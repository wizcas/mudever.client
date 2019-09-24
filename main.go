package main

import (
	"github.com/wizcas/mudever.svc/data"
	"github.com/wizcas/mudever.svc/nvt"
	"go.uber.org/zap"
)

var (
	game = data.MudGame{
		Name:   "pkuxkx",
		Server: data.NewServer("mud.pkuxkx.net", 8080),
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
	logger = logger.Named("client")
	zap.ReplaceGlobals(logger)
	return logger
}

func main() {
	logger := initLogger(false)
	defer logger.Sync()
	client := NewClient(nvt.EncodingGB18030)
	if err := client.Connect(game.Server); err != nil {
		logger.Sugar().Fatalf("unexpected exit: %v", err)
	}
	logger.Info("good bye ;)")
}
