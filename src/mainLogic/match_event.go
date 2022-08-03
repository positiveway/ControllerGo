package mainLogic

import "github.com/positiveway/gofuncs"

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
	event.print()

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
		gofuncs.Panic("Event dropped")
	}
}
