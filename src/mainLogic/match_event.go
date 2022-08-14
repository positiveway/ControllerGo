package mainLogic

import "github.com/positiveway/gofuncs"

func axisChanged() {
	switch Event.btnOrAxis {
	case AxisLeftPadX:
		LeftPad.SetX()
	case AxisLeftPadY:
		LeftPad.SetY()
	case AxisRightPadStickX:
		RightPadStick.SetX()
	case AxisRightPadStickY:
		RightPadStick.SetY()
	}
	TypeLetter()
}

func padReleased() {
	switch Event.btnOrAxis {
	case AxisLeftPadX, AxisLeftPadY:
		LeftPad.Reset()
	case AxisRightPadStickX, AxisRightPadStickY:
		RightPadStick.Reset()
	}
}

func gamepadDisconnected() {
	LeftPad.Reset()
	RightPadStick.Reset()
	LeftStick.Reset()

	releaseAll()
}

func matchEvent() {
	//Event.print()

	switch Event.eventType {
	case EvAxisChanged:
		axisChanged()
	case EvPadReleased:
		padReleased()
	case EvButtonChanged:
		buttonChanged(Event.btnOrAxis, Event.value)
	case EvDisconnected:
		gamepadDisconnected()
		gofuncs.Print("Gamepad disconnected")
	case EvConnected:
		gofuncs.Print("Gamepad connected")
	case EvDropped:
		gofuncs.Panic("Event dropped")
	default:
		gofuncs.Panic("Unsupported event type: %v", Event.eventType)
	}
}
