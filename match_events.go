package main

import (
	"github.com/bendahl/uinput"
	"time"
)

const mouseInterval = 25 * time.Millisecond
const mouseMaxMove float64 = 20

const scrollMinValue = 35
const scrollMaxValue int64 = 200

const scrollMaxRange = float64(scrollMaxValue - scrollMinValue)

var mouseMovement = Coords{}
var scrollMovement = Coords{}
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
			case AxisRightStickX:
				scrollMovement.setX(event.value)
			case AxisRightStickY:
				scrollMovement.setY(event.value)
			}
		case EvButtonChanged:
			println("hui")
		}
	}

}
