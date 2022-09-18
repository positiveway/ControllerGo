package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"reflect"
	"strings"
)

type CodeTypeT string

const (
	CTAbs CodeTypeT = "ABS"
	CTKey CodeTypeT = "KEY"
)

type CodeT uint

const (
	//CodeStickXSC CodeT = 0
	//CodeStickYSC   CodeT = 1
	//CodeLeftPadXSC CodeT = 16
	//CodeLeftPadYSC  CodeT = 17
	//CodeRightPadXSC CodeT = 3
	//CodeRightPadYSC CodeT = 4
	CodeLeftWingSC  CodeT = 336
	CodeRightWingSC CodeT = 337
)

type BtnOrAxisT string

const HoldSuffix = "_hold" //lower case

func addHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(string(btn) + HoldSuffix)
}

func removeHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(strings.TrimSuffix(string(btn), HoldSuffix))
}

func InitCurStickButton() *BtnOrAxisT {
	//required for full reset
	return new(BtnOrAxisT)
}

type BtnAxisMapT map[uint8]BtnOrAxisT

func (allBtnAxis *AllBtnAxis) InitBtnAxisMap() BtnAxisMapT {
	mapping := BtnAxisMapT{}
	for k, v := range allBtnAxis.initAxisMap() {
		gofuncs.AssignWithDuplicateKeyValueCheck(mapping, k, v, true)
	}
	for k, v := range allBtnAxis.initBtnMap() {
		gofuncs.AssignWithDuplicateKeyCheck(mapping, k, v)
	}
	return mapping
}

const (
	AxisUnknown BtnOrAxisT = "AxisUnknown"
)

func (allBtnAxis *AllBtnAxis) initAxisMap() BtnAxisMapT {
	return BtnAxisMapT{
		'u': allBtnAxis.AxisLeftStickX,
		'v': allBtnAxis.AxisLeftStickY,
		'w': allBtnAxis.AxisLeftStickZ,
		'x': allBtnAxis.AxisRightPadStickX,
		'y': allBtnAxis.AxisRightPadStickY,
		'z': allBtnAxis.AxisRightPadStickZ,
		'0': allBtnAxis.AxisLeftPadX,
		'1': allBtnAxis.AxisLeftPadY,
		'2': AxisUnknown,
	}
}

func (allBtnAxis *AllBtnAxis) initPadAndStickAxes() []BtnOrAxisT {
	return []BtnOrAxisT{
		allBtnAxis.AxisLeftPadX,
		allBtnAxis.AxisLeftPadY,
		allBtnAxis.AxisRightPadStickX,
		allBtnAxis.AxisRightPadStickY,
		allBtnAxis.AxisLeftStickX,
		allBtnAxis.AxisLeftStickY,
	}
}

func (allBtnAxis *AllBtnAxis) initConfigDependent(cfg *ConfigsT) {
	switch cfg.ControllerInUse {
	case SteamController:
		//axis
		allBtnAxis.AxisLeftPadX = "LeftPadX"
		allBtnAxis.AxisLeftPadY = "LeftPadY"

		allBtnAxis.AxisLeftStickX = "StickX"
		allBtnAxis.AxisLeftStickY = "StickY"
		allBtnAxis.AxisLeftStickZ = "StickZ"

		allBtnAxis.AxisRightPadStickX = "RightPadX"
		allBtnAxis.AxisRightPadStickY = "RightPadY"
		allBtnAxis.AxisRightPadStickZ = "RightPadZ"

		//buttons
		allBtnAxis.BtnLeftPad = "LeftPad"
		allBtnAxis.BtnLeftStick = "Stick"
		allBtnAxis.BtnRightPadStick = "RightPad"

		allBtnAxis.BtnLeftWingSC = "LeftWing"
		allBtnAxis.BtnRightWingSC = "RightWing"

		allBtnAxis.BtnStickUpSC = "StickUp"
		allBtnAxis.BtnStickDownSC = "StickDown"
		allBtnAxis.BtnStickLeftSC = "StickLeft"
		allBtnAxis.BtnStickRightSC = "StickRight"

		allBtnAxis.BtnDPadUp = allBtnAxis.BtnLeftPad
		allBtnAxis.BtnDPadDown = allBtnAxis.BtnLeftPad
		allBtnAxis.BtnDPadLeft = allBtnAxis.BtnLeftPad
		allBtnAxis.BtnDPadRight = allBtnAxis.BtnLeftPad
	case DualShock:
		//axis
		allBtnAxis.AxisLeftStickX = "LeftStickX"
		allBtnAxis.AxisLeftStickY = "LeftStickY"
		allBtnAxis.AxisLeftStickZ = "LeftStickZ"

		allBtnAxis.AxisRightPadStickX = "RightStickX"
		allBtnAxis.AxisRightPadStickY = "RightStickY"
		allBtnAxis.AxisRightPadStickZ = "RightStickZ"

		//buttons
		allBtnAxis.BtnLeftStick = "LeftStick"
		allBtnAxis.BtnRightPadStick = "RightStick"

		allBtnAxis.BtnDPadUp = "DPadUp"
		allBtnAxis.BtnDPadDown = "DPadDown"
		allBtnAxis.BtnDPadLeft = "DPadLeft"
		allBtnAxis.BtnDPadRight = "DPadRight"
	}
}

