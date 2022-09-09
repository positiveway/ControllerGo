package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"math"
	"strconv"
	"strings"
)

type EventT struct {
	dependentVars *DependentVariablesT
	eventType     EventTypeT
	btnOrAxis     BtnOrAxisT
	value         float64
	codeType      CodeTypeT
	code          CodeT

	update func(msg string)
	transformStickToDPadSC,
	transformToWingsSC,
	transformToPadReleasedEvent,
	transformAndFilter,
	fullReset,
	axisChanged,
	padReleased,
	match func()
}

func MakeEvent(dependentVars *DependentVariablesT) *EventT {
	event := &EventT{}
	event.dependentVars = dependentVars

	event.update = event.GetUpdateFunc()

	event.transformStickToDPadSC = event.GetTransformStickSCFunc()
	event.transformToWingsSC = event.GetTransformToWingsSCFunc()
	event.transformToPadReleasedEvent = event.GetTransformToPadReleasedFunc()
	event.transformAndFilter = event.GetTransformAndFilterFunc()
	event.fullReset = event.GetFullResetFunc()

	event.axisChanged = event.GetAxisChangedFunc()
	event.padReleased = event.GetPadReleasedFunc()
	event.match = event.GetMatchFunc()

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
	dependentVars := event.dependentVars

	curPressedStickButton := CurPressedStickButtonSC
	boundariesMap := dependentVars.cfg.PadsSticks.Stick.BoundariesMapSC
	zoneToBtnMap := initStickZoneBtnMap()

	LeftStick := dependentVars.LeftStick
	buttons := dependentVars.Buttons

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
				buttons.releaseButton(*curPressedStickButton)
				*curPressedStickButton = ""
			}
			if LeftStick.zoneCanBeUsed {
				*curPressedStickButton = gofuncs.GetOrPanic(zoneToBtnMap, LeftStick.zone)
				buttons.pressButton(*curPressedStickButton)
			}
		}
		event.btnOrAxis = BtnUnknown
	}
}

func (event *EventT) applyDeadzoneDS() bool {
	stickDeadzone := event.dependentVars.cfg.PadsSticks.Stick.DeadzoneDS

	applyDeadzone := func(value float64) float64 {
		if gofuncs.IsNotInit(value) {
			return value
		}
		if math.Abs(value) <= stickDeadzone {
			value = 0
		}
		return value
	}

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

func (event *EventT) GetTransformAndFilterFunc() func() {
	controllerInUse := event.dependentVars.cfg.ControllerInUse

	return func() {
		//gofuncs.Print("Before: ")
		//Event.print()

		switch event.eventType {
		case EvButtonPressed, EvButtonReleased:
			return
		}

		event.transformToPadReleasedEvent()

		switch controllerInUse {
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
}

func (event *EventT) GetFullResetFunc() func() {
	dependentVars := event.dependentVars
	cfg := dependentVars.cfg

	RightPadStick := dependentVars.RightPadStick
	LeftStick := dependentVars.LeftStick
	LeftPad := dependentVars.LeftPad

	controllerInUse := cfg.ControllerInUse
	highPrecisionMode := dependentVars.HighPrecisionMode
	padsSticksMode := cfg.PadsSticks.Mode
	buttons := dependentVars.Buttons

	return func() {
		RightPadStick.Reset()
		LeftStick.Reset()

		switch controllerInUse {
		case SteamController:
			LeftPad.Reset()
		}

		*CurPressedStickButtonSC = ""

		padsSticksMode.SetToDefault()
		highPrecisionMode.Disable()

		buttons.releaseAll("")
	}
}

func (event *EventT) GetUpdateFunc() func(msg string) {
	EventTypeMap := initEventTypeMap()
	BtnAxisMap := InitBtnAxisMap()

	split := func(str, sep string) []string {
		if str == "" {
			gofuncs.PanicIsEmpty()
		}
		return strings.Split(str, sep)
	}

	strToFloat := func(str string) float64 {
		value, err := strconv.ParseFloat(str, 32)
		gofuncs.CheckErr(err)
		return value
	}

	strToCode := func(str string) CodeT {
		return CodeT(gofuncs.StrToInt(str))
	}

	startsWith := strings.HasSuffix

	var found bool

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
				valueAndCode := split(msg, ";")

				event.value = strToFloat(valueAndCode[0])
				if startsWith(msg, ";") {
					return
				}
				typeAndCode := split(valueAndCode[1], "(")
				event.codeType = CodeTypeT(typeAndCode[0])

				code := typeAndCode[1]
				event.code = strToCode(code[:len(code)-1])
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
