package bootstrap

import (
	"context"
	"go-clean-api/api/controller"
	"go-clean-api/api/middleware"
	"go-clean-api/api/repositories"
	"go-clean-api/api/routes"
	"go-clean-api/api/service"

	"go-clean-api/infrastructure"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	infrastructure.Module,
	routes.Module,
	controller.Module,
	service.Module,
	repositories.Module,
	middleware.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	routes routes.Routes,
	middlewares middleware.Middleware,
	migrations infrastructure.Migration,
	handler infrastructure.RequestHandler,
	env infrastructure.Env,
	logger infrastructure.Logger,
	database infrastructure.Database,
) {
	conn, _ := database.DB.DB()
	_, cancel := context.WithCancel(context.Background())
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("---------------------------")
			logger.Zap.Info("----- go-clean-API ♾️  -----")
			logger.Zap.Info("---------------------------")
			conn.SetMaxOpenConns(10)

			go func() {
				logger.Zap.Info("Migrating DB schema...")
				migrations.Migrate()

				// logger.Zap.Info("Attaching firestore listeners")
				// listeners.Attach(ctx)
				logger.Zap.Info("Middleware setup...")
				middlewares.SetUp()

				logger.Zap.Info(" Routes setup...")
				routes.Setup()

				// logger.Zap.Info("🌱 seeding data...")

				// seeds.Run()

				// logger.Zap.Info(" firebase setup...")
				// firebase.

				if env.ENVIRONMENT != "local" {
					_ = handler.Gin.Run()
				} else {
					_ = handler.Gin.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Zap.Info("Stopping Application")
			_ = conn.Close()
			cancel()
			return nil
		},
	})
}
