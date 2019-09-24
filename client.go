package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/wizcas/mudever.svc/data"
	"github.com/wizcas/mudever.svc/nvt"
	"github.com/wizcas/mudever.svc/nvt/common"
	"github.com/wizcas/mudever.svc/stream"
	"go.uber.org/zap"
)

// Client is the session where MUD game runs.
// It communicates with game server with a built-in telnet NVT, and
// interacts with the front-end UI. MUD gaming tools (such as alias,
// trigger, script runner, etc.) operate in the client for consistency.
type Client struct {
	conn     net.Conn
	reader   *stream.Reader
	writer   *stream.Writer
	terminal *nvt.Terminal
}

// NewClient creates a mud client with given encoding
func NewClient(encoding nvt.Encoding) *Client {
	return &Client{
		terminal: nvt.NewTerminal(encoding),
	}
}

// Connect the client to given server
func (c *Client) Connect(server data.Server) error {
	conn, err := dial(server)
	if err != nil {
		return err
	}
	if conn == nil {
		return fmt.Errorf("cannot establish connection to %s", server)
	}
	c.conn = conn
	c.reader = stream.NewReader(conn)
	c.writer = stream.NewWriter(conn)
	defer c.close()
	return c.run()
}

func dial(server data.Server) (net.Conn, error) {
	conn, err := net.Dial("tcp", server.Addr())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Client) close() {
	c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.writer = nil
}

func (c *Client) run() error {
	chSig := make(chan os.Signal)
	chErr := make(chan nvt.TerminalError)
	signal.Notify(chSig, os.Interrupt)
	defer signal.Stop(chSig)
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.terminal.Start(rootCtx, c.reader, c.writer, chErr)
	for {
		select {
		case <-chSig:
			common.Logger().Info("gracefully stopping...")
			cancel()
		case err := <-chErr:
			if err.Panic() {
				return err.RawErr()
			}
			zap.S().Error(err)
		case <-c.terminal.Stopped():
			common.Logger().Info("client stopped")
			return nil
		}
	}
}
