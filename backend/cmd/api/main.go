package main

import (
	"log"
	"os"

	"travel-itinerary-d1/backend/internal/config"
	"travel-itinerary-d1/backend/internal/handler"
	"travel-itinerary-d1/backend/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db := config.ConnectDB()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Itinerary{},
		&model.Spot{},
		&model.SpotNote{},
	); err != nil {
		log.Fatal("Gagal migrasi tabel: ", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{allowedOrigin},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Server Travel Itinerary D1 berjalan",
		})
	})

	itineraryHandler := handler.NewItineraryHandler(db)
	itineraries := api.Group("/itineraries")
	itineraries.Get("/", itineraryHandler.GetAll)
	itineraries.Get("/:id", itineraryHandler.GetByID)
	itineraries.Post("/", itineraryHandler.Create)
	itineraries.Put("/:id", itineraryHandler.Update)
	itineraries.Delete("/:id", itineraryHandler.Delete)

	spotNoteHandler := handler.NewSpotNoteHandler(db)
	spotNotes := api.Group("/spot-notes")
	spotNotes.Get("/", spotNoteHandler.GetAll)
	spotNotes.Get("/:id", spotNoteHandler.GetByID)
	spotNotes.Post("/", spotNoteHandler.Create)
	spotNotes.Put("/:id", spotNoteHandler.Update)
	spotNotes.Delete("/:id", spotNoteHandler.Delete)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}