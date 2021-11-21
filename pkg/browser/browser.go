package browser

// These constants represent an individual browser as a string value. The
// values themselves represents the accepted value for the `-b/--browser` flag
// for the command line tool.
const (
	Brave   = "brave"
	Chrome  = "chrome"
	Firefox = "firefox"
	Safari  = "safari"
)

// GetBraveBookmarks returns a list of all the bookmarks collected from the
// Brave browser.
func GetBraveBookmarks() ([]Bookmark, error) {
	return newBraveBrowser().Bookmarks()
}

// GetChromeBookmarks returns a list of all the bookmarks collected from the
// Google Chrome browser.
func GetChromeBookmarks() ([]Bookmark, error) {
	return newGoogleChrome().Bookmarks()
}

// GetFirefoxBookmarks returns a list of all the bookmarks collected from the
// Firefox browser. The bookmarks are collected for the default user profile.
func GetFirefoxBookmarks() ([]Bookmark, error) {
	return newFirefox().Bookmarks()
}

// GetFirefoxProfileBookmarks is similar to GetFirefoxBookmarks, but this
// returns the bookmarks for the given profile name.
func GetFirefoxProfileBookmarks(profile string) ([]Bookmark, error) {
	b := newFirefox()
	b.SetProfile(profile)
	return b.Bookmarks()
}

// GetSafariBookmarks returns a list of all the bookmarks collected from the
// Safari browser.
func GetSafariBookmarks() ([]Bookmark, error) {
	return newSafari().Bookmarks()
}

// Default returns the name of the default browser. The returned string is one
// of the constants provided by this package (`Brave`, `Chrome`, `Firefox`,
// `Safari`).
//
// This is determined based on the operating system. For MacOS, the launch
// services plist file is queried. For Linux, the `xdg-mime` command is used to
// query the default handler for the HTTP protocol.
func Default() (string, error) {
	return defaultBrowser()
}
