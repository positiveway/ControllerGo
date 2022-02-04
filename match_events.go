package main

import (
	"github.com/bendahl/uinput"
)

var mouse uinput.Mouse
var keyboard uinput.Keyboard

func matchEvents(events []Event) {
	for _, event := range events {
		switch event.eventType {
		case EvAxisChanged:
			switch typingMode.get() {
			case false:
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
			case true:
				println("jopa")
			}

		case EvButtonChanged:
			detectTriggers(event)
		case EvButtonPressed:
			buttonPressed(event.btnOrAxis)
		case EvButtonReleased:
			buttonReleased(event.btnOrAxis)
		}
	}

}
