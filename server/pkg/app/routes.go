package app

import (
	"time"

	"github.com/gin-contrib/cors"
)

// Routes add routes
func (s *Server) Routes() {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		// AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))
	s.router.POST("/register", s.CreateUser())
	s.router.POST("/login", s.LoginUser())

	user := s.router.Group("/user")
	user.Use(CreateAuthMiddleware(s.userService))
	{
		user.GET("/list", s.ListAllUsers())
		user.GET("/conversation_history/:id", s.GetConversationHistory())
	}

	// Takes a http request and upgrade it to a ws connection.
	s.router.GET("/ws", s.WSUpgrader())
}
