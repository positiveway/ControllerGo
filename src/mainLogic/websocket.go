package mainLogic

import (
	"fmt"
	"net"
)

const gamepadConnectedMsg = "gamepadConnected"

func RunWebSocket() {
	addr := net.UDPAddr{
		Port: SocketPort,
		IP:   net.ParseIP(SocketIP),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panicMsg("Listen err %v\n", err)
	}
	fmt.Printf("Listen at %v\n", addr.String())

	p := make([]byte, 32)
	event := Event{}

	for {
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}

		//fmt.Println(nn)
		//fmt.Println(string(p[:nn]))

		event.update(string(p[:nn]))
		matchEvent(&event)
	}
}
