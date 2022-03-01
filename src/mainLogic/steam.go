package mainLogic

func detectDPadButton() {

}

func eventChangedSteam(event *Event) {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
		case false:
			switch event.btnOrAxis {
			case AxisDPadX:
				scrollMovement.setX(event.value)
			case AxisDPadY:
				scrollMovement.setY(event.value)
			case AxisRightStickX:
				mouseMovement.setX(event.value)
			case AxisRightStickY:
				mouseMovement.setY(event.value)
			}
		case true:
			switch event.btnOrAxis {
			case AxisDPadX:
				joystickTyping.leftCoords.setDirectlyX(event.value)
				joystickTyping.updateLeftZone()
			case AxisDPadY:
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
		case AxisDPadX:
			movementCoords.setX(event.value)
		case AxisDPadY:
			movementCoords.setY(event.value)
		case AxisRightStickX:
			mouseMovement.setX(event.value)
		case AxisRightStickY:
			mouseMovement.setY(event.value)
		}
	}
}
