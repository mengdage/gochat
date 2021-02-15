package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"githum.com/mengdage/gochat/pkg/api"
)

// Server represents a server
type Server struct {
	router          *gin.Engine
	userService     api.UserService
	wsService       api.WSService
	convService     api.ConversationService
	clientManager   *api.ClientManager
	messageOperator *api.MessageOperator
}

// NewServer creates a new server
func NewServer(router *gin.Engine, userService api.UserService, wsService api.WSService, convService api.ConversationService, cm *api.ClientManager, mo *api.MessageOperator) *Server {
	return &Server{
		clientManager:   cm,
		router:          router,
		userService:     userService,
		wsService:       wsService,
		messageOperator: mo,
		convService:     convService,
	}
}

// Run starts a server
func (s *Server) Run(addr string) error {
	s.Routes()

	httpAddr := "0.0.0.0:" + viper.GetString("app.http_port")
	if addr != "" {
		httpAddr = "0.0.0.0" + addr

	}

	err := s.router.Run(httpAddr)
	if err != nil {
		log.Printf("Error while running server: %v", err)
		return err
	}
	return nil
}
