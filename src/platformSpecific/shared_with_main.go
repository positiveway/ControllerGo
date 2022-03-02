package platformSpecific

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const LeftMouse = -3
const RightMouse = -4
const MiddleMouse = -5

const EventServerNotRunningMsg = "Event server is not running"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sPrint(message string, variables ...any) string {
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	return fmt.Sprintf(message, variables...)
}

func print(message string, variables ...any) {
	fmt.Print(sPrint(message, variables...))
}

func GetCurFileDir() string {
	ex, err := os.Executable()
	CheckErr(err)
	exPath := filepath.Dir(ex)
	print("Exec path: %s", exPath)
	return exPath
}
