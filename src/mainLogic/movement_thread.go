package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"sync"
	"time"
)

type IntervalRangeT struct {
	slowest, fastest float64
}

func (intervalRange *IntervalRangeT) checkIntervals() {
	if intervalRange.fastest >= intervalRange.slowest {
		gofuncs.Panic("Fastest interval can't be greater than slowest")
	}
}

func MakeIntervalRange(slowest, fastest float64) *IntervalRangeT {
	intervalRange := &IntervalRangeT{}
	intervalRange.slowest, intervalRange.fastest = slowest, fastest
	intervalRange.checkIntervals()
	return intervalRange
}

type ModeT string

const (
	TypingMode ModeT = "Typing"
	MouseMode  ModeT = "Mouse"
	GamingMode ModeT = "Gaming"
)

type PadsSticksModeT struct {
	currentMode, defaultMode ModeT
	lock                     sync.Mutex
}

func MakePadsSticksMode(defaultMode ModeT) *PadsSticksModeT {
	return &PadsSticksModeT{currentMode: defaultMode, defaultMode: defaultMode}
}

func (mode *PadsSticksModeT) SwitchMode() {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	if mode.currentMode == mode.defaultMode {
		mode.currentMode = TypingMode
	} else {
		mode.currentMode = mode.defaultMode
	}
}

func (mode *PadsSticksModeT) GetMode() ModeT {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.currentMode
}

type HighPrecisionModeT struct {
	isActive bool
	lock     sync.Mutex
}

func MakeHighPrecisionMode() *HighPrecisionModeT {
	return &HighPrecisionModeT{}
}

func (mode *HighPrecisionModeT) SwitchMode() bool {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.isActive
}

type IntervalTimerT struct {
	repeatInterval, intervalLeft float64
}

func MakeIntervalTimer(repeatInterval float64) *IntervalTimerT {
	intervalTimer := &IntervalTimerT{}
	intervalTimer.SetInterval(repeatInterval)
	return intervalTimer
}

func (i *IntervalTimerT) SetInterval(repeatInterval float64) {
	gofuncs.PanicAnyNotInitOrEmpty(repeatInterval)

	i.repeatInterval = repeatInterval
	i.intervalLeft = repeatInterval
}

func (i *IntervalTimerT) ResetInterval() bool {
	if i.intervalLeft <= 0 {
		i.intervalLeft = i.repeatInterval
		return true
	}
	return false
}

func (i *IntervalTimerT) DecreaseInterval() bool {
	i.intervalLeft -= Cfg.tickerInterval
	return i.ResetInterval()
}

type IntervalTimers2T struct {
	X, Y *IntervalTimerT
}

func MakeIntervalTimers2() *IntervalTimers2T {
	return &IntervalTimers2T{
		X: MakeIntervalTimer(0),
		Y: MakeIntervalTimer(0),
	}
}

func RunGlobalEventsThread() {
	ticker := time.NewTicker(gofuncs.NumberToMillis(Cfg.tickerInterval))

	mouseIntervals := MakeIntervalTimers2()
	scrollIntervals := MakeIntervalTimers2()

	for range ticker.C {
		switch Cfg.PadsSticksMode.GetMode() {
		case MouseMode:
			MoveScroll(scrollIntervals)
			fallthrough
		case GamingMode:
			switch Cfg.ControllerInUse {
			case DualShock:
				MoveMouse(mouseIntervals)
			}
		}

		//should be placed last to not interfere with GetMode
		//and MoveMouse has higher priority
		RepeatCommand()
	}
}
