package raknet

import (
	"encoding/binary"
	"testing"

	"github.com/jnaraujo/mcprotocol/api/uuid"
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

func TestNegativeVarInt(t *testing.T) {
	buf := NewBuffer()
	err := buf.WriteVarInt(-6)
	assert.Nil(t, err)

	actual, err := buf.ReadVarInt()
	assert.Nil(t, err)
	assert.Equal(t, int32(-6), actual)
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

func TestNegativeVarLong(t *testing.T) {
	buf := NewBuffer()
	err := buf.WriteVarLong(-63412337812637)
	assert.Nil(t, err)

	actual, err := buf.ReadVarLong()
	assert.Nil(t, err)
	assert.Equal(t, int64(-63412337812637), actual)
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

func TestWriteUUID(t *testing.T) {
	expected := []byte{137, 206, 23, 145, 218, 176, 75, 42, 151, 217, 238, 114, 160, 205, 193, 253}
	buf := NewBuffer()
	actual, err := uuid.UUIDFromString("89ce1791-dab0-4b2a-97d9-ee72a0cdc1fd")
	assert.Nil(t, err)

	err = buf.WriteUUID(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected, buf.Bytes())
}

func TestReadUUID(t *testing.T) {
	expected := "89ce1791-dab0-4b2a-97d9-ee72a0cdc1fd"
	buf := NewBuffer()

	uuid, err := uuid.UUIDFromString("89ce1791-dab0-4b2a-97d9-ee72a0cdc1fd")
	assert.Nil(t, err)

	err = buf.WriteUUID(uuid)
	assert.Nil(t, err)

	actualUUID, err := buf.ReadUUID()
	assert.Nil(t, err)

	assert.Equal(t, expected, actualUUID.String())
}
