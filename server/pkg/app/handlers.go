package app

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"githum.com/mengdage/gochat/pkg/api"
)

// CreateUser handles the request to create a new user
func (s *Server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUserRequest api.NewUserRequest
		err := c.ShouldBindJSON(&newUserRequest)

		if err != nil {
			log.Printf("Error while parsing new user request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		user, err := s.userService.New(c, newUserRequest)
		if err != nil {
			log.Printf("error while creating a new user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}

// LoginUser handles the request to create a new user
func (s *Server) LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest api.LoginRequest
		err := c.ShouldBindJSON(&loginRequest)

		if err != nil {
			log.Printf("Error while parsing login request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		user, err := s.userService.Login(c, loginRequest)
		if err != nil {
			log.Printf("error while logging in a new user: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": user.Name,
			"user":  user,
		})
	}
}

func (s *Server) ListAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := getCurrentUser(c)
		users, err := s.userService.GetAllUsers(c, currentUser.Name)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}
		// return equalASCIIFold(u.Host, r.Host)
		log.Printf("origin: %s; host: %s\n", u.Host, r.Host)
		return true
	},
}

func (s *Server) WSUpgrader() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Printf("Error while upgrading to a ws connection: %v", err)
		}

		err = s.wsService.New(c, conn, s.clientManager, s.messageOperator)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

	}
}

// GetConversationHistory returns the conversation history between the current user and the target user
// /user/:id/conversation_history
func (s *Server) GetConversationHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := getCurrentUser(c)
		targetUser := c.Param("id")

		msgs, err := s.convService.GetConversationHistory(c, currentUser.Name, currentUser.Name, targetUser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, msgs)
	}
}
