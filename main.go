package main

import (
	"github.com/wizcas/mudever.svc/telnet"
	"github.com/wizcas/mudever.svc/telnet/nvt"
	"go.uber.org/zap"
)

// MudGame contains the profile of a mud server
type MudGame struct {
	Name   string
	Server telnet.Server
}

var (
	game = MudGame{
		Name:   "pkuxkx",
		Server: telnet.NewServer("mud.pkuxkx.net", 8080),
	}
)

func initLogger(production bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger.Named("svc"))
	return logger
}

func main() {
	logger := initLogger(false)
	defer logger.Sync()
	client := telnet.NewClient(nvt.EncodingGB18030)
	if err := client.Connect(game.Server); err != nil {
		logger.Sugar().Fatalw("unrecoverable error",
			"message", err,
		)
	}
}
