package telnet

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/wizcas/mudever.svc/telnet/nvt"
	"github.com/wizcas/mudever.svc/telnet/nvt/common"
	"github.com/wizcas/mudever.svc/telnet/stream"
	"go.uber.org/zap"
)

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

func (s Server) Addr() string {
	return s.String()
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
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
	conn, err := net.Dial("tcp", server.Addr())
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
		return fmt.Errorf("cannot establish connection to %s", server)
	}
	c.conn = conn
	c.reader = stream.NewReader(conn)
	c.writer = stream.NewWriter(conn)
	defer c.Close()
	return c.run()
}

func (c *Client) Close() {
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
			common.Logger().Info("gracefully closing client...")
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
