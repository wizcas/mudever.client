package main

import (
	"github.com/wizcas/mudever.svc/telnet"
	"github.com/wizcas/mudever.svc/telnet/nvt"
	"github.com/wizcas/mudever.svc/utils"
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

func main() {
	utils.InitLogger("main", false)
	defer utils.Logger().Sync()
	client := telnet.NewClient(nvt.EncodingGB18030)
	if err := client.Connect(game.Server); err != nil {
		utils.Logger().Fatalw("unrecoverable error",
			"message", err,
		)
	}
}
