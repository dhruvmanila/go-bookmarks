package browser

import (
	"os/exec"
)

// defaultBrowser returns the bookmarks Collector for the default browser.
func defaultBrowser() (string, error) {
	output, err := exec.Command("xdg-mime", "query", "default", "x-scheme-handler/http").Output()
	if err != nil {
		return "", err
	}

	switch browserName := string(output); browserName {
	case "brave-browser.desktop":
		return Brave, nil
	case "google-chrome.desktop":
		return Chrome, nil
	case "firefox.desktop":
		return Firefox, nil
	default:
		return "", &ErrUnsupportedBrowser{Name: browserName}
	}
}
