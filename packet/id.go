package packet

type PacketID byte

const (
	IDServerIdentification PacketID = 0
	IDPing                 PacketID = 1
)
