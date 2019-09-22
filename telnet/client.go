package telnet

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/wizcas/mudever.svc/telnet/nvt"
	"github.com/wizcas/mudever.svc/telnet/stream"
)

type Server struct {
	Host string
	Port uint16
}

type Client struct {
	conn     net.Conn
	reader   *stream.Reader
	writer   *stream.Writer
	terminal *nvt.Terminal
	stopFn   context.CancelFunc
}

func NewClient(encoding nvt.Encoding) *Client {
	return &Client{
		terminal: nvt.NewTerminal(encoding),
	}
}

func dial(server Server) (net.Conn, error) {
	if len(strings.TrimSpace(server.Host)) == 0 {
		server.Host = "127.0.0.1"
	}
	if server.Port == 0 {
		server.Port = 23
	}
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Client) Connect(server Server) error {
	conn, err := dial(server)
	if err != nil {
		return err
	}
	if conn == nil {
		panic("empty connection")
	}
	c.conn = conn
	c.reader = stream.NewReader(conn)
	c.writer = stream.NewWriter(conn)
	return c.run()
}

func (c *Client) Close() {
	c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.writer = nil
}

func (c *Client) run() error {
	defer c.Close()
	return c.terminal.Start(c.reader, c.writer)
}
