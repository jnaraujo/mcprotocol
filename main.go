package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jnaraujo/mcprotocol/protocol"
)

func main() {
	listener, err := net.Listen("tcp", ":25565")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	slog.Info("New connection", "addr", conn.RemoteAddr().String())
	defer conn.Close()

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		slog.Error("Error reading from connection", "err", err.Error())
		return
	}
	handshakePkt, err := protocol.ReceiveHandshakePacket(buf[:n])
	if err != nil {
		slog.Error("Error reading handshake", "err", err.Error())
		return
	}

	// show motd
	if handshakePkt.NextState == protocol.HandshakeNextStateStatus {
		srp, err := protocol.StatusResponsePacket()
		if err != nil {
			slog.Error("Error creating status response packet", "err", err.Error())
			return
		}

		srpBytes, err := srp.MarshalBinary()
		if err != nil {
			slog.Error("Error marshalling status response packet", "err", err.Error())
			return
		}
		conn.Write(srpBytes)

		// PING REQUEST
		conn.Read(buf) // drop

		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("Error reading from connection", "err", err.Error())
			return
		}

		prp, err := protocol.ReceivePingRequestPacket(buf[:n])
		if err != nil {
			slog.Error("Error receiving ping request packet", "err", err.Error())
			return
		}

		fmt.Println(prp.Payload)

		prPkt, err := protocol.PingResponsePacket(prp.Payload)
		if err != nil {
			slog.Error("Error creating ping response packet", "err", err.Error())
			return
		}

		prBytes, err := prPkt.MarshalBinary()
		if err != nil {
			slog.Error("Error marshalling ping response packet", "err", err.Error())
			return
		}

		conn.Write(prBytes)
	}
}
