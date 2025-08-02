package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trueegorletov/analabit/service/idmsu/cache"
	"github.com/trueegorletov/analabit/service/idmsu/handler"
	"github.com/trueegorletov/analabit/service/idmsu/resolver"
)

func main() {
	log.Println("Starting idmsu service...")

	// Initialize cache
	cacheClient, err := cache.NewPostgresCache()
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}
	defer cacheClient.Close()

	// Ensure database tables exist with proper schema
	if err := cacheClient.EnsureTables(); err != nil {
		log.Fatalf("Failed to ensure database tables: %v", err)
	}

	// Initialize resolver
	resolver := resolver.NewMSUResolver(cacheClient)

	// Initialize handler
	handler := handler.NewHandler(resolver)

	// Setup router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/resolve", handler.ResolveBatch)
		v1.GET("/health", handler.Health)
		v1.GET("/ready", handler.Ready)
	}

	// Start background fetcher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go resolver.StartBackgroundFetcher(ctx)

	// Start server
	port := os.Getenv("IDMSU_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("idmsu service listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down idmsu service...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("idmsu service stopped")
}
