package bootstrap

import (
	"julo-backend/pkg/logruslogger"
	api "julo-backend/server/handler"
	"julo-backend/server/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RegisterRoutes ...
func (boot *Bootup) RegisterRoutes() {
	handlerType := api.Handler{
		DB:         boot.DB,
		EnvConfig:  boot.EnvConfig,
		Validate:   boot.Validator,
		Translator: boot.Translator,
		ContractUC: &boot.ContractUC,
	}
	mJwt := middleware.VerifyMiddlewareInit{
		ContractUC: &boot.ContractUC,
	}

	boot.R.Route("/api", func(r chi.Router) {
		// Define a limit rate to 1000 requests per IP per request.
		rate, _ := limiter.NewRateFromFormatted("1000-S")
		store, _ := sredis.NewStoreWithOptions(boot.ContractUC.Redis, limiter.StoreOptions{
			Prefix:   "limiter_global",
			MaxRetry: 3,
		})
		rateMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
		r.Use(rateMiddleware.Handler)

		// Logging setup
		r.Use(chimiddleware.RequestID)
		r.Use(logruslogger.NewStructuredLogger(boot.EnvConfig["LOG_FILE_PATH"], boot.EnvConfig["LOG_DEFAULT"], boot.ContractUC.ReqID))
		r.Use(chimiddleware.Recoverer)

		// API
		r.Route("/v1", func(r chi.Router) {
			walletHandler := api.WalletHandler{Handler: handlerType}
			r.Group(func(r chi.Router) {
				r.Use(mJwt.VerifyBasicAuth)
				r.Post("/init", walletHandler.InitHandler)
			})

			balanceHandler := api.BalanceHandler{Handler: handlerType}
			r.Route("/wallet", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyTokenCredential)
					r.Post("/", walletHandler.EnableHandler)
					r.Get("/", walletHandler.GetWalletHandler)
					r.Post("/deposits", balanceHandler.DepositHandler)
					r.Post("/withdrawals", balanceHandler.WithdrawalHandler)
					r.Patch("/", walletHandler.DisableHandler)
				})
			})
		})
	})
}
