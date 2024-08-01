package protocol

import (
	"encoding/json"

	"github.com/jnaraujo/mcprotocol/packet"
)

type StatusResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}
type StatusResponsePlayers struct {
	Max    int `json:"max"`
	Online int `json:"online"`
	Sample []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"sample"`
}
type StatusResponseDescription struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Version            StatusResponseVersion     `json:"version,omitempty"`
	Players            StatusResponsePlayers     `json:"players,omitempty"`
	Description        StatusResponseDescription `json:"description,omitempty"`
	Favicon            string                    `json:"favicon,omitempty"`
	EnforcesSecureChat bool                      `json:"enforcesSecureChat,omitempty"`
}

func CreateStatusResponsePacket(response StatusResponse) (*packet.Packet, error) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	p := packet.NewPacket(packet.IDServerIdentification)
	p.Buffer().WriteString(string(respBytes))

	return p, nil
}
