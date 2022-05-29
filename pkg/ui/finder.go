package ui

import (
	"fmt"
	"strings"

	"github.com/dhruvmanila/go-bookmarks/pkg/browser"
	"github.com/ktr0731/go-fuzzyfinder"
)

func FindBookmarksMulti(browserName string, bookmarks []browser.Bookmark) ([]int, error) {
	return fuzzyfinder.FindMulti(bookmarks, func(i int) string {
		// TODO: format each entry by dividing the available space into two
		// parts where the first one will contain the path and the other will
		// contain the URL. The URL part will be greyed and italics.
		return bookmarks[i].Path
	}, fuzzyfinder.WithPromptString(fmt.Sprintf("%s Bookmarks > ", strings.Title(browserName))))
}
