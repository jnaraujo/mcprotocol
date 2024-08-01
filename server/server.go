package server

import (
	"io"
	"log/slog"
	"net"

	"github.com/jnaraujo/mcprotocol/api/uuid"
	"github.com/jnaraujo/mcprotocol/auth"
	"github.com/jnaraujo/mcprotocol/fsm"
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/protocol"
)

type Server struct {
	addr           string
	statusResponse protocol.StatusResponse

	crypto *auth.Crypto
}

func NewServer(addr string) *Server {
	crypto, err := auth.NewCrypto()
	if err != nil {
		panic(err)
	}

	return &Server{
		addr:   addr,
		crypto: crypto,
		statusResponse: protocol.StatusResponse{
			Version: protocol.StatusResponseVersion{
				Name:     "1.7.10",
				Protocol: 5,
			},
			Description: protocol.StatusResponseDescription{
				Text: "Hello, world!",
			},
			Players: protocol.StatusResponsePlayers{
				Online: 0,
				Max:    20,
			},
		},
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
	defer func() {
		conn.Close()
		slog.Info("Connection closed")
	}()

	slog.Info("New connection", "addr", conn.RemoteAddr().String())

	buf := make([]byte, packet.MaxPacketSizeInBytes)
	state := fsm.NewFSM()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			slog.Error("Error reading from connection", "err", err.Error())
			continue
		}

		pkt := new(packet.Packet)
		err = pkt.UnmarshalBinary(buf[:n])
		if err != nil {
			slog.Error("Error unmarshalling packet", "err", err.Error())
			return
		}

		slog.Info("New packet!", "id", pkt.ID(), "state", state.State())

		switch state.State() {
		case fsm.FSMStateHandshake:
			s.handleHandshakeState(conn, pkt, state)
		case fsm.FSMStateStatus:
			s.handleStatusState(conn, pkt)
		case fsm.FSMStateLogin:
			s.handleLoginState(conn, pkt, state)
		default:
			slog.Error("State not implemented", "state", state.State())
		}
	}
}

func (s *Server) handleHandshakeState(conn *net.TCPConn, pkt *packet.Packet, state *fsm.FSM) {
	handshakePkt, err := protocol.ReceiveHandshakePacket(pkt)
	if err != nil {
		slog.Error("Error reading handshake", "err", err.Error())
		return
	}

	slog.Info("New HandShake Packet", "nextState", handshakePkt.NextState)

	switch handshakePkt.NextState {
	case protocol.HandshakeNextStateStatus:
		state.SetState(fsm.FSMStateStatus)
		// show motd
		statusRespPkt, err := protocol.CreateStatusResponsePacket(s.statusResponse)
		if err != nil {
			slog.Error("Error creating status response packet", "err", err.Error())
			return
		}

		err = s.sendPacket(conn, statusRespPkt)
		if err != nil {
			slog.Error("Error sending status response bytes")
			return
		}
	case protocol.HandshakeNextStateLogin:
		state.SetState(fsm.FSMStateLogin)
	default:
		slog.Error("next state not implemented", "nextState", handshakePkt.NextState)
	}
}

func (s *Server) handleStatusState(conn *net.TCPConn, pkt *packet.Packet) {
	if pkt.Buffer().Len() == 0 {
		return
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

	err = s.sendPacket(conn, pingRespPkt)
	if err != nil {
		slog.Error("Error sending ping response bytes")
		return
	}
}

func (s *Server) handleLoginState(conn *net.TCPConn, pkt *packet.Packet, state *fsm.FSM) {
	slog.Info("New Login Packet", "id", pkt.ID())

	switch pkt.ID() {
	case 0x00: // login start
		loginStartPkt, err := protocol.ReceiveLoginStartPacket(pkt)
		if err != nil {
			slog.Error("error receiving login start packet", "err", err.Error())
			return
		}

		slog.Info("Hello, Player!", "name", loginStartPkt.Name)

		// TODO: implement encryption!!!

		loginSuccessPkt, err := protocol.CreateLoginSuccessPacket(uuid.GenerateUUID(), loginStartPkt.Name)
		if err != nil {
			slog.Error("error creating login success packet", "err", err.Error())
			return
		}
		err = s.sendPacket(conn, loginSuccessPkt)
		if err != nil {
			slog.Error("Error sending login success bytes")
			return
		}

		joinGamePkt, err := protocol.CreateJoinGamePacket()
		if err != nil {
			slog.Error("error creating join game packet", "err", err.Error())
			return
		}

		err = s.sendPacket(conn, joinGamePkt)
		if err != nil {
			slog.Error("Error sending join game bytes")
			return
		}

		// change the state to game mode
		state.SetState(fsm.FSMStatePlay)
	case 0x01:
		slog.Error("login encryption response not implemented yet")
	case 0x03:
		slog.Info("Login success!")
	default:
		slog.Error("login id not implemented", "id", pkt.ID())
	}

}

func (s *Server) sendPacket(conn net.Conn, pkt *packet.Packet) error {
	pktBytes, err := pkt.MarshalBinary()
	if err != nil {
		slog.Error("error marshalling packet", "err", err.Error())
		return err
	}

	_, err = conn.Write(pktBytes)
	return err
}
