package raknet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/jnaraujo/mcprotocol/api/uuid"
)

const (
	segmentBits  = 0x7F // the value itself
	continueBits = 0x80 // if there is more bytes after current byte
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

func (buf *Buffer) ReadBytes(n int) ([]byte, error) {
	readBuf := make([]byte, n)
	_, err := buf.data.Read(readBuf)
	if err != nil {
		return []byte{}, err
	}
	return readBuf, nil
}

func (buf *Buffer) WriteVarInt(value int32) error {
	uValue := uint32(value)

	for {
		if (uValue & ^uint32(segmentBits)) == 0 {
			return buf.WriteByte(byte(uValue))
		}
		err := buf.WriteByte(byte((uValue & segmentBits) | continueBits))
		if err != nil {
			return err
		}
		uValue >>= 7
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
		val |= int32(currentByte&segmentBits) << pos
		// if there is no byte after the current
		if currentByte&continueBits == 0 {
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

func (buf *Buffer) ReadShort() (int16, error) {
	b := make([]byte, 2)
	_, err := buf.data.Read(b)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b)), nil
}

func (buf *Buffer) WriteShort(value int16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(value))
	_, err := buf.data.Write(b)
	return err
}

func (buf *Buffer) WriteLong(value int64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err := buf.data.Write(b)
	return err
}

func (buf *Buffer) ReadLong() (int64, error) {
	b := make([]byte, 8)
	_, err := buf.data.Read(b)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

func (buf *Buffer) ReadVarLong() (int64, error) {
	val := int64(0)
	pos := uint8(0)
	for {
		currentByte, err := buf.ReadByte()
		if err != nil {
			return 0, err
		}
		val |= int64(currentByte&segmentBits) << pos
		if currentByte&continueBits == 0 {
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
	uValue := uint64(value)

	for {
		if uValue & ^uint64(segmentBits) == 0 {
			return buf.WriteByte(byte(uValue))
		}
		err := buf.WriteByte(byte((uValue & segmentBits) | continueBits))
		if err != nil {
			return err
		}
		uValue >>= 7
	}
}

func (buf *Buffer) ReadUUID() (uuid.UUID, error) {
	b := make([]byte, 16)
	_, err := buf.data.Read(b)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.UUIDFromBytes(b)
}

func (buf *Buffer) WriteBool(v bool) error {
	if v {
		return buf.WriteByte(0x01)
	}
	return buf.WriteByte(0x00)
}

func (buf *Buffer) ReadBool() (bool, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return false, err
	}
	return b == 0x01, nil
}

func (buf *Buffer) WriteUUID(uuid uuid.UUID) error {
	_, err := buf.data.Write(uuid[:])
	return err
}

func (buf *Buffer) WriteBytes(p []byte) (n int, err error) {
	return buf.data.Write(p)
}

// Signed 32-bit integer, two's complement
func (buf *Buffer) WriteInt(value int32) error {
	va := make([]byte, 4)
	binary.BigEndian.PutUint32(va, uint32(value))
	_, err := buf.data.Write(va)
	return err
}

func (buf *Buffer) ReadInt() (int32, error) {
	b := make([]byte, 4)
	_, err := buf.data.Read(b)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

func (buf *Buffer) WriteDouble(d float64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(d))
	_, err := buf.data.Write(b)
	return err
}

func (buf *Buffer) ReadDouble() (float64, error) {
	b := make([]byte, 8)
	_, err := buf.data.Read(b)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.BigEndian.Uint64(b)), nil
}

func (buf *Buffer) Len() int {
	return buf.data.Len()
}

func (buf *Buffer) Bytes() []byte {
	return buf.data.Bytes()
}
