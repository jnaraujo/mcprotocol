package protocol

import (
	"encoding/json"

	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/raknet"
)

func StatusResponsePacket() (*packet.Packet, error) {
	response := map[string]any{
		"version": map[string]any{
			"name":     "1.21",
			"protocol": 767,
		},
		"players": map[string]any{
			"max":    99999,
			"online": 99999,
		},
		"description": map[string]string{
			"text": "Hello world",
		},
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	buf := raknet.NewBuffer()
	buf.WriteString(string(respBytes))

	p := packet.NewPacket(buf, packet.IDServerIdentification)
	return p, nil
}
