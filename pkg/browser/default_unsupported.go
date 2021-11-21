//go:build !darwin && !linux && !windows

package browser

import (
	"fmt"
	"runtime"
)

func defaultBrowser() (string, error) {
	return "", ErrUnsupportedOS
}
