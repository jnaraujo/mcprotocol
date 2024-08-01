package packet

import "fmt"

type PacketID byte

const (
	IDServerIdentification PacketID = 0
	IDPing                 PacketID = 1
)

func (pID PacketID) String() string {
	return fmt.Sprintf("0x%x", byte(pID))
}
