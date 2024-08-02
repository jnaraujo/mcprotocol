package protocol

import "github.com/jnaraujo/mcprotocol/packet"

type ClientSettings struct {
	Locale       string
	ViewDistance byte
	ChatFlags    byte
	ChatColours  bool
	Difficulty   byte
	ShowCape     bool
}

func ReceiveClientSettings(pkt *packet.Packet) (*ClientSettings, error) {
	settings := &ClientSettings{}

	var err error

	settings.Locale, err = pkt.Buffer().ReadString()
	if err != nil {
		return nil, err
	}
	settings.ViewDistance, err = pkt.Buffer().ReadByte()
	if err != nil {
		return nil, err
	}
	settings.ChatFlags, err = pkt.Buffer().ReadByte()
	if err != nil {
		return nil, err
	}
	settings.ChatColours, err = pkt.Buffer().ReadBool()
	if err != nil {
		return nil, err
	}
	settings.Difficulty, err = pkt.Buffer().ReadByte()
	if err != nil {
		return nil, err
	}
	settings.ShowCape, err = pkt.Buffer().ReadBool()
	if err != nil {
		return nil, err
	}
	return settings, nil
}
