package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"path"
	"sync"
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
	linesParts := Cfg.ReadLayoutFile(path.Join(Cfg.LayoutInUse, "buttons.csv"), 2)
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
		if isEmptyCmd(codes) {
			gofuncs.Panic("Empty command mapping for button %s", btn)
		}
		if codes[0] == NoAction {
			continue
		}
		pressLayout[btn] = MakeCommandInfo(codes, gofuncs.NaN())
	}
	return pressLayout
}

type Command []int

type ButtonToCommand map[BtnOrAxisT]*CommandInfo

var pressCommandsLayout ButtonToCommand

var buttonsToRelease = gofuncs.MakeMap[BtnOrAxisT, *CommandInfo]()

var ButtonsLock sync.Mutex

func PutButton(btn BtnOrAxisT, commandInfo *CommandInfo) bool {
	if _, exist := buttonsToRelease.CheckAndGet(btn); exist {
		return false
	}
	buttonsToRelease.Put(btn, commandInfo)
	return true
}

type Interval struct {
	repeatInterval, intervalLeft float64
}

func MakeInterval(repeatInterval float64) *Interval {
	interval := &Interval{}
	interval.SetInterval(repeatInterval)
	return interval
}

func (i *Interval) SetInterval(repeatInterval float64) {
	gofuncs.PanicAnyNotInitOrEmpty(repeatInterval)

	i.repeatInterval = repeatInterval
	i.intervalLeft = repeatInterval
}

func (i *Interval) ResetInterval() bool {
	if i.intervalLeft <= 0 {
		i.intervalLeft = i.repeatInterval
		return true
	}
	return false
}

func (i *Interval) DecreaseInterval(tickerInterval float64) bool {
	i.intervalLeft -= tickerInterval
	return i.ResetInterval()
}

type CommandInfo struct {
	Interval
	command              Command
	specialCaseIsHandled bool
}

func MakeCommandInfo(command Command, repeatInterval float64) *CommandInfo {
	if gofuncs.IsNotInit(repeatInterval) {
		repeatInterval = Cfg.holdRepeatInterval
	}

	commandInfo := &CommandInfo{command: command}
	commandInfo.SetInterval(repeatInterval)
	return commandInfo
}

func MakeUndeterminedCommandInfo() *CommandInfo {
	return MakeCommandInfo(nil, Cfg.holdingStateThreshold)
}

func (c *CommandInfo) GetCopy() *CommandInfo {
	return MakeCommandInfo(c.command, c.repeatInterval)
}

func (c *CommandInfo) CopyFromOther(other *CommandInfo) {
	c.command = other.command
	c.SetInterval(other.repeatInterval)
}

func isEmptyCommandInfo(commandInfo *CommandInfo) bool {
	return commandInfo == nil
}

func isEmptyCommandForButton(btn BtnOrAxisT, hold bool) bool {
	return isEmptyCommandInfo(getCommandInfo(btn, hold))
}

func isEmptyCmd(command Command) bool {
	return gofuncs.IsEmptySlice(command)
}

func commandNotExists(btn BtnOrAxisT) bool {
	return isEmptyCommandForButton(btn, false) &&
		isEmptyCommandForButton(btn, true)
}

func hasHoldCommand(btn BtnOrAxisT) bool {
	return !isEmptyCommandForButton(btn, true)
}

func getCommandInfo(btn BtnOrAxisT, hold bool) *CommandInfo {
	if hold {
		btn = addHoldSuffix(btn)
	}
	//Have only one point of access. Don't forget to copy
	commandInfo := pressCommandsLayout[btn]
	//nil can't be copied
	if !isEmptyCommandInfo(commandInfo) {
		commandInfo = commandInfo.GetCopy()
	}
	return commandInfo
}

func getFirstCmdSymbol(command Command) int {
	return command[0]
}

func isSwitchModeCmd(command Command) bool {
	if isEmptyCmd(command) {
		return false
	}
	return getFirstCmdSymbol(command) == SwitchMode
}

