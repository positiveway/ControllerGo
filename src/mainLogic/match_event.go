package mainLogic

func matchEvent() {
	//event.print()

	switch event.eventType {
	case EvAxisChanged:
		switch SteamController {
		case true:
			eventChangedSteam()
		case false:
			eventChangedDS()
		}
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
