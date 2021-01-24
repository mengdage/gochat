package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"githum.com/mengdage/gochat/server/websocket"
)

func ListUsers(c *gin.Context) {
	users := websocket.GetAllUsers()
	log.Println(users)

	c.JSON(http.StatusOK, users)
}
