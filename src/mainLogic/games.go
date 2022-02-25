package mainLogic

import "time"

var movementCoords = Coords{}

var fastKeyRepeatInterval time.Duration = 40
var slowKeyRepeatInterval time.Duration = 100

func RunMovementThread() {
	//shiftCode := getCodeFromLetter("shift")
	//shiftPressed := false

	for {
		//metrics := movementCoords.getMetrics()
		time.Sleep(DefaultRefreshInterval)
	}
}
