package raknet

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

type Buffer struct {
	data *bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{
		data: new(bytes.Buffer),
	}
}

func NewBufferFrom(data []byte) *Buffer {
	return &Buffer{
		data: bytes.NewBuffer(data),
	}
}

func (buf *Buffer) WriteByte(value byte) error {
	return buf.data.WriteByte(value)
}

func (buf *Buffer) ReadByte() (byte, error) {
	currentByte, err := buf.data.ReadByte()
	if err != nil {
		return 0, err
	}
	return currentByte, nil
}

func (buf *Buffer) WriteVarInt(value int32) error {
	assert.Assert(value >= 0, "value should the greater than 0")
	for {
		if (value & ^segment_bits) == 0 {
			return buf.WriteByte(byte(value))
		}
		err := buf.WriteByte(byte((value & segment_bits) | continue_bits))
		if err != nil {
			return err
		}
		value >>= 7
	}
}

func (buf *Buffer) ReadVarInt() (int32, error) {
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

func (buf *Buffer) WriteString(str string) error {
	err := buf.WriteVarInt(int32(len(str)))
	if err != nil {
		return err
	}
	_, err = buf.data.WriteString(str)
	return err
}

func (buf *Buffer) ReadString() (string, error) {
	length, err := buf.ReadVarInt()
	if err != nil {
		return "", err
	}

	strBytes := make([]byte, length)
	_, err = buf.data.Read(strBytes)
	if err != nil {
		return "", err
	}

	return string(strBytes), nil
}

func (buf *Buffer) ReadUShort() (uint16, error) {
	b := make([]byte, 2)
	_, err := buf.data.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func (buf *Buffer) WriteUShort(value uint16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, value)
	_, err := buf.data.Write(b)
	return err
}

func (buf *Buffer) ReadVarLong() (int64, error) {
	val := int64(0)
	pos := uint8(0)
	for {
		currentByte, err := buf.ReadByte()
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

func (buf *Buffer) WriteVarLong(value int64) error {
	for {
		if value & ^segment_bits == 0 {
			return buf.WriteByte(byte(value))
		}
		err := buf.WriteByte(byte((value & segment_bits) | continue_bits))
		if err != nil {
			return err
		}
		value >>= 7
	}
}

func (buf *Buffer) WriteBytes(p []byte) (n int, err error) {
	return buf.data.Write(p)
}

func (buf *Buffer) Len() int {
	return buf.data.Len()
}

func (buf *Buffer) Bytes() []byte {
	return buf.data.Bytes()
}
