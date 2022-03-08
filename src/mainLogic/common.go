package mainLogic

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StrToBool(value string) bool {
	res, err := strconv.ParseBool(value)
	CheckErr(err)
	return res
}

func StrToInt(value string) int {
	res, err := strconv.Atoi(value)
	CheckErr(err)
	return res
}

func StrToMillis(value string) time.Duration {
	number := StrToFloat(value)
	return numberToMillis(number)
}

func StrToIntToFloat(value string) float64 {
	return float64(StrToInt(value))
}

func StrToFloat(value string) float64 {
	res, err := strconv.ParseFloat(value, 64)
	CheckErr(err)
	return res
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

func StartsWithAnyOf(str string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func EndsWithAnyOf(str string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

func TrimAnyPrefix(str string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return strings.TrimPrefix(str, prefix)
		}
	}
	return str
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
		if line == "" || StartsWithAnyOf(line, ";", "//") {
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

func sPrint(message string, variables ...any) string {
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	return fmt.Sprintf(message, variables...)
}

func print(message string, variables ...any) {
	fmt.Print(sPrint(message, variables...))
}

func panicMsg(message string, variables ...any) {
	releaseAll()
	panic(sPrint(message, variables...))
}

func PanicMisspelled(str any) {
	panicMsg("Probably misspelled: %v", str)
}

func CheckErr(err error) {
	if err != nil {
		panicMsg("%v", err)
	}
}

func pop[K comparable, V any](m map[K]V, key K) V {
	value := m[key]
	delete(m, key)
	return value
}

func AssignWithDuplicateCheck[K comparable, V any](m map[K]V, key K, val V) {
	if _, found := m[key]; found {
		panic("duplicate position")
	}
	m[key] = val
}

func getOrDefault[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if val, found := m[key]; found {
		return val
	} else {
		return defaultVal
	}
}

type Int interface {
	int | int32 | int64
}

type Number interface {
	Int | float64 | float32
}

type BasicType interface {
	Number | string | bool | rune
}

func applySign(sign bool, val float64) float64 {
	if sign {
		val *= -1
	}
	return val
}

func getSignAndAbs(val float64) (bool, float64) {
	sign := math.Signbit(val)
	val = math.Abs(val)
	return sign, val
}

func floatToInt(value float64) int {
	return int(math.Round(value))
}

func floatToInt32(value float64) int32 {
	return int32(math.Round(value))
}

func floatToInt64(value float64) int64 {
	return int64(math.Round(value))
}

func numberToMillis[T Number](value T) time.Duration {
	return time.Duration(float64(value)*1000) * time.Microsecond
}

func trunc(number float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Trunc(number*multiplier) / multiplier
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

func isEmpty[T BasicType](seq []T) bool {
	return len(seq) == 0
}

func reverse[T BasicType](seq []T) []T {
	var res []T
	for i := len(seq) - 1; i >= 0; i-- {
		res = append(res, seq[i])
	}
	return res
}

func contains[T comparable](s []T, e T) bool {
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

type ThreadSafeMap struct {
	mapping map[string]any
	mutex   sync.Mutex
}

func (threadMap *ThreadSafeMap) set(key string, value any) {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	threadMap.mapping[key] = value
}

func (threadMap *ThreadSafeMap) get(key string) any {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	return threadMap.mapping[key]
}

func (threadMap *ThreadSafeMap) pop(key string) any {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	value := threadMap.mapping[key]
	delete(threadMap.mapping, key)
	return value
}
