package mainLogic

func eventChanged() {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
		case false:
			switch event.btnOrAxis {
			case AxisLeftPadX:
				scrollMovement.setX()
			case AxisLeftPadY:
				scrollMovement.setY()
			case AxisRightPadX:
				//print("x: %v", event.value)
				mousePad.setX()
			case AxisRightPadY:
				//print("y: %v", event.value)
				mousePad.setY()
			}
			//scrollMovement.printCurState()
		case true:
			switch event.btnOrAxis {
			case AxisLeftPadX:
				joystickTyping.leftCoords.setDirectlyX()
				joystickTyping.updateLeftZone()
			case AxisLeftPadY:
				joystickTyping.leftCoords.setDirectlyY()
				joystickTyping.updateLeftZone()
			case AxisRightPadX:
				joystickTyping.rightCoords.setDirectlyX()
				joystickTyping.updateRightZone()
			case AxisRightPadY:
				joystickTyping.rightCoords.setDirectlyY()
				joystickTyping.updateRightZone()
			}
		}
	case true:
		switch event.btnOrAxis {
		case AxisRightPadX:
			movementCoords.setX()
		case AxisRightPadY:
			movementCoords.setY()
		case AxisLeftPadX:
			mousePad.setX()
		case AxisLeftPadY:
			mousePad.setY()
		}
	}
}

func padReleased() {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
		case false:
			switch event.btnOrAxis {
			case AxisLeftPadX, AxisLeftPadY:
				scrollMovement.reset()
			case AxisRightPadX, AxisRightPadY:
				mousePad.reset()
			}
		case true:
			switch event.btnOrAxis {
			case AxisLeftPadX, AxisLeftPadY:
				joystickTyping.leftCoords.reset()
				joystickTyping.updateLeftZone()
			case AxisRightPadX, AxisRightPadY:
				joystickTyping.rightCoords.reset()
				joystickTyping.updateRightZone()
			}
		}
	case true:
		switch event.btnOrAxis {
		case AxisLeftPadX, AxisLeftPadY:
			mousePad.reset()
		case AxisRightPadX, AxisRightPadY:
			movementCoords.reset()
		}
	}
}

func matchEvent() {
	switch event.eventType {
	case EvAxisChanged:
		eventChanged()
	case EvPadFirstTouched:
		return
	case EvPadReleased:
		padReleased()
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
