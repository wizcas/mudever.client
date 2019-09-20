package main

import (
	"log"

	"github.com/wizcas/mudever.svc/telnet"
	"github.com/wizcas/mudever.svc/telnet/nvt"
)

// MudGame contains the profile of a mud server
type MudGame struct {
	Name   string
	Server telnet.Server
}

var (
	game = MudGame{
		Name:   "pkuxkx",
		Server: telnet.Server{"mud.pkuxkx.net", 8080},
	}
)

func main() {
	client := telnet.NewClient(nvt.EncodingGB18030)
	if err := client.Connect(game.Server); err != nil {
		log.Fatalf("[FATAL ERROR]: %v\n", err)
	}
}
