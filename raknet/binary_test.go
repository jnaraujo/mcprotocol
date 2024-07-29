package raknet

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteByte(t *testing.T) {
	buf := NewBuffer()

	buf.WriteByte(byte(123))
	assert.Equal(t, byte(123), buf.data.Bytes()[0])

	buf.WriteByte(byte(255))
	assert.Equal(t, byte(255), buf.data.Bytes()[1])
}

func TestReadByte(t *testing.T) {
	buf := NewBuffer()

	buf.WriteByte(byte(123))

	actual, err := buf.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(123), actual)

	buf.WriteByte(byte(255))
	buf.WriteByte(byte(154))

	actual, err = buf.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(255), actual)

	actual, err = buf.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(154), actual)
}

func TestReadVarInt(t *testing.T) {
	buf := NewBuffer()
	buf.WriteByte(0x80)
	buf.WriteByte(0x01)

	actual, err := buf.ReadVarInt()
	assert.Nil(t, err)
	assert.Equal(t, int32(128), actual)
}

func TestWriteVarInt(t *testing.T) {
	buf := NewBuffer()

	buf.WriteVarInt(0)
	buf.WriteVarInt(128)
	buf.WriteVarInt(2147483647)

	assert.Equal(t, []byte{0x00}, buf.data.Bytes()[0:1])
	assert.Equal(t, []byte{0x80, 0x01}, buf.data.Bytes()[1:3])
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0x07}, buf.data.Bytes()[3:8])
}

func TestWriteString(t *testing.T) {
	expected := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a elit ex."
	buf := NewBuffer()
	buf.WriteString(expected)
	length, _ := buf.ReadVarInt()

	assert.Equal(t, int32(75), length)
	assert.Equal(t, expected, buf.data.String())
}

func TestReadString(t *testing.T) {
	expected := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a elit ex."
	buf := NewBuffer()
	buf.WriteString(expected)

	actual, err := buf.ReadString()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadUShort(t *testing.T) {
	expected := make([]byte, 2)
	binary.BigEndian.PutUint16(expected, 21312)

	buf := NewBuffer()
	buf.data.Write(expected)

	value, err := buf.ReadUShort()
	assert.Nil(t, err)
	assert.Equal(t, uint16(21312), value)
}

func TestWriteUShort(t *testing.T) {
	expected := make([]byte, 2)
	binary.BigEndian.PutUint16(expected, 13456)

	buf := NewBuffer()
	buf.WriteUShort(13456)

	assert.Equal(t, expected, buf.data.Bytes())
}

func TestReadVarLong(t *testing.T) {
	buf := NewBuffer()
	buf.data.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})

	actual, err := buf.ReadVarLong()
	assert.Nil(t, err)
	assert.Equal(t, int64(9223372036854775807), actual)
}

func TestWriteVarLong(t *testing.T) {
	buf := NewBuffer()

	buf.WriteVarLong(0)
	buf.WriteVarLong(9223372036854775807)

	assert.Equal(t, []byte{0x00}, buf.data.Bytes()[0:1])
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, buf.data.Bytes()[1:10])
}

func TestReadLong(t *testing.T) {
	expected := make([]byte, 8)
	binary.BigEndian.PutUint64(expected, 9223372036854775807)

	buf := NewBuffer()
	buf.data.Write(expected)

	value, err := buf.ReadLong()
	assert.Nil(t, err)
	assert.Equal(t, int64(9223372036854775807), value)
}

func TestWriteLong(t *testing.T) {
	expected := make([]byte, 8)
	binary.BigEndian.PutUint64(expected, 9223372036854775807)

	buf := NewBuffer()
	buf.WriteLong(9223372036854775807)

	assert.Equal(t, expected, buf.data.Bytes())
}
