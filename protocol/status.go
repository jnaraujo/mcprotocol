package protocol

import (
	"encoding/json"
	"net"

	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/raknet"
)

func SendStatusResponse(conn net.Conn) error {
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
		return err
	}

	buf := raknet.NewBuffer()
	buf.WriteString(string(respBytes))

	p := packet.NewPacket(buf, packet.IDConnectedPing)
	b, err := p.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	return err
}
