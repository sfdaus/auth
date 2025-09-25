package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"prakarsa-app/utils/crypto"
	"prakarsa-app/utils/jwt"

	_ "prakarsa-app/docs"
	"prakarsa-app/utils"

	"prakarsa-app/config"
	httpDelivery "prakarsa-app/delivery/http"
	appMiddleware "prakarsa-app/delivery/middleware"
	"prakarsa-app/infrastructure/datastore"
	filestorage "prakarsa-app/infrastructure/filestorage"
	pgsqlRepository "prakarsa-app/repository/pgsql"
	"prakarsa-app/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Auth API
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization
func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	fileStorageInstance, err := filestorage.NewFileStorage(configApp)
	utils.PanicIfNeeded(err)

	// cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	// redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	userRepo := pgsqlRepository.NewPgsqlUserRepository(dbInstance)
	authTokenRepo := pgsqlRepository.NewPgsqlAuthTokenRepository(dbInstance)
	profileRepo := pgsqlRepository.NewPgsqlProfileRepository(dbInstance)

	// Setup Service
	cryptoSvc := crypto.NewCryptoService()
	jwtSvc := jwt.NewJWTService(configApp.JWTSecretKey)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	signUpUC := usecase.SignUpUsecase(userRepo, authTokenRepo, profileRepo, cryptoSvc, jwtSvc, ctxTimeout)
	signInUC := usecase.SignInUsecase(userRepo, cryptoSvc, jwtSvc, ctxTimeout)
	profileUC := usecase.ProfileUsecase(profileRepo, ctxTimeout, fileStorageInstance)

	// Setup app middleware
	appMiddleware := appMiddleware.NewMiddleware(jwtSvc)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	// e.Use(appMiddleware.Logger(nil))
	e.Use(appMiddleware.CustomLogger())
	e.Logger.Info("🚀 Server is alive and running")

	// Setup handler
	e.GET("/api/v1/auth/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	httpDelivery.NewAuthHandler(e, appMiddleware, signUpUC, signInUC, profileUC)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configApp.ContextTimeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
