package server

import (
	"errors"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net"
	"syscall"
	"time"

	"github.com/jnaraujo/mcprotocol/api/uuid"
	"github.com/jnaraujo/mcprotocol/auth"
	"github.com/jnaraujo/mcprotocol/fsm"
	"github.com/jnaraujo/mcprotocol/packet"
	"github.com/jnaraujo/mcprotocol/player"
	"github.com/jnaraujo/mcprotocol/protocol"
)

type Server struct {
	addr           string
	statusResponse protocol.StatusResponse

	crypto  *auth.Crypto
	players map[string]*player.Player
}

func NewServer(addr string) *Server {
	crypto, err := auth.NewCrypto()
	if err != nil {
		panic(err)
	}

	return &Server{
		addr:    addr,
		crypto:  crypto,
		players: make(map[string]*player.Player),
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

	go func() {
		for {
			time.Sleep(10 * time.Second)
			rndId := rand.Int31n(math.MaxInt32)

			slog.Info("Sending KeepAlive packets", "id", rndId)

			pkt := packet.NewPacket(packet.IDServerKeepAlive)
			pkt.Buffer().WriteInt(rndId)
			for addr, plr := range s.players {
				if !plr.IsLoggedIn {
					continue
				}
				err := plr.SendPacket(pkt)
				if err != nil {
					slog.Error("error sending keep alive packet", "addr", addr, "err", err.Error())
				}
			}
		}
	}()

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

	plr, exists := s.players[conn.RemoteAddr().String()]
	if !exists {
		plr = &player.Player{
			Conn: conn,
		}
		s.players[conn.RemoteAddr().String()] = plr
	}

	// close player connection
	defer s.closeConn(plr)

	buf := make([]byte, packet.MaxPacketSizeInBytes)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			switch {
			case errors.Is(err, net.ErrClosed),
				errors.Is(err, io.EOF),
				errors.Is(err, syscall.EPIPE):
				return
			default:
				slog.Error("Error reading from connection", "err", err.Error())
				continue
			}
		}

		pkt := new(packet.Packet)
		err = pkt.UnmarshalBinary(buf[:n])
		if err != nil {
			slog.Error("Error unmarshalling packet", "err", err.Error())
			return
		}

		switch plr.State.State() {
		case fsm.FSMStateHandshake:
			s.handleHandshakeState(plr, pkt)
		case fsm.FSMStateStatus:
			s.handleStatusState(plr, pkt)
		case fsm.FSMStateLogin:
			s.handleLoginState(plr, pkt)
		case fsm.FSMStatePlay:
			s.handlePlayState(plr, pkt)
		default:
			slog.Error("State not implemented", "id", pkt.ID(), "size", n, "state", plr.State.State())
		}
	}
}

func (s *Server) handleHandshakeState(plr *player.Player, pkt *packet.Packet) {
	handshakePkt, err := protocol.ReceiveHandshakePacket(pkt)
	if err != nil {
		slog.Error("Error reading handshake", "err", err.Error())
		return
	}

	slog.Info("New HandShake Packet", "nextState", handshakePkt.NextState)

	switch handshakePkt.NextState {
	case protocol.HandshakeNextStateStatus:
		plr.State.SetState(fsm.FSMStateStatus)
		// show motd
		statusRespPkt, err := protocol.CreateStatusResponsePacket(s.statusResponse)
		if err != nil {
			slog.Error("Error creating status response packet", "err", err.Error())
			return
		}

		err = plr.SendPacket(statusRespPkt)
		if err != nil {
			slog.Error("Error sending status response bytes", "err", err.Error())
			return
		}

	case protocol.HandshakeNextStateLogin:
		plr.State.SetState(fsm.FSMStateLogin)
	default:
		slog.Error("next state not implemented", "state", handshakePkt.NextState)
	}
}

func (s *Server) handleStatusState(plr *player.Player, pkt *packet.Packet) {
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

	err = plr.SendPacket(pingRespPkt)
	if err != nil {
		slog.Error("error sending ping response packet", "err", err.Error())
		return
	}
}

