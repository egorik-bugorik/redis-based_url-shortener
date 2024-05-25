package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

var Ctx = context.Background()

func NewClient(dbNum int) *redis.Client {
	println(os.Getenv("DB_PORT"))
	println(os.Getenv("DB_PASS"))
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_PORT"),
		DB:       dbNum,
		Password: os.Getenv("DB_PASS"),
	})

}
