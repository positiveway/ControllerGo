package main

import (
	"github.com/bendahl/uinput"
	"sync"
	"time"
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
var needToSwitchBackLang = false

var pressMouseOrKey = func(key int) {
	switch key {
	case LeftMouse:
		mouse.LeftPress()
	case RightMouse:
		mouse.RightPress()
	case MiddleMouse:
		mouse.MiddlePress()
	default:
		keyboard.KeyDown(key)
	}
}

var releaseMouseOrKey = func(key int) {
	switch key {
	case LeftMouse:
		mouse.LeftRelease()
	case RightMouse:
		mouse.RightRelease()
	case MiddleMouse:
		mouse.MiddleRelease()
	default:
		keyboard.KeyUp(key)
	}
}

//var altTab = []int{uinput.KeyLeftalt, uinput.KeyTab}

const NoAction = -1
const LeftMouse = -2
const RightMouse = -3
const MiddleMouse = -4
const SwitchToTyping = -5

var commonCmdMapping = map[string]int{
	"LeftMouse":      LeftMouse,
	"RightMouse":     RightMouse,
	"MiddleMouse":    MiddleMouse,
	"SwitchToTyping": SwitchToTyping,
}

var holdStartTime = map[string]time.Time{}

var commandsMap = map[string][]int{
	BtnSouth:         {uinput.KeyLeftctrl, uinput.KeyZ},
	BtnEast:          {uinput.KeyBackspace},
	BtnNorth:         {uinput.KeySpace},
	BtnWest:          {uinput.KeyEnter},
	BtnWestHold:      {uinput.KeyLeftctrl, uinput.KeyLeftalt, uinput.KeyL},
	BtnNorthHold:     {uinput.KeyLeftctrl, uinput.KeyLeftalt, uinput.KeyB},
	BtnC:             {NoAction},
	BtnZ:             {NoAction},
	BtnLeftTrigger:   {SwitchToTyping},
	BtnLeftTrigger2:  {RightMouse},
	BtnRightTrigger:  {uinput.KeyRightalt},
	BtnRightTrigger2: {LeftMouse},
	BtnSelect:        {uinput.KeyLeftmeta},
	BtnStart:         {uinput.KeyEsc},
	BtnMode:          {NoAction},
	BtnLeftThumb:     {uinput.KeyLeftctrl, uinput.KeyC},
	BtnRightThumb:    {uinput.KeyLeftctrl, uinput.KeyV},
	BtnDPadUp:        {uinput.KeyUp},
	BtnDPadDown:      {uinput.KeyDown},
	BtnDPadLeft:      {uinput.KeyLeft},
	BtnDPadRight:     {uinput.KeyRight},
	BtnUnknown:       {NoAction},
}

func press(seq []int) {
	switch seq[0] {
	case NoAction:
		return
	case SwitchToTyping:
		typingMode.switchMode()
		return
	}
	//if len(seq) > 1 && seq[0] == controlKey {
	//	locale := getLocale()
	//	println(locale)
	//}
	for _, el := range seq {
		pressMouseOrKey(el)
	}
}

func release(seq []int) {
	switch seq[0] {
	case NoAction:
		return
	case SwitchToTyping:
		return
	}
	for _, el := range reverse(seq) {
		releaseMouseOrKey(el)
	}
}

const TriggerThreshold float64 = 0.3

var triggersPressed = map[string]bool{
	BtnLeftTrigger2:  false,
	BtnRightTrigger2: false,
}

func isTriggerBtn(btn string) bool {
	return btn == BtnLeftTrigger2 || btn == BtnRightTrigger2
}

func detectTriggers(event Event) {
	btn := event.btnOrAxis
	if !isTriggerBtn(btn) {
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
	if isTriggerBtn(btn) {
		return
	}
	holdBtn := btn + HoldSuffix
	if _, found := commandsMap[holdBtn]; found {
		holdStartTime[holdBtn] = time.Now()
		return
	}
	command := commandsMap[btn]
	press(command)
}

const holdThreshold time.Duration = 400 * time.Millisecond

func buttonReleased(btn string) {
	if isTriggerBtn(btn) {
		return
	}
	holdBtn := btn + HoldSuffix
	if _, found := commandsMap[holdBtn]; found {
		startTime := holdStartTime[holdBtn]
		holdDuration := time.Now().Sub(startTime)
		//fmt.Printf("duration: %v\n", holdDuration)
		if holdDuration > holdThreshold {
			btn = holdBtn
		}
		press(commandsMap[btn])
	}
	command := commandsMap[btn]
	release(command)
}
