package packet

import (
	"github.com/jnaraujo/mcprotocol/raknet"
)

type Packet struct {
	id     PacketID
	buffer *raknet.Buffer
}

func NewPacket(buffer *raknet.Buffer, id PacketID) *Packet {
	return &Packet{
		buffer: buffer,
		id:     id,
	}
}

func (p *Packet) ID() PacketID {
	return p.id
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	dataBuf := raknet.NewBuffer()
	err := dataBuf.WriteByte(byte(p.id))
	if err != nil {
		return []byte{}, err
	}

	_, err = dataBuf.WriteBytes(p.Bytes())
	if err != nil {
		return []byte{}, err
	}

	packetBuf := raknet.NewBuffer()
	err = packetBuf.WriteVarInt(int32(dataBuf.Len()))
	if err != nil {
		return []byte{}, err
	}

	_, err = packetBuf.WriteBytes(dataBuf.Bytes())
	if err != nil {
		return []byte{}, err
	}

	return packetBuf.Bytes(), nil
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	buf := raknet.NewBufferFrom(data)

	_, err := buf.ReadVarInt() // packet length
	if err != nil {
		return err
	}

	packetID, err := buf.ReadByte()
	if err != nil {
		return err
	}

	p.id = PacketID(packetID)
	p.buffer = buf

	return nil
}

func (p *Packet) Bytes() []byte {
	return p.buffer.Bytes()
}

func (p *Packet) Buffer() *raknet.Buffer {
	return p.buffer
}
