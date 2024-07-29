package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/jnaraujo/mcprotocol/packet"
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
		defer conn.Close()

		fmt.Printf("New connection from %s\n", conn.RemoteAddr())

		buf := make([]byte, 1024*1024)

		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
				continue
			}
			fmt.Println("eof")
			break
		}

		fmt.Println(n)
		fmt.Println(buf[:n])

		var packet packet.Packet
		err = packet.UnmarshalBinary(buf[:n])
		if err != nil {
			slog.Error("Error unmarshalling packet", "err", err.Error())
			continue
		}

		protocolVersion, _ := packet.Data().ReadVarInt()
		serverAddr, _ := packet.Data().ReadString()
		serverPort, _ := packet.Data().ReadUShort()
		nextState, _ := packet.Data().ReadVarInt()

		fmt.Println("data", packet.Id(), protocolVersion, serverAddr, serverPort, nextState)

		switch nextState {
		case 1:
			protocol.SendStatusResponse(conn)
			fmt.Println("status")
		case 2:
			fmt.Println("login")
		case 3:
			fmt.Println("transfer")
		default:
			fmt.Println("not implemented")
		}
	}
}
