package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	"strings"
)

var LetterToCodes = map[string]int{
	"Grave":        uinput.KeyGrave,
	"Minus":        uinput.KeyMinus,
	"Equal":        uinput.KeyEqual,
	"Leftbrace":    uinput.KeyLeftbrace,
	"Rightbrace":   uinput.KeyRightbrace,
	"Comma":        uinput.KeyComma,
	"Dot":          uinput.KeyDot,
	"Slash":        uinput.KeySlash,
	"Backslash":    uinput.KeyBackslash,
	"Semicolon":    uinput.KeySemicolon,
	"Apostrophe":   uinput.KeyApostrophe,
	"Esc":          uinput.KeyEsc,
	"Backspace":    uinput.KeyBackspace,
	"Delete":       uinput.KeyDelete,
	"Space":        uinput.KeySpace,
	"Enter":        uinput.KeyEnter,
	"Tab":          uinput.KeyTab,
	"Windows":      uinput.KeyLeftmeta,
	"LeftControl":  uinput.KeyLeftctrl,
	"RightControl": uinput.KeyRightctrl,
	"LeftAlt":      uinput.KeyLeftalt,
	"RightAlt":     uinput.KeyRightalt,
	"LeftShift":    uinput.KeyLeftshift,
	"RightShift":   uinput.KeyRightshift,
	"Up":           uinput.KeyUp,
	"Down":         uinput.KeyDown,
	"Left":         uinput.KeyLeft,
	"Right":        uinput.KeyRight,

	"A": uinput.KeyA,
	"B": uinput.KeyB,
	"C": uinput.KeyC,
	"D": uinput.KeyD,
	"E": uinput.KeyE,
	"F": uinput.KeyF,
	"G": uinput.KeyG,
	"H": uinput.KeyH,
	"I": uinput.KeyI,
	"J": uinput.KeyJ,
	"K": uinput.KeyK,
	"L": uinput.KeyL,
	"M": uinput.KeyM,
	"N": uinput.KeyN,
	"O": uinput.KeyO,
	"P": uinput.KeyP,
	"Q": uinput.KeyQ,
	"R": uinput.KeyR,
	"S": uinput.KeyS,
	"T": uinput.KeyT,
	"U": uinput.KeyU,
	"V": uinput.KeyV,
	"W": uinput.KeyW,
	"X": uinput.KeyX,
	"Y": uinput.KeyY,
	"Z": uinput.KeyZ,

	"0": uinput.Key0,
	"1": uinput.Key1,
	"2": uinput.Key2,
	"3": uinput.Key3,
	"4": uinput.Key4,
	"5": uinput.Key5,
	"6": uinput.Key6,
	"7": uinput.Key7,
	"8": uinput.Key8,
	"9": uinput.Key9,
}

func getCodeFromLetter(letter string) int {
	letter = strings.ToLower(letter)
	if code, found := LetterToCodes[letter]; found {
		return code
	} else {
		panic(fmt.Sprintf("No such letter in mapping %s\n", letter))
	}
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
	}
	toLowerMap(synonyms)
	for k, v := range synonyms {
		synonyms[k] = strings.ToLower(v)
	}

	toLowerMap(LetterToCodes)
	for orig, synonym := range synonyms {
		if code, found := LetterToCodes[orig]; found {
			LetterToCodes[synonym] = code
		}
	}
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
