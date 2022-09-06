package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func MakeConfigs() *ConfigsT {
	c := &ConfigsT{}

	c.setConfigConstants()
	c.InitBasePath()
	c.loadLayoutDir()
	c.loadConfigs()
	c.setConfigVars()

	return c
}

func RunFreshInitSequence() {
	Cfg = MakeConfigs()
	Cfg.initDependentOnCfg()

	initEventTypes()
	initCodeMapping()
	initTyping()
	initCommands()
}

func RunMain() {
	//run as maximum priority process

	RunFreshInitSequence()

	osSpec.InitInput()
	defer osSpec.CloseInputResources()
	defer releaseAll("")

	runtime.GC()
	debug.SetGCPercent(Cfg.System.GCPercent)

	switch Cfg.ControllerInUse {
	//go thread should always come first
	case DualShock:
		go RunWebSocket()
		RunGlobalEventsThread()
	case SteamController:
		go RunGlobalEventsThread()
		RunWebSocket()
	}
}

func (c *ConfigsT) InitBasePath() {
	BaseDir := func() string {
		if c.System.RunFromTerminal {
			return filepath.Dir(filepath.Dir(gofuncs.GetCurFileDir()))
		} else {
			return osSpec.DefaultProjectDir
		}
	}()

	c.Path.AllLayoutsDir = gofuncs.JoinPathCheckIfExists(BaseDir, "layouts")
}

func (c *ConfigsT) loadLayoutDir() {
	layoutInUse := gofuncs.ReadFile[string](filepath.Join(c.Path.AllLayoutsDir, "layout_to_use.txt"))
	layoutInUse = strings.TrimSpace(layoutInUse)

	c.Path.CurLayoutDir = gofuncs.JoinPathCheckIfExists(c.Path.AllLayoutsDir, layoutInUse)
}

func (c *ConfigsT) loadConfigs() {
	_RawCfg = &RawConfigsT{}

	gofuncs.ReadJson(_RawCfg, []string{c.Path.CurLayoutDir, "configs.json"})
}
