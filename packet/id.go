package packet

type PacketID byte

const (
	IDConnectedPing                  PacketID = 0x00
	IDUnconnectedPing                PacketID = 0x01
	IDUnconnectedPingOpenConnections PacketID = 0x02
	IDConnectedPong                  PacketID = 0x03
	IDDetectLostConnections          PacketID = 0x04
	IDOpenConnectionRequest1         PacketID = 0x05
	IDOpenConnectionReply1           PacketID = 0x06
	IDOpenConnectionRequest2         PacketID = 0x07
	IDOpenConnectionReply2           PacketID = 0x08
	IDConnectionRequest              PacketID = 0x09
	IDConnectionRequestAccepted      PacketID = 0x10
	IDNewIncomingConnection          PacketID = 0x13
	IDDisconnectNotification         PacketID = 0x15

	IDIncompatibleProtocolVersion PacketID = 0x19

	IDUnconnectedPong PacketID = 0x1c
)
