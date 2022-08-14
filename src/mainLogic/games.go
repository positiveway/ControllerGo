package mainLogic

import "time"

var fastKeyRepeatInterval time.Duration = 40
var slowKeyRepeatInterval time.Duration = 100

func RunGameMovementThread() {
	//shiftCode := getCodeFromLetter("shift")
	//shiftPressed := false

	for {
		keyRepeatInterval := fastKeyRepeatInterval

		if Cfg.PadsSticksMode.GetMode() != MouseMode {
			time.Sleep(keyRepeatInterval)
			continue
		}

		time.Sleep(keyRepeatInterval)
	}
}
