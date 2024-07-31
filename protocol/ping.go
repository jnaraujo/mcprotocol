package protocol

import (
	"github.com/jnaraujo/mcprotocol/packet"
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
	pkt := packet.NewPacket(packet.IDPing)

	err := pkt.Buffer().WriteLong(payload)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}
