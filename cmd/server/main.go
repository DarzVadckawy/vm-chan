package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vm-chan/internal/config"
	"vm-chan/internal/domain"
	"vm-chan/internal/handler"
	"vm-chan/internal/middleware"
	"vm-chan/internal/repository"
	"vm-chan/internal/service"

	_ "vm-chan/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title VM-Chan Text Analysis API
// @version 1.0
// @description VM-Chan is a production-ready microservice for analyzing text sentences and extracting linguistic statistics.
//
// @description ## Features
// @description - **Text Analysis**: Extract word count, vowel count, and consonant count from sentences
// @description - **JWT Authentication**: Secure API access with Bearer token authentication
// @description - **Production Ready**: Includes logging, metrics, health checks, and graceful shutdown
// @description - **OpenAPI Documentation**: Interactive API documentation with examples
//
// @description ## Quick Start
// @description 1. **Authenticate**: POST to `/auth/login` with credentials `{"username": "admin", "password": "password"}`
// @description 2. **Get Token**: Extract the JWT token from the response
// @description 3. **Analyze Text**: POST to `/api/v1/analyze` with Authorization header `Bearer <token>`
//
// @description ## Example Usage
// @description ```bash
// @description # Login and get token
// @description curl -X POST http://localhost:8080/auth/login \
// @description   -H "Content-Type: application/json" \
// @description   -d '{"username": "admin", "password": "password"}'
// @description
// @description # Use token to analyze text
// @description curl -X POST http://localhost:8080/api/v1/analyze \
// @description   -H "Authorization: Bearer <your-token>" \
// @description   -H "Content-Type: application/json" \
// @description   -d '{"sentence": "Hello world!"}'
// @description ```
//
// @description ## Response Format
// @description All responses are in JSON format. Error responses include an error field with a descriptive message.
//
// @host localhost:8080
// @BasePath /
//
// @contact.name API Support
// @contact.url https://github.com/your-username/vm-chan
// @contact.email support@vm-chan.io
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @termsOfService http://swagger.io/terms/
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token
//
// @tag.name auth
// @tag.description Authentication operations
//
// @tag.name Text Analysis
// @tag.description Text analysis operations
//
// @externalDocs.description GitHub Repository
// @externalDocs.url https://github.com/your-username/vm-chan

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	logger := initLogger(cfg.Logging.Level)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	logger.Info("Starting VM-Chan microservice", zap.String("version", "1.0.0"))

	// Initialize dependencies
	userRepo := repository.NewUserRepository(logger)
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, logger)
	textAnalysisService := service.NewTextAnalysisService(logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, logger)
	textAnalysisHandler := handler.NewTextAnalysisHandler(textAnalysisService, logger)

	// Setup router
	router := setupRouter(cfg, logger, authHandler, textAnalysisHandler, authService)

	// Setup server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func setupRouter(cfg *config.Config, logger *zap.Logger, authHandler *handler.AuthHandler,
	textAnalysisHandler *handler.TextAnalysisHandler, authService domain.AuthService) *gin.Engine {

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityMiddleware())

	if cfg.Metrics.Enabled {
		router.Use(middleware.MetricsMiddleware())
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "vm-chan",
			"version":   "1.0.0",
		})
	})

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Protected API routes
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(middleware.AuthMiddleware(authService, logger))
	{
		apiGroup.POST("/analyze", textAnalysisHandler.AnalyzeText)
	}

	return router
}

func initLogger(level string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	return logger
}
