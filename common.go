package main

import "math"

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

type Number interface {
	int64 | float64 | int | int32 | float32
}

type BasicType interface {
	Number | string | bool
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