func initEventTypeMap() map[uint8]EventTypeT {
	return map[uint8]EventTypeT{
		'a': EvAxisChanged,
		'b': EvButtonChanged,
		'c': EvButtonReleased,
		'd': EvButtonPressed,
		'e': EvButtonRepeated,
		'f': EvConnected,
		'g': EvDisconnected,
		'h': EvDropped,
	}
}

func (allBtnAxis *AllBtnAxis) initUnknownCodesMapSC() map[CodeT]BtnOrAxisT {
	return map[CodeT]BtnOrAxisT{
		//CodeStickXSC:    allBtnAxis.AxisLeftStickX,
		//CodeStickYSC:    allBtnAxis.AxisLeftStickY,
		//CodeLeftPadXSC:  allBtnAxis.AxisLeftPadX,
		//CodeLeftPadYSC:  allBtnAxis.AxisLeftPadY,
		//CodeRightPadXSC: allBtnAxis.AxisRightPadStickX,
		//CodeRightPadYSC: allBtnAxis.AxisRightPadStickY,
		CodeLeftWingSC:  allBtnAxis.BtnLeftWingSC,
		CodeRightWingSC: allBtnAxis.BtnRightWingSC,
	}
}

type ZoneToBtnMapT map[ZoneT]BtnOrAxisT

func (allBtnAxis *AllBtnAxis) initStickZoneBtnMap() ZoneToBtnMapT {
	return ZoneToBtnMapT{
		ZoneRight: allBtnAxis.BtnStickRightSC,
		ZoneUp:    allBtnAxis.BtnStickUpSC,
		ZoneLeft:  allBtnAxis.BtnStickLeftSC,
		ZoneDown:  allBtnAxis.BtnStickDownSC,
	}
}

type AvailableButtonsT []BtnOrAxisT

func (allBtnAxis *AllBtnAxis) initAvailableButtons() AvailableButtonsT {
	_availableButtons := AvailableButtonsT{
		allBtnAxis.BtnLeftWingSC,
		allBtnAxis.BtnRightWingSC,
		allBtnAxis.BtnA,
		allBtnAxis.BtnB,
		allBtnAxis.BtnY,
		allBtnAxis.BtnX,
		allBtnAxis.BtnC,
		allBtnAxis.BtnZ,
		allBtnAxis.BtnLeftButton,
		allBtnAxis.BtnLeftTrigger,
		allBtnAxis.BtnRightButton,
		allBtnAxis.BtnRightTrigger,
		allBtnAxis.BtnLeftSpecial,
		allBtnAxis.BtnRightSpecial,
		allBtnAxis.BtnCentralSpecial,

		allBtnAxis.BtnLeftStick,
		allBtnAxis.BtnRightPadStick,
		allBtnAxis.BtnLeftPad,

		allBtnAxis.BtnDPadUp,
		allBtnAxis.BtnDPadDown,
		allBtnAxis.BtnDPadLeft,
		allBtnAxis.BtnDPadRight,

		allBtnAxis.BtnStickUpSC,
		allBtnAxis.BtnStickDownSC,
		allBtnAxis.BtnStickLeftSC,
		allBtnAxis.BtnStickRightSC,
	}

	//filter empty
	var NonEmptyAvailableButtons AvailableButtonsT
	for _, button := range _availableButtons {
		if !gofuncs.IsEmptyStripStr(string(button)) {
			NonEmptyAvailableButtons = append(NonEmptyAvailableButtons, button)
		}
	}

	return NonEmptyAvailableButtons
}

func (allBtnAxis *AllBtnAxis) initBtnMap() BtnAxisMapT {
	return BtnAxisMapT{
		'a': allBtnAxis.BtnA,
		'b': allBtnAxis.BtnB,
		'c': allBtnAxis.BtnY,
		'd': allBtnAxis.BtnX,
		'e': allBtnAxis.BtnC,
		'f': allBtnAxis.BtnZ,
		'g': allBtnAxis.BtnLeftButton,
		'h': allBtnAxis.BtnLeftTrigger,
		'i': allBtnAxis.BtnRightButton,
		'j': allBtnAxis.BtnRightTrigger,
		'k': allBtnAxis.BtnLeftSpecial,
		'l': allBtnAxis.BtnRightSpecial,
		'm': allBtnAxis.BtnCentralSpecial,
		'n': allBtnAxis.BtnLeftStick,
		'o': allBtnAxis.BtnRightPadStick,
		'p': allBtnAxis.BtnDPadUp,
		'q': allBtnAxis.BtnDPadDown,
		'r': allBtnAxis.BtnDPadLeft,
		's': allBtnAxis.BtnDPadRight,
		't': allBtnAxis.BtnUnknown,
	}
}

