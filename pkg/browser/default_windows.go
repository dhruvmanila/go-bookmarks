package browser

import (
	"errors"
)

// defaultBrowser returns the bookmarks Collector for the default browser.
func defaultBrowser() (string, error) {
	return "", errors.New("not implemented")
}
