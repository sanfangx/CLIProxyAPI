package tui

import (
	"os/exec"
	"runtime"
)

// openBrowser opens the specified URL in the user's default browser.
func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin", "ios":
		if _, err := exec.LookPath("uiopen"); err == nil {
			return exec.Command("uiopen", url).Start()
		}
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