type AllBtnAxis struct {
	BtnB,
	BtnY,
	BtnX,
	BtnA,
	BtnC,
	BtnZ,
	BtnLeftButton,
	BtnLeftTrigger,
	BtnRightButton,
	BtnRightTrigger,
	BtnLeftSpecial,
	BtnRightSpecial,
	BtnCentralSpecial,
	BtnUnknown,

	BtnLeftPad,
	BtnLeftStick,
	BtnRightPadStick,

	BtnLeftWingSC,
	BtnRightWingSC,

	BtnStickUpSC,
	BtnStickDownSC,
	BtnStickLeftSC,
	BtnStickRightSC,

	BtnDPadUp,
	BtnDPadDown,
	BtnDPadLeft,
	BtnDPadRight,

	AxisLeftPadX,
	AxisLeftPadY,

	AxisLeftStickX,
	AxisLeftStickY,
	AxisLeftStickZ,

	AxisRightPadStickX,
	AxisRightPadStickY,
	AxisRightPadStickZ BtnOrAxisT
}

func MakeAllBtnAxis(cfg *ConfigsT) *AllBtnAxis {
	allBtnAxis := &AllBtnAxis{
		BtnB:              "B",
		BtnY:              "Y",
		BtnX:              "X",
		BtnA:              "A",
		BtnC:              "BtnC",
		BtnZ:              "BtnZ",
		BtnLeftButton:     "LB",
		BtnLeftTrigger:    "LT",
		BtnRightButton:    "RB",
		BtnRightTrigger:   "RT",
		BtnLeftSpecial:    "LeftSpecial",
		BtnRightSpecial:   "RightSpecial",
		BtnCentralSpecial: "CentralSpecial",
		BtnUnknown:        "BtnUnknown",
	}

	allBtnAxis.initConfigDependent(cfg)
	allBtnAxis.ToLower()

	return allBtnAxis
}

func (allBtnAxis *AllBtnAxis) ToLower() {
	v := reflect.Indirect(reflect.ValueOf(allBtnAxis))

	for i := 0; i < v.NumField(); i++ {
		strFieldValue := gofuncs.ToLowerPanicIfSpaces(v.Field(i).String())
		v.Field(i).Set(reflect.ValueOf(BtnOrAxisT(strFieldValue)))
	}
}

type StrOrBtn interface {
	BtnOrAxisT | string
}

type SynonymsT[T StrOrBtn] map[T]T

func ToLower[T StrOrBtn](strOrBtn T) T {
	return T(gofuncs.ToLowerAndStripPanicIfEmpty(strOrBtn))
}

func ToLowerSynonyms[T StrOrBtn](synonyms SynonymsT[T]) SynonymsT[T] {
	lowered := SynonymsT[T]{}
	for synonym, orig := range synonyms {
		lowered[ToLower(synonym)] = ToLower(orig)
	}
	return lowered
}

func (allBtnAxis *AllBtnAxis) genBtnSynonyms() SynonymsT[BtnOrAxisT] {
	synonyms := SynonymsT[BtnOrAxisT]{
		"LeftButton":                allBtnAxis.BtnLeftButton,
		addHoldSuffix("LeftButton"): addHoldSuffix(allBtnAxis.BtnLeftButton),

		"RightButton":                allBtnAxis.BtnRightButton,
		addHoldSuffix("RightButton"): addHoldSuffix(allBtnAxis.BtnRightButton),

		"LeftTrigger":  allBtnAxis.BtnLeftTrigger,
		"RightTrigger": allBtnAxis.BtnRightTrigger,
	}

	return ToLowerSynonyms(synonyms)
}

type EventTypeT string

const (
	EvAxisChanged    EventTypeT = "AxisChanged"
	EvButtonChanged  EventTypeT = "ButtonChanged"
	EvButtonReleased EventTypeT = "ButtonReleased"
	EvButtonPressed  EventTypeT = "ButtonPressed"
	EvButtonRepeated EventTypeT = "ButtonRepeated"
	EvConnected      EventTypeT = "Connected"
	EvDisconnected   EventTypeT = "Disconnected"
	EvDropped        EventTypeT = "Dropped"
)
