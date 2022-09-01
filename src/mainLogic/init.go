package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
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
	Cfg.initTouchpads()

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
	debug.SetGCPercent(Cfg.GCPercent)

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
	if c.RunFromTerminal {
		c.BaseDir = filepath.Dir(filepath.Dir(GetCurFileDir()))
	} else {
		c.BaseDir = osSpec.DefaultProjectDir
	}
	c.LayoutsDir = filepath.Join(c.BaseDir, "layouts")
}

func (c *ConfigsT) loadLayoutDir() {
	c.LayoutInUse = gofuncs.ReadFile(filepath.Join(c.LayoutsDir, "layout_to_use.txt"))
	c.LayoutInUse = strings.TrimSpace(c.LayoutInUse)

	curLayoutDir := path.Join(c.LayoutsDir, c.LayoutInUse)
	if _, err := os.Stat(curLayoutDir); os.IsNotExist(err) {
		gofuncs.Panic("Layout folder with such name doesn't exist: %s", c.LayoutInUse)
	}
}

func (c *ConfigsT) loadConfigs() {
	c.RawStrConfigs = map[string]string{}

	linesParts := c.ReadLayoutFile(path.Join(c.LayoutInUse, "configs.csv"), 0)
	for _, parts := range linesParts {
		constName := parts[0]
		constValue := parts[1]

		constName = strings.ToLower(constName)
		gofuncs.AssignWithDuplicateCheck(c.RawStrConfigs, constName, constValue)
	}
}

func GetCurFileDir() string {
	ex, err := os.Executable()
	gofuncs.CheckErr(err)
	exPath := filepath.Dir(ex)
	gofuncs.Print("Exec path: %s", exPath)
	return exPath
}

func (c *ConfigsT) ReadLayoutFile(pathFromLayoutsDir string, skipLines int) [][]string {
	file := filepath.Join(c.LayoutsDir, pathFromLayoutsDir)
	lines := gofuncs.ReadLines(file)
	lines = lines[skipLines:]

	var linesParts [][]string
	for _, line := range lines {
		line = gofuncs.Strip(line)
		if gofuncs.IsEmptyStripStr(line) || gofuncs.StartsWithAnyOf(line, ";", "//") {
			continue
		}
		parts := gofuncs.SplitByAnyOf(line, "&|>:,=")
		for ind, part := range parts {
			parts[ind] = gofuncs.Strip(part)
		}
		linesParts = append(linesParts, parts)
	}
	return linesParts
}

func (c *ConfigsT) getConfig(constName string) string {
	constName = strings.ToLower(constName)
	return gofuncs.GetOrPanic(c.RawStrConfigs, constName, "No such name in config")
}

func (c *ConfigsT) toBoolConfig(name string) bool {
	return gofuncs.StrToBool(c.getConfig(name))
}

func (c *ConfigsT) toIntConfig(name string) int {
	return gofuncs.StrToInt(c.getConfig(name))
}

func (c *ConfigsT) toMillisConfig(name string) time.Duration {
	return gofuncs.StrToMillis(c.getConfig(name))
}

func (c *ConfigsT) toFloatConfig(name string) float64 {
	return gofuncs.StrToFloat(c.getConfig(name))
}

func (c *ConfigsT) toIntToFloatConfig(name string) float64 {
	return gofuncs.StrToIntToFloat(c.getConfig(name))
}

func (c *ConfigsT) toPctConfig(name string) float64 {
	return gofuncs.StrToPct(c.getConfig(name))
}