func isEscLetterCode(command Command) bool {
	if isEmptyCmd(command) {
		return false
	}
	return getFirstCmdSymbol(command) == EscLetterCode
}

func pressSequence(btn BtnOrAxisT, commandInfo *CommandInfo) {
	command := commandInfo.command

	if !commandInfo.specialCaseIsHandled {
		commandInfo.specialCaseIsHandled = true

		if isSwitchModeCmd(command) {
			// don't do release all
			Cfg.PadsSticksMode.SwitchMode()
			return
		} else if isEscLetterCode(command) {
			releaseAll(btn)
		}
	}

	for _, el := range command {
		osSpec.PressKeyOrMouse(el)
	}
}

func releaseSequence(command Command) {
	if isSwitchModeCmd(command) {
		return
	}

	for _, el := range gofuncs.Reverse(command) {
		osSpec.ReleaseKeyOrMouse(el)
	}
}

func pressImmediately(btn BtnOrAxisT) {
	commandInfo := getCommandInfo(btn, false)

	if PutButton(btn, commandInfo) {
		pressSequence(btn, commandInfo)
	}
}

func pressButton(btn BtnOrAxisT) {
	if hasHoldCommand(btn) {
		PutButton(btn, MakeUndeterminedCommandInfo())
	} else {
		pressImmediately(btn)
	}
}

func releaseButton(btn BtnOrAxisT) {
	commandInfo := buttonsToRelease.Pop(btn)

	if isEmptyCommandInfo(commandInfo) {
		return
	}

	if isEmptyCmd(commandInfo.command) {
		//has hold command but no "immediately press" command
		//and not enough time have passed for hold command to be triggered
		if isEmptyCommandForButton(btn, false) {
			return
		}
		commandInfo.CopyFromOther(getCommandInfo(btn, false))
		pressSequence(btn, commandInfo)
	}

	releaseSequence(commandInfo.command)
}

func RunRepeatCommandThread() {
	var tickerInterval float64 = 1
	ticker := time.NewTicker(gofuncs.NumberToMillis(tickerInterval))
	for range ticker.C {
		ButtonsLock.Lock()
		//Esc button's releaseAll will break state (changing map over iteration)
		//RangeOverCopy prevents this: states will be restored, esc command executed,
		//and then states will be properly released
		buttonsToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfo) {
			if isEmptyCmd(commandInfo.command) {
				//if hold state Interval passed
				if commandInfo.DecreaseInterval(tickerInterval) {
					//assign hold command, reset interval
					commandInfo.CopyFromOther(getCommandInfo(btn, true))
					pressSequence(btn, commandInfo)
				}
			} else { //if command already assigned to hold or immediate
				if commandInfo.DecreaseInterval(tickerInterval) {
					pressSequence(btn, commandInfo)
				}
			}
		})
		ButtonsLock.Unlock()
	}
}

func releaseAll(curButton BtnOrAxisT) {
	buttonsToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfo) {
		//current button should stay in map
		if btn != curButton {
			releaseButton(btn)
		}
	})
}

func isTriggerBtn(btn BtnOrAxisT) bool {
	return btn == BtnLeftTrigger || btn == BtnRightTrigger
}

func handleTriggers(btn BtnOrAxisT, value float64) {
	if value > Cfg.TriggerThreshold {
		pressImmediately(btn)
	} else if value < Cfg.TriggerThreshold {
		releaseButton(btn)
	}
}

func buttonChanged(btn BtnOrAxisT, value float64) {
	ButtonsLock.Lock()
	defer ButtonsLock.Unlock()

	if commandNotExists(btn) {
		return
	}

	if isTriggerBtn(btn) {
		handleTriggers(btn, value)
	} else {
		switch value {
		case 1:
			pressButton(btn)
		case 0:
			releaseButton(btn)
		default:
			gofuncs.Panic("Unsupported value: \"%s\" %v", btn, value)
		}
	}
}
