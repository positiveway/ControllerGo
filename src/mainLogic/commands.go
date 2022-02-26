package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"path"
	"sync"
	"time"
)

type CommandsMode struct {
	_mode bool
}

func (c *CommandsMode) switchMode() {
	c._mode = !c._mode
	mouseMovement.reset()
	scrollMovement.reset()
}

func (c *CommandsMode) get() bool {
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

type HoldStartTime = map[string]time.Time
type Command = []int
type ButtonToCommand = map[string]Command

var holdStartTime = HoldStartTime{}
var buttonsMutex = sync.Mutex{}
var buttonsToRelease = ButtonToCommand{}

var commandsLayout ButtonToCommand

func loadCommandsLayout() ButtonToCommand {
	layout := ButtonToCommand{}
	linesParts := ReadLayoutFile(path.Join(LayoutInUse, "commands.csv"), 2)
	for _, parts := range linesParts {
		btn := parts[0]
		keys := parts[1:]

		if btnSynonym, found := BtnSynonyms[btn]; found {
			btn = btnSynonym
		}
		if !contains(AllOriginalButtons, removeHoldSuffix(btn)) {
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
			panicMsg("Empty command mapping for button %s\n", btn)
		}
		if codes[0] == NoAction {
			continue
		}
		layout[btn] = codes
	}
	return layout
}

func getCommand(btn string, hold bool) Command {
	if hold {
		return commandsLayout[addHoldSuffix(btn)]
	} else {
		return commandsLayout[btn]
	}
}

func press(btn string, hold bool) {
	command := getCommand(btn, hold)

	if isEmpty(command) {
		return
	}
	switch command[0] {
	case SwitchToTyping:
		typingMode.switchMode()
		return
	case EscLetter:
		releaseAll()
	}

	buttonsToRelease[btn] = command
	//if len(command) > 1 && command[0] == controlKey {
	//	locale := osSpecific.GetLocale()
	//	println(locale)
	//}
	for _, el := range command {
		osSpecific.PressKeyOrMouse(el)
	}
}

func release(btn string) {
	command := pop(buttonsToRelease, btn)
	if isEmpty(command) {
		return
	}

	for _, el := range reverse(command) {
		osSpecific.ReleaseKeyOrMouse(el)
	}
}

func releaseAll() {
	var buttonsCopy []string
	for btn := range buttonsToRelease {
		buttonsCopy = append(buttonsCopy, btn)
	}
	for _, btn := range buttonsCopy {
		release(btn)
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
	buttonsMutex.Lock()
	defer buttonsMutex.Unlock()

	if event.value > TriggerThreshold && !triggersPressed[btn] {
		triggersPressed[btn] = true
		press(btn, false)
	} else if event.value < TriggerThreshold && triggersPressed[btn] {
		triggersPressed[btn] = false
		release(btn)
	}
}

func buttonPressed(btn string) {
	if isTriggerBtn(btn) {
		return
	}
	buttonsMutex.Lock()
	defer buttonsMutex.Unlock()

	if _, found := commandsLayout[addHoldSuffix(btn)]; found {
		holdStartTime[btn] = time.Now()
	} else {
		press(btn, false)
	}
}

const holdRefreshInterval = 15 * time.Millisecond

func RunReleaseHoldThread() {
	for {
		buttonsMutex.Lock()
		for btn, startTime := range holdStartTime {
			holdDuration := time.Now().Sub(startTime)
			//fmt.Printf("duration: %v\n", holdDuration)
			if holdDuration > holdThreshold {
				press(btn, true)
				delete(holdStartTime, btn)
			}
		}
		buttonsMutex.Unlock()

		time.Sleep(holdRefreshInterval)
	}
}

func buttonReleased(btn string) {
	if isTriggerBtn(btn) {
		return
	}
	buttonsMutex.Lock()
	defer buttonsMutex.Unlock()

	if _, found := holdStartTime[btn]; found {
		press(btn, false)
		delete(holdStartTime, btn)
	}
	release(btn)
}
