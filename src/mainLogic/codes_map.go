package mainLogic

import (
	"ControllerGo/src/osSpec"
	"strings"
)

func getCodeFromLetter(letter string) int {
	letter = strings.ToLower(letter)
	if code, found := osSpec.LetterToCodes[letter]; found {
		return code
	} else {
		panicMsg("No such letter in mapping %s", letter)
	}
	return 0
}

func toLowerMap[V any](m map[string]V) {
	for k, v := range m {
		if k != strings.TrimSpace(k) {
			panic("Mapping identifiers check failed")
		}
		delete(m, k)
		k = strings.ToLower(k)
		m[k] = v
	}
}

func convertLetterToCodeMapping() {
	synonyms := map[string]string{
		"LeftControl": "Control",
		"LeftAlt":     "Alt",
		"LeftShift":   "Shift",
		"Backspace":   "BS",
		"Delete":      "Del",
		"CapsLock":    "Caps",
	}
	toLowerMap(synonyms)
	for k, v := range synonyms {
		synonyms[k] = strings.ToLower(v)
	}

	toLowerMap(osSpec.LetterToCodes)
	for orig, synonym := range synonyms {
		if code, found := osSpec.LetterToCodes[orig]; found {
			osSpec.LetterToCodes[synonym] = code
		}
	}
}

var EscLetter int

func initCodeMapping() {
	convertLetterToCodeMapping()
	EscLetter = getCodeFromLetter("Esc")
}

type ZoneT string

const (
	ZoneRight     ZoneT = "Right"
	ZoneUpRight   ZoneT = "UpRight"
	ZoneUp        ZoneT = "Up"
	ZoneUpLeft    ZoneT = "UpLeft"
	ZoneLeft      ZoneT = "Left"
	ZoneDownLeft  ZoneT = "DownLeft"
	ZoneDown      ZoneT = "Down"
	ZoneDownRight ZoneT = "DownRight"
)

var AllZones = []ZoneT{
	ZoneRight,
	ZoneUpRight,
	ZoneUp,
	ZoneUpLeft,
	ZoneLeft,
	ZoneDownLeft,
	ZoneDown,
	ZoneDownRight,
}
