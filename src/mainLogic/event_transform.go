package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"strconv"
)

var PadsAxesOffsetMap = map[BtnOrAxisT]float64{
	AxisLeftPadX:  0,
	AxisLeftPadY:  0,
	AxisRightPadX: 0,
	AxisRightPadY: 0,
}

var PadAndStickAxes = []BtnOrAxisT{
	AxisLeftPadX,
	AxisLeftPadY,
	AxisRightPadX,
	AxisRightPadY,
	AxisStickX,
	AxisStickY,
}

type Event struct {
	eventType EventTypeT
	btnOrAxis BtnOrAxisT
	value     float64
	codeType  CodeTypeT
	code      CodeT
}

func (event *Event) applyPadAxesAdjustments() {
	if event.eventType == EvAxisChanged {
		event.value -= gofuncs.GetOrDefault(PadsAxesOffsetMap, event.btnOrAxis, 0)
	}
}

func (event *Event) fixButtonNamesForSteamController() {
	switch event.btnOrAxis {
	case BtnY:
		event.btnOrAxis = BtnX
	case BtnX:
		event.btnOrAxis = BtnY
	}
}

func (event *Event) transformToPadReleasedEvent() {
	if gofuncs.Contains(PadAndStickAxes, event.btnOrAxis) &&
		event.eventType == EvAxisChanged && event.value == 0.0 {

		event.eventType = EvPadReleased
		event.code = 0
		event.codeType = ""
	}
}

func (event *Event) transformToWings() {
	if event.btnOrAxis == BtnUnknown && event.codeType == CTKey {
		if btn, found := UnknownCodesResolvingMap[event.code]; found {
			event.btnOrAxis = btn
		}
	}
}

var StickZoneToBtnMap = map[Zone]BtnOrAxisT{
	ZoneRight: BtnStickRight,
	ZoneUp:    BtnStickUp,
	ZoneLeft:  BtnStickLeft,
	ZoneDown:  BtnStickDown,
}

var curPressedStickButton BtnOrAxisT

func (event *Event) transformStickToDPad() {
	isStickEvent := (event.eventType == EvPadReleased || event.eventType == EvAxisChanged) &&
		(event.btnOrAxis == AxisStickX || event.btnOrAxis == AxisStickY)
	if !isStickEvent {
		return
	}
	switch event.eventType {
	case EvPadReleased:
		Stick.Reset()
	case EvAxisChanged:
		switch event.btnOrAxis {
		case AxisStickX:
			Stick.SetX()
		case AxisStickY:
			Stick.SetY()
		}
	}

	Stick.ReCalculateZone(Cfg.StickBoundariesMap)

	if Stick.zoneChanged {
		if curPressedStickButton != "" {
			event.btnOrAxis = curPressedStickButton
			event.eventType = EvButtonReleased
			curPressedStickButton = ""
			matchEvent()
		}
		if Stick.zoneCanBeUsed {
			stickBtn := StickZoneToBtnMap[Stick.zone]
			curPressedStickButton = stickBtn

			event.btnOrAxis = stickBtn
			event.eventType = EvButtonPressed
			matchEvent()
		}
	}
	event.btnOrAxis = BtnUnknown
}

func (event *Event) transformAndFilter() {
	//printDebug("Before: ")
	//event.print()

	if Cfg.ControllerInUse.SteamController {
		event.fixButtonNamesForSteamController()
		event.transformToWings()
	}

	event.transformToPadReleasedEvent()

	if Cfg.ControllerInUse.SteamController {
		event.applyPadAxesAdjustments()
		event.transformStickToDPad()
	}

	if event.btnOrAxis == BtnUnknown {
		return
	}

	//printDebug("After: ")

	matchEvent()
}

func (event *Event) update(msg string) {
	var found bool
	var err error

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

func (event *Event) print() {
	gofuncs.PrintDebug("%s \"%s\": %.2f",
		gofuncs.TrimAnyPrefix(string(event.eventType), "Ev"),
		gofuncs.TrimAnyPrefix(string(event.btnOrAxis), "Btn", "Axis"),
		event.value)
}
