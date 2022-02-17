//go:build !windows

package osSpecific

func RunOsLogic() {
	GetLocaleExecPath = GetCurFileDir() + "/getLocale.sh"
}
