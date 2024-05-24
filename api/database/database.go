package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

var Ctx = context.Background()

func NewClient(dbNum int) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		DB:       dbNum,
		Password: os.Getenv("DB_PASS"),
	})

}
