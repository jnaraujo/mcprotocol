package protocol

import (
	"crypto/rand"

	"github.com/jnaraujo/mcprotocol/api/uuid"
	"github.com/jnaraujo/mcprotocol/auth"
	"github.com/jnaraujo/mcprotocol/packet"
)

type LoginStartPacket struct {
	Name string
	UUID uuid.UUID
}

func ReceiveLoginStartPacket(pkt *packet.Packet) (*LoginStartPacket, error) {
	name, err := pkt.Buffer().ReadString()
	if err != nil {
		return nil, err
	}

	uuid, err := pkt.Buffer().ReadUUID()
	if err != nil {
		return nil, err
	}

	return &LoginStartPacket{
		Name: name,
		UUID: uuid,
	}, nil
}

type EncryptionResponsePacket struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func ReceiveEncryptionResponsePacket(pkt *packet.Packet) (*EncryptionResponsePacket, error) {
	sharedSecretLength, err := pkt.Buffer().ReadVarInt()
	if err != nil {
		return nil, err
	}

	sharedSecret, err := pkt.Buffer().ReadBytes(int(sharedSecretLength))
	if err != nil {
		return nil, err
	}

	verifyTokenLength, err := pkt.Buffer().ReadVarInt()
	if err != nil {
		return nil, err
	}

	verifyToken, err := pkt.Buffer().ReadBytes(int(verifyTokenLength))
	if err != nil {
		return nil, err
	}

	return &EncryptionResponsePacket{
		SharedSecret: sharedSecret,
		VerifyToken:  verifyToken,
	}, nil
}

func CreateEncryptionRequestPacket(crypto *auth.Crypto) (*packet.Packet, error) {
	pkt := packet.NewPacket(0x01)

	rndBytes := make([]byte, 4)
	rand.Read(rndBytes)

	err := pkt.Buffer().WriteString("")
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteVarInt(int32(len(crypto.PublicKeyBytes())))
	if err != nil {
		return nil, err
	}
	_, err = pkt.Buffer().WriteBytes(crypto.PublicKeyBytes())
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteVarInt(4)
	if err != nil {
		return nil, err
	}
	_, err = pkt.Buffer().WriteBytes(rndBytes)
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteBool(false)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func CreateLoginSuccessPacket(playerUUID uuid.UUID, playerUsername string) (*packet.Packet, error) {
	pkt := packet.NewPacket(0x02)

	err := pkt.Buffer().WriteUUID(playerUUID)
	if err != nil {
		return nil, err
	}

	err = pkt.Buffer().WriteString(playerUsername)
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteVarInt(0)
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteBool(false)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}
