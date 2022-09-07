package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"strconv"
)

type EventT struct {
	eventType EventTypeT
	btnOrAxis BtnOrAxisT
	value     float64
	codeType  CodeTypeT
	code      CodeT

	update func(msg string)
	transformStickToDPadSC,
	transformToWingsSC,
	transformToPadReleasedEvent func()
}

func MakeEvent() *EventT {
	event := &EventT{}

	event.update = event.GetUpdateFunc()
	event.transformStickToDPadSC = event.GetTransformStickSCFunc()
	event.transformToWingsSC = event.GetTransformToWingsSCFunc()
	event.transformToPadReleasedEvent = event.GetTransformToPadReleasedFunc()

	return event
}

func (event *EventT) fixButtonNamesForSteamController() {
	switch event.btnOrAxis {
	case BtnY:
		event.btnOrAxis = BtnX
	case BtnX:
		event.btnOrAxis = BtnY
	}
}

func (event *EventT) GetTransformToPadReleasedFunc() func() {
	PadAndStickAxes := initPadAndStickAxes()

	return func() {
		if gofuncs.Contains(PadAndStickAxes, event.btnOrAxis) &&
			event.eventType == EvAxisChanged && event.value == 0 {

			event.eventType = EvPadReleased
			event.code = 0
			event.codeType = ""
		}
	}
}

func (event *EventT) GetTransformToWingsSCFunc() func() {
	UnknownCodesResolvingMapSC := initUnknownCodesMapSC()

	return func() {
		if event.btnOrAxis == BtnUnknown && event.codeType == CTKey {
			if btn, found := UnknownCodesResolvingMapSC[event.code]; found {
				event.btnOrAxis = btn
			}
		}
	}
}

func (event *EventT) GetTransformStickSCFunc() func() {
	curPressedStickButton := CurPressedStickButtonSC
	boundariesMap := Cfg.PadsSticks.Stick.BoundariesMapSC
	zoneToBtnMap := initStickZoneBtnMap()

	return func() {
		isStickEvent := (event.eventType == EvPadReleased || event.eventType == EvAxisChanged) &&
			(event.btnOrAxis == AxisLeftStickX || event.btnOrAxis == AxisLeftStickY)
		if !isStickEvent {
			return
		}
		switch event.eventType {
		case EvPadReleased:
			LeftStick.Reset()
		case EvAxisChanged:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				LeftStick.SetX(event.value)
			case AxisLeftStickY:
				LeftStick.SetY(event.value)
			}
		}

		LeftStick.ReCalculateZone(boundariesMap)

		if LeftStick.zoneChanged {
			if *curPressedStickButton != "" {
				releaseButton(*curPressedStickButton)
				*curPressedStickButton = ""
			}
			if LeftStick.zoneCanBeUsed {
				*curPressedStickButton = gofuncs.GetOrPanic(zoneToBtnMap, LeftStick.zone)
				pressButton(*curPressedStickButton)
			}
		}
		event.btnOrAxis = BtnUnknown
	}
}

func (event *EventT) applyDeadzoneDS() bool {
	if event.eventType == EvAxisChanged {
		switch event.btnOrAxis {
		case AxisLeftStickX, AxisLeftStickY, AxisRightPadStickX, AxisRightPadStickY:
			event.value = applyDeadzone(event.value)
			if event.value == 0 {
				return true
			}
		}
	}
	return false
}

func (event *EventT) transformAndFilter() {
	//gofuncs.Print("Before: ")
	//Event.print()

	switch event.eventType {
	case EvButtonPressed, EvButtonReleased:
		return
	}

	event.transformToPadReleasedEvent()

	switch Cfg.ControllerInUse {
	case SteamController:
		event.fixButtonNamesForSteamController()
		event.transformToWingsSC()
		event.transformStickToDPadSC()
	case DualShock:
		if event.applyDeadzoneDS() {
			return
		}
	}

	if event.btnOrAxis == BtnUnknown {
		return
	}

	event.match()
}

func fullReset() {
	RightPadStick.Reset()
	LeftStick.Reset()

	switch Cfg.ControllerInUse {
	case SteamController:
		LeftPad.Reset()
	}

	*CurPressedStickButtonSC = ""

	releaseAll("")
}

func (event *EventT) GetUpdateFunc() func(msg string) {
	EventTypeMap := initEventTypeMap()
	BtnAxisMap := InitBtnAxisMap()
	var found bool
	var err error

	return func(msg string) {
		event.eventType, found = EventTypeMap[msg[0]]
		if !found {
			gofuncs.PanicMisspelled(string(msg[0]))
		}
		if event.eventType != EvConnected && event.eventType != EvDisconnected && event.eventType != EvDropped {
			event.btnOrAxis, found = BtnAxisMap[msg[1]]
			if !found {
				gofuncs.PanicMisspelled(string(msg[1]))
			}
			if event.eventType == EvAxisChanged || event.eventType == EvButtonChanged {
				msg = msg[2:]
				valueAndCode := gofuncs.Split(msg, ";")

				event.value, err = strconv.ParseFloat(valueAndCode[0], 32)
				gofuncs.CheckErr(err)

				if gofuncs.StartsWith(msg, ";") {
					return
				}
				typeAndCode := gofuncs.Split(valueAndCode[1], "(")
				event.codeType = CodeTypeT(typeAndCode[0])

				code := typeAndCode[1]
				codeNum := gofuncs.StrToInt(code[:len(code)-1])
				event.code = CodeT(codeNum)
			}
		}
		event.transformAndFilter()
	}
}

func (event *EventT) print() {
	gofuncs.Print("%s \"%s\": %.2f; %s: %v",
		gofuncs.TrimAnyPrefix(string(event.eventType), "Ev"),
		gofuncs.TrimAnyPrefix(string(event.btnOrAxis), "Btn", "Axis"),
		event.value, event.codeType, event.code)
}
