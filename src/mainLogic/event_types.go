package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"strings"
)

type CodeTypeT string

const (
	CTAbs CodeTypeT = "ABS"
	CTKey CodeTypeT = "KEY"
)

type CodeT int

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

const HoldSuffix = "_Hold"

func addHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(string(btn) + HoldSuffix)
}

func removeHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(strings.TrimSuffix(string(btn), HoldSuffix))
}

type BtnAxisMapT map[uint8]BtnOrAxisT

func genBtnAxisMap() BtnAxisMapT {
	mapping := BtnAxisMapT{}
	for k, v := range _AxisMap {
		gofuncs.AssignWithDuplicateCheck(mapping, k, v)
	}
	for k, v := range _BtnMap {
		gofuncs.AssignWithDuplicateCheck(mapping, k, v)
	}
	return mapping
}

const (
	AxisUnknown BtnOrAxisT = "Unknown"
)

func initAxisMap() {
	_AxisMap = BtnAxisMapT{
		'u': AxisLeftStickX,
		'v': AxisLeftStickY,
		'w': AxisLeftStickZ,
		'x': AxisRightPadStickX,
		'y': AxisRightPadStickY,
		'z': AxisRightPadStickZ,
		'0': AxisLeftPadX,
		'1': AxisLeftPadY,
		'2': AxisUnknown,
	}
}

func initPadAndStickAxes() {
	PadAndStickAxes = []BtnOrAxisT{
		AxisLeftPadX,
		AxisLeftPadY,
		AxisRightPadStickX,
		AxisRightPadStickY,
		AxisLeftStickX,
		AxisLeftStickY,
	}
}

func initEventTypes() {
	switch Cfg.ControllerInUse {
	case SteamController:
		//axis
		AxisLeftPadX = "LeftPadX"
		AxisLeftPadY = "LeftPadY"

		AxisLeftStickX = "StickX"
		AxisLeftStickY = "StickY"
		AxisLeftStickZ = "StickZ"

		AxisRightPadStickX = "RightPadX"
		AxisRightPadStickY = "RightPadY"
		AxisRightPadStickZ = "RightPadZ"

		//buttons
		BtnLeftPad = "LeftPad"
		BtnLeftStick = "Stick"
		BtnRightPadStick = "RightPad"

		BtnLeftWingSC = "LeftWing"
		BtnRightWingSC = "RightWing"

		BtnStickUpSC = "StickUp"
		BtnStickDownSC = "StickDown"
		BtnStickLeftSC = "StickLeft"
		BtnStickRightSC = "StickRight"

		BtnDPadUp = BtnLeftPad
		BtnDPadDown = BtnLeftPad
		BtnDPadLeft = BtnLeftPad
		BtnDPadRight = BtnLeftPad
	case DualShock:
		//axis
		AxisLeftStickX = "LeftStickX"
		AxisLeftStickY = "LeftStickY"
		AxisLeftStickZ = "LeftStickZ"

		AxisRightPadStickX = "RightStickX"
		AxisRightPadStickY = "RightStickY"
		AxisRightPadStickZ = "RightStickZ"

		//buttons
		BtnLeftStick = "LeftStick"
		BtnRightPadStick = "RightStick"

		BtnDPadUp = "DPadUp"
		BtnDPadDown = "DPadDown"
		BtnDPadLeft = "DPadLeft"
		BtnDPadRight = "DPadRight"
	}

	initButtonsAndAxesFullSequence()
}

func initButtonsAndAxesFullSequence() {
	//event types
	initEventTypeMap()

	//axis
	initAxisMap()
	initPadAndStickAxes()

	//buttons
	BtnSynonyms = genBtnSynonyms()
	initStickZoneBtnMap()
	initAvailableButtons()
	initBtnMap()
	initUnknownCodesMapSC()

	//axis and buttons
	BtnAxisMap = genBtnAxisMap()

}

