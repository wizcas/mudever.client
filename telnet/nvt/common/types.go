package common

import "github.com/wizcas/mudever.svc/telnet/packet"

type PacketSender interface {
	Send(p packet.Packet) error
}

type OnError func(err error)
