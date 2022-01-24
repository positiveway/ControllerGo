package main

import (
	"fmt"
	"net"
	"os"
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
		fmt.Printf("Bytes: %v; Event: %s Host: %v\n", nn, msg, raddr)
		//fmt.Printf("Max: %v\n", maxSize)

	}
}
