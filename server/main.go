package main

import (
	"julo-backend/pkg/aes"
	"julo-backend/pkg/env"
	"julo-backend/pkg/interfacepkg"
	"julo-backend/pkg/jwe"
	"julo-backend/pkg/jwt"
	"julo-backend/pkg/logruslogger"
	"julo-backend/pkg/pg"
	"julo-backend/pkg/str"
	boot "julo-backend/server/bootstrap"
	"julo-backend/usecase"

	"github.com/rs/xid"

	"github.com/rs/cors"

	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v7"
	validator "gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	idTranslations "gopkg.in/go-playground/validator.v9/translations/id"
)

var (
	_, b, _, _      = runtime.Caller(0)
	basepath        = filepath.Dir(b)
	debug           = false
	host            string
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	envConfig       map[string]string
	corsDomainList  []string
)

// Init first time running function
func init() {
	// Load env variable from .env file
	envConfig = env.NewEnvConfig("../.env")

	host = envConfig["APP_HOST"]
	if str.StringToBool(envConfig["APP_DEBUG"]) {
		debug = true
		log.Printf("Running on Debug Mode: On at host [%v]", host)
	}
}

func main() {
	ctx := "main"

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig["REDIS_HOST"],
		Password: envConfig["REDIS_PASSWORD"],
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	// Postgre DB connection
	dbInfo := pg.Connection{
		Host:    envConfig["DATABASE_HOST"],
		DB:      envConfig["DATABASE_DB"],
		User:    envConfig["DATABASE_USER"],
		Pass:    envConfig["DATABASE_PASSWORD"],
		Port:    str.StringToInt(envConfig["DATABASE_PORT"]),
		SslMode: envConfig["DATABASE_SSL_MODE"],
	}
	db, err := dbInfo.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// JWT credential
	jwtCredential := jwt.Credential{
		Secret:           envConfig["TOKEN_SECRET"],
		ExpSecret:        str.StringToInt(envConfig["TOKEN_EXP_SECRET"]),
		RefreshSecret:    envConfig["TOKEN_REFRESH_SECRET"],
		RefreshExpSecret: str.StringToInt(envConfig["TOKEN_EXP_REFRESH_SECRET"]),
	}

	// JWE credential
	jweCredential := jwe.Credential{
		KeyLocation: envConfig["APP_PRIVATE_KEY_LOCATION"],
		Passphrase:  envConfig["APP_PRIVATE_KEY_PASSPHRASE"],
	}

	// AES credential
	aesCredential := aes.Credential{
		Key: envConfig["AES_KEY"],
	}

	// Validator initialize
	validatorInit()

	// Load contract struct
	contractUC := usecase.ContractUC{
		ReqID:     xid.New().String(),
		DB:        db,
		Redis:     redisClient,
		EnvConfig: envConfig,
		Jwt:       jwtCredential,
		Jwe:       jweCredential,
		Aes:       aesCredential,
	}

	r := chi.NewRouter()
	// Cors setup
	r.Use(cors.New(cors.Options{
		AllowedOrigins: corsDomainList,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler)

	// load application bootstrap
	bootApp := boot.Bootup{
		R:          r,
		CorsDomain: corsDomainList,
		EnvConfig:  envConfig,
		DB:         db,
		Redis:      redisClient,
		Validator:  validatorDriver,
		Translator: translator,
		ContractUC: contractUC,
	}

	// register middleware
	bootApp.RegisterMiddleware()

	// register routes
	bootApp.RegisterRoutes()

	// Log start server
	startBody := map[string]interface{}{
		"Host":     host,
		"Location": str.DefaultData(envConfig["APP_DEFAULT_LOCATION"], "Asia/Jakarta"),
	}
	logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshall(startBody), ctx, "server_start", "")

	// Run the app
	http.ListenAndServe(host, r)
}

func validatorInit() {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch envConfig["APP_LOCALE"] {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}
