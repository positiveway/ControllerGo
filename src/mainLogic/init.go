package mainLogic

import (
	"ControllerGo/src/osSpec"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const RunFromTerminal = false

func GetCurFileDir() string {
	ex, err := os.Executable()
	checkErr(err)
	exPath := filepath.Dir(ex)
	print("Exec path: %s", exPath)
	return exPath
}

func InitPath() {
	if RunFromTerminal {
		BaseDir = filepath.Dir(filepath.Dir(GetCurFileDir()))
	} else {
		BaseDir = osSpec.DefaultProjectDir
	}
	LayoutsDir = filepath.Join(BaseDir, "layouts")
}

func setLayoutDir() {
	LayoutInUse = ReadFile(filepath.Join(LayoutsDir, "layout_to_use.txt"))
	LayoutInUse = strings.TrimSpace(LayoutInUse)

	curLayoutDir := path.Join(LayoutsDir, LayoutInUse)
	if _, err := os.Stat(curLayoutDir); os.IsNotExist(err) {
		panicMsg("Layout folder with such name doesn't exist: %s", LayoutInUse)
	}
}

func InitSettings() {
	InitPath()
	setLayoutDir()
	setConfigVars()
	initCodeMapping()
	initTyping()
	initCommands()
}

func RunMain() {
	InitSettings()

	osSpec.InitInput()
	defer osSpec.CloseInputResources()
	defer releaseAll()

	go RunMouseThread()
	go RunScrollThread()
	go RunReleaseHoldThread()
	go RunGameMovementThread()

	runtime.GC()
	//debug.SetGCPercent(100)

	RunWebSocket()
}
