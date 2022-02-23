package mainLogic

import "time"

var movementCoords = Coords{}

var fastKeyRepeatInterval time.Duration = 40
var slowKeyRepeatInterval time.Duration = 100

type TimeSinceLastPress struct {
	vertical, horizontal time.Duration
}

func (t *TimeSinceLastPress) update() {
	t.vertical += RefreshInterval
	t.horizontal += RefreshInterval
}

func runMovementThread() {
	//shiftCode := getCodeFromLetter("shift")
	//shiftPressed := false

	timeSinceLastPress := TimeSinceLastPress{}

	for {
		timeSinceLastPress.update()

		time.Sleep(RefreshInterval)
	}
}
