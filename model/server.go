package model

import "fmt"

// Server represents a server.
type Server struct {
	IP   string `json:"ip"`
	Port string `json:"port`
}

// NewServer returns a new Server instance
func NewServer(ip string, port string) *Server {
	return &Server{ip, port}
}

func (s *Server) String() string {
	return fmt.Sprintf("%s:%s", s.IP, s.Port)
}
