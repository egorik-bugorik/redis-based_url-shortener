package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
	"redis_based_url_shortener/database"
	"redis_based_url_shortener/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)

	}
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)
	c := database.NewClient(0)
	_ = c
	if err != nil {
		panic(err)
	}
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
