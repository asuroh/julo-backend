package bootstrap

import (
	"julo-backend/pkg/str"
	appMW "julo-backend/server/middleware"
)

// RegisterMiddleware ...
func (boot *Bootup) RegisterMiddleware() {
	recovery := appMW.RecoverInit{
		Debug: str.StringToBool(boot.EnvConfig["APP_DEBUG"]),
	}
	boot.R.Use(recovery.Recoverer)

	boot.R.Use(appMW.NotfoundMiddleware)
}
