package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"log"
	"redis_based_url_shortener/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	panic(err)
	//
	//}
	//app := fiber.New()
	//
	//app.Use(logger.New())
	//
	//setupRoutes(app)
	//c := database.NewClient(0)
	//_ = c
	//if err != nil {
	//	panic(err)
	//}
	//log.Fatal(app.Listen(os.Getenv("APP_PORT")))
	log.Println("HELO")
	cli := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	res, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("It's all right!!!!!")
	log.Println(res)
}
