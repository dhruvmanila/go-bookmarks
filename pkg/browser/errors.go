package browser

import (
	"fmt"
	"runtime"
)

// ErrUnsupportedOS is returned when the current operating system is not
// supported.
var ErrUnsupportedOS = fmt.Errorf("unsupported operating system: %q", runtime.GOOS)

// ErrUnsupportedBrowser is returned when user provided browser or the default
// browser name is not supported.
type ErrUnsupportedBrowser struct {
	// Name is the browser name.
	Name string
}

func (e *ErrUnsupportedBrowser) Error() string {
	return fmt.Sprintf("unsupported browser: %q", e.Name)
}
