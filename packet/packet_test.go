package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteByte(t *testing.T) {
	buf := new(bytes.Buffer)

	WriteByte(buf, byte(123))
	assert.Equal(t, byte(123), buf.Bytes()[0])

	WriteByte(buf, byte(255))
	assert.Equal(t, byte(255), buf.Bytes()[1])
}

func TestReadByte(t *testing.T) {
	buf := new(bytes.Buffer)

	WriteByte(buf, byte(123))

	actual, err := ReadByte(buf)
	assert.Nil(t, err)
	assert.Equal(t, byte(123), actual)

	WriteByte(buf, byte(255))
	WriteByte(buf, byte(154))

	actual, err = ReadByte(buf)
	assert.Nil(t, err)
	assert.Equal(t, byte(255), actual)

	actual, err = ReadByte(buf)
	assert.Nil(t, err)
	assert.Equal(t, byte(154), actual)
}

func TestReadVarInt(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteByte(buf, 0x80)
	WriteByte(buf, 0x01)

	actual, err := ReadVarInt(buf)
	assert.Nil(t, err)
	assert.Equal(t, int32(128), actual)
}

func TestWriteVarInt(t *testing.T) {
	buf := new(bytes.Buffer)

	WriteVarInt(buf, 0)
	WriteVarInt(buf, 128)
	WriteVarInt(buf, 2147483647)

	assert.Equal(t, []byte{0x00}, buf.Bytes()[0:1])
	assert.Equal(t, []byte{0x80, 0x01}, buf.Bytes()[1:3])
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0x07}, buf.Bytes()[3:8])
}

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
