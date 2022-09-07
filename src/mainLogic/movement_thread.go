package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"sync"
	"time"
)

type PrecisionsT[T any] struct {
	Normal        T `json:"Normal"`
	HighPrecision T `json:"HighPrecision"`
}

type PrecisionsSpeedT PrecisionsT[float64]

func (precisionsSpeed *PrecisionsSpeedT) Validate() {
	gofuncs.PanicAnyNotPositive(precisionsSpeed.Normal, precisionsSpeed.HighPrecision)
}

type PrecisionsIntervalsT PrecisionsT[IntervalRangeT]

func (precisionsIntervals *PrecisionsIntervalsT) Validate() {
	precisionsIntervals.Normal.Validate()
	precisionsIntervals.HighPrecision.Validate()
}

type IntervalRangeT struct {
	Slowest float64 `json:"SlowestMs"`
	Fastest float64 `json:"FastestMs"`
}

func (intervalRange *IntervalRangeT) Validate() {
	gofuncs.PanicAnyNotPositive(intervalRange.Slowest, intervalRange.Fastest)

	if intervalRange.Fastest >= intervalRange.Slowest {
		gofuncs.Panic("Fastest interval can't be greater than Slowest")
	}
}

func MakeIntervalRange(slowest, fastest float64) *IntervalRangeT {
	intervalRange := &IntervalRangeT{}

	intervalRange.Slowest = slowest
	intervalRange.Fastest = fastest

	intervalRange.Validate()
	return intervalRange
}

type ModeT string

const (
	TypingMode ModeT = "Typing"
	MouseMode  ModeT = "Mouse"
	GamingMode ModeT = "Gaming"
)

type PadsSticksModeT struct {
	CurrentMode, defaultMode ModeT
	lock                     sync.Mutex
}

func MakePadsSticksMode(defaultMode ModeT) *PadsSticksModeT {
	return &PadsSticksModeT{CurrentMode: defaultMode, defaultMode: defaultMode}
}

func (mode *PadsSticksModeT) SwitchMode() {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	if mode.CurrentMode == mode.defaultMode {
		mode.CurrentMode = TypingMode
	} else {
		mode.CurrentMode = mode.defaultMode
	}
}

func (mode *PadsSticksModeT) GetMode() ModeT {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.CurrentMode
}

type HighPrecisionModeT struct {
	lock sync.Mutex

	isActive bool

	curMouseIntervals, curScrollIntervals *IntervalRangeT
	curMouseSpeed                         float64

	ctrlVirtualButton BtnOrAxisT
	ctrlCommandInfo   *CommandInfoT
}

func MakeHighPrecisionMode() *HighPrecisionModeT {
	mode := &HighPrecisionModeT{}

	CtrlCommand := []int{getCodeFromLetter("Ctrl")}
	mode.ctrlVirtualButton, mode.ctrlCommandInfo = CreateVirtualButton(CtrlCommand)

	mode.setSpeedValues()
	return mode
}

func (mode *HighPrecisionModeT) PressCtrl() {
	if mode.isActive {
		pressIfNotAlready(mode.ctrlVirtualButton, mode.ctrlCommandInfo)
	}
}

func (mode *HighPrecisionModeT) ReleaseCtrl() {
	releaseButton(mode.ctrlVirtualButton)
}

func (mode *HighPrecisionModeT) IsActive() bool {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.isActive
}

func (mode *HighPrecisionModeT) setSpeedValues() {
	if mode.isActive {
		mode.curScrollIntervals = &Cfg.Scroll.Intervals.HighPrecision
		mode.curMouseIntervals = &Cfg.Mouse.Intervals.HighPrecision
		mode.curMouseSpeed = Cfg.Mouse.Speed.HighPrecision
	} else {
		mode.ReleaseCtrl()

		mode.curScrollIntervals = &Cfg.Scroll.Intervals.Normal
		mode.curMouseIntervals = &Cfg.Mouse.Intervals.Normal
		mode.curMouseSpeed = Cfg.Mouse.Speed.Normal
	}
}

func (mode *HighPrecisionModeT) SwitchMode() {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	mode.isActive = !mode.isActive
	mode.setSpeedValues()
}

type IntervalTimerT struct {
	repeatInterval, intervalLeft, tickerInterval float64
}

func MakeIntervalTimer(repeatInterval float64) *IntervalTimerT {
	intervalTimer := &IntervalTimerT{}
	intervalTimer.InitIntervalTimer(repeatInterval)
	return intervalTimer
}

func (i *IntervalTimerT) InitIntervalTimer(repeatInterval float64) {
	i.tickerInterval = Cfg.System.TickerInterval
	gofuncs.PanicAnyNotPositive(i.tickerInterval)

	i.SetInterval(repeatInterval)

}

func (i *IntervalTimerT) SetInterval(repeatInterval float64) {
	gofuncs.PanicAnyNotPositive(repeatInterval)

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
	i.intervalLeft -= i.tickerInterval
	return i.ResetInterval()
}

type IntervalTimers2T struct {
	X, Y *IntervalTimerT
}

func MakeIntervalTimers2() *IntervalTimers2T {
	//Very Small Initial Interval That Will Be Immediately Reset
	initIntervalToReset := Cfg.Math.FloatEqualityMargin
	return &IntervalTimers2T{
		X: MakeIntervalTimer(initIntervalToReset),
		Y: MakeIntervalTimer(initIntervalToReset),
	}
}

func RunGlobalEventsThread() {
	tickerInterval := Cfg.System.TickerInterval
	ticker := time.NewTicker(gofuncs.NumberToMillis(tickerInterval))

	padsSticksMode := Cfg.PadsSticks.Mode
	highPrecisionMode := Cfg.PadsSticks.HighPrecisionMode

	mousePadStick := Cfg.PadsSticks.MousePS
	mousePosition := mousePadStick.transformedPos
	moveMouseInInterval := GetMoveInInterval(mousePadStick, mousePosition,
		GetMouseMoveFunc(), nil)

	scrollPadStick := Cfg.PadsSticks.ScrollPS
	scrollPosition := scrollPadStick.transformedPos
	moveScrollInInterval := GetMoveInInterval(scrollPadStick, scrollPosition,
		GetScrollMoveFunc(), GetScrollFilterFunc())

	for range ticker.C {
		switch padsSticksMode.CurrentMode {
		case MouseMode:
			moveScrollInInterval(highPrecisionMode.curScrollIntervals)

			if scrollPosition.x == 0 && scrollPosition.y == 0 {
				highPrecisionMode.ReleaseCtrl()
			}

			fallthrough
		case GamingMode:
			switch Cfg.ControllerInUse {
			case DualShock:
				moveMouseInInterval(highPrecisionMode.curMouseIntervals)
			}
		}

		//should be placed last to not interfere with GetMode
		//and MoveMouse has higher priority
		RepeatCommand()
	}
}