func (s *Server) handleLoginState(plr *player.Player, pkt *packet.Packet) {
	slog.Info("New Login Packet", "id", pkt.ID())

	switch pkt.ID() {
	case 0x00: // login start
		loginStartPkt, err := protocol.ReceiveLoginStartPacket(pkt)
		if err != nil {
			slog.Error("error receiving login start packet", "err", err.Error())
			return
		}

		slog.Info("Hello, Player!", "name", loginStartPkt.Name)

		plr.Name = loginStartPkt.Name
		plr.UUID = uuid.GenerateUUID() // generating a random UUID for now
		plr.IsLoggedIn = true

		// TODO: implement encryption!!!

		loginSuccessPkt, err := protocol.CreateLoginSuccessPacket(plr.UUID, loginStartPkt.Name)
		if err != nil {
			slog.Error("error creating login success packet", "err", err.Error())
			return
		}
		err = plr.SendPacket(loginSuccessPkt)
		if err != nil {
			slog.Error("error sending login success packet", "err", err.Error())
			return
		}

		joinGamePkt, err := protocol.CreateJoinGamePacket()
		if err != nil {
			slog.Error("error creating join game packet", "err", err.Error())
			return
		}
		err = plr.SendPacket(joinGamePkt)
		if err != nil {
			slog.Error("Error sending join game packet")
			return
		}

		// change the state to game mode
		plr.State.SetState(fsm.FSMStatePlay)

		// send the spawn position
		spawnPositionPkt, err := protocol.CreateSpawnPositionPacket()
		if err != nil {
			slog.Error("error creating spawn position packet", "err", err.Error())
		}
		err = plr.SendPacket(spawnPositionPkt)
		if err != nil {
			slog.Error("error sending spawn position")
			return
		}
	default:
		slog.Error("login id not implemented", "id", pkt.ID())
	}

}

func (s *Server) handlePlayState(plr *player.Player, pkt *packet.Packet) {
	switch pkt.ID() {
	case packet.IDClientKeepAlive:
		rndId, _ := pkt.Buffer().ReadInt()
		slog.Info("Client sent KeepAlive packet!", "id", rndId)
	case packet.IDClientPlayer:
		// This packet is used to indicate whether the player is on ground (walking/swimming), or airborne (jumping/falling).
		plr.Position.OnGround, _ = pkt.Buffer().ReadBool()
	case packet.IDClientClientSettings: // Sent when the player connects, or when settings are changed.
		clientSettings, err := protocol.ReceiveClientSettings(pkt)
		if err != nil {
			slog.Error("error receiving client settings packet", "err", err.Error())
			return
		}
		_ = clientSettings
		// I don't know exactly how to deal with it yet <-<
	case packet.IDClientPluginMessage: // Plugin Message
		pluginMessage, err := protocol.ReceivePluginMessage(pkt)
		if err != nil {
			slog.Error("error receiving plugin message packet", "err", err.Error())
			return
		}
		if pluginMessage.Channel == "MC|Brand" {
			pluginMessagePkt, err := protocol.CreatePluginMessagePacket(pluginMessage)
			if err != nil {
				slog.Error("error creating plugin message packet", "err", err.Error())
				return
			}

			err = plr.SendPacket(pluginMessagePkt)
			if err != nil {
				slog.Error("error sending plugin message packet", "err", err.Error())
				return
			}
		}
	case packet.IDClientPlayerPosition:
		playerPosition, err := protocol.ReceivePlayerPosition(pkt)
		if err != nil {
			slog.Error("error receiving player position", "err", err.Error())
			return
		}

		_ = playerPosition
	default:
		slog.Error("Play State not implemented yet", "id", pkt.ID())
	}
}

func (s *Server) closeConn(plr *player.Player) error {
	addr := plr.Conn.RemoteAddr().String()
	_, exists := s.players[addr]
	if exists {
		delete(s.players, addr)
	}
	slog.Info("Connection Closed", "name", plr.Name, "addr", addr)
	return plr.Conn.Close()
}
