package player

import (
	"errors"
	"net"

	"github.com/jnaraujo/mcprotocol/api/uuid"
	"github.com/jnaraujo/mcprotocol/fsm"
	"github.com/jnaraujo/mcprotocol/packet"
)

type Position struct {
	X        float64
	FeetY    float64
	HeadY    float64
	Z        float64
	OnGround bool
}

type Player struct {
	UUID uuid.UUID
	Name string

	IsLoggedIn bool
	IsAlive    bool
	Position   Position

	// connection stuff
	Conn  *net.TCPConn
	State fsm.FSM
}

func (p *Player) SendPacket(pkt *packet.Packet) error {
	pktBytes, err := pkt.MarshalBinary()
	if err != nil {
		return nil
	}

	if p.Conn == nil {
		return errors.New("conn was not set")
	}

	_, err = p.Conn.Write(pktBytes)
	if err != nil {
		return err
	}
	return nil
}
