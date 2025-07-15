package main

import (
	"analabit/service/api/config"
	"analabit/service/api/handlers"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/trueegorletov/analabit/core/database"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/migrate"
	"github.com/trueegorletov/analabit/core/metrics"
	"github.com/trueegorletov/analabit/core/migrations"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database Connection
	cfg := config.AppConfig
	var connStr string

	// Use legacy connection string if provided, otherwise build from individual fields
	if cfg.PostgresConnStrings != "" {
		connStr = cfg.PostgresConnStrings
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseDBName, cfg.DatabaseSSLMode)
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Initialize metrics
	metrics.InitMetrics()
	log.Println("Metrics initialized")

	// Run database migrations
	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Create database client wrapper
	dbClient, err := database.NewClient(client)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}

	// Run custom migrations
	migrationRunner := migrations.NewMigrationRunner(dbClient)
	if err := migrationRunner.Run(context.Background()); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	// Create Prometheus metrics
	requestTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"method", "endpoint"},
	)

	// Register metrics with Prometheus
	prometheus.MustRegister(requestTotal)
	prometheus.MustRegister(requestDuration)

	// Initialize a new Fiber app
	app := fiber.New()

	// Custom Prometheus middleware
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()

		// Continue to next handler
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()
		endpoint := c.Route().Path

		requestTotal.WithLabelValues(method, endpoint, statusCode).Inc()
		requestDuration.WithLabelValues(method, endpoint).Observe(duration)

		return err
	})

	// Prometheus metrics endpoint
	app.Get("/metrics", func(c fiber.Ctx) error {
		registry := prometheus.DefaultGatherer
		metricFamilies, err := registry.Gather()
		if err != nil {
			return c.Status(500).SendString("Error gathering metrics")
		}

		c.Set("Content-Type", string(expfmt.FmtText))

		encoder := expfmt.NewEncoder(c, expfmt.FmtText)
		for _, mf := range metricFamilies {
			if err := encoder.Encode(mf); err != nil {
				return c.Status(500).SendString("Error encoding metrics")
			}
		}

		return nil
	})

	// Enable CORS (allow all origins by default; adjust via env if needed)
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1 Group
	api := app.Group("/v1")

	// Routes
	api.Get("/varsities", handlers.GetVarsities(client))
	api.Get("/headings", handlers.GetHeadings(client))
	api.Get("/headings/:id", handlers.GetHeadingByID(client))
	api.Get("/applications", handlers.GetApplications(client))
	api.Get("/students/:id", handlers.GetStudentByID(client))
	api.Get("/results", handlers.GetResults(client))



	// Start the server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
