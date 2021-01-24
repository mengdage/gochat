package cache

import (
	"encoding/json"
	"fmt"
	"log"

	"githum.com/mengdage/gochat/lib/redislib"
	"githum.com/mengdage/gochat/model"
)

const (
	userOnlinePrefix    = "user:online:"
	userOnlineCacheTime = 24 * 60 * 60
)

func getUserOnlineKey(userKey string) string {
	return fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
}

func GetUserOnlineInfo(userID string) (*model.UserOnline, error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userID)

	bs, err := redisClient.Get(key).Bytes()
	if err != nil {
		log.Printf("User %s does not exist\n", userID)
		return nil, err
	}

	userOnline := &model.UserOnline{}
	if err = json.Unmarshal(bs, userOnline); err != nil {
		log.Printf("Unabled to unmarshall data for user %s\n", userID)
		return nil, err
	}

	log.Printf("Online info for the user %s: %v\n", userID, userOnline)

	return userOnline, nil
}

func SetUserOnlineInfo(userOnline *model.UserOnline) (err error) {
	redisClient := redislib.GetClient()

	cacheKey := getUserOnlineKey(userOnline.UserID)

	bs, err := json.Marshal(userOnline)
	if err != nil {
		return err
	}

	_, err = redisClient.Do("setEx", cacheKey, userOnlineCacheTime, string(bs)).Result()
	if err != nil {
		return
	}

	return
}
