package initialize

import (
	"github.com/go-redis/redis"
	"xmlt/global"
)

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ping := client.Ping()
	err := ping.Err()
	if err != nil {
		panic(err)
	}
	global.Redis = client

}
