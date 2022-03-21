package mainLogic

import (
	"ControllerGo/src/osSpec"
	"path"
	"sync"
	"time"
)

type CommandsMode struct {
	mode bool
}

func (c *CommandsMode) switchMode() {
	c.mode = !c.mode
	movementCoords.reset()
	scrollMovement.reset()
	mousePad.reset()
}

var typingMode = CommandsMode{}
var needToSwitchBackLang = false

const NoAction = -1
const SwitchToTyping = -2

var commonCmdMapping = map[string]int{
	"NoAction":       NoAction,
	"LeftMouse":      osSpec.LeftMouse,
	"RightMouse":     osSpec.RightMouse,
	"MiddleMouse":    osSpec.MiddleMouse,
	"SwitchToTyping": SwitchToTyping,
}

type HoldStartTime map[BtnOrAxisT]time.Time
type Command []int
type ButtonToCommand map[BtnOrAxisT]Command

var holdStartTime = HoldStartTime{}
var buttonsMutex = sync.Mutex{}
var buttonsToRelease = ButtonToCommand{}

var pressCommandsLayout, releaseCommandsLayout ButtonToCommand

func loadCommandsLayout() (ButtonToCommand, ButtonToCommand) {
	pressLayout := ButtonToCommand{}
	releaseLayout := ButtonToCommand{}
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
		releaseLayout[btn] = reverse(codes)
	}
	return pressLayout, releaseLayout
}

func getPressCommand(btn BtnOrAxisT, hold bool) Command {
	if hold {
		btn = addHoldSuffix(btn)
	}
	return pressCommandsLayout[btn]
}

func getReleaseCommand(btn BtnOrAxisT, hold bool) Command {
	if hold {
		btn = addHoldSuffix(btn)
	}
	return releaseCommandsLayout[btn]
}

func pressCommand(command Command) {
	for _, el := range command {
		osSpec.PressKeyOrMouse(el)
	}
}

func press(btn BtnOrAxisT, hold bool) {
	command := getPressCommand(btn, hold)

	switch command[0] {
	case SwitchToTyping:
		typingMode.switchMode()
		return
	case EscLetter:
		releaseAll()
	}

	buttonsToRelease[btn] = getReleaseCommand(btn, hold)
	//if len(command) > 1 && command[0] == controlKey {
	//	locale := osSpecific.GetLocale()
	//	print(locale)
	//}
	pressCommand(command)
}

func releaseCommand(command Command) {
	for _, el := range command {
		osSpec.ReleaseKeyOrMouse(el)
	}
}

func release(btn BtnOrAxisT) {
	command := pop(buttonsToRelease, btn)
	releaseCommand(command)
}

func releaseAll() {
	for _, command := range buttonsToRelease {
		releaseCommand(command)
	}
	buttonsToRelease = ButtonToCommand{}
}

var triggersPressed = map[BtnOrAxisT]bool{
	BtnLeftTrigger2:  false,
	BtnRightTrigger2: false,
}

func isTriggerBtn(btn BtnOrAxisT) bool {
	return btn == BtnLeftTrigger2 || btn == BtnRightTrigger2
}

func detectTriggers() {
	btn := event.btnOrAxis
	if !isTriggerBtn(btn) {
		return
	}
	value := event.value
	triggerPressed := triggersPressed[btn]

	if value > TriggerThreshold && !triggerPressed {
		triggersPressed[btn] = true
		pressCommand(pressCommandsLayout[btn])
	} else if value < TriggerThreshold && triggerPressed {
		triggersPressed[btn] = false
		releaseCommand(releaseCommandsLayout[btn])
	}
}

func buttonPressed() {
	btn := event.btnOrAxis
	if isTriggerBtn(btn) {
		return
	}
	buttonsMutex.Lock()
	defer buttonsMutex.Unlock()

	if _, found := pressCommandsLayout[addHoldSuffix(btn)]; found {
		holdStartTime[btn] = time.Now()
	} else {
		press(btn, false)
	}
}

const holdRefreshInterval = 15 * time.Millisecond

func RunReleaseHoldThread() {
	var holdDuration time.Duration
	for {
		buttonsMutex.Lock()
		for btn, startTime := range holdStartTime {
			holdDuration = time.Now().Sub(startTime)
			//print("duration: %v", holdDuration)
			if holdDuration > holdThreshold {
				press(btn, true)
				delete(holdStartTime, btn)
			}
		}
		buttonsMutex.Unlock()

		time.Sleep(holdRefreshInterval)
	}
}

func buttonReleased() {
	btn := event.btnOrAxis
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
