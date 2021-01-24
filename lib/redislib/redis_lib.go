package redislib

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
)

// InitClient initializes the redis client.
func InitClient() {
	if client != nil {
		return
	}

	client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolsize"),
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
	})
}

// GetClient returns the redis client.
func GetClient() *redis.Client {
	if client == nil {
		InitClient()
	}

	return client
}
