package mainLogic

import (
	"fmt"
	"github.com/jinzhu/copier"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func strToBool(value string) bool {
	res, err := strconv.ParseBool(value)
	checkErr(err)
	return res
}

func strToInt(value string) int {
	res, err := strconv.Atoi(value)
	checkErr(err)
	return res
}

func strToMillis(value string) time.Duration {
	number := strToFloat(value)
	return numberToMillis(number)
}

func strToIntToFloat(value string) float64 {
	return float64(strToInt(value))
}

func strToPct(value string) float64 {
	return strToIntToFloat(value) / 100
}

func strToFloat(value string) float64 {
	res, err := strconv.ParseFloat(value, 64)
	checkErr(err)
	return res
}

func splitByAnyOf(str string, separators string) []string {
	if separators == "" {
		panicMsg("Empty separator")
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

func startsWith(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func endsWith(s string, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func StartsWithAnyOf(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if startsWith(s, prefix) {
			return true
		}
	}
	return false
}

func endsWithAnyOf(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if endsWith(s, suffix) {
			return true
		}
	}
	return false
}

func trimAnyPrefix(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if startsWith(s, prefix) {
			return strings.TrimPrefix(s, prefix)
		}
	}
	return s
}

func strip(s string) string {
	return strings.TrimSpace(s)
}

func filterEmptyStrings(slice []string) []string {
	var filtered []string
	for _, s := range slice {
		if !isEmptyStr(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func split(s string, sep ...string) []string {

	panicIfEmptyStr(s)

	separator := func(separator []string) string {
		switch len(separator) {
		case 0:
			return ""
		case 1:
			return sep[0]
		default:
			panicMsg("Invalid separator for split")
		}
		panic("")
	}(sep)

	res := strings.Split(s, separator)
	res = filterEmptyStrings(res)
	return res
}

func words(s string) []string {
	return split(s)
}

func firstWord(s string) string {
	return words(s)[0]
}

func lastWord(s string) string {
	return lastElem(words(s))
}

func lastElem[T BasicType](slice []T) T {
	return slice[len(slice)-1]
}

func ReadFile(file string) string {
	dat, err := os.ReadFile(file)
	checkErr(err)
	return string(dat)
}

func ReadLines(file string) []string {
	content := ReadFile(file)
	return split(content, "\n")
}

func ReadLayoutFile(pathFromLayoutsDir string, skipLines int) [][]string {
	file := filepath.Join(LayoutsDir, pathFromLayoutsDir)
	lines := ReadLines(file)
	lines = lines[skipLines:]

	var linesParts [][]string
	for _, line := range lines {
		line = strip(line)
		if isEmptyStr(line) || StartsWithAnyOf(line, ";", "//") {
			continue
		}
		parts := splitByAnyOf(line, "&|>:,=")
		for ind, part := range parts {
			parts[ind] = strip(part)
		}
		linesParts = append(linesParts, parts)
	}
	return linesParts
}

func sPrint(message string, variables ...any) string {
	if !startsWith(message, "\n") {
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

func checkErr(err error) {
	if err != nil {
		panicMsg("%v", err)
	}
}

func Pop[K comparable, V any](m map[K]V, key K) V {
	value := m[key]
	delete(m, key)
	return value
}

func AssignWithDuplicateCheck[K comparable, V any](m map[K]V, key K, val V) {
	if _, found := m[key]; found {
		panicMsg("duplicate position")
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

func getPanicMsg(message []string, defaultMsg string) string {
	switch len(message) {
	case 0:
		return defaultMsg
	case 1:
		return message[0]
	default:
		panicMsg("Only one message can be specified")
	}
	panic("")
}

func getOrPanic[K comparable, V any](m map[K]V, key K, msg ...string) V {
	if val, found := m[key]; found {
		return val
	}
	message := getPanicMsg(msg, "No such key in map")

	panicMsg(message+": \"%v\"", key)
	panic("")
}

type Int interface {
	int | int32 | int64
}

type float = float64

type Float interface {
	float32 | float64
}

type Number interface {
	Int | Float
}

type BasicType interface {
	Number | string | bool | rune
}

func NaN() float64 {
	return math.NaN()
}

func isEmptyStr(s string) bool {
	return strip(s) == ""
}

func panicIfEmptyStr(s string) {
	if isEmptyStr(s) {
		panicMsg("String is empty")
	}
}

func stripOrPanicIfEmpty(s string) string {
	panicIfEmptyStr(s)
	return strip(s)
}

func _isNotInitSingleValue(value any) bool {
	switch v := value.(type) {
	case float64:
		return math.IsNaN(v)
	case string:
		return v == ""
	}
	panicMsg("Type is not supported")
	return false
}

func isNotInit(values ...any) bool {
	for _, value := range values {
		if _isNotInitSingleValue(value) {
			return true
		}
	}
	return false
}

func panicNotInit() {
	panicMsg("Value is not initialized")
}

func panicIfNotInit(values ...any) {
	if isNotInit(values...) {
		panicNotInit()
	}
}

func anyCmp[T Number](pairs [][]T, cmp func(val1, val2 T) bool) bool {
	for _, pair := range pairs {
		if len(pair) > 2 {
			panicMsg("Pair can only have 2 elements")
		}
		if cmp(pair[0], pair[1]) {
			return true
		}
	}
	return false
}

func anyGreater[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 > val2 })
}

func anyGreaterOrEqual[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 >= val2 })
}

func anyLess[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 < val2 })
}

func anyLessOrEqual[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 <= val2 })
}

func anyEqual[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 == val2 })
}

