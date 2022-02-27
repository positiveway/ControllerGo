package mainLogic

import "fmt"

func matchEvent(event *Event) {
	switch event.eventType {
	case EvAxisChanged:
		switch GamesModeOn {
		case false:
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
					joystickTyping.leftCoords.setDirectlyX(event.value)
					joystickTyping.updateLeftZone()
				case AxisLeftStickY:
					joystickTyping.leftCoords.setDirectlyY(event.value)
					joystickTyping.updateLeftZone()
				case AxisRightStickX:
					joystickTyping.rightCoords.setDirectlyX(event.value)
					joystickTyping.updateRightZone()
				case AxisRightStickY:
					joystickTyping.rightCoords.setDirectlyY(event.value)
					joystickTyping.updateRightZone()
				}
			}
		case true:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				movementCoords.setX(event.value)
			case AxisLeftStickY:
				movementCoords.setY(event.value)
			case AxisRightStickX:
				mouseMovement.setX(event.value)
			case AxisRightStickY:
				mouseMovement.setY(event.value)
			}
		}

	case EvButtonChanged:
		detectTriggers(event)
	case EvButtonPressed:
		buttonPressed(event.btnOrAxis)
	case EvButtonReleased:
		buttonReleased(event.btnOrAxis)
	case EvDisconnected:
		fmt.Printf("Gamepad disconnected\n")
	case EvConnected:
		fmt.Printf("Gamepad connected\n")
	case EvDropped:
		panicMsg("Event dropped\n")
	}
}
