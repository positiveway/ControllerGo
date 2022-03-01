package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"strings"
)

func getCodeFromLetter(letter string) int {
	letter = strings.ToLower(letter)
	if code, found := osSpecific.LetterToCodes[letter]; found {
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

	toLowerMap(osSpecific.LetterToCodes)
	for orig, synonym := range synonyms {
		if code, found := osSpecific.LetterToCodes[orig]; found {
			osSpecific.LetterToCodes[synonym] = code
		}
	}
}

var EscLetter int

func initCodeMapping() {
	convertLetterToCodeMapping()
	EscLetter = getCodeFromLetter("Esc")
}

const (
	ZoneRight     string = "Right"
	ZoneUpRight          = "UpRight"
	ZoneUp               = "Up"
	ZoneUpLeft           = "UpLeft"
	ZoneLeft             = "Left"
	ZoneDownLeft         = "DownLeft"
	ZoneDown             = "Down"
	ZoneDownRight        = "DownRight"
)

var AllZones = []string{
	ZoneRight,
	ZoneUpRight,
	ZoneUp,
	ZoneUpLeft,
	ZoneLeft,
	ZoneDownLeft,
	ZoneDown,
	ZoneDownRight,
}
