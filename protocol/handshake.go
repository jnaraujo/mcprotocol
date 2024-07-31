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

func ReceiveHandshakePacket(pkt *packet.Packet) (*HandshakePacket, error) {
	protocolVersion, err := pkt.Buffer().ReadVarInt()
	if err != nil {
		return nil, err
	}
	addr, err := pkt.Buffer().ReadString()
	if err != nil {
		return nil, err
	}
	port, err := pkt.Buffer().ReadUShort()
	if err != nil {
		return nil, err
	}
	nextState, err := pkt.Buffer().ReadVarInt()
	if err != nil {
		return nil, err
	}

	return &HandshakePacket{
		ProtocolVersion: protocolVersion,
		Addr:            addr,
		Port:            port,
		NextState:       NextState(nextState),
	}, nil
}
