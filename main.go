package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"

	"github.com/jnaraujo/mcprotocol/fsm"
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

	state := fsm.NewFSM()
	buf := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("Error reading from connection", "err", err.Error())
			return
		}
		data := buf[:n]

		switch state.State() {
		case fsm.FSMStateHandshake:
			pkt, err := protocol.ReceiveHandshakePacket(data)
			if err != nil {
				slog.Error("Error reading handshake", "err", err.Error())
				return
			}

			switch pkt.NextState {
			case protocol.HandshakeNextStateStatus:
				state.SetState(fsm.FSMStateStatus)
				// show motd
				statusRespPkt, err := protocol.CreateStatusResponsePacket()
				if err != nil {
					slog.Error("Error creating status response packet", "err", err.Error())
					return
				}
				statusRespBytes, err := statusRespPkt.MarshalBinary()
				if err != nil {
					slog.Error("Error marshalling status response packet", "err", err.Error())
					return
				}
				_, err = conn.Write(statusRespBytes)
				if err != nil {
					slog.Error("Error writing status response bytes")
					return
				}
			case protocol.HandshakeNextStateLogin:
				state.SetState(fsm.FSMStateLogin)
			default:
				slog.Error("next state not implemented", "nextState", pkt.NextState)
			}
		case fsm.FSMStateStatus:
			if bytes.Equal(data, []byte{1, 0}) {
				continue
			}
			pingReqPkt, err := protocol.ReceivePingRequestPacket(data)
			if err != nil {
				slog.Error("Error receiving ping request packet", "err", err.Error())
				return
			}

			pingRespPkt, err := protocol.PingResponsePacket(pingReqPkt.Payload)
			if err != nil {
				slog.Error("Error creating ping response packet", "err", err.Error())
				return
			}

			pingRespBytes, err := pingRespPkt.MarshalBinary()
			if err != nil {
				slog.Error("Error marshalling ping response packet", "err", err.Error())
				return
			}
			_, err = conn.Write(pingRespBytes)
			if err != nil {
				slog.Error("Error writing ping response bytes")
				return
			}
		default:
			slog.Error("State not implemented", "state", state.State())
		}
	}

}
