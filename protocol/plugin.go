package protocol

import "github.com/jnaraujo/mcprotocol/packet"

type PluginMessage struct {
	Channel string
	Length  int16
	Data    []byte
}

func ReceivePluginMessage(pkt *packet.Packet) (*PluginMessage, error) {
	pm := &PluginMessage{}

	var err error
	pm.Channel, err = pkt.Buffer().ReadString()
	if err != nil {
		return nil, err
	}
	pm.Length, err = pkt.Buffer().ReadShort()
	if err != nil {
		return nil, err
	}
	pm.Data, err = pkt.Buffer().ReadBytes(int(pm.Length))
	if err != nil {
		return nil, err
	}
	return pm, err
}

func CreatePluginMessagePacket(pm *PluginMessage) (*packet.Packet, error) {
	pkt := packet.NewPacket(0x3F)

	err := pkt.Buffer().WriteString(pm.Channel)
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteShort(pm.Length)
	if err != nil {
		return nil, err
	}
	_, err = pkt.Buffer().WriteBytes(pm.Data)
	if err != nil {
		return nil, err
	}
	return pkt, err
}
