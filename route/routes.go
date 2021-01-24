package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"githum.com/mengdage/gochat/controller/user"
	"githum.com/mengdage/gochat/server/websocket"
)

func Init(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]string{"msg": "hello world!"})
		})
		userRouter.GET("/list", user.ListUsers)
	}

	router.GET("/ws", func(c *gin.Context) {
		websocket.WsConnect(c.Writer, c.Request)
	})
}
