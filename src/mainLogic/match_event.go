package mainLogic

func eventChanged() {
	switch event.btnOrAxis {
	case AxisLeftPadX:
		LeftPad.SetX()
	case AxisLeftPadY:
		LeftPad.SetY()
	case AxisRightPadX:
		RightPad.SetX()
	case AxisRightPadY:
		RightPad.SetY()
	}
	TypeLetter()
}

func padReleased() {
	switch event.btnOrAxis {
	case AxisLeftPadX, AxisLeftPadY:
		LeftPad.Reset()
	case AxisRightPadX, AxisRightPadY:
		RightPad.Reset()
	}
}

func gamepadDisconnected() {
	LeftPad.Reset()
	RightPad.Reset()
	Stick.Reset()

	releaseAll()
}

func matchEvent() {
	if PrintTypingDebugInfo {
		print("%v \"%v\": %.2f", event.eventType, event.btnOrAxis, event.value)
	}

	switch event.eventType {
	case EvAxisChanged:
		eventChanged()
	case EvPadReleased:
		padReleased()
	case EvButtonChanged:
		detectTriggers()
	case EvButtonPressed:
		buttonPressed()
	case EvButtonReleased:
		buttonReleased()
	case EvDisconnected:
		gamepadDisconnected()
		print("Gamepad disconnected")
	case EvConnected:
		print("Gamepad connected")
	case EvDropped:
		panicMsg("Event dropped")
	}
}
