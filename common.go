package main

import (
	"fmt"
	"math"
)

func panicMisspelled(str string) {
	panic(fmt.Sprintf("Probably misspelled: %s\n", str))
}

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

func getOrDefault[k comparable, v any](m map[k]v, key k, defaultVal v) v {
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
	Number | string | bool
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
