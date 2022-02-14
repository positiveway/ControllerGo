package main

import "fmt"

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
				switch event.btnOrAxis {
				case AxisLeftStickX:
					joystickTyping.leftCoords.setX(event.value)
					joystickTyping.updateZoneLeft()
				case AxisLeftStickY:
					joystickTyping.leftCoords.setY(event.value)
					joystickTyping.updateZoneLeft()
				case AxisRightStickX:
					joystickTyping.rightCoords.setX(event.value)
					joystickTyping.updateZoneRight()
				case AxisRightStickY:
					joystickTyping.rightCoords.setY(event.value)
					joystickTyping.updateZoneRight()
				}
			}

		case EvButtonChanged:
			detectTriggers(event)
		case EvButtonPressed:
			buttonPressed(event.btnOrAxis)
		case EvButtonReleased:
			buttonReleased(event.btnOrAxis)
		case EvDisconnected:
			fmt.Printf("Gamepad %v: disconnected\n", event.deviceID)
		case EvConnected:
			fmt.Printf("Gamepad %v: connected\n", event.deviceID)
		}
	}
}
