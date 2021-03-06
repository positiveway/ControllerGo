package mainLogic

import (
	"ControllerGo/src/osSpec"
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
	linesParts := ReadLayoutFile(path.Join(LayoutInUse, "commands.csv"), 2)
	for _, parts := range linesParts {
		btn := BtnOrAxisT(parts[0])
		keys := parts[1:]

		if btnSynonym, found := BtnSynonyms[btn]; found {
			btn = btnSynonym
		}
		if !contains(AllAvailableButtons, removeHoldSuffix(btn)) {
			PanicMisspelled(btn)
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
			panicMsg("Empty command mapping for button %s", btn)
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
	for _, el := range reverse(command) {
		osSpec.ReleaseKeyOrMouse(el)
	}
}

type Command []int

type ButtonToCommand map[BtnOrAxisT]Command

var pressCommandsLayout ButtonToCommand

var buttonsToRelease = MakeThreadSafeMap[BtnOrAxisT, CommandToReleaseWithHoldStartTime]()

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

func hasHoldCommand(button BtnOrAxisT) bool {
	_, found := pressCommandsLayout[addHoldSuffix(button)]
	return found
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
	if isEmpty(command) {
		return
	}

	switch command[0] {
	case SwitchMode:
		//releaseAll()
		padsMode.SwitchMode()
		return
	case EscLetter:
		releaseAll()
	}

	PutButton(btn, command, true)
	pressSequence(command)
}

func buttonPressed() {
	btn := event.btnOrAxis
	if isTriggerBtn(btn) {
		return
	}

	if hasHoldCommand(btn) {
		PutButton(btn, nil, false)
	} else { //press immediately
		pressButton(btn, false)
	}
}

const holdRefreshInterval = 15 * time.Millisecond

func RunReleaseHoldThread() {
	ticker := time.NewTicker(holdRefreshInterval)
	for range ticker.C {
		buttonsToRelease.RangeOverCopy(func(btn BtnOrAxisT, cmdWithTime CommandToReleaseWithHoldStartTime) {
			holdDuration := time.Now().Sub(cmdWithTime.holdStartTime)
			//print("duration: %v", holdDuration)
			if holdDuration > holdingThreshold {
				pressButton(btn, true)
			}
		})
	}
}

func buttonReleased() {
	btn := event.btnOrAxis
	if isTriggerBtn(btn) {
		return
	}

	pressButton(btn, false)
	releaseButton(btn)
}

func releaseButton(btn BtnOrAxisT) {
	cmdWithTime := buttonsToRelease.Pop(btn)
	command := cmdWithTime.command

	if isEmpty(command) {
		return
	}
	releaseSequence(command)
}

func releaseAll() {
	buttonsToRelease.RangeOverCopy(func(btn BtnOrAxisT, cmdWithTime CommandToReleaseWithHoldStartTime) {
		releaseButton(btn)
	})
}

func isTriggerBtn(btn BtnOrAxisT) bool {
	return btn == BtnLeftTrigger || btn == BtnRightTrigger
}

func detectTriggers() {
	btn := event.btnOrAxis
	if !isTriggerBtn(btn) {
		return
	}

	value := event.value

	if value > TriggerThreshold {
		pressButton(btn, false)
	} else if value < TriggerThreshold {
		releaseButton(btn)
	}
}