func initEventTypeMap() {
	EventTypeMap = map[uint8]EventTypeT{
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

func initUnknownCodesMapSC() {
	UnknownCodesResolvingMapSC = map[CodeT]BtnOrAxisT{
		//CodeStickXSC:    AxisLeftStickX,
		//CodeStickYSC:    AxisLeftStickY,
		//CodeLeftPadXSC:  AxisLeftPadX,
		//CodeLeftPadYSC:  AxisLeftPadY,
		//CodeRightPadXSC: AxisRightPadStickX,
		//CodeRightPadYSC: AxisRightPadStickY,
		CodeLeftWingSC:  BtnLeftWingSC,
		CodeRightWingSC: BtnRightWingSC,
	}
}

func initStickZoneBtnMap() {
	StickZoneToBtnMapSC = map[Zone]BtnOrAxisT{
		ZoneRight: BtnStickRightSC,
		ZoneUp:    BtnStickUpSC,
		ZoneLeft:  BtnStickLeftSC,
		ZoneDown:  BtnStickDownSC,
	}
}

func initAvailableButtons() {
	_availableButtons := []BtnOrAxisT{
		BtnLeftWingSC,
		BtnRightWingSC,
		BtnA,
		BtnB,
		BtnY,
		BtnX,
		BtnC,
		BtnZ,
		BtnLeftButton,
		BtnLeftTrigger,
		BtnRightButton,
		BtnRightTrigger,
		BtnSelect,
		BtnStart,
		BtnMode,

		BtnLeftStick,
		BtnRightPadStick,
		BtnLeftPad,

		BtnDPadUp,
		BtnDPadDown,
		BtnDPadLeft,
		BtnDPadRight,

		BtnStickUpSC,
		BtnStickDownSC,
		BtnStickLeftSC,
		BtnStickRightSC,
	}

	for _, button := range _availableButtons {
		if !gofuncs.IsEmptyStripStr(string(button)) {
			AllAvailableButtons = append(AllAvailableButtons, button)
		}
	}
}

func initBtnMap() {
	_BtnMap = BtnAxisMapT{
		'a': BtnA,
		'b': BtnB,
		'c': BtnY,
		'd': BtnX,
		'e': BtnC,
		'f': BtnZ,
		'g': BtnLeftButton,
		'h': BtnLeftTrigger,
		'i': BtnRightButton,
		'j': BtnRightTrigger,
		'k': BtnSelect,
		'l': BtnStart,
		'm': BtnMode,
		'n': BtnLeftStick,
		'o': BtnRightPadStick,
		'p': BtnDPadUp,
		'q': BtnDPadDown,
		'r': BtnDPadLeft,
		's': BtnDPadRight,
		't': BtnUnknown,
	}
}

const (
	BtnB            BtnOrAxisT = "B"
	BtnY            BtnOrAxisT = "Y"
	BtnX            BtnOrAxisT = "X"
	BtnA            BtnOrAxisT = "A"
	BtnC            BtnOrAxisT = "BtnC"
	BtnZ            BtnOrAxisT = "BtnZ"
	BtnLeftButton   BtnOrAxisT = "LB"
	BtnLeftTrigger  BtnOrAxisT = "LT"
	BtnRightButton  BtnOrAxisT = "RB"
	BtnRightTrigger BtnOrAxisT = "RT"
	BtnSelect       BtnOrAxisT = "Select"
	BtnStart        BtnOrAxisT = "Start"
	BtnMode         BtnOrAxisT = "Mode"

	BtnUnknown BtnOrAxisT = "BtnUnknown"
)

type Synonyms map[BtnOrAxisT]BtnOrAxisT

func genBtnSynonyms() Synonyms {
	synonyms := Synonyms{
		"LeftButton":   BtnLeftButton,
		"LeftTrigger":  BtnLeftTrigger,
		"RightButton":  BtnRightButton,
		"RightTrigger": BtnRightTrigger,
	}
	for key, val := range synonyms {
		synonyms[addHoldSuffix(key)] = addHoldSuffix(val)
	}
	return synonyms
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
	EvPadReleased    EventTypeT = "PadReleased"
)
