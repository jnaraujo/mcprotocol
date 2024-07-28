package packet

import (
	"bytes"
	"errors"

	"github.com/jnaraujo/mcprotocol/assert"
)

var (
	ErrTooBig = errors.New("too big")
)

func WriteByte(buf *bytes.Buffer, value byte) error {
	return buf.WriteByte(value)
}

func ReadByte(buf *bytes.Buffer) (byte, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

const (
	segment_bits  = 0x7F // the value itself
	continue_bits = 0x80 // if there is more bytes after current byte
)

func WriteVarInt(buf *bytes.Buffer, value int32) {
	assert.Assert(value >= 0, "value should the greater than 0")
	for {
		if (value & ^segment_bits) == 0 {
			buf.WriteByte(byte(value))
			return
		}
		buf.WriteByte(byte((value & segment_bits) | continue_bits))
		value >>= 7
	}
}

func ReadVarInt(buf *bytes.Buffer) (int32, error) {
	val := int32(0)
	pos := int32(0)

	for {
		currentByte, err := buf.ReadByte()
		if err != nil {
			return 0, err
		}
		val |= int32(currentByte&segment_bits) << pos

		// if there is no byte after the current
		if currentByte&continue_bits == 0 {
			break
		}

		pos += 7
		if pos >= 32 {
			return 0, ErrTooBig
		}
	}
	return val, nil
}

func WriteString(buf *bytes.Buffer, str string) {
	WriteVarInt(buf, int32(len(str)))
	buf.WriteString(str)
}

func ReadString(buf *bytes.Buffer) (string, error) {
	length, err := ReadVarInt(buf)
	if err != nil {
		return "", err
	}

	strBytes := make([]byte, length)
	_, err = buf.Read(strBytes)
	if err != nil {
		return "", err
	}

	return string(strBytes), nil
}
