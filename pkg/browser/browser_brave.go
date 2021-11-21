package browser

import (
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

type braveBrowser struct {
	name string
}

func newBraveBrowser() *braveBrowser {
	return &braveBrowser{name: Brave}
}

// Name returns the browser name as given by the Brave constant.
func (b *braveBrowser) Name() string { return b.name }

func (b *braveBrowser) bookmarkPath() (string, error) {
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
			"BraveSoftware",
			"Brave-Browser",
			"Default",
			"Bookmarks",
		), nil
	case "linux":
		return filepath.Join(
			home,
			".config",
			"BraveSoftware",
			"Brave-Browser",
			"Default",
			"Bookmarks",
		), nil
	case "windows":
		return filepath.Join(
			home,
			"AppData",
			"Local",
			"BraveSoftware",
			"Brave-Browser",
			"User Data",
			"Default",
			"Bookmarks",
		), nil
	default:
		return "", ErrUnsupportedOS
	}
}

func (b *braveBrowser) Bookmarks() ([]Bookmark, error) {
	path, err := b.bookmarkPath()
	if err != nil {
		return nil, err
	}
	return chromeBasedBookmarks(path)
}
