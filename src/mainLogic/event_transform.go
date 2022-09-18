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
	applyDeadzoneDS,
	fixButtonNamesForSC,
	transformStickToDPadSC,
	transformToWingsSC,
	transformAndFilter,
	fullReset,
	axisChanged,
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
	event.transformAndFilter = event.GetTransformAndFilterFunc()
	event.fullReset = event.GetFullResetFunc()

	event.axisChanged = event.GetAxisChangedFunc()
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

		isStickEvent := eventType == EvAxisChanged &&
			(btnOrAxis == AxisLeftStickX || btnOrAxis == AxisLeftStickY)
		if !isStickEvent {
			return
		}
		switch btnOrAxis {
		case AxisLeftStickX:
			stick.SetX(event.value)
		case AxisLeftStickY:
			stick.SetY(event.value)
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

func (event *EventT) GetApplyDeadzoneFunc() func() {
	dependentVars := event.dependentVars
	allBtnAxis := dependentVars.allBtnAxis

	AxisLeftStickX := allBtnAxis.AxisLeftStickX
	AxisLeftStickY := allBtnAxis.AxisLeftStickY

	AxisRightPadStickX := allBtnAxis.AxisRightPadStickX
	AxisRightPadStickY := allBtnAxis.AxisRightPadStickY

	stickDeadzone := dependentVars.cfg.PadsSticks.Stick.DeadzoneDS

	return func() {
		if event.eventType == EvAxisChanged {
			switch event.btnOrAxis {
			case AxisLeftStickX, AxisLeftStickY, AxisRightPadStickX, AxisRightPadStickY:
				if math.Abs(event.value) <= stickDeadzone {
					event.value = 0
				}
			}
		}
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

		switch controllerInUse {
		case SteamController:
			event.fixButtonNamesForSC()
			event.transformToWingsSC()
			event.transformStickToDPadSC()
		case DualShock:
			event.applyDeadzoneDS()
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
