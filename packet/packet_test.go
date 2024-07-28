package packet

import (
	"bytes"
	"encoding/binary"
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

func TestWriteString(t *testing.T) {
	expected := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a elit ex."
	buf := new(bytes.Buffer)
	WriteString(buf, expected)
	length, _ := ReadVarInt(buf)

	assert.Equal(t, int32(75), length)
	assert.Equal(t, expected, buf.String())
}

func TestReadString(t *testing.T) {
	expected := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a elit ex."
	buf := new(bytes.Buffer)
	WriteString(buf, expected)

	actual, err := ReadString(buf)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadUShort(t *testing.T) {
	expected := make([]byte, 2)
	binary.BigEndian.PutUint16(expected, 21312)

	buf := new(bytes.Buffer)
	buf.Write(expected)

	value, err := ReadUShort(buf)
	assert.Nil(t, err)
	assert.Equal(t, uint16(21312), value)
}

func TestWriteUShort(t *testing.T) {
	expected := make([]byte, 2)
	binary.BigEndian.PutUint16(expected, 13456)

	buf := new(bytes.Buffer)
	WriteUShort(buf, 13456)

	assert.Equal(t, expected, buf.Bytes())
}
