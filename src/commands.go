package main

import (
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
}

func (c *CommandsMode) get() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._mode
}

var typingMode = CommandsMode{}
var needToSwitchBackLang = false

const NoAction = -1
const LeftMouse = -2
const RightMouse = -3
const MiddleMouse = -4
const SwitchToTyping = -5

var commonCmdMapping = map[string]int{
	"NoAction":       NoAction,
	"LeftMouse":      LeftMouse,
	"RightMouse":     RightMouse,
	"MiddleMouse":    MiddleMouse,
	"SwitchToTyping": SwitchToTyping,
}

var holdStartTime = map[string]time.Time{}

type CommandsLayout = map[string][]int

func genEmptyCommandsLayout() CommandsLayout {
	layout := CommandsLayout{}
	NoActionSlice := []int{NoAction}
	for _, btn := range AllButtons {
		layout[btn] = NoActionSlice
	}
	return layout
}

func loadCommandsLayout() CommandsLayout {
	layout := genEmptyCommandsLayout()
	linesParts := readLayoutFile("cmd_layout.csv")
	for _, parts := range linesParts {
		btn := parts[0]
		keys := parts[1:]

		if btnSynonym, found := BtnSynonyms[btn]; found {
			btn = btnSynonym
		}
		if _, found := layout[btn]; !found {
			panicMisspelled(btn)
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
		layout[btn] = codes
	}
	return layout
}

var commandsLayout CommandsLayout

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
		pressKeyOrMouse(el)
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
		releaseKeyOrMouse(el)
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
	command := commandsLayout[btn]

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
	if _, found := commandsLayout[holdBtn]; found {
		holdStartTime[holdBtn] = time.Now()
		return
	}
	command := commandsLayout[btn]
	press(command)
}

const holdThreshold time.Duration = 400 * time.Millisecond

func buttonReleased(btn string) {
	if isTriggerBtn(btn) {
		return
	}
	holdBtn := btn + HoldSuffix
	if _, found := commandsLayout[holdBtn]; found {
		startTime := holdStartTime[holdBtn]
		holdDuration := time.Now().Sub(startTime)
		//fmt.Printf("duration: %v\n", holdDuration)
		if holdDuration > holdThreshold {
			btn = holdBtn
		}
		press(commandsLayout[btn])
	}
	command := commandsLayout[btn]
	release(command)
}