func anyNotEqual[T Number](pairs [][]T) bool {
	return anyCmp(pairs, func(val1, val2 T) bool { return val1 != val2 })
}

func swap[T any](value1, value2 *T) {
	*value1, *value2 = *value2, *value1
}

func abs[T Number](val T) T {
	return T(math.Abs(float64(val)))
}

func sign[T Number](val T) T {
	if val != 0 {
		val /= abs(val)
	}
	return val
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

func floatToInt64(value float64) int64 {
	return int64(math.Round(value))
}

func numberToMillis[T Number](value T) time.Duration {
	return time.Duration(float64(value)*1000) * time.Microsecond
}

func sqr[T Number](x T) T {
	return x * x
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
func isSlicesEqual[T BasicType](a, b []T) bool {
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

type ComparableFields [][2]any

func isFieldsEqual(fields ComparableFields) bool {
	for _, pair := range fields {
		if pair[0] != pair[1] {
			return false
		}
	}
	return true
}

//type SafeMap[K comparable, V any] struct {
//	mapping map[K]V
//}
//
//func (threadMap *SafeMap[K, V]) Put(key K, value V) {
//	threadMap.mapping[key] = value
//}
//
//func (threadMap *SafeMap[K, V]) CheckAndGet(key K) (V, bool) {
//	value, present := threadMap.mapping[key]
//	return value, present
//}
//
//func (threadMap *SafeMap[K, V]) Get(key K) V {
//	value, _ := threadMap.CheckAndGet(key)
//	return value
//}
//
//func (threadMap *SafeMap[K, V]) RangeOverCopy(elementHandler func(key K, value V)) {
//	copiedMap := map[K]V{}
//	err := copier.Copy(&copiedMap, &(threadMap.mapping))
//	if err != nil {
//		panic(err)
//	}
//
//	for k, v := range copiedMap {
//		elementHandler(k, v)
//	}
//}
//
//func (threadMap *SafeMap[K, V]) Pop(key K) V {
//	value := threadMap.mapping[key]
//	delete(threadMap.mapping, key)
//	return value
//}

type ThreadSafeMap[K comparable, V any] struct {
	mapping map[K]V
	mutex   sync.Mutex
}

func MakeThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	tsMap := &ThreadSafeMap[K, V]{}
	tsMap.mapping = map[K]V{}
	return tsMap
}

func (threadMap *ThreadSafeMap[K, V]) Put(key K, value V) {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	threadMap.mapping[key] = value
}

func (threadMap *ThreadSafeMap[K, V]) CheckAndGet(key K) (V, bool) {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	value, present := threadMap.mapping[key]
	return value, present
}

func (threadMap *ThreadSafeMap[K, V]) Get(key K) V {
	value, _ := threadMap.CheckAndGet(key)
	return value
}

func (threadMap *ThreadSafeMap[K, V]) RangeOverCopy(elementHandler func(key K, value V)) {
	threadMap.mutex.Lock()

	copiedMap := map[K]V{}
	err := copier.Copy(&copiedMap, &(threadMap.mapping))
	if err != nil {
		panic(err)
	}

	threadMap.mutex.Unlock()

	for k, v := range copiedMap {
		elementHandler(k, v)
	}
}

func (threadMap *ThreadSafeMap[K, V]) Pop(key K) V {
	threadMap.mutex.Lock()
	defer threadMap.mutex.Unlock()

	value := threadMap.mapping[key]
	delete(threadMap.mapping, key)
	return value
}
