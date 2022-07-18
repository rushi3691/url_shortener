package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rushi3691/url_shortener/shortener"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/:url", shortener.ResolveUrl)
	app.Post("/api", shortener.ShortenUrl)
}
