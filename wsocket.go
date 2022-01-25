package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func checkBytesAmount(bytesNumStr string, nn int) (bytesAmount int) {
	bytesAmount, err := strconv.Atoi(bytesNumStr)
	check_err(err)
	bytesAmount += len(bytesNumStr + ";")
	if bytesAmount > 1000 {
		panic("buffer overflow")
	}
	if nn != bytesAmount {
		panic("Incorrect buffer size")
	}
	return
}

func mainWS() {
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("0.0.0.0"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Listen err %v\n", err)
		os.Exit(-1)
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
		//if nn > maxSize {
		//	maxSize = nn
		//}

		msgStr := string(msg)

		if msgStr[len(msgStr)-1] != ';' {
			panic("Invalid message format")
		}
		rawEvents := strings.Split(msgStr, ";")

		bytesAmount := checkBytesAmount(rawEvents[0], nn)

		rawEvents = rawEvents[1 : len(rawEvents)-1]
		batchStr := ""
		for _, rawEvent := range rawEvents {
			rawSlice := strings.Split(rawEvent, ",")
			id, eventType, btnOrAxis, value := rawSlice[0], rawSlice[1], rawSlice[2], rawSlice[3]
			event := makeEvent(id, eventType, btnOrAxis, value)
			batchStr += fmt.Sprintf("%v;", event)
		}

		//fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, msg, raddr)
		fmt.Printf("Bytes: %v; Event: %s Host: %v\n", bytesAmount, batchStr, raddr)
		//fmt.Printf("Max: %v\n", maxSize)

	}
}
