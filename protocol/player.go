package protocol

import (
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/player"
)

func ReceivePlayerPosition(pkt *packet.Packet) (*player.Position, error) {
	pp := new(player.Position)

	var err error
	pp.X, err = pkt.Buffer().ReadDouble()
	if err != nil {
		return nil, err
	}
	pp.FeetY, err = pkt.Buffer().ReadDouble()
	if err != nil {
		return nil, err
	}
	pp.HeadY, err = pkt.Buffer().ReadDouble()
	if err != nil {
		return nil, err
	}
	pp.Z, err = pkt.Buffer().ReadDouble()
	if err != nil {
		return nil, err
	}
	pp.OnGround, err = pkt.Buffer().ReadBool()
	if err != nil {
		return nil, err
	}

	return pp, nil
}
