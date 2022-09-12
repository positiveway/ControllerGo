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

	applyDeadzoneDS func() bool
	update          func(msg string)
	fixButtonNamesForSC,
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

	event.applyDeadzoneDS = event.GetApplyDeadzoneFunc()
	event.fixButtonNamesForSC = event.GetFixButtonNamesForSCFunc()
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

func (event *EventT) GetFixButtonNamesForSCFunc() func() {
	allBtnAxis := event.dependentVars.allBtnAxis
	BtnX := allBtnAxis.BtnX
	BtnY := allBtnAxis.BtnY

	return func() {
		switch event.btnOrAxis {
		case BtnY:
			event.btnOrAxis = BtnX
		case BtnX:
			event.btnOrAxis = BtnY
		}
	}
}

func (event *EventT) GetTransformToPadReleasedFunc() func() {
	allBtnAxis := event.dependentVars.allBtnAxis
	PadAndStickAxes := allBtnAxis.initPadAndStickAxes()

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
	allBtnAxis := event.dependentVars.allBtnAxis
	UnknownCodesResolvingMapSC := allBtnAxis.initUnknownCodesMapSC()
	BtnUnknown := allBtnAxis.BtnUnknown

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
	allBtnAxis := dependentVars.allBtnAxis
	BtnUnknown := allBtnAxis.BtnUnknown
	AxisLeftStickX := allBtnAxis.AxisLeftStickX
	AxisLeftStickY := allBtnAxis.AxisLeftStickY

	curPressedStickButtonSC := dependentVars.CurPressedStickButtonSC
	boundariesMap := dependentVars.cfg.PadsSticks.Stick.BoundariesMapSC
	zoneToBtnMap := allBtnAxis.initStickZoneBtnMap()

	stick := dependentVars.LeftStick
	buttons := dependentVars.Buttons

	return func() {
		eventType := event.eventType
		btnOrAxis := event.btnOrAxis

		isStickEvent := (eventType == EvPadReleased || eventType == EvAxisChanged) &&
			(btnOrAxis == AxisLeftStickX || btnOrAxis == AxisLeftStickY)
		if !isStickEvent {
			return
		}
		switch eventType {
		case EvPadReleased:
			stick.Reset()
		case EvAxisChanged:
			switch btnOrAxis {
			case AxisLeftStickX:
				stick.SetX(event.value)
			case AxisLeftStickY:
				stick.SetY(event.value)
			}
		}

		stick.ReCalculateZone(boundariesMap)

		if stick.zoneChanged {
			if *curPressedStickButtonSC != "" {
				buttons.releaseButton(*curPressedStickButtonSC)
				*curPressedStickButtonSC = ""
			}
			if stick.zoneCanBeUsed {
				*curPressedStickButtonSC = gofuncs.GetOrPanic(zoneToBtnMap, stick.zone)
				buttons.pressImmediately(*curPressedStickButtonSC)
			}
		}
		event.btnOrAxis = BtnUnknown
	}
}

func (event *EventT) GetApplyDeadzoneFunc() func() bool {
	dependentVars := event.dependentVars
	allBtnAxis := dependentVars.allBtnAxis

	AxisLeftStickX := allBtnAxis.AxisLeftStickX
	AxisLeftStickY := allBtnAxis.AxisLeftStickY

	AxisRightPadStickX := allBtnAxis.AxisRightPadStickX
	AxisRightPadStickY := allBtnAxis.AxisRightPadStickY

	stickDeadzone := dependentVars.cfg.PadsSticks.Stick.DeadzoneDS

	applyDeadzone := func(value float64) float64 {
		if gofuncs.IsNotInit(value) {
			return value
		}
		if math.Abs(value) <= stickDeadzone {
			value = 0
		}
		return value
	}

	return func() bool {
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
}

func (event *EventT) GetTransformAndFilterFunc() func() {
	dependentVars := event.dependentVars
	allBtnAxis := dependentVars.allBtnAxis
	BtnUnknown := allBtnAxis.BtnUnknown

	controllerInUse := dependentVars.cfg.ControllerInUse

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
			event.fixButtonNamesForSC()
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

	curPressedStickButtonSC := dependentVars.CurPressedStickButtonSC
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

		*curPressedStickButtonSC = ""

		padsSticksMode.SetToDefault()
		highPrecisionMode.Disable()

		buttons.releaseAll("")
	}
}

func (event *EventT) GetUpdateFunc() func(msg string) {
	dependentVars := event.dependentVars
	allBtnAxis := dependentVars.allBtnAxis

	BtnAxisMap := allBtnAxis.InitBtnAxisMap()
	EventTypeMap := initEventTypeMap()

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
