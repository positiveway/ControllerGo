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

func ToLowerMap[V any](mapping map[string]V) {
	for k, v := range mapping {
		if k != strings.TrimSpace(k) {
			gofuncs.Panic("Mapping identifiers check failed")
		}
		delete(mapping, k)
		k = strings.ToLower(k)
		mapping[k] = v
	}
}

func convertLetterToCodeMapping() {
	synonyms := map[string]string{
		"Control": "LeftControl",
		"Ctrl":    "LeftControl",
		"Alt":     "LeftAlt",
		"Shift":   "LeftShift",
		"BS":      "Backspace",
		"Del":     "Delete",
		"Caps":    "CapsLock",
	}
	ToLowerMap(synonyms)
	for synonym, orig := range synonyms {
		synonyms[synonym] = strings.ToLower(orig)
	}

	ToLowerMap(osSpec.LetterToCodes)
	for synonym, orig := range synonyms {
		if code, found := osSpec.LetterToCodes[orig]; found {
			gofuncs.AssignWithDuplicateKeyCheck(osSpec.LetterToCodes, synonym, code)
		} else {
			gofuncs.Panic("No such button name: %v", orig)
		}
	}
}

var EscLetterCode int

func initCodeMapping() {
	convertLetterToCodeMapping()
	EscLetterCode = getCodeFromLetter("Esc")
}
