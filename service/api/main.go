package main

import (
	"analabit/service/api/config"
	"analabit/service/api/handlers"
	"fmt"
	"log"

	"analabit/core/ent"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	_ "github.com/lib/pq"
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

	// Initialize a new Fiber app
	app := fiber.New()

	// Enable CORS (allow all origins by default; adjust via env if needed)
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1 Group
	api := app.Group("/api")

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
