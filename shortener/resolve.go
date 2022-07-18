package shortener

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rushi3691/url_shortener/database"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")
	rdb0 := database.CreateClient(0)
	defer rdb0.Close()

	value, err := rdb0.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found on database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",
		})
	}
	return c.Redirect(value, 301)
}
