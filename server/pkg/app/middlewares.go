package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"githum.com/mengdage/gochat/pkg/api"
)

type MiddlewareService interface {
	GetUserByName(ctx context.Context, userID string) (*api.User, error)
}

var (
	currentUserKey = "currentUser"
)

func getAuthToken(req *http.Request) string {
	v := req.Header.Get("Authorization")
	return v
}

func getCurrentUser(c *gin.Context) *api.User {
	userVal, err := c.Get(currentUserKey)
	if !err {
		log.Println("currentUser does not exist in the context")
		return nil
	}

	user, _ := userVal.(*api.User)

	return user
}

// CreateAuthMiddleware creates a auth middleware which validate the request.
func CreateAuthMiddleware(s MiddlewareService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAuthToken(c.Request)
		if token == "" {
			c.AbortWithError(http.StatusNonAuthoritativeInfo, errors.New("no Authorization info"))
			return
		}
		log.Printf("With Authorization: %s", token)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user, err := s.GetUserByName(ctx, token)
		if err != nil {
			c.AbortWithError(http.StatusNonAuthoritativeInfo, errors.New("invalid Authorization info"))
			return
		}
		log.Printf("Find user: %v", user)

		c.Set(currentUserKey, user)
	}
}
