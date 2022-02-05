package main

import (
	"encoding/json"
	"flag"
	"julo-backend/model"
	"julo-backend/pkg/aes"
	amqpPkg "julo-backend/pkg/amqp"
	"julo-backend/pkg/amqpconsumer"
	"julo-backend/pkg/env"
	"julo-backend/pkg/logruslogger"
	"julo-backend/pkg/pg"
	"julo-backend/pkg/str"
	"julo-backend/usecase"
	"julo-backend/usecase/viewmodel"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/streadway/amqp"
)

var (
	uri          *string
	formURL      = flag.String("form_url", "http://localhost", "The URL that requests are sent to")
	logFile      = flag.String("log_file", "system.log", "The file where errors are logged")
	threads      = flag.Int("threads", 1, "The max amount of go routines that you would like the process to use")
	maxprocs     = flag.Int("max_procs", 1, "The max amount of processors that your application should use")
	paymentsKey  = flag.String("payments_key", "secret", "Access key")
	exchange     = flag.String("exchange", amqpPkg.UpdateBalanceExchange, "The exchange we will be binding to")
	exchangeType = flag.String("exchange_type", "direct", "Type of exchange we are binding to | topic | direct| etc..")
	queue        = flag.String("queue", amqpPkg.UpdateBalance, "Name of the queue that you would like to connect to")
	routingKey   = flag.String("routing_key", amqpPkg.UpdateBalanceDeadLetter, "queue to route messages to")
	workerName   = flag.String("worker_name", "worker.name", "name to identify worker by")
	verbosity    = flag.Bool("verbos", false, "Set true if you would like to log EVERYTHING")

	// Hold consumer so our go routine can listen to
	// it's done error chan and trigger reconnects
	// if it's ever returned
	conn      *amqpconsumer.Consumer
	envConfig map[string]string
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*maxprocs)
	envConfig = env.NewEnvConfig("../.env")
	uri = flag.String("uri", envConfig["AMQP_URL"], "The rabbitmq endpoint")
}

func main() {
	file := false
	// Open a system file to start logging to
	if file {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			log.Printf("error opening file: %v", err.Error())
		}
		log.SetOutput(f)
	}

	conn := amqpconsumer.NewConsumer(*workerName, *uri, *exchange, *exchangeType, *queue)

	if err := conn.Connect(); err != nil {
		log.Printf("Error: %v", err)
	}

	deliveries, err := conn.AnnounceQueue(*queue, *routingKey)
	if err != nil {
		log.Printf("Error when calling AnnounceQueue(): %v", err.Error())
	}

	// Postgre DB connection
	dbInfo := pg.Connection{
		Host:    envConfig["DATABASE_HOST"],
		DB:      envConfig["DATABASE_DB"],
		User:    envConfig["DATABASE_USER"],
		Pass:    envConfig["DATABASE_PASSWORD"],
		Port:    str.StringToInt(envConfig["DATABASE_PORT"]),
		SslMode: "disable",
	}
	db, err := dbInfo.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig["REDIS_HOST"],
		Password: envConfig["REDIS_PASSWORD"],
		DB:       0,
	})
	pong, err := redisClient.Ping().Result()
	log.Println("Redis ping status: "+pong, err)

	// AES credential
	aesCredential := aes.Credential{
		Key: envConfig["AES_KEY"],
	}

	cUC := usecase.ContractUC{
		DB:        db,
		Redis:     redisClient,
		EnvConfig: envConfig,
		Aes:       aesCredential,
	}

	conn.Handle(deliveries, handler, *threads, *queue, *routingKey, cUC)
}

func handler(deliveries <-chan amqp.Delivery, uc *usecase.ContractUC) {
	var (
		ctx  = "UpdateBalanceListener"
		rand = rand.Intn(5-1) + 1
	)

	for d := range deliveries {
		var formData map[string]interface{}

		err := json.Unmarshal(d.Body, &formData)
		if err != nil {
			log.Printf("Error unmarshaling data: %s", err.Error())
		}

		uc.ReqID = formData["qid"].(string)
		time.Sleep(time.Second * time.Duration(rand))

		tx := model.SQLDBTx{DB: uc.DB}
		txFunc, err := tx.TxBegin()
		txDB := txFunc.DB
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "rejected", formData["qid"].(string))
			d.Reject(false)
		}
		walletUc := usecase.WalletUC{ContractUC: uc}
		err = walletUc.AddBalance(viewmodel.SendQueue{
			Amount:    int(formData["amount"].(float64)),
			Type:      formData["type"].(string),
			OwnedBy:   formData["owned_by"].(string),
			BalanceID: formData["balance_id"].(string),
		})
		if err != nil {
			txDB.Rollback()
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "err", formData["qid"].(string))

			// Get fail counter from redis
			failCounter := amqpconsumer.FailCounter{}
			err = uc.GetFromRedis("amqpFail"+formData["qid"].(string), &failCounter)
			if err != nil {
				failCounter = amqpconsumer.FailCounter{
					Counter: 1,
				}
			}

			if failCounter.Counter > amqpconsumer.MaxFailCounter {
				logruslogger.Log(logruslogger.WarnLevel, strconv.Itoa(failCounter.Counter), ctx, "rejected", formData["qid"].(string))
				d.Reject(false)
			} else {
				// Save the new counter to redis
				failCounter.Counter++
				err = uc.StoreToRedisExp("amqpFail"+formData["qid"].(string), failCounter, "10m")

				logruslogger.Log(logruslogger.WarnLevel, strconv.Itoa(failCounter.Counter), ctx, "failed", formData["qid"].(string))
				d.Nack(false, true)
			}
		} else {
			txDB.Commit()
			logruslogger.Log(logruslogger.InfoLevel, string(d.Body), ctx, "success", formData["qid"].(string))
			d.Ack(false)
		}
	}

	return
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
