package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"strings"
)

type CodesFromLetterFuncT func(letter string) int

func GetGetCodesFromLetterFunc() CodesFromLetterFuncT {
	letterToCodes := osSpec.InitLetterToCodes()
	initLetterToCodesMapping(letterToCodes)

	return func(letter string) int {
		letter = strings.ToLower(letter)
		return gofuncs.GetOrPanic(letterToCodes, letter, "No such letter in mapping")
	}
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

func initLetterToCodesMapping(letterToCodes osSpec.LetterToCodesT) {
	ToLowerMap(letterToCodes)

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

	for synonym, orig := range synonyms {
		if code, found := letterToCodes[orig]; found {
			gofuncs.AssignWithDuplicateKeyCheck(letterToCodes, synonym, code)
		} else {
			gofuncs.Panic("No such button name: %v", orig)
		}
	}
}
