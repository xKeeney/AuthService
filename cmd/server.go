package cmd

import (
	"github.com/xKeeney/httpForge"
	"github.com/xKeeney/httpForge/httpLogger"
	"github.com/xKeeney/httpForge/httpMiddlewares"
)

func StartServer() {
	// Init logger
	appLogger := httpLogger.NewHttpLogger("app.log", httpLogger.TRACE)

	// Init app
	app := httpForge.NewHttpApp(":8080", appLogger)

	// Init base middlewares
	baseMiddlewares := httpMiddlewares.InitBaseMiddlewares(appLogger)

	app.AddMiddleware(baseMiddlewares.InfoMiddleware)
	app.AddMiddleware(baseMiddlewares.RequestsLoggerMiddleware)

	// Start server
	app.ListenAndServe()
}
