package main

import (
	"github.com/bendahl/uinput"
	"time"
)

const INTERVAL = 25 * time.Millisecond

var mouseMovement Movement = makeMovement(mouseAccelStep, mouseUpdateEveryN)

func matchEvents(events []Event, mouse uinput.Mouse) {
	for _, event := range events {
		switch event.eventType {
		case EvAxisChanged:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				mouseMovement.SetX(event.value)
			case AxisLeftStickY:
				mouseMovement.SetY(event.value)
				//case AxisRightStickX:
				//	scrollMovement.SetX(event.value)
				//case AxisRightStickY:
				//	scrollMovement.SetY(event.value)
			}
		case EvButtonChanged:
			println("hui")
		}
	}

	mouseMovement.UpdateAccel()
	xForce, yForce := mouseMovement.CalcForces()
	err := mouse.Move(xForce, yForce)
	check_err(err)
}
