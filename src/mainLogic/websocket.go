package mainLogic

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func checkBytesAmount(bytesNumStr string, nn int) int {
	bytesAmount, err := strconv.Atoi(bytesNumStr)
	CheckErr(err)
	bytesAmount += len(bytesNumStr + ";")
	if bytesAmount > 1000 {
		panic("buffer overflow")
	}
	if nn != bytesAmount {
		panic("Incorrect buffer size")
	}
	return bytesAmount
}

func printEvents(events []Event, bytesAmount int, raddr *net.UDPAddr) {
	return
	batchStr := ""
	for _, event := range events {
		batchStr += fmt.Sprintf("%v;", event)
	}
	fmt.Printf("Bytes: %v; Event: %s Host: %v\n", bytesAmount, batchStr, raddr)
}

func convertToEvents(rawEvents []string) []Event {
	var events []Event

	rawEvents = rawEvents[1:]
	for _, rawEvent := range rawEvents {
		rawSlice := strings.Split(rawEvent, ",")
		id, eventType, btnOrAxis, value := rawSlice[0], rawSlice[1], rawSlice[2], rawSlice[3]
		event := makeEvent(id, eventType, btnOrAxis, value)
		events = append(events, event)
	}
	return events
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

	//maxSize := 300
	for {
		p := make([]byte, 1024)
		nn, raddr, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}

		msg := p[:nn]
		msgStr := string(msg)

		if msgStr[len(msgStr)-1] != ';' {
			panic("Invalid message format")
		}
		msgStr = msgStr[:len(msgStr)-1]
		rawEvents := strings.Split(msgStr, ";")

		bytesAmount := checkBytesAmount(rawEvents[0], nn)
		//if bytesAmount > maxSize {
		//	maxSize = bytesAmount
		//}

		if rawEvents[1] == gamepadConnectedMsg {
			fmt.Println("Gamepad connected")
			continue
		}

		events := convertToEvents(rawEvents)
		matchEvents(events)
		printEvents(events, bytesAmount, raddr)

		//fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, msg, raddr)
		//fmt.Printf("Max: %v\n", maxSize)
	}
}
