package mainLogic

func matchEvent(event *Event) {
	//event.print()

	switch event.eventType {
	case EvAxisChanged:
		switch SteamController {
		case true:
			eventChangedSteam(event)
		case false:
			eventChangedDS(event)
		}
	case EvButtonChanged:
		detectTriggers(event)
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
