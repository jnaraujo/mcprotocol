package packet

import (
	"errors"
	"io"

	"github.com/jnaraujo/mcprotocol/assert"
)

var (
	ErrTooBig = errors.New("too big")
)

type Packet struct {
	data []byte
	n    int
}

func NewPacket(n int) *Packet {
	return &Packet{
		data: make([]byte, 0, n),
		n:    n,
	}
}

func (p *Packet) WriteByte(b byte) error {
	p.data = append(p.data, b)
	return nil
}

func (p *Packet) ReadByte() (byte, error) {
	if len(p.data) == 0 {
		return 0, io.EOF
	}

	b := p.data[0]
	p.data = p.data[1:]
	return b, nil
}

const (
	segment_bits  = 0x7F // the value itself
	continue_bits = 0x80 // if there is more bytes after current byte
)

func (p *Packet) WriteVarInt(value int32) {
	assert.Assert(value >= 0, "value should the greater than 0")
	for {
		if (value & ^segment_bits) == 0 {
			p.WriteByte(byte(value))
			return
		}
		p.WriteByte(byte((value & segment_bits) | continue_bits))
		value >>= 7
	}
}

func (p *Packet) ReadVarInt() (int32, error) {
	val := int32(0)
	pos := int32(0)

	for {
		currentByte, err := p.ReadByte()
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
