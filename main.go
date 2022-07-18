package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rushi3691/url_shortener/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
