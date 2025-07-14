package main

import (
	"analabit/service/api/config"
	"analabit/service/api/handlers"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/migrate"

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

	// Run database migrations
	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Create materialized view for application flags
	createViewQuery := `
CREATE MATERIALIZED VIEW IF NOT EXISTS application_flags AS
SELECT
  a.id AS application_id,
  a.run_id,
  a.student_id,
  a.heading_id,
  a.priority,
  a.original_submitted,
  (SELECT COUNT(*) FROM application a2
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.priority < a.priority
     AND a2.heading_id != a.heading_id
     AND EXISTS (SELECT 1 FROM calculation c
                 WHERE c.student_id = a2.student_id
                   AND c.heading_id = a2.heading_id
                   AND c.run_id = a2.run_id)) AS passing_to_more_priority,
  EXISTS (SELECT 1 FROM calculation c
          WHERE c.student_id = a.student_id
            AND c.heading_id = a.heading_id
            AND c.run_id = a.run_id) AS passing_now,
  (SELECT COUNT(*) FROM application a2
   JOIN heading h2 ON a2.heading_id = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND h2.varsity_id != (SELECT varsity_id FROM heading h3 WHERE h3.id = a.heading_id)) AS another_varsities_count
FROM application a;

CREATE UNIQUE INDEX IF NOT EXISTS application_flags_pkey ON application_flags (application_id);
`
	if _, err := client.ExecContext(context.Background(), createViewQuery); err != nil {
		log.Fatalf("failed creating materialized view: %v", err)
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
