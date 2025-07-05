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
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database Connection
	dbCfg := config.AppConfig.Database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.SSLMode)

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Initialize a new Fiber app
	app := fiber.New()

	// Enable CORS (allow all origins by default; adjust via env if needed)
	app.Use(cors.New())

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
	log.Fatal(app.Listen(":" + config.AppConfig.Server.Port))
}
