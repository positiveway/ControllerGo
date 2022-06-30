package mainLogic

import (
	"strconv"
)

const (
	LeftAdjustX float64 = -0.04
	LeftAdjustY float64 = -0.14

	RightAdjustX float64 = 0
	RightAdjustY float64 = 0
)

var PadAxes = []BtnOrAxisT{
	AxisRightPadX,
	AxisRightPadY,
	AxisLeftPadX,
	AxisLeftPadY,
}

type Event struct {
	eventType EventTypeT
	btnOrAxis BtnOrAxisT
	value     float64
	codeType  CodeTypeT
	code      CodeT
}

func (event *Event) applyAdjustments() {
	if event.eventType == EvAxisChanged {
		switch event.btnOrAxis {
		case AxisLeftPadX:
			event.value -= LeftAdjustX
		case AxisLeftPadY:
			event.value -= LeftAdjustY
		case AxisRightPadX:
			event.value -= RightAdjustX
		case AxisRightPadY:
			event.value -= RightAdjustY
		}
	}
}

func (event *Event) fixButtonNames() {
	switch event.btnOrAxis {
	case BtnY:
		event.btnOrAxis = BtnX
	case BtnX:
		event.btnOrAxis = BtnY
	}
}

func (event *Event) transformToPadEvent() {
	if event.btnOrAxis == BtnUnknown && event.codeType == CTAbs {
		if axis, found := UnknownCodesResolvingMap[event.code]; found {
			switch event.eventType {
			case EvButtonPressed:
				event.eventType = EvPadFirstTouched
			case EvButtonReleased:
				event.eventType = EvPadReleased
			default:
				return
			}
			event.btnOrAxis = axis
			event.value = 0
			event.code = 0
			event.codeType = ""
		}
	}
}

func (event *Event) transformToWings() {
	if event.btnOrAxis == BtnUnknown && event.codeType == CTKey {
		if btn, found := UnknownCodesResolvingMap[event.code]; found {
			event.btnOrAxis = btn
		}
	}
}

var stickCoords = makeCoords()

func (event *Event) transformStickToButtons() {

	if event.eventType == EvAxisChanged {
		stickCoords.mu.Lock()
		defer stickCoords.mu.Unlock()

		switch event.btnOrAxis {
		case AxisStickX:
			stickCoords.setDirectlyX()
		case AxisStickY:
			stickCoords.setDirectlyY()
		default:
			return
		}

		stickCoords.updateValues()
		stickCoords.updateAngle()
	}
}

func (event *Event) transformAndFilter() {
	//fmt.Printf("Before: ")
	//event.print()

	event.fixButtonNames()

	if event.eventType == EvAxisChanged &&
		contains(PadAxes, event.btnOrAxis) &&
		event.value == 0.0 {
		return
	}

	//event.applyAdjustments()

	event.transformToPadEvent()
	event.transformToWings()

	if event.btnOrAxis == BtnUnknown {
		return
	}

	//fmt.Printf("After: ")
	//event.print()

	matchEvent()
}

func (event *Event) update(msg string) {
	var found bool
	var err error

	event.eventType, found = EventTypeMap[msg[0]]
	if !found {
		PanicMisspelled(string(msg[0]))
	}
	if event.eventType != EvConnected && event.eventType != EvDisconnected && event.eventType != EvDropped {
		event.btnOrAxis, found = BtnAxisMap[msg[1]]
		if !found {
			PanicMisspelled(string(msg[1]))
		}
		if event.eventType == EvAxisChanged || event.eventType == EvButtonChanged {
			msg = msg[2:]
			valueAndCode := split(msg, ";")

			event.value, err = strconv.ParseFloat(valueAndCode[0], 32)
			checkErr(err)

			if startsWith(msg, ";") {
				return
			}
			typeAndCode := split(valueAndCode[1], "(")
			event.codeType = CodeTypeT(typeAndCode[0])

			code := typeAndCode[1]
			codeNum := strToInt(code[:len(code)-1])
			event.code = CodeT(codeNum)
		}
	}
	event.transformAndFilter()
}

func (event *Event) print() {
	print("%s %s %s %v %0.2f",
		trimAnyPrefix(string(event.eventType), "Ev"),
		trimAnyPrefix(string(event.btnOrAxis), "Btn", "Axis"),
		event.codeType,
		event.code,
		event.value)
}
