package mainLogic

import (
	"fmt"
	"net"
	"strings"
)

func convertToEvent(rawEvent string) Event {
	rawSlice := strings.Split(rawEvent, ",")
	eventType, btnOrAxis, value := rawSlice[0], rawSlice[1], rawSlice[2]
	return makeEvent(eventType, btnOrAxis, value)
}

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

	for {
		p := make([]byte, 64)
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		//fmt.Println(nn)

		msg := p[:nn]
		msgStr := string(msg)

		event := convertToEvent(msgStr)
		matchEvent(event)

		//fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, msgStr, raddr)
	}
}
