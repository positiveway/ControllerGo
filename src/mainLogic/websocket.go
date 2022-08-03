package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"net"
)

const gamepadConnectedMsg = "gamepadConnected"

var event Event

func RunWebSocket() {
	addr := net.UDPAddr{
		Port: SocketPort,
		IP:   net.ParseIP(SocketIP),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		gofuncs.Panic("Client is already running: %v", err)
	}
	print("Listening at %v", addr.String())

	p := make([]byte, 32)

	for {
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			print("Read err  %v", err)
			continue
		}

		//print(nn)
		//print(string(p[:nn]))

		event.update(string(p[:nn]))
	}
}
