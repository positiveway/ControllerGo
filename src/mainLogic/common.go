package mainLogic

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ToInt(value string) int {
	res, err := strconv.Atoi(value)
	CheckErr(err)
	return res
}

func ToMilliseconds(value string) time.Duration {
	valInt := ToInt(value)
	return time.Duration(valInt) * time.Millisecond
}

func ToIntToFloat(value string) float64 {
	return float64(ToInt(value))
}

func ToFloat(value string) float64 {
	res, err := strconv.ParseFloat(value, 64)
	CheckErr(err)
	return res
}

func AssignWithDuplicateCheck[K comparable, V any](m map[K]V, key K, val V) {
	if _, found := m[key]; found {
		panic("duplicate position")
	}
	m[key] = val
}

func SplitByAnyOf(str string, separators string) []string {
	if separators == "" {
		panic("Empty separator")
	}
	var res []string
	prevSplitInd := 0
	for ind, symbol := range str {
		if strings.ContainsRune(separators, symbol) {
			res = append(res, str[prevSplitInd:ind])
			prevSplitInd = ind + 1
		}
	}
	prevSplitInd = min(len(str), prevSplitInd)
	res = append(res, str[prevSplitInd:])
	return res
}

func StartsWithAnyOf(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func ReadFile(file string) string {
	dat, err := os.ReadFile(file)
	CheckErr(err)
	return string(dat)
}

func ReadLines(file string) []string {
	content := ReadFile(file)
	return strings.Split(content, "\n")
}

func ReadLayoutFile(pathFromLayoutsDir string, skipLines int) [][]string {
	file := filepath.Join(LayoutsDir, pathFromLayoutsDir)
	lines := ReadLines(file)
	lines = lines[skipLines:]

	var linesParts [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || StartsWithAnyOf(line, []string{";", "//"}) {
			continue
		}
		parts := SplitByAnyOf(line, "&|>:,=")
		for ind, part := range parts {
			parts[ind] = strings.TrimSpace(part)
		}
		linesParts = append(linesParts, parts)
	}
	return linesParts
}

func PanicMisspelled(str string) {
	panic(fmt.Sprintf("Probably misspelled: %s\n", str))
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getOrDefault[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if val, found := m[key]; found {
		return val
	} else {
		return defaultVal
	}
}

type Number interface {
	int64 | float64 | int | int32 | float32
}

type BasicType interface {
	Number | string | bool | rune
}

func applySign(sign bool, val *float64) {
	if sign {
		*val *= -1
	}
}

func getSignMakeAbs(val *float64) bool {
	sign := math.Signbit(*val)
	*val = math.Abs(*val)
	return sign
}

func round(number float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(number*multiplier) / multiplier
}

func max[T Number](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func reverse[T BasicType](seq []T) []T {
	var res []T
	for _, el := range seq {
		res = append(res, el)
	}
	return res
}

func contains[T BasicType](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//A nil argument is equivalent to an empty slice
func equal[T BasicType](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func isGreater[T Number](oldValue, newValue T) bool {
	return math.Abs(float64(newValue)) > math.Abs(float64(oldValue))
}