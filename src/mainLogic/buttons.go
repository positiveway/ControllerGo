package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"sync"
)

const (
	NoAction                = -10
	SwitchPadStickMode      = -11
	SwitchHighPrecisionMode = -12
)

func initButtons() {
	VirtualButtonCounter = MakeVirtualButtonCounter()
	pressCommandsLayout = loadCommandsLayout()
}

func initCommonCmdMapping() map[string]int {
	mapping := map[string]int{
		"NoAction":      NoAction,
		"LeftMouse":     osSpec.LeftMouse,
		"RightMouse":    osSpec.RightMouse,
		"MiddleMouse":   osSpec.MiddleMouse,
		"SwitchMode":    SwitchPadStickMode,
		"HighPrecision": SwitchHighPrecisionMode,
	}
	gofuncs.PanicIfDuplicateValueInMap(mapping, false)
	return mapping
}

func loadCommandsLayout() ButtonToCommandT {
	commonCmdMapping := initCommonCmdMapping()
	BtnSynonyms := genBtnSynonyms()
	AllAvailableButtons := initAvailableButtons()

	pressLayout := ButtonToCommandT{}
	linesParts := gofuncs.ReadLayoutFile(2,
		[]string{Cfg.Path.CurLayoutDir, "buttons.csv"})

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
				//Be careful! It probably works because variable was reassigned
				//and original map key isn't broken
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
		pressLayout[btn] = MakeEmptyCommandInfo(codes)
	}
	return pressLayout
}

type CommandT []int

type ButtonToCommandT map[BtnOrAxisT]*CommandInfoT

var pressCommandsLayout ButtonToCommandT

var buttonsToRelease = gofuncs.MakeMap[BtnOrAxisT, *CommandInfoT]()

var ButtonsLock sync.Mutex

type VirtualButtonCounterT struct {
	counter uint
}

var VirtualButtonCounter *VirtualButtonCounterT

func MakeVirtualButtonCounter() *VirtualButtonCounterT {
	return &VirtualButtonCounterT{}
}

func (v *VirtualButtonCounterT) GetButton() BtnOrAxisT {
	ButtonsLock.Lock()
	defer ButtonsLock.Unlock()

	//if uint overflows it will be zero
	v.counter += 1

	virtualButton := gofuncs.Format("VirtualButton_%v", v.counter)
	return BtnOrAxisT(virtualButton)
}

func PutButton(btn BtnOrAxisT, commandInfo *CommandInfoT) bool {
	if _, exist := buttonsToRelease.CheckAndGet(btn); exist {
		return false
	}
	buttonsToRelease.Put(btn, commandInfo)
	return true
}

type CommandInfoT struct {
	IntervalTimerT
	command              CommandT
	specialCaseIsHandled bool
}

func MakeCommandInfo(command CommandT, repeatInterval float64) *CommandInfoT {
	commandInfo := &CommandInfoT{command: command}
	commandInfo.InitIntervalTimer(repeatInterval)
	return commandInfo
}

func MakeEmptyCommandInfo(command CommandT) *CommandInfoT {
	return MakeCommandInfo(command, Cfg.Buttons.HoldRepeatInterval)
}

func MakeUndeterminedCommandInfo() *CommandInfoT {
	return MakeCommandInfo(nil, Cfg.Buttons.HoldingStateThreshold)
}

func (c *CommandInfoT) GetCopy() *CommandInfoT {
	return MakeCommandInfo(c.command, c.repeatInterval)
}

func (c *CommandInfoT) CopyFromOther(other *CommandInfoT) {
	c.command = other.command
	c.SetInterval(other.repeatInterval)
}

func isEmptyCommandInfo(commandInfo *CommandInfoT) bool {
	return commandInfo == nil
}

func isEmptyCommandForButton(btn BtnOrAxisT, hold bool) bool {
	return isEmptyCommandInfo(getCommandInfo(btn, hold))
}

func isEmptyCmd(command CommandT) bool {
	return gofuncs.IsEmptySlice(command)
}

func commandNotExists(btn BtnOrAxisT) bool {
	return isEmptyCommandForButton(btn, false) &&
		isEmptyCommandForButton(btn, true)
}

func hasHoldCommand(btn BtnOrAxisT) bool {
	return !isEmptyCommandForButton(btn, true)
}

func getCommandInfo(btn BtnOrAxisT, hold bool) *CommandInfoT {
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

func getFirstCmdSymbol(command CommandT) int {
	return command[0]
}

func isSwitchModeCmd(command CommandT) bool {
	return getFirstCmdSymbol(command) == SwitchPadStickMode
}

func pressSequence(btn BtnOrAxisT, commandInfo *CommandInfoT) {
	command := commandInfo.command

	if !commandInfo.specialCaseIsHandled {
		commandInfo.specialCaseIsHandled = true

		switch getFirstCmdSymbol(command) {
		case SwitchPadStickMode:
			// don't do release all
			Cfg.PadsSticks.Mode.SwitchMode()
			return
		case SwitchHighPrecisionMode:
			Cfg.PadsSticks.HighPrecisionMode.SwitchMode()
			return
		case EscLetterCode:
			releaseAll(btn)
		}
	}

	for _, el := range command {
		osSpec.PressKeyOrMouse(el)
	}
}

func releaseSequence(command CommandT) {
	if isSwitchModeCmd(command) {
		return
	}

	for _, el := range gofuncs.Reverse(command) {
		osSpec.ReleaseKeyOrMouse(el)
	}
}

func pressIfNotAlready(btn BtnOrAxisT, commandInfo *CommandInfoT) {
	if PutButton(btn, commandInfo) {
		pressSequence(btn, commandInfo)
	}
}

func pressImmediately(btn BtnOrAxisT) {
	commandInfo := getCommandInfo(btn, false)

	pressIfNotAlready(btn, commandInfo)
}

func CreateVirtualButton(command CommandT) (BtnOrAxisT, *CommandInfoT) {
	virtualButton := VirtualButtonCounter.GetButton()
	commandInfo := MakeEmptyCommandInfo(command)

	return virtualButton, commandInfo
}

func PressVirtualButton(btn BtnOrAxisT, command CommandT) {
	commandInfo := MakeEmptyCommandInfo(command)
	pressIfNotAlready(btn, commandInfo)
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

func RepeatCommand() {
	ButtonsLock.Lock()
	//Esc button's releaseAll will break state (changing map over iteration)
	//RangeOverCopy prevents this: states will be restored, esc command executed,
	//and then states will be properly released
	buttonsToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
		if isEmptyCmd(commandInfo.command) {
			//if hold state Interval passed
			if commandInfo.DecreaseInterval() {
				//assign hold command, reset interval
				commandInfo.CopyFromOther(getCommandInfo(btn, true))
				pressSequence(btn, commandInfo)
			}
		} else { //if command already assigned to hold or immediate
			if commandInfo.DecreaseInterval() {
				pressSequence(btn, commandInfo)
			}
		}
	})
	ButtonsLock.Unlock()
}

func releaseAll(curButton BtnOrAxisT) {
	buttonsToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
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
	if value > Cfg.Buttons.TriggerThreshold {
		pressImmediately(btn)
	} else if value < Cfg.Buttons.TriggerThreshold {
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
