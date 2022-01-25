package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

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
		if nn > 1000 {
			panic("buffer overflow")
		}

		msg := p[:nn]
		//if nn > maxSize {
		//	maxSize = nn
		//}
		msg_str := string(msg)
		raw_events := strings.Split(msg_str, ";")
		raw_events = raw_events[:len(raw_events)-1]
		batch_str := ""
		for _, raw_event := range raw_events {
			rawSlice := strings.Split(raw_event, ",")
			id, eventType, btnOrAxis, value := rawSlice[0], rawSlice[1], rawSlice[2], rawSlice[3]
			event := makeEvent(id, eventType, btnOrAxis, value)
			batch_str += fmt.Sprintf("%v;", event)
		}

		fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, msg, raddr)
		fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, batch_str, raddr)
		//fmt.Printf("Max: %v\n", maxSize)

	}
}
