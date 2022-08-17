package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"strings"
)

func getCodeFromLetter(letter string) int {
	letter = strings.ToLower(letter)
	return gofuncs.GetOrPanic(osSpec.LetterToCodes, letter, "No such letter in mapping")
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

var EscLetterCode int

func initCodeMapping() {
	convertLetterToCodeMapping()
	EscLetterCode = getCodeFromLetter("Esc")
}
