package browser

import (
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

type chrome struct {
	name string
}

func newGoogleChrome() *chrome {
	return &chrome{name: Chrome}
}

// Name returns the browser name as given by the Chrome constant.
func (b *chrome) Name() string { return b.name }

func (b *chrome) bookmarkPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(
			home,
			"Library",
			"Application Support",
			"Google",
			"Chrome",
			"Default",
			"Bookmarks",
		), nil
	case "linux":
		return filepath.Join(
			home,
			".config",
			"google-chrome",
			"Default",
			"Bookmarks",
		), nil
	case "windows":
		return filepath.Join(
			home,
			"AppData",
			"Local",
			"Google",
			"Chrome",
			"User Data",
			"Default",
			"Bookmarks",
		), nil
	default:
		return "", ErrUnsupportedOS
	}
}

func (b *chrome) Bookmarks() ([]Bookmark, error) {
	path, err := b.bookmarkPath()
	if err != nil {
		return nil, err
	}
	return chromeBasedBookmarks(path)
}
