package main

import (
	"fmt"
	"github.com/karalabe/hid"
	dualshock "github.com/kvartborg/go-dualshock"
	"log"
)

func mainDS() {
	vendorID, productID := uint16(1356), uint16(1476)
	devices := hid.Enumerate(vendorID, productID)

	if len(devices) == 0 {
		log.Fatal("no dualshock controller where found")
	}

	device, err := devices[0].Open()

	if err != nil {
		log.Fatal(err)
	}

	controller := dualshock.New(device)

	controller.Listen(func(state dualshock.State) {
		fmt.Println(state.Analog.L2)
	})
}
