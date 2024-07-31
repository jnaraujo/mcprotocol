package server

import (
	"log/slog"
	"net"

	"github.com/jnaraujo/mcprotocol/fsm"
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/protocol"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn *net.TCPConn) {
	slog.Info("New connection", "addr", conn.RemoteAddr().String())
	defer conn.Close()

	defer conn.Close()

	buf := make([]byte, packet.MaxPacketSizeInBytes)
	state := fsm.NewFSM()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("Error reading from connection", "err", err.Error())
			return
		}

		pkt := new(packet.Packet)
		err = pkt.UnmarshalBinary(buf[:n])
		if err != nil {
			slog.Error("Error unmarshalling packet", "err", err.Error())
			return
		}

		switch state.State() {
		case fsm.FSMStateHandshake:
			handshakePkt, err := protocol.ReceiveHandshakePacket(pkt)
			if err != nil {
				slog.Error("Error reading handshake", "err", err.Error())
				return
			}

			switch handshakePkt.NextState {
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
				slog.Error("next state not implemented", "nextState", handshakePkt.NextState)
			}
		case fsm.FSMStateStatus:
			if pkt.Bytes() == nil {
				continue
			}
			pingReqPkt, err := protocol.ReceivePingRequestPacket(pkt)
			if err != nil {
				slog.Error("Error receiving ping request packet", "err", err.Error())
				return
			}

			pingRespPkt, err := protocol.CreatePingResponsePacket(pingReqPkt.Payload)
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
