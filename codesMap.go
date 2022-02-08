package main

import (
	"github.com/bendahl/uinput"
	"strings"
)

var LetterToCodes = map[string]int{
	"`":         uinput.KeyGrave,
	"-":         uinput.KeyMinus,
	"=":         uinput.KeyEqual,
	"[":         uinput.KeyLeftbrace,
	"]":         uinput.KeyRightbrace,
	"Comma":     uinput.KeyComma,
	"Dot":       uinput.KeyDot,
	"/":         uinput.KeySlash,
	"\\":        uinput.KeyBackslash,
	"Semicolon": uinput.KeySemicolon,
	"'":         uinput.KeyApostrophe,
	"Esc":       uinput.KeyEsc,
	"Backspace": uinput.KeyBackspace,
	"BS":        uinput.KeyBackspace,
	"Delete":    uinput.KeyDelete,
	"Del":       uinput.KeyDelete,
	"Space":     uinput.KeySpace,
	"Enter":     uinput.KeyEnter,
	"Tab":       uinput.KeyTab,
	"Windows":   uinput.KeyLeftmeta,
	"Control":   uinput.KeyLeftctrl,
	"Alt":       uinput.KeyLeftalt,
	"Shift":     uinput.KeyLeftshift,

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

func addLowercaseLetters() {
	for k, v := range LetterToCodes {
		k = strings.ToLower(k)
		LetterToCodes[k] = v
	}
}
