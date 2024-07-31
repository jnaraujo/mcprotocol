package protocol

import (
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/raknet"
)

type PingRequestPacket struct {
	Payload int64
}

func ReceivePingRequestPacket(pkt *packet.Packet) (*PingRequestPacket, error) {
	payload, err := pkt.Buffer().ReadLong()
	if err != nil {
		return nil, err
	}

	return &PingRequestPacket{
		Payload: payload,
	}, nil
}

func CreatePingResponsePacket(payload int64) (*packet.Packet, error) {
	buf := raknet.NewBuffer()
	err := buf.WriteLong(payload)
	if err != nil {
		return nil, err
	}
	return packet.NewPacket(buf, packet.IDPing), nil
}
