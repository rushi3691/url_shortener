package shortener

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rushi3691/url_shortener/database"
	"github.com/rushi3691/url_shortener/helpers"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	RateRemaining  int           `json:"rate_limit"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx) error {
	reqBody := new(Request)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	rdb1 := database.CreateClient(1)
	defer rdb1.Close()
	val, err := rdb1.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = rdb1.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		return c.JSON(err)
	} else {
		intVal, _ := strconv.Atoi(val)
		if intVal <= 0 {
			limit, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}
	if !govalidator.IsURL(reqBody.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}
	if !helpers.RemoveDomainError(reqBody.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "can't shorten this domain",
		})
	}

	reqBody.URL = helpers.EnforceHTTP(reqBody.URL)
	var id string
	if reqBody.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = reqBody.CustomShort
	}

	rdb0 := database.CreateClient(0)
	defer rdb0.Close()

	val, _ = rdb0.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL short already in use",
		})
	}

	if reqBody.Expiry == 0 {
		reqBody.Expiry = 24
	}

	err = rdb0.Set(database.Ctx, id, reqBody.URL, reqBody.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}
	resp := response{
		URL:            reqBody.URL,
		CustomShort:    "",
		Expiry:         reqBody.Expiry,
		RateRemaining:  10,
		RateLimitReset: 30,
	}
	newVal, _ := rdb1.Decr(database.Ctx, c.IP()).Result()
	resp.RateRemaining = int(newVal)
	ttl, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
	resp.RateLimitReset = ttl / time.Nanosecond / time.Minute
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
