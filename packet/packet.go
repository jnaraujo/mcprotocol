package packet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/jnaraujo/mcprotocol/assert"
)

const (
	segment_bits  = 0x7F // the value itself
	continue_bits = 0x80 // if there is more bytes after current byte
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

func WriteVarInt(buf *bytes.Buffer, value int32) error {
	assert.Assert(value >= 0, "value should the greater than 0")
	for {
		if (value & ^segment_bits) == 0 {
			return WriteByte(buf, byte(value))
		}
		err := WriteByte(buf, byte((value&segment_bits)|continue_bits))
		if err != nil {
			return err
		}
		value >>= 7
	}
}

func ReadVarInt(buf *bytes.Buffer) (int32, error) {
	val := int32(0)
	pos := int32(0)
	for {
		currentByte, err := ReadByte(buf)
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

func WriteString(buf *bytes.Buffer, str string) error {
	err := WriteVarInt(buf, int32(len(str)))
	if err != nil {
		return err
	}
	_, err = buf.WriteString(str)
	return err
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

func ReadUShort(buf *bytes.Buffer) (uint16, error) {
	b := make([]byte, 2)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func WriteUShort(buf *bytes.Buffer, value uint16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, value)
	_, err := buf.Write(b)
	return err
}

func ReadVarLong(buf *bytes.Buffer) (int64, error) {
	val := int64(0)
	pos := uint8(0)
	for {
		currentByte, err := ReadByte(buf)
		if err != nil {
			return 0, err
		}
		val |= int64(currentByte&segment_bits) << pos
		if currentByte&continue_bits == 0 {
			break
		}
		pos += 7
		if pos >= 64 {
			return 0, ErrTooBig
		}
	}
	return val, nil
}

func WriteVarLong(buf *bytes.Buffer, value int64) error {
	for {
		if value & ^segment_bits == 0 {
			return WriteByte(buf, byte(value))
		}
		err := WriteByte(buf, byte((value&segment_bits)|continue_bits))
		if err != nil {
			return err
		}
		value >>= 7
	}
}
