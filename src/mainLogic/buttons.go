package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

const (
	NoAction                = -10
	SwitchPadStickMode      = -11
	SwitchHighPrecisionMode = -12
)

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

func (buttons *ButtonsT) loadCommandsLayout() ButtonToCommandT {
	allBtnAxis := buttons.allBtnAxis

	commonCmdMapping := initCommonCmdMapping()
	BtnSynonyms := allBtnAxis.genBtnSynonyms()
	AllAvailableButtons := allBtnAxis.initAvailableButtons()

	pressLayout := ButtonToCommandT{}
	linesParts := gofuncs.ReadLayoutFile(2,
		[]string{buttons.cfg.Path.CurLayoutDir, "buttons.csv"})

	for _, parts := range linesParts {
		btn := BtnOrAxisT(parts[0])
		keys := parts[1:]

		btn = ToLower(btn)

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
				code := buttons.getCodeFromLetter(key)
				codes = append(codes, code)
			}
		}
		if isEmptyCmd(codes) {
			gofuncs.Panic("Empty command mapping for button %s", btn)
		}
		if codes[0] == NoAction {
			continue
		}
		pressLayout[btn] = MakeEmptyCommandInfo(codes, buttons.cfg)
	}
	return pressLayout
}

type CommandT []int

type ButtonToCommandT map[BtnOrAxisT]*CommandInfoT

type CommandInfoT struct {
	CfgStruct
	RepeatedTimerT
	command              CommandT
	specialCaseIsHandled bool
}

func MakeCommandInfo(command CommandT, repeatInterval float64, cfg *ConfigsT) *CommandInfoT {
	commandInfo := &CommandInfoT{command: command}
	commandInfo.Init(cfg)
	commandInfo.InitIntervalTimer(repeatInterval, cfg)
	return commandInfo
}

func MakeEmptyCommandInfo(command CommandT, cfg *ConfigsT) *CommandInfoT {
	return MakeCommandInfo(command, cfg.Buttons.HoldRepeatInterval, cfg)
}

func MakeUndeterminedCommandInfo(cfg *ConfigsT) *CommandInfoT {
	return MakeCommandInfo(nil, cfg.Buttons.HoldingStateThreshold, cfg)
}

func (c *CommandInfoT) GetCopy() *CommandInfoT {
	return MakeCommandInfo(c.command, c.repeatInterval, c.cfg)
}

func (c *CommandInfoT) CopyFromOther(other *CommandInfoT) {
	c.command = other.command
	c.SetInterval(other.repeatInterval)
}

type ButtonsT struct {
	CfgLockStruct
	allBtnAxis           *AllBtnAxis
	highPrecisionMode    *HighPrecisionModeT
	ToRelease            *gofuncs.Map[BtnOrAxisT, *CommandInfoT]
	ToCommandLayout      ButtonToCommandT
	EscLetterCode        int
	virtualButtonCounter uint
	getCodeFromLetter    CodesFromLetterFuncT
	isTriggerBtn         IsTriggerBtnFuncT
	handleTriggers       func(btn BtnOrAxisT, value float64)
	pressSequence        func(btn BtnOrAxisT, commandInfo *CommandInfoT)
}

func (buttons *ButtonsT) Init(cfg *ConfigsT, highPrecisionMode *HighPrecisionModeT, allBtnAxis *AllBtnAxis) {
	buttons.CfgLockStruct.Init(cfg)
	buttons.allBtnAxis = allBtnAxis
	buttons.highPrecisionMode = highPrecisionMode

	buttons.getCodeFromLetter = GetGetCodesFromLetterFunc()
	//should come before other functions initialization
	buttons.EscLetterCode = buttons.getCodeFromLetter("Esc")

	buttons.isTriggerBtn = buttons.GetIsTriggerBtnFunc()
	buttons.handleTriggers = buttons.GetHandleTriggersFunc()
	buttons.pressSequence = buttons.GetPressSequenceFunc()

	buttons.ToRelease = gofuncs.MakeMap[BtnOrAxisT, *CommandInfoT]()
	buttons.ToCommandLayout = buttons.loadCommandsLayout()
}

func (buttons *ButtonsT) GetVirtualButton() BtnOrAxisT {
	buttons.Lock()
	defer buttons.Unlock()

	//if uint overflows it will be zero
	buttons.virtualButtonCounter += 1

	virtualButton := gofuncs.Format("VirtualButton_%v",
		buttons.virtualButtonCounter)
	return BtnOrAxisT(virtualButton)
}

func (buttons *ButtonsT) PutButton(btn BtnOrAxisT, commandInfo *CommandInfoT) bool {
	if _, exist := buttons.ToRelease.CheckAndGet(btn); exist {
		return false
	}
	buttons.ToRelease.Put(btn, commandInfo)
	return true
}

func isEmptyCommandInfo(commandInfo *CommandInfoT) bool {
	return commandInfo == nil
}

func (buttons *ButtonsT) isEmptyCommandForButton(btn BtnOrAxisT, hold bool) bool {
	return isEmptyCommandInfo(buttons.getCommandInfo(btn, hold))
}

func isEmptyCmd(command CommandT) bool {
	return gofuncs.IsEmptySlice(command)
}

func (buttons *ButtonsT) commandNotExists(btn BtnOrAxisT) bool {
	return buttons.isEmptyCommandForButton(btn, false) &&
		buttons.isEmptyCommandForButton(btn, true)
}

func (buttons *ButtonsT) hasHoldCommand(btn BtnOrAxisT) bool {
	return !buttons.isEmptyCommandForButton(btn, true)
}

