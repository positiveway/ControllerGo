package mainLogic

import "github.com/positiveway/gofuncs"

func axisChanged() {
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
	event.print()

	switch event.eventType {
	case EvAxisChanged:
		axisChanged()
	case EvPadReleased:
		padReleased()
	case EvButtonChanged:
		buttonChanged(event.btnOrAxis, event.value)
	case EvDisconnected:
		gamepadDisconnected()
		gofuncs.Print("Gamepad disconnected")
	case EvConnected:
		gofuncs.Print("Gamepad connected")
	case EvDropped:
		gofuncs.Panic("Event dropped")
	default:
		gofuncs.Panic("Unsupported event type: %v", event.eventType)
	}
}
