package protocol

import (
	"crypto/rand"

	"github.com/jnaraujo/mcprotocol/api/uuid"
	"github.com/jnaraujo/mcprotocol/auth"
	"github.com/jnaraujo/mcprotocol/packet"
)

type LoginStartPacket struct {
	Name string
}

func ReceiveLoginStartPacket(pkt *packet.Packet) (*LoginStartPacket, error) {
	name, err := pkt.Buffer().ReadString()
	if err != nil {
		return nil, err
	}
	return &LoginStartPacket{
		Name: name,
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

	err := pkt.Buffer().WriteString(playerUUID.String())
	if err != nil {
		return nil, err
	}

	err = pkt.Buffer().WriteString(playerUsername)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func CreateJoinGamePacket() (*packet.Packet, error) {
	pkt := packet.NewPacket(0x01)

	// entity id
	err := pkt.Buffer().WriteInt(0)
	if err != nil {
		return nil, err
	}

	// game mode
	// 0: survival, 1: creative, 2: adventure. Bit 3 (0x8) is the hardcore flag
	err = pkt.Buffer().WriteByte(1)
	if err != nil {
		return nil, err
	}
	// Dimension
	// -1: nether, 0: overworld, 1: end
	err = pkt.Buffer().WriteByte(0)
	if err != nil {
		return nil, err
	}
	// Difficulty
	// 0 thru 3 for Peaceful, Easy, Normal, Hard
	err = pkt.Buffer().WriteByte(1)
	if err != nil {
		return nil, err
	}
	// Max Players
	err = pkt.Buffer().WriteByte(2)
	if err != nil {
		return nil, err
	}
	// level type
	// default, flat, largeBiomes, amplified, default_1_1
	err = pkt.Buffer().WriteString("default")
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func CreateSpawnPositionPacket() (*packet.Packet, error) {
	pkt := packet.NewPacket(packet.IDServerSpawnPosition)

	err := pkt.Buffer().WriteInt(0) // x
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteInt(200) // y
	if err != nil {
		return nil, err
	}
	err = pkt.Buffer().WriteInt(0) // z
	if err != nil {
		return nil, err
	}
	return pkt, nil
}
