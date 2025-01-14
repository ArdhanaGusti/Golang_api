package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	err := RDB.Ping().Err()
	if err != nil {
		panic("Failed to connect redis because " + err.Error())
	}
	fmt.Println("Connected to Redis!")
}
