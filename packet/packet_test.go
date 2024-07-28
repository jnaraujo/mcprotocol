package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteByte(t *testing.T) {
	p := NewPacket(2)
	p.WriteByte(byte(123))
	assert.Equal(t, byte(123), p.data[0])

	p.WriteByte(byte(255))
	assert.Equal(t, byte(255), p.data[1])
}

func TestReadByte(t *testing.T) {
	p := NewPacket(2)
	p.WriteByte(byte(123))
	actual, err := p.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(123), actual)

	p.WriteByte(byte(255))
	p.WriteByte(byte(154))
	actual, err = p.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(255), actual)

	actual, err = p.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(154), actual)
}

func TestReadVarInt(t *testing.T) {
	p := NewPacket(2)
	p.WriteByte(0x80)
	p.WriteByte(0x01)

	actual, err := p.ReadVarInt()
	assert.Nil(t, err)
	assert.Equal(t, int32(128), actual)
}

func TestWriteVarInt(t *testing.T) {
	p := NewPacket(10)

	p.WriteVarInt(0)
	p.WriteVarInt(128)
	p.WriteVarInt(2147483647)

	assert.Equal(t, []byte{0x00}, p.data[0:1])
	assert.Equal(t, []byte{0x80, 0x01}, p.data[1:3])
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0x07}, p.data[3:8])

	actual, _ := p.ReadVarInt()
	assert.Equal(t, int32(0), actual)

	actual, _ = p.ReadVarInt()
	assert.Equal(t, int32(128), actual)

	actual, _ = p.ReadVarInt()
	assert.Equal(t, int32(2147483647), actual)
}
