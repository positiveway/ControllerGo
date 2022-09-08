package mainLogic

import "github.com/positiveway/gofuncs"

func (event *EventT) GetAxisChangedFunc() func() {
	dependentVars := event.dependentVars
	controllerInUse := dependentVars.cfg.ControllerInUse

	RightPadStick := dependentVars.RightPadStick
	LeftPad := dependentVars.LeftPad
	LeftStick := dependentVars.LeftStick

	typing := dependentVars.Typing

	return func() {
		switch event.btnOrAxis {
		case AxisRightPadStickX:
			RightPadStick.SetX(event.value)
		case AxisRightPadStickY:
			RightPadStick.SetY(event.value)
		}

		switch controllerInUse {
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

		typing.typeLetter()
	}
}

func (event *EventT) GetPadReleasedFunc() func() {
	dependentVars := event.dependentVars
	controllerInUse := dependentVars.cfg.ControllerInUse

	RightPadStick := dependentVars.RightPadStick
	LeftPad := dependentVars.LeftPad
	LeftStick := dependentVars.LeftStick

	return func() {
		switch event.btnOrAxis {
		case AxisRightPadStickX, AxisRightPadStickY:
			RightPadStick.Reset()
		}

		switch controllerInUse {
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
}

func (event *EventT) GetMatchFunc() func() {
	buttons := event.dependentVars.Buttons

	return func() {
		//gofuncs.Print("After: ")
		//event.print()

		switch event.eventType {
		case EvAxisChanged:
			event.axisChanged()
		case EvPadReleased:
			event.padReleased()
		case EvButtonChanged:
			buttons.buttonChanged(event.btnOrAxis, event.value)
		case EvDisconnected:
			event.fullReset()
			gofuncs.Print("Gamepad disconnected")
		case EvConnected:
			gofuncs.Print("Gamepad connected")
		case EvDropped:
			gofuncs.Panic("Event dropped")
		default:
			gofuncs.Panic("Unsupported event type: %v", event.eventType)
		}
	}
}
