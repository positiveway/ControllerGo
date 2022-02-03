package main

import (
	"github.com/bendahl/uinput"
	"time"
)

const INTERVAL = 25 * time.Millisecond

var mouseMovement = Coords{}
var mouse uinput.Mouse

func matchEvents(events []Event) {
	for _, event := range events {
		switch event.eventType {
		case EvAxisChanged:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				mouseMovement.setX(event.value)
			case AxisLeftStickY:
				mouseMovement.setY(event.value)
				//case AxisRightStickX:
				//	scrollMovement.SetX(event.value)
				//case AxisRightStickY:
				//	scrollMovement.SetY(event.value)
			}
		case EvButtonChanged:
			println("hui")
		}
	}

}
