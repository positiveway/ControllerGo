package mainLogic

func eventChanged() {
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
				//print("x: %v", event.value)
				mousePad.x = event.value
				moveMouse()
			case AxisRightStickY:
				//print("y: %v", event.value)
				mousePad.y = event.value
				moveMouse()
			}
			//scrollMovement.printCurState()
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
			mousePad.x = event.value
			moveMouse()
		case AxisRightStickY:
			mousePad.y = event.value
			moveMouse()
		}
	}
}

func matchEvent() {
	//event.print()

	switch event.eventType {
	case EvAxisChanged:
		eventChanged()
	case EvButtonChanged:
		detectTriggers()
	case EvButtonPressed:
		buttonPressed(event.btnOrAxis)
	case EvButtonReleased:
		buttonReleased(event.btnOrAxis)
	case EvDisconnected:
		panicMsg("Gamepad disconnected")
	case EvConnected:
		print("Gamepad connected")
	case EvDropped:
		panicMsg("Event dropped")
	}
}
