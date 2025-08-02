package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trueegorletov/analabit/service/idmsu/cache"
	"github.com/trueegorletov/analabit/service/idmsu/handler"
	"github.com/trueegorletov/analabit/service/idmsu/resolver"
)

func main() {
	slog.Info("Starting sophisticated IDMSU service with persistent caching...")

	// Build database connection string
	dbHost := getEnvOrDefault("DATABASE_HOST", "localhost")
	dbPort := getEnvOrDefault("DATABASE_PORT", "5433")
	dbUser := getEnvOrDefault("DATABASE_USER", "postgres")
	dbPassword := getEnvOrDefault("DATABASE_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DATABASE_DBNAME", "postgres")
	dbSSLMode := getEnvOrDefault("DATABASE_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Initialize database store
	dbStore, err := cache.NewDatabaseStore(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database store: %v", err)
	}
	defer dbStore.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := dbStore.Health(ctx); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	slog.Info("Database connection established")

	// Initialize layered cache with database persistence
	memCache := cache.NewMemoryCache(30 * time.Minute)
	layered := cache.NewLayeredCache(memCache, dbStore)

	// Initialize resolver with both caches
	res := resolver.NewMSUResolver(layered, dbStore)
	h := handler.NewHandler(res, res.HasRecentData)

	// Start background fetcher
	fetchCtx, fetchCancel := context.WithCancel(context.Background())
	defer fetchCancel()
	res.StartBackgroundFetcher(fetchCtx)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add request logging middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[IDMSU] %s - %s \"%s %s %s\" %d %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
		)
	}))

	// Wire API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/resolve", h.ResolveBatch)
		v1.GET("/health", h.Health)
		v1.GET("/ready", h.Ready)
		v1.GET("/wait", h.Wait)
	}

	port := getEnvOrDefault("IDMSU_PORT", "8081")
	slog.Info("Starting IDMSU server", "port", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
