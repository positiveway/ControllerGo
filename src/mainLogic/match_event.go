package mainLogic

func eventChanged() {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
		case false:
			switch event.btnOrAxis {
			case AxisDPadX:
				scrollMovement.setX()
			case AxisDPadY:
				scrollMovement.setY()
			case AxisRightStickX:
				//print("x: %v", event.value)
				mousePad.setX()
			case AxisRightStickY:
				//print("y: %v", event.value)
				mousePad.setY()
			}
			//scrollMovement.printCurState()
		case true:
			switch event.btnOrAxis {
			case AxisDPadX:
				joystickTyping.leftCoords.setDirectlyX()
				joystickTyping.updateLeftZone()
			case AxisDPadY:
				joystickTyping.leftCoords.setDirectlyY()
				joystickTyping.updateLeftZone()
			case AxisRightStickX:
				joystickTyping.rightCoords.setDirectlyX()
				joystickTyping.updateRightZone()
			case AxisRightStickY:
				joystickTyping.rightCoords.setDirectlyY()
				joystickTyping.updateRightZone()
			}
		}
	case true:
		switch event.btnOrAxis {
		case AxisRightStickX:
			movementCoords.setX()
		case AxisRightStickY:
			movementCoords.setY()
		case AxisDPadX:
			mousePad.setX()
		case AxisDPadY:
			mousePad.setY()
		}
	}
}

func matchEvent() {
	event.print()

	switch event.eventType {
	case EvAxisChanged:
		eventChanged()
	case EvButtonChanged:
		detectTriggers()
	case EvButtonPressed:
		buttonPressed()
	case EvButtonReleased:
		buttonReleased()
	case EvDisconnected:
		panicMsg("Gamepad disconnected")
	case EvConnected:
		print("Gamepad connected")
	case EvDropped:
		panicMsg("Event dropped")
	}
}
