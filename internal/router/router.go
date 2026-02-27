package router

import (
	"auth_service/internal/auth"

	"github.com/xKeeney/httpForge"
	"github.com/xKeeney/httpForge/httpLogger"
	"gorm.io/gorm"
)

func AddRoutes(app *httpForge.HttpApp, db *gorm.DB, appLogger *httpLogger.HttpLogger) {
	/* AUTH */
	authRepo := auth.InitAuthRepository(db)
	authService := auth.InitAuthService(authRepo)
	authHandler := auth.InitAuthHandler(authService, appLogger)

	// routes
	auth := app.NewRouter("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
}
