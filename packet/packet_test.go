package packet

import (
	"fmt"
	"testing"

	"github.com/jnaraujo/mcprotocol/raknet"
	"github.com/stretchr/testify/assert"
)

func TestPacket(t *testing.T) {
	p := NewPacket(nil, 2)
	assert.Equal(t, PacketID(2), p.ID())
}

func TestPacketMarshalBinary(t *testing.T) {
	buf := raknet.NewBuffer()
	buf.WriteUShort(1200)
	p := NewPacket(buf, 5)
	b, err := p.MarshalBinary()

	fmt.Println(b)

	assert.Nil(t, err)
	assert.Equal(t, byte(3), b[0])         // packet size
	assert.Equal(t, byte(5), b[1])         // packet id
	assert.Equal(t, []byte{4, 176}, b[2:]) // packet data
}

func TestPacketUnmarshalBinary(t *testing.T) {
	buf := raknet.NewBuffer()
	buf.WriteUShort(1200)

	p := NewPacket(buf, 12)
	b, err := p.MarshalBinary()
	assert.Nil(t, err)

	var p2 Packet
	err = p2.UnmarshalBinary(b)
	assert.Nil(t, err)
	assert.Equal(t, PacketID(12), p2.ID())
	assert.Equal(t, []byte{4, 176}, p2.data.Bytes())
}
