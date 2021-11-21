package browser

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
	ini "gopkg.in/ini.v1"
)

type firefoxProfile struct {
	Name       string `ini:"Name"`
	IsRelative bool   `ini:"IsRelative"`
	Path       string `ini:"Path"`
	Default    bool   `ini:"Default"`
}

type firefox struct {
	name    string
	profile string
}

func newFirefox() *firefox {
	return &firefox{name: Firefox}
}

// Name returns the browser name as given by the Firefox constant.
func (b *firefox) Name() string { return b.name }

// SetProfile is used to set the firefox profile name for which to retrieve
// the bookmarks.
func (b *firefox) SetProfile(profile string) {
	b.profile = profile
}

// Bookmarks returns a list of Bookmark either for the user provided profile
// or the default profile.
func (b *firefox) Bookmarks() ([]Bookmark, error) {
	profileDir, err := b.getProfileDir()
	if err != nil {
		return nil, err
	}

	// Firefox bookmarks are stored in a sqlite3 database.
	placesDb := filepath.Join(profileDir, "places.sqlite")
	db, err := sql.Open("sqlite3", placesDb)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// fk¹ is the item's place id in moz_places table.
	// type² is the item's type constant (1 -> bookmark, 2 -> folder, 3 -> separator).
	//
	// ¹https://hg.mozilla.org/mozilla-central/file/tip/toolkit/components/places/nsPlacesTables.h#l93
	// ²https://hg.mozilla.org/mozilla-central/file/tip/toolkit/components/places/Bookmarks.jsm#l107
	rows, err := db.Query("SELECT fk,parent,title FROM moz_bookmarks WHERE type=1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []Bookmark
	for rows.Next() {
		// Step 1: Extract the bookmark row.
		var placeId, parentId int
		var title string
		if err = rows.Scan(&placeId, &parentId, &title); err != nil {
			return nil, err
		}

		// Step 2: Extract the URL for the bookmark.
		var url string
		row := db.QueryRow("SELECT url FROM moz_places WHERE id=?", placeId)
		if err = row.Scan(&url); err != nil {
			return nil, err
		}

		// Step 3: Extract the full path to the bookmark.
		var path []string
		for {
			var folder string
			row := db.QueryRow("SELECT title,parent FROM moz_bookmarks WHERE id=?", parentId)
			if err = row.Scan(&folder, &parentId); err != nil {
				// If there are no rows, then we've reached the root.
				if errors.Is(err, sql.ErrNoRows) {
					break
				}
				return nil, err
			}
			// Prepend the value as we're travelling up.
			path = append([]string{folder}, path...)
		}

		bookmarks = append(bookmarks, Bookmark{
			Name: title,
			Path: strings.Join(path, "/") + "/" + title,
			Url:  url,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookmarks, nil
}

// getProfileDir returns the absolute path to the firefox profile directory.
// The profile name will either be the one provided by the user, which is set
// using the `SetProfile` method or the default one.
func (b *firefox) getProfileDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	var defaultDir string
	switch runtime.GOOS {
	case "darwin":
		defaultDir = filepath.Join(home, "Library", "Application Support", "Firefox")
	case "linux":
		defaultDir = filepath.Join(home, ".mozilla", "firefox")
	case "windows":
		defaultDir = filepath.Join(home, "AppData", "Roaming", "Mozilla", "Firefox")
	default:
		return "", ErrUnsupportedOS
	}

	cfg, err := ini.Load(filepath.Join(defaultDir, "profiles.ini"))
	if err != nil {
		return "", err
	}

	var profileDir string
	for _, section := range cfg.Sections() {
		if strings.HasPrefix(section.Name(), "Profile") {
			// We need a new object for every profile.
			fp := new(firefoxProfile)
			if err = section.StrictMapTo(fp); err != nil {
				return "", err
			}
			// First, check if the user provided profile name match and then
			// fallback to using the default profile.
			if b.profile == fp.Name || (fp.Default && b.profile == "") {
				if fp.IsRelative {
					profileDir = filepath.Join(defaultDir, fp.Path)
				} else {
					profileDir = fp.Path
				}
				break
			}
		}
	}
	if profileDir == "" {
		if b.profile == "" {
			return "", errors.New("unable to find the default profile name")
		}
		return "", fmt.Errorf("unable to find the profile name: %q", b.profile)
	}
	return profileDir, nil
}
