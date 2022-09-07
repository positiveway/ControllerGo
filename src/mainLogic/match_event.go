package mainLogic

import "github.com/positiveway/gofuncs"

func (event *EventT) axisChanged() {
	switch event.btnOrAxis {
	case AxisRightPadStickX:
		RightPadStick.SetX(event.value)
	case AxisRightPadStickY:
		RightPadStick.SetY(event.value)
	}

	switch Cfg.ControllerInUse {
	case SteamController:
		switch event.btnOrAxis {
		case AxisLeftPadX:
			LeftPad.SetX(event.value)
		case AxisLeftPadY:
			LeftPad.SetY(event.value)
		}
	case DualShock:
		switch event.btnOrAxis {
		case AxisLeftStickX:
			LeftStick.SetX(event.value)
		case AxisLeftStickY:
			LeftStick.SetY(event.value)
		}
	}

	TypeLetter()
}

func (event *EventT) padReleased() {
	switch event.btnOrAxis {
	case AxisRightPadStickX, AxisRightPadStickY:
		RightPadStick.Reset()
	}

	switch Cfg.ControllerInUse {
	case SteamController:
		switch event.btnOrAxis {
		case AxisLeftPadX, AxisLeftPadY:
			LeftPad.Reset()
		}
	case DualShock:
		switch event.btnOrAxis {
		case AxisLeftStickX, AxisLeftStickY:
			LeftStick.Reset()
		}
	}
}

func (event *EventT) match() {
	//gofuncs.Print("After: ")
	//event.print()

	switch event.eventType {
	case EvAxisChanged:
		event.axisChanged()
	case EvPadReleased:
		event.padReleased()
	case EvButtonChanged:
		buttonChanged(event.btnOrAxis, event.value)
	case EvDisconnected:
		fullReset()
		gofuncs.Print("Gamepad disconnected")
	case EvConnected:
		gofuncs.Print("Gamepad connected")
	case EvDropped:
		gofuncs.Panic("Event dropped")
	default:
		gofuncs.Panic("Unsupported event type: %v", event.eventType)
	}
}
