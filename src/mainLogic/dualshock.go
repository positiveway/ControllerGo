package mainLogic

func eventChangedDS(event *Event) {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
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
}
