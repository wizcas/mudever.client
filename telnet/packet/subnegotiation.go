package packet

type SubnegotiationPacket struct {
	Option    byte
	Parameter []byte
}

func (p *SubnegotiationPacket) GetKind() Kind {
	return KindSubnegotiation
}
