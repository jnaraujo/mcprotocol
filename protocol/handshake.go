package protocol

import (
	"github.com/jnaraujo/mcprotocol/packet"
)

type NextState int32

const (
	HandshakeNextStateStatus NextState = 1
	HandshakeNextStateLogin  NextState = 2
)

type HandshakePacket struct {
	ProtocolVersion int32
	NextState       NextState
	Addr            string
	Port            uint16
}

func ReceiveHandshakePacket(data []byte) (*HandshakePacket, error) {
	var pkt packet.Packet
	err := pkt.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	protocolVersion, _ := pkt.Data().ReadVarInt()
	addr, _ := pkt.Data().ReadString()
	port, _ := pkt.Data().ReadUShort()
	nextState, _ := pkt.Data().ReadVarInt()

	return &HandshakePacket{
		ProtocolVersion: protocolVersion,
		Addr:            addr,
		Port:            port,
		NextState:       NextState(nextState),
	}, nil
}
