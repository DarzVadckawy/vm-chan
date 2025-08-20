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

	_ "vm-chan/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger := initLogger(cfg.Logging.Level)
	defer func(logger *zap.Logger) {
		if syncErr := logger.Sync(); syncErr != nil {
			fmt.Printf("Error syncing logger: %v\n", syncErr)
		}
	}(logger)

	logger.Info("Starting VM-Chan microservice", zap.String("version", "1.0.0"))

	userRepo := repository.NewUserRepository(logger)
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, logger)
	textAnalysisService := service.NewTextAnalysisService(logger)

	authHandler := handler.NewAuthHandler(authService, logger)
	textAnalysisHandler := handler.NewTextAnalysisHandler(textAnalysisService, logger)

	router := setupRouter(cfg, logger, authHandler, textAnalysisHandler, authService)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func setupRouter(
	cfg *config.Config,
	logger *zap.Logger,
	authHandler *handler.AuthHandler,
	textAnalysisHandler *handler.TextAnalysisHandler,
	authService domain.AuthService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityMiddleware())

	if cfg.Metrics.Enabled {
		router.Use(middleware.MetricsMiddleware())
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "vm-chan",
			"version":   "1.0.0",
		})
	})

	if cfg.Metrics.Enabled {
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	authGroup.POST("/login", authHandler.Login)

	apiGroup := router.Group("/api/v1")
	apiGroup.Use(middleware.AuthMiddleware(authService, logger))
	apiGroup.POST("/analyze", textAnalysisHandler.AnalyzeText)

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
