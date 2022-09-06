package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"net"
)

const gamepadConnectedMsg = "gamepadConnected"

var Event EventT

func RunWebSocket() {
	addr := net.UDPAddr{
		Port: Cfg.WebSocket.Port,
		IP:   net.ParseIP(Cfg.WebSocket.IP),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		gofuncs.Panic("Client is already running: %v", err)
	}
	gofuncs.Print("Listening at %v", addr.String())

	p := make([]byte, 32)

	for {
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			gofuncs.Print("Read err  %v", err)
			continue
		}

		//gofuncs.Print(nn)
		//gofuncs.Print(string(p[:nn]))

		Event.update(string(p[:nn]))
	}
}
