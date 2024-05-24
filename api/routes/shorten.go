package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"os"
	"redis_based_url_shortener/database"
	"redis_based_url_shortener/helpers"
	"strconv"
	"time"
)

type Response struct {
	Url             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemain     int           `json:"x_rate_remain"`
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

type Request struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

func ShortenURL(ctx fiber.Ctx) error {
	body := new(Request)
	databaseCtx := database.Ctx
	if err := ctx.Bind().Body(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "couldn't parse json"})
	}

	r2 := database.NewClient(1)
	defer r2.Close()
	result, err := r2.Get(databaseCtx, ctx.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(databaseCtx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error in reddis"})

	}
	i, _ := strconv.Atoi(result)
	if i <= 0 {
		limit := r2.TTL(databaseCtx, ctx.IP())
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate time exceded", "rate_limit__rest": limit})

	}

	if !govalidator.IsURL(body.Url) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})

	}

	if !helpers.RemoveDomainError(body.Url) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "unavailable!!!"})

	}

	body.Url = helpers.EnforceHTTP(body.Url)

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r1 := database.NewClient(0)
	defer r1.Close()

	val, _ := r1.Get(databaseCtx, id).Result()
	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "already in use"})

	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = r1.Set(databaseCtx, id, body.Url, body.Expiry).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unaelr to connect to server "})

	}

	r2.Decr(databaseCtx, ctx.IP())

	resp := Response{
		Url:             body.Url,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemain:     10,
		XRateLimitReset: 30,
	}

	ttl, _ := r1.TTL(databaseCtx, ctx.IP()).Result()
	resp.XRateLimitReset = ttl

	val, _ = r1.Get(databaseCtx, ctx.IP()).Result()
	resp.XRateRemain, _ = strconv.Atoi(val)

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(resp)

}
