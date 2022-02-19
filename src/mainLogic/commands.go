package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"path"
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
}

func (c *CommandsMode) get() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._mode
}

var typingMode = CommandsMode{}
var needToSwitchBackLang = false

const NoAction = -1
const SwitchToTyping = -2

var commonCmdMapping = map[string]int{
	"NoAction":       NoAction,
	"LeftMouse":      osSpecific.LeftMouse,
	"RightMouse":     osSpecific.RightMouse,
	"MiddleMouse":    osSpecific.MiddleMouse,
	"SwitchToTyping": SwitchToTyping,
}

var holdTimeMutex = sync.Mutex{}
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
	linesParts := ReadLayoutFile(path.Join(LayoutInUse, "cmd_layout.csv"))
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
		if len(codes) == 0 {
			panic(fmt.Sprintf("Empty command mapping for button %s\n", btn))
		}
		if codes[0] == NoAction {
			continue
		}
		layout[btn] = codes
	}
	return layout
}

var commandsLayout CommandsLayout

func press(seq []int) {
	if seq[0] == SwitchToTyping {
		typingMode.switchMode()
		return
	}
	//if len(seq) > 1 && seq[0] == controlKey {
	//	locale := osSpecific.GetLocale()
	//	println(locale)
	//}
	for _, el := range seq {
		osSpecific.PressKeyOrMouse(el)
	}
}

func release(seq []int) {
	if seq[0] == SwitchToTyping {
		return
	}
	for _, el := range reverse(seq) {
		osSpecific.ReleaseKeyOrMouse(el)
	}
}

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
		holdTimeMutex.Lock()
		holdStartTime[holdBtn] = time.Now()
		holdTimeMutex.Unlock()
		return
	}
	press(commandsLayout[btn])
}

func RunReleaseHoldThread() {
	for {
		holdTimeMutex.Lock()
		for holdBtn, startTime := range holdStartTime {
			holdDuration := time.Now().Sub(startTime)
			//fmt.Printf("duration: %v\n", holdDuration)
			if holdDuration > holdThreshold {
				press(commandsLayout[holdBtn])
				delete(holdStartTime, holdBtn)
			}
		}
		holdTimeMutex.Unlock()

		time.Sleep(DefaultWaitInterval)
	}
}

func buttonReleased(btn string) {
	if isTriggerBtn(btn) {
		return
	}
	holdBtn := btn + HoldSuffix
	if _, found := commandsLayout[holdBtn]; found {
		holdTimeMutex.Lock()
		if _, exist := holdStartTime[holdBtn]; exist {
			press(commandsLayout[btn])
			delete(holdStartTime, holdBtn)
		} else {
			btn = holdBtn
		}
		holdTimeMutex.Unlock()
	}
	release(commandsLayout[btn])
}