func (buttons *ButtonsT) getCommandInfo(btn BtnOrAxisT, hold bool) *CommandInfoT {
	if hold {
		btn = addHoldSuffix(btn)
	}
	//Have only one point of access. Don't forget to copy
	commandInfo := buttons.ToCommandLayout[btn]
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

func (buttons *ButtonsT) GetPressSequenceFunc() func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
	padsSticksMode := buttons.cfg.PadsSticks.Mode
	highPrecisionMode := buttons.highPrecisionMode
	EscLetterCode := buttons.EscLetterCode

	return func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
		command := commandInfo.command

		if !commandInfo.specialCaseIsHandled {
			commandInfo.specialCaseIsHandled = true

			switch getFirstCmdSymbol(command) {
			case SwitchPadStickMode:
				// don't do release all
				padsSticksMode.SwitchMode()
				return
			case SwitchHighPrecisionMode:
				highPrecisionMode.SwitchMode()
				return
			case EscLetterCode:
				buttons.releaseAll(btn)
			}
		}

		for _, el := range command {
			osSpec.PressKeyOrMouse(el)
		}
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

func (buttons *ButtonsT) pressIfNotAlready(btn BtnOrAxisT, commandInfo *CommandInfoT) {
	if isEmptyCommandInfo(commandInfo) {
		gofuncs.Panic("CommandInfo can't be empty at this point")
	}
	if buttons.PutButton(btn, commandInfo) {
		buttons.pressSequence(btn, commandInfo)
	}
}

func (buttons *ButtonsT) pressImmediately(btn BtnOrAxisT) {
	commandInfo := buttons.getCommandInfo(btn, false)

	buttons.pressIfNotAlready(btn, commandInfo)
}

func (buttons *ButtonsT) CreateVirtualButton(command CommandT) (BtnOrAxisT, *CommandInfoT) {
	virtualButton := buttons.GetVirtualButton()
	commandInfo := MakeEmptyCommandInfo(command, buttons.cfg)

	return virtualButton, commandInfo
}

func (buttons *ButtonsT) CreateAndPressVirtualButton(command CommandT) BtnOrAxisT {
	virtualButton, commandInfo := buttons.CreateVirtualButton(command)
	buttons.pressIfNotAlready(virtualButton, commandInfo)
	return virtualButton
}

func (buttons *ButtonsT) pressButton(btn BtnOrAxisT) {
	if buttons.hasHoldCommand(btn) {
		buttons.PutButton(btn, MakeUndeterminedCommandInfo(buttons.cfg))
	} else {
		buttons.pressImmediately(btn)
	}
}

func (buttons *ButtonsT) releaseButton(btn BtnOrAxisT) {
	commandInfo := buttons.ToRelease.Pop(btn)

	if isEmptyCommandInfo(commandInfo) {
		return
	}

	if isEmptyCmd(commandInfo.command) {
		//has hold command but no "immediately press" command
		//and not enough time have passed for hold command to be triggered
		if buttons.isEmptyCommandForButton(btn, false) {
			return
		}
		commandInfo.CopyFromOther(buttons.getCommandInfo(btn, false))
		buttons.pressSequence(btn, commandInfo)
	}

	releaseSequence(commandInfo.command)
}

func (buttons *ButtonsT) RepeatCommand() {
	buttons.Lock()
	//Esc button's releaseAll will break state (changing map over iteration)
	//RangeOverCopy prevents this: states will be restored, esc command executed,
	//and then states will be properly released
	buttons.ToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
		if isEmptyCmd(commandInfo.command) {
			//if hold state Interval passed
			if commandInfo.DecreaseInterval() {
				//assign hold command, reset interval
				commandInfo.CopyFromOther(buttons.getCommandInfo(btn, true))
				buttons.pressSequence(btn, commandInfo)
			}
		} else { //if command already assigned to hold or immediate
			if commandInfo.DecreaseInterval() {
				buttons.pressSequence(btn, commandInfo)
			}
		}
	})
	buttons.Unlock()
}

func (buttons *ButtonsT) releaseAll(curButton BtnOrAxisT) {
	buttons.ToRelease.RangeOverShallowCopy(func(btn BtnOrAxisT, commandInfo *CommandInfoT) {
		//current button should stay in map
		if btn != curButton {
			buttons.releaseButton(btn)
		}
	})
}

type IsTriggerBtnFuncT func(btn BtnOrAxisT) bool

func (buttons *ButtonsT) GetIsTriggerBtnFunc() IsTriggerBtnFuncT {
	allBtnAxis := buttons.allBtnAxis
	BtnLeftTrigger := allBtnAxis.BtnLeftTrigger
	BtnRightTrigger := allBtnAxis.BtnRightTrigger

	return func(btn BtnOrAxisT) bool {
		return btn == BtnLeftTrigger || btn == BtnRightTrigger
	}
}

func (buttons *ButtonsT) GetHandleTriggersFunc() func(btn BtnOrAxisT, value float64) {
	triggerThreshold := buttons.cfg.Buttons.TriggerThreshold

	return func(btn BtnOrAxisT, value float64) {
		if value > triggerThreshold {
			buttons.pressImmediately(btn)
		} else if value < triggerThreshold {
			buttons.releaseButton(btn)
		}
	}
}

func (buttons *ButtonsT) buttonChanged(btn BtnOrAxisT, value float64) {
	buttons.Lock()
	defer buttons.Unlock()

	if buttons.commandNotExists(btn) {
		return
	}

	if buttons.isTriggerBtn(btn) {
		buttons.handleTriggers(btn, value)
	} else {
		switch value {
		case 1:
			buttons.pressButton(btn)
		case 0:
			buttons.releaseButton(btn)
		default:
			gofuncs.Panic("Unsupported value: \"%s\" %v", btn, value)
		}
	}
}
