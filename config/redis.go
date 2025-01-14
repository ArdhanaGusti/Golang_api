package config

import (
	"fmt"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err := RDB.Ping().Err()
	if err != nil {
		panic("Failed to connect redis because " + err.Error())
	}
	fmt.Println("Connected to Redis!")
}
