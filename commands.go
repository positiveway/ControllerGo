package main

import (
	"github.com/bendahl/uinput"
	"sync"
)

type CommandsMode struct {
	_mode bool
	mu    sync.Mutex
}

func (c *CommandsMode) switchMode() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._mode = !c._mode

	mouseMovement.reset()
	scrollMovement.reset()
	//for k := range triggersPressed {
	//	triggersPressed[k] = false
	//}
}

func (c *CommandsMode) get() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._mode
}

var typingMode = CommandsMode{}

var TRIGGERS = []string{BtnLeftTrigger2, BtnRightTrigger2}

var NoAction = []int{-1}

var controlKey = uinput.KeyLeftctrl
var copyCmd = []int{controlKey, uinput.KeyC}
var pasteCmd = []int{controlKey, uinput.KeyV}
var cutCmd = []int{controlKey, uinput.KeyX}
var undoCmd = []int{controlKey, uinput.KeyZ}

var SwitchLang = []int{uinput.KeyRightalt}
var altTab = []int{uinput.KeyLeftalt, uinput.KeyTab}

var LeftMouse = -2
var RightMouse = -3
var MiddleMouse = -4
var SwitchToTyping = []int{-5}

var commandsMap = map[string][]int{
	BtnSouth:         undoCmd,
	BtnEast:          {uinput.KeyBackspace},
	BtnNorth:         {uinput.KeySpace},
	BtnWest:          {uinput.KeyEnter},
	BtnC:             NoAction,
	BtnZ:             NoAction,
	BtnLeftTrigger:   SwitchToTyping,
	BtnLeftTrigger2:  {RightMouse},
	BtnRightTrigger:  SwitchLang,
	BtnRightTrigger2: {LeftMouse},
	BtnSelect:        {uinput.KeyLeftmeta},
	BtnStart:         NoAction,
	BtnMode:          NoAction,
	BtnLeftThumb:     copyCmd,
	BtnRightThumb:    pasteCmd,
	BtnDPadUp:        {uinput.KeyUp},
	BtnDPadDown:      {uinput.KeyDown},
	BtnDPadLeft:      {uinput.KeyLeft},
	BtnDPadRight:     {uinput.KeyRight},
	BtnUnknown:       NoAction,
}

func press(seq []int) {
	if equal(seq, NoAction) {
		return
	}
	if equal(seq, SwitchToTyping) {
		typingMode.switchMode()
		return
	}
	for _, el := range seq {
		switch el {
		case LeftMouse:
			mouse.LeftPress()
		case RightMouse:
			mouse.RightPress()
		case MiddleMouse:
			mouse.MiddlePress()
		default:
			keyboard.KeyDown(el)
		}
	}
}

func release(seq []int) {
	if equal(seq, NoAction) {
		return
	}
	if equal(seq, SwitchToTyping) {
		return
	}
	for _, el := range reverse(seq) {
		switch el {
		case LeftMouse:
			mouse.LeftRelease()
		case RightMouse:
			mouse.RightRelease()
		case MiddleMouse:
			mouse.MiddleRelease()
		default:
			keyboard.KeyUp(el)
		}
	}
}

const TriggerThreshold = 0.3

var triggersPressed = map[string]bool{
	BtnLeftTrigger2:  false,
	BtnRightTrigger2: false,
}

func detectTriggers(event Event) {
	btn := event.btnOrAxis
	if !contains(TRIGGERS, btn) {
		return
	}
	command := commandsMap[btn]

	if event.value > TriggerThreshold && !triggersPressed[btn] {
		triggersPressed[btn] = true
		press(command)
	} else if event.value < TriggerThreshold && triggersPressed[btn] {
		triggersPressed[btn] = false
		release(command)
	}
}

func buttonPressed(btn string) {
	if contains(TRIGGERS, btn) {
		return
	}
	command := commandsMap[btn]
	press(command)
}

func buttonReleased(btn string) {
	if contains(TRIGGERS, btn) {
		return
	}
	command := commandsMap[btn]
	release(command)
}
