package telnet

import (
	"fmt"
	"net"
	"strings"
)

type Server struct {
	Host string
	Port uint16
}

type Client struct {
	conn     net.Conn
	reader   *reader
	writer   *writer
	terminal *Terminal
}

func NewClient(encoding TermEncoding) *Client {
	return &Client{
		terminal: NewTerminal(encoding),
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
	c.reader = newReader(conn)
	c.writer = newWriter(conn)
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
	return c.terminal.proc(c.reader, c.writer)
}
