package browser

// Bookmark contains information regarding an individual browser bookmark.
type Bookmark struct {
	// Name is the name of a bookmark. This name is only the final element of
	// the path (the base name), not the entire path.
	// For example, Name would be "baz" and not "foo/bar/baz".
	Name string `json:"name"`

	// Path is the full path to a bookmark. In case of nested structure, where
	// there are folders within folders, this will correspond to the entire
	// path separated by a forward slash (/).
	Path string `json:"path"`

	// Url is the URL of a bookmark.
	Url string `json:"url"`
}
