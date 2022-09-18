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

type LockStruct struct {
	_lock        sync.Mutex
	Lock, Unlock func()
}

func (base *LockStruct) Init() {
	base.Lock = base._lock.Lock
	base.Unlock = base._lock.Unlock
}

type CfgStruct struct {
	cfg *ConfigsT
}

func (base *CfgStruct) Init(cfg *ConfigsT) {
	base.cfg = cfg
}

type ButtonsStruct struct {
	buttons *ButtonsT
}

func (base *ButtonsStruct) Init(buttons *ButtonsT) {
	base.buttons = buttons
}

type CfgButtonsStruct struct {
	CfgStruct
	ButtonsStruct
}

func (base *CfgButtonsStruct) Init(cfg *ConfigsT, buttons *ButtonsT) {
	base.CfgStruct.Init(cfg)
	base.ButtonsStruct.Init(buttons)
}

type CfgButtonsLockStruct struct {
	LockStruct
	CfgButtonsStruct
}

func (base *CfgButtonsLockStruct) Init(cfg *ConfigsT, buttons *ButtonsT) {
	base.LockStruct.Init()
	base.CfgButtonsStruct.Init(cfg, buttons)
}

type CfgLockStruct struct {
	LockStruct
	CfgStruct
}

func (base *CfgLockStruct) Init(cfg *ConfigsT) {
	base.LockStruct.Init()
	base.CfgStruct.Init(cfg)
}

type PadsSticksModeT struct {
	CurrentMode, defaultMode ModeT
	LockStruct
}

func MakePadsSticksMode(defaultMode ModeT) *PadsSticksModeT {
	mode := &PadsSticksModeT{CurrentMode: defaultMode, defaultMode: defaultMode}
	mode.Init()
	return mode
}

func (mode *PadsSticksModeT) SetToDefault() {
	mode.Lock()
	defer mode.Unlock()

	mode.CurrentMode = mode.defaultMode
}

func (mode *PadsSticksModeT) SwitchMode() {
	mode.Lock()
	defer mode.Unlock()

	if mode.CurrentMode == mode.defaultMode {
		mode.CurrentMode = TypingMode
	} else {
		mode.CurrentMode = mode.defaultMode
	}
}

func (mode *PadsSticksModeT) GetMode() ModeT {
	mode.Lock()
	defer mode.Unlock()

	return mode.CurrentMode
}

type HighPrecisionModeT struct {
	CfgButtonsLockStruct

	isActive bool

	curMouseIntervals, curScrollIntervals *IntervalRangeT
	curMouseSpeed                         float64

	ctrlVirtualButton BtnOrAxisT
	ctrlCommandInfo   *CommandInfoT

	setSpeedValues func()
}

func (mode *HighPrecisionModeT) Init(cfg *ConfigsT, buttons *ButtonsT) {
	mode.CfgButtonsLockStruct.Init(cfg, buttons)

	mode.setSpeedValues = mode.GetSetSpeedValuesFunc()

	CtrlCommand := CommandT{buttons.getCodeFromLetter("Ctrl")}
	mode.ctrlVirtualButton, mode.ctrlCommandInfo = buttons.CreateVirtualButton(CtrlCommand)

	mode.setSpeedValues()
}

func (mode *HighPrecisionModeT) PressCtrl() {
	if mode.isActive {
		mode.buttons.pressIfNotAlready(
			mode.ctrlVirtualButton,
			mode.ctrlCommandInfo)
	}
}

func (mode *HighPrecisionModeT) ReleaseCtrl() {
	mode.buttons.releaseButton(mode.ctrlVirtualButton)
}

func (mode *HighPrecisionModeT) IsActive() bool {
	mode.Lock()
	defer mode.Unlock()

	return mode.isActive
}

func (mode *HighPrecisionModeT) GetSetSpeedValuesFunc() func() {
	scrollIntervals := mode.cfg.Scroll.Intervals
	mouseIntervals := mode.cfg.Mouse.Intervals
	mouseSpeed := mode.cfg.Mouse.Speed

	return func() {
		if mode.isActive {
			mode.curScrollIntervals = &scrollIntervals.HighPrecision
			mode.curMouseIntervals = &mouseIntervals.HighPrecision
			mode.curMouseSpeed = mouseSpeed.HighPrecision
		} else {
			mode.ReleaseCtrl()

			mode.curScrollIntervals = &scrollIntervals.Normal
			mode.curMouseIntervals = &mouseIntervals.Normal
			mode.curMouseSpeed = mouseSpeed.Normal
		}
	}
}

