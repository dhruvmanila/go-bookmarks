package browser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"howett.net/plist"
)

type safariBookmarkNode struct {
	Children        []safariBookmarkNode `plist:"Children"`
	WebBookmarkType string               `plist:"WebBookmarkType"`
	Title           string               `plist:"Title"`
	URLString       string               `plist:"URLString"`
	URIDictionary   struct {
		Title string `plist:"title"`
	} `plist:"URIDictionary"`
}

type safari struct {
	name string
}

func newSafari() *safari {
	return &safari{name: Safari}
}

// Name returns the browser name as given by the Safari constant.
func (b *safari) Name() string { return b.name }

func (b *safari) bookmarkPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Safari", "Bookmarks.plist"), nil
	default:
		return "", ErrUnsupportedOS
	}
}

func (b *safari) Bookmarks() ([]Bookmark, error) {
	path, err := b.bookmarkPath()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root safariBookmarkNode
	if _, err = plist.Unmarshal(content, &root); err != nil {
		return nil, err
	}

	var bookmarks []Bookmark
	var insertBookmark func([]string, safariBookmarkNode)

	insertBookmark = func(path []string, n safariBookmarkNode) {
		switch n.WebBookmarkType {
		case "WebBookmarkTypeList":
			path = append(path, n.Title)
			for _, child := range n.Children {
				insertBookmark(path, child)
			}
		case "WebBookmarkTypeLeaf":
			path = append(path, n.URIDictionary.Title)
			bookmarks = append(bookmarks, Bookmark{
				Name: n.URIDictionary.Title,
				Path: strings.Join(path, "/"),
				Url:  n.URLString,
			})
		}
	}

	categories := map[string]bool{
		"BookmarksBar":          true,
		"BookmarksMenu":         true,
		"com.apple.ReadingList": true,
	}
	for _, categoryNode := range root.Children {
		if _, exist := categories[categoryNode.Title]; exist {
			for _, child := range categoryNode.Children {
				insertBookmark(nil, child)
			}
		}
	}

	return bookmarks, nil
}
