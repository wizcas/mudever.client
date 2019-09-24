package data

import (
	"fmt"
	"strings"
)

// Server contains the connection parameters of a game server.
type Server struct {
	Host string
	Port uint16
}

// NewServer returns a specified server config.
// If host is not given, "127.0.0.1" will be set as default.
// If port is not given, 23 will be set as default.
func NewServer(host string, port uint16) Server {
	host = strings.TrimSpace(host)
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 23
	}
	return Server{host, port}
}

// Addr returns the address string for tcp connection.
func (s Server) Addr() string {
	return s.String()
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
