package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func MakeConfigs() (*ConfigsT, *RawConfigsT) {
	cfg := &ConfigsT{}

	cfg.setConfigConstants()
	cfg.InitBasePath()
	cfg.loadLayoutDir()
	rawCfg := cfg.loadConfigs()
	cfg.setConfigVars(rawCfg)

	return cfg, rawCfg
}

func RunFreshInitSequence() *DependentVariablesT {
	cfg, rawCfg := MakeConfigs()

	initEventTypes(cfg)

	return MakeDependentVariables(rawCfg, cfg)
}

func RunMain() {
	//run as maximum priority process
	dependentVars := RunFreshInitSequence()

	osSpec.InitInput()
	defer osSpec.CloseInputResources()
	defer dependentVars.Buttons.releaseAll("")

	runtime.GC()
	debug.SetGCPercent(dependentVars.cfg.System.GCPercent)

	switch dependentVars.cfg.ControllerInUse {
	//go thread should always come first
	case DualShock:
		go dependentVars.RunWebSocket()
		dependentVars.RunGlobalEventsThread()
	case SteamController:
		go dependentVars.RunGlobalEventsThread()
		dependentVars.RunWebSocket()
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

func (c *ConfigsT) loadConfigs() *RawConfigsT {
	rawCfg := &RawConfigsT{}

	gofuncs.ReadJson(rawCfg, []string{c.Path.CurLayoutDir, "configs.json"})
	return rawCfg
}
