package routes

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"redis_based_url_shortener/database"
)

var ctx = context.Background()

func ResolveURL(ctx fiber.Ctx) error {

	url := ctx.Params("url")

	if url == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "empty url"})
	}

	c := database.NewClient(0)
	defer c.Close()
	result, err := c.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "coudln't get url short"})

	}
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "coudln't get url short internal error "})

	}

	rInr := database.NewClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	return ctx.Redirect().To(result)

}
