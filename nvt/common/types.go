package common

import "github.com/wizcas/mudever.svc/packet"

type PacketSender interface {
	Send(p packet.Packet) error
}

type OnError func(err error)
