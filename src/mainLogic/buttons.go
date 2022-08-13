package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"path"
	"time"
)

const NoAction = -1
const SwitchMode = -2

var commonCmdMapping = map[string]int{
	"NoAction":    NoAction,
	"LeftMouse":   osSpec.LeftMouse,
	"RightMouse":  osSpec.RightMouse,
	"MiddleMouse": osSpec.MiddleMouse,
	"SwitchMode":  SwitchMode,
}

func initCommands() {
	pressCommandsLayout = loadCommandsLayout()
}

func loadCommandsLayout() ButtonToCommand {
	pressLayout := ButtonToCommand{}
	linesParts := Cfg.ReadLayoutFile(path.Join(Cfg.LayoutInUse, "commands.csv"), 2)
	for _, parts := range linesParts {
		btn := BtnOrAxisT(parts[0])
		keys := parts[1:]

		if btnSynonym, found := BtnSynonyms[btn]; found {
			btn = btnSynonym
		}
		if !gofuncs.Contains(AllAvailableButtons, removeHoldSuffix(btn)) {
			gofuncs.PanicMisspelled(btn)
		}
		var codes []int
		for _, key := range keys {
			if code, found := commonCmdMapping[key]; found {
				codes = append(codes, code)
			} else {
				code := getCodeFromLetter(key)
				codes = append(codes, code)
			}
		}
		if len(codes) == 0 {
			gofuncs.Panic("Empty command mapping for button %s", btn)
		}
		if codes[0] == NoAction {
			continue
		}
		pressLayout[btn] = codes
	}
	return pressLayout
}

func pressSequence(command Command) {
	for _, el := range command {
		osSpec.PressKeyOrMouse(el)
	}
}

func releaseSequence(command Command) {
	for _, el := range gofuncs.Reverse(command) {
		osSpec.ReleaseKeyOrMouse(el)
	}
}

type Command []int

type ButtonToCommand map[BtnOrAxisT]Command

var pressCommandsLayout ButtonToCommand

var buttonsToRelease = gofuncs.MakeThreadSafeMap[BtnOrAxisT, CommandToReleaseWithHoldStartTime]()

func PutButton(btn BtnOrAxisT, command Command, alreadyPressed bool) {
	buttonsToRelease.Put(btn,
		CommandToReleaseWithHoldStartTime{
			command:        command,
			holdStartTime:  time.Now(),
			alreadyPressed: alreadyPressed,
		})
}

type CommandToReleaseWithHoldStartTime struct {
	command        Command
	holdStartTime  time.Time
	alreadyPressed bool
}

func isEmptyCmd(command Command) bool {
	return gofuncs.IsEmptySlice(command)
}

func commandNotExists(button BtnOrAxisT) bool {
	return isEmptyCmd(pressCommandsLayout[button]) &&
		isEmptyCmd(pressCommandsLayout[addHoldSuffix(button)])
}

func hasHoldCommand(button BtnOrAxisT) bool {
	return !isEmptyCmd(pressCommandsLayout[addHoldSuffix(button)])
}

func getPressCommand(btn BtnOrAxisT, hold bool) Command {
	if hold {
		btn = addHoldSuffix(btn)
	}
	return pressCommandsLayout[btn]
}

func pressButton(btn BtnOrAxisT, hold bool) {
	if cmdWithTime, found := buttonsToRelease.CheckAndGet(btn); found {
		if cmdWithTime.alreadyPressed {
			return
		}
	}

	command := getPressCommand(btn, hold)
	//if isEmptyCmd(command) {
	//	return
	//}

	firstCmdSymbol := command[0]
	switch firstCmdSymbol {
	case SwitchMode:
		PutButton(btn, nil, true)
		// don't do release all
		Cfg.padsMode.SwitchMode()
		return

	case EscLetter:
		releaseAll()
	}

	PutButton(btn, command, true)
	pressSequence(command)
}

func buttonPressed(btn BtnOrAxisT) {
	if hasHoldCommand(btn) {
		PutButton(btn, nil, false)
	} else { //press immediately
		pressButton(btn, false)
	}
}

func buttonReleased(btn BtnOrAxisT) {
	pressButton(btn, false)
	releaseButton(btn)
}

func releaseButton(btn BtnOrAxisT) {
	cmdWithTime := buttonsToRelease.Pop(btn)
	command := cmdWithTime.command

	//if isEmptyCmd(command) {
	//	return
	//}
	releaseSequence(command)
}

func RunReleaseHoldThread() {
	ticker := time.NewTicker(Cfg.holdRefreshInterval)
	for range ticker.C {
		buttonsToRelease.RangeOverCopy(func(btn BtnOrAxisT, cmdWithTime CommandToReleaseWithHoldStartTime) {
			holdDuration := time.Now().Sub(cmdWithTime.holdStartTime)
			//gofuncs.Print("duration: %v", holdDuration)
			if holdDuration > Cfg.holdingThreshold {
				pressButton(btn, true)
			}
		})
	}
}

func releaseAll() {
	buttonsToRelease.RangeOverCopy(func(btn BtnOrAxisT, cmdWithTime CommandToReleaseWithHoldStartTime) {
		releaseButton(btn)
	})
}

func isTriggerBtn(btn BtnOrAxisT) bool {
	return btn == BtnLeftTrigger || btn == BtnRightTrigger
}

func handleTriggers(btn BtnOrAxisT, value float64) {
	if value > Cfg.TriggerThreshold {
		pressButton(btn, false)
	} else if value < Cfg.TriggerThreshold {
		releaseButton(btn)
	}
}

func buttonChanged(btn BtnOrAxisT, value float64) {
	if commandNotExists(btn) {
		return
	}

	if isTriggerBtn(btn) {
		handleTriggers(btn, value)
	} else {
		switch value {
		case 1:
			buttonPressed(btn)
		case 0:
			buttonReleased(btn)
		default:
			gofuncs.Panic("Unsupported value: \"%s\" %v", btn, value)
		}
	}
}
