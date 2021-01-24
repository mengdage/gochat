package cache

import (
	"fmt"

	"githum.com/mengdage/gochat/lib/redislib"
	"githum.com/mengdage/gochat/model"
)

const (
	serverHashKey        = "acc:hash:server"
	serversHashCacheTime = 2 * 60 * 60
	serversHashTimeout   = 3 * 60
)

// SetServerInfo stores the server info in the cache.
func SetServerInfo(server *model.Server, currentTime int64) (err error) {
	redisClient := redislib.GetClient()
	num, err := redisClient.Do("hSet", serverHashKey, server.String(), currentTime).Int()
	if err != nil {
		fmt.Printf("Failed to set server info: %s\n", err.Error())
		return
	}

	if num != 1 {
		return
	}

	redisClient.Do("Expire", serverHashKey, serversHashTimeout)
	return
}
