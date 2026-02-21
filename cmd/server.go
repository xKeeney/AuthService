package cmd

import (
	"auth_service/internal/config"
	"auth_service/internal/database"
	"fmt"
	"log"

	"github.com/xKeeney/httpForge"
	"github.com/xKeeney/httpForge/httpLogger"
	"github.com/xKeeney/httpForge/httpMiddlewares"
)

func StartServer() {
	// Init configs
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	appLogger := httpLogger.NewHttpLogger(cfg.Logger.LogFile, cfg.Logger.LogLevel)

	// Init db
	db, err := database.InitGormPostgresql(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	if err != nil {
		appLogger.Fatal(err)
	}

	// Migrations
	if err := database.StartMigrations(db); err != nil {
		appLogger.Fatal(err)
	}

	// Init app
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	app := httpForge.NewHttpApp(addr, appLogger)

	// Init base middlewares
	baseMiddlewares := httpMiddlewares.InitBaseMiddlewares(appLogger)

	app.AddMiddleware(baseMiddlewares.InfoMiddleware)
	app.AddMiddleware(baseMiddlewares.RequestsLoggerMiddleware)

	// Start server
	app.ListenAndServe()
}
