package browser

import (
	"encoding/json"
	"os"
	"strings"
)

type chromeBasedBookmarkNode struct {
	Children []chromeBasedBookmarkNode `json:"children"`
	Name     string                    `json:"name"`
	Type     string                    `json:"type"`
	Url      string                    `json:"url"`
}

type chromeBasedBookmarkRoot struct {
	Category struct {
		BookmarkBar chromeBasedBookmarkNode `json:"bookmark_bar"`
		Synced      chromeBasedBookmarkNode `json:"synced"`
		Other       chromeBasedBookmarkNode `json:"other"`
	} `json:"roots"`
}

func chromeBasedBookmarks(path string) ([]Bookmark, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root chromeBasedBookmarkRoot
	if err = json.Unmarshal(content, &root); err != nil {
		return nil, err
	}

	var bookmarks []Bookmark
	var insertBookmark func([]string, chromeBasedBookmarkNode)

	insertBookmark = func(path []string, n chromeBasedBookmarkNode) {
		path = append(path, n.Name)
		switch n.Type {
		case "folder":
			for _, child := range n.Children {
				insertBookmark(path, child)
			}
		case "url":
			bookmarks = append(bookmarks, Bookmark{
				Name: n.Name,
				Path: strings.Join(path, "/"),
				Url:  n.Url,
			})
		}
	}

	for _, categoryNode := range []chromeBasedBookmarkNode{
		root.Category.BookmarkBar,
		root.Category.Synced,
		root.Category.Other,
	} {
		for _, child := range categoryNode.Children {
			insertBookmark(nil, child)
		}
	}

	return bookmarks, nil
}