func (mode *HighPrecisionModeT) Disable() {
	mode.Lock()
	defer mode.Unlock()

	mode.isActive = false
	mode.setSpeedValues()
}

func (mode *HighPrecisionModeT) SwitchMode() {
	if mode.cfg.PadsSticks.Mode.CurrentMode == TypingMode {
		return
	}

	mode.Lock()
	defer mode.Unlock()

	mode.isActive = !mode.isActive
	mode.setSpeedValues()
}

type RepeatedTimerT struct {
	repeatInterval, intervalLeft, tickerInterval float64
}

func MakeIntervalTimer(repeatInterval float64, cfg *ConfigsT) *RepeatedTimerT {
	intervalTimer := &RepeatedTimerT{}
	intervalTimer.InitIntervalTimer(repeatInterval, cfg)
	return intervalTimer
}

func (t *RepeatedTimerT) InitIntervalTimer(repeatInterval float64, cfg *ConfigsT) {
	t.tickerInterval = cfg.System.TickerInterval
	gofuncs.PanicAnyNotPositive(t.tickerInterval)

	t.SetInterval(repeatInterval)

}

func (t *RepeatedTimerT) SetInterval(repeatInterval float64) {
	gofuncs.PanicAnyNotPositive(repeatInterval)

	t.repeatInterval = repeatInterval
	t.intervalLeft = repeatInterval
}

func (t *RepeatedTimerT) ResetInterval() bool {
	if t.intervalLeft <= 0 {
		t.intervalLeft = t.repeatInterval
		return true
	}
	return false
}

func (t *RepeatedTimerT) DecreaseInterval() bool {
	t.intervalLeft -= t.tickerInterval
	return t.ResetInterval()
}

type IntervalTimers2T struct {
	X, Y *RepeatedTimerT
}

func MakeIntervalTimers2(cfg *ConfigsT) *IntervalTimers2T {
	//Very Small Initial Interval That Will Be Immediately Reset
	initIntervalToReset := cfg.Math.FloatEqualityMargin
	return &IntervalTimers2T{
		X: MakeIntervalTimer(initIntervalToReset, cfg),
		Y: MakeIntervalTimer(initIntervalToReset, cfg),
	}
}

func (dependentVars *DependentVariablesT) RunGlobalEventsThread() {
	cfg := dependentVars.cfg
	buttons := dependentVars.Buttons

	tickerInterval := cfg.System.TickerInterval
	ticker := time.NewTicker(gofuncs.NumberToMillis(tickerInterval))

	padsSticksMode := cfg.PadsSticks.Mode
	highPrecisionMode := dependentVars.HighPrecisionMode
	controllerInUse := cfg.ControllerInUse

	mousePadStick := dependentVars.MousePS
	mousePosition := mousePadStick.transformedPos

	moveMouseInInterval := GetMoveInInterval(cfg, mousePadStick, mousePosition,
		GetMouseMoveFunc(), nil)

	scrollPadStick := dependentVars.ScrollPS
	scrollPosition := scrollPadStick.transformedPos
	moveScrollInInterval := GetMoveInInterval(cfg, scrollPadStick, scrollPosition,
		dependentVars.GetScrollMoveFunc(), dependentVars.GetScrollFilterFunc())

	for range ticker.C {
		switch padsSticksMode.CurrentMode {
		case MouseMode:
			moveScrollInInterval(highPrecisionMode.curScrollIntervals)

			if scrollPosition.x == 0 && scrollPosition.y == 0 {
				highPrecisionMode.ReleaseCtrl()
			}

			fallthrough
		case GamingMode:
			switch controllerInUse {
			case DualShock:
				moveMouseInInterval(highPrecisionMode.curMouseIntervals)
			}
		}

		//should be placed last to not interfere with GetMode
		//and MoveMouse has higher priority
		buttons.RepeatCommand()
	}
}
