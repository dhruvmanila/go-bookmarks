package browser

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"howett.net/plist"
)

// launchServicesPlist is the absolute path to the launch services plist file
// on MacOS.
const launchServicesPlist = "~/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

// lsPreferences is the structure of the launch services plist file. This
// contains only the required fields.
type lsPreferences struct {
	LSHandlers []struct {
		LSHandlerRoleAll   string `plist:"LSHandlerRoleAll"`
		LSHandlerURLScheme string `plist:"LSHandlerURLScheme"`
	} `plist:"LSHandlers"`
}

// defaultBrowser returns the name of the default browser on MacOS.
func defaultBrowser() (string, error) {
	path, err := homedir.Expand(launchServicesPlist)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var prefs lsPreferences
	if _, err = plist.Unmarshal(content, &prefs); err != nil {
		return "", err
	}

	var bundleId string
	for _, handler := range prefs.LSHandlers {
		if handler.LSHandlerURLScheme == "http" {
			bundleId = handler.LSHandlerRoleAll
			break
		}
	}
	if bundleId == "" {
		return "", fmt.Errorf("URL scheme handler for %q absent from plist file: %q", "http", path)
	}

	switch bundleId {
	case "com.brave.browser":
		return Brave, nil
	case "com.google.chrome":
		return Chrome, nil
	case "org.mozilla.firefox":
		return Firefox, nil
	case "com.apple.safari":
		return Safari, nil
	default:
		return "", &ErrUnsupportedBrowser{Name: bundleId}
	}
}
