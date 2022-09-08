package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"net"
)

const gamepadConnectedMsg = "gamepadConnected"

func (dependentVars *DependentVariablesT) RunWebSocket() {
	addr := net.UDPAddr{
		Port: dependentVars.cfg.WebSocket.Port,
		IP:   net.ParseIP(dependentVars.cfg.WebSocket.IP),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		gofuncs.Panic("Client is already running: %v", err)
	}
	gofuncs.Print("Listening at %v", addr.String())

	p := make([]byte, 32)

	event := MakeEvent(dependentVars)

	for {
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			gofuncs.Print("Read err  %v", err)
			continue
		}

		//gofuncs.Print(nn)
		//gofuncs.Print(string(p[:nn]))

		event.update(string(p[:nn]))
	}
}
