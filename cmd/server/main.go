// cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/shravanirajulu2004/go-user-api/config"
	"github.com/shravanirajulu2004/go-user-api/internal/handler"
	"github.com/shravanirajulu2004/go-user-api/internal/logger"
	"github.com/shravanirajulu2004/go-user-api/internal/middleware"
	"github.com/shravanirajulu2004/go-user-api/internal/repository"
	"github.com/shravanirajulu2004/go-user-api/internal/routes"
	"github.com/shravanirajulu2004/go-user-api/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting application",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.Port),
	)

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Log.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Log.Info("Successfully connected to database")

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, logger.Log)
	userHandler := handler.NewUserHandler(userService, logger.Log)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "User API",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(cors.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.RecoveryMiddleware(logger.Log))
	app.Use(middleware.LoggerMiddleware(logger.Log))

	// Setup routes
	routes.SetupRoutes(app, userHandler)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		logger.Log.Info("Server starting", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Log.Error("Server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Server stopped")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}