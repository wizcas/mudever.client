package main

import (
	"log"

	"github.com/wizcas/mudever.svc/telnet"
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
	//@TODO: Configure the TLS connection here, if you need to.
	// tlsConfig := &tls.Config{}

	// caller := comm.NewEncodedCaller(comm.GB18030)
	// // call
	// addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
	// if err := telnet.DialToAndCall(addr, caller); err != nil {
	// 	log.Fatal(err)
	// }

	client := telnet.NewClient(telnet.TermEncodingGB18030)
	if err := client.Connect(game.Server); err != nil {
		log.Fatalf("[FATAL ERROR]: %v\n", err)
	}
}
