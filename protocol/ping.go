package protocol

import (
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/raknet"
)

type PingRequestPacket struct {
	Payload int64
}

func ReceivePingRequestPacket(data []byte) (*PingRequestPacket, error) {
	var pkt packet.Packet
	err := pkt.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	payload, err := pkt.Data().ReadLong()
	if err != nil {
		return nil, err
	}

	return &PingRequestPacket{
		Payload: payload,
	}, nil
}

func PingResponsePacket(payload int64) (*packet.Packet, error) {
	buf := raknet.NewBuffer()
	err := buf.WriteLong(payload)
	if err != nil {
		return nil, err
	}
	return packet.NewPacket(buf, packet.IDPing), nil
}
