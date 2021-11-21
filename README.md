# Go Bookmarks

Go library and command-line tool to list your browser bookmarks.

## Installation

```bash
go install github.com/dhruvmanila/go-bookmarks/cmd/bookmarks@latest
```

### Usage

```
Usage: bookmarks [options]

Bookmarks lists your browser bookmarks.

If no browser name is provided, Bookmarks will collect the bookmarks
for the default browser.

Options:
  -b, --browser <name>  Browser to collect the bookmarks for
                        Available: 'brave', 'chrome', 'firefox', 'safari'
  -j, --json            Output in JSON format
  -p, --profile <name>  Firefox profile name
  -h, --help            Show this help message
  -v, --version         Show version
```

## Library

[![Go Reference](https://pkg.go.dev/badge/github.com/dhruvmanila/go-bookmarks/pkg/browser.svg)](https://pkg.go.dev/github.com/dhruvmanila/go-bookmarks/pkg/browser)

Bookmarks also exposes a library with a minimal API. See
[documentation](https://pkg.go.dev/github.com/dhruvmanila/go-bookmarks/pkg/browser) for
more details.

## License

This project is licensed under the MIT License.

See [LICENSE](./LICENSE) for details.
