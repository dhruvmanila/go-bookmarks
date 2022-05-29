package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dhruvmanila/go-bookmarks/pkg/browser"
	"github.com/dhruvmanila/go-bookmarks/pkg/ui"
	webbrowser "github.com/pkg/browser"
)

var (
	Version   = "dev"
	GitCommit = ""
)

// Command-line options
var (
	browserName        string
	jsonOutput         bool
	firefoxProfileName string
	showHelp           bool
	showVersion        bool
)

func init() {
	flag.StringVar(&browserName, "b", "", "")
	flag.StringVar(&browserName, "browser", "", "")
	flag.BoolVar(&jsonOutput, "j", false, "")
	flag.BoolVar(&jsonOutput, "json", false, "")
	flag.StringVar(&firefoxProfileName, "p", "", "")
	flag.StringVar(&firefoxProfileName, "profile", "", "")
	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
}

func main() {
	log.SetPrefix("bookmarks: ")
	log.SetFlags(0)

	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		return
	case showVersion:
		fmt.Printf("bookmarks %s (%s)\n", Version, GitCommit)
		return
	}

	var err error
	if browserName == "" {
		browserName, err = browser.Default()
		if err != nil {
			log.Fatal(fmt.Errorf("error: unable to get the default browser: %w", err))
		}
	}

	if browserName != browser.Firefox && firefoxProfileName != "" {
		log.Fatal("error: -p/--profile can only be used with firefox browser")
	}

	var bookmarks []browser.Bookmark
	switch browserName {
	case browser.Brave:
		bookmarks, err = browser.GetBraveBookmarks()
	case browser.Chrome:
		bookmarks, err = browser.GetChromeBookmarks()
	case browser.Firefox:
		if firefoxProfileName == "" {
			bookmarks, err = browser.GetFirefoxBookmarks()
		} else {
			bookmarks, err = browser.GetFirefoxProfileBookmarks(firefoxProfileName)
		}
	case browser.Safari:
		bookmarks, err = browser.GetSafariBookmarks()
	default:
		log.Fatalf("unsupported browser: %q", browserName)
	}

	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", browserName, err))
	}

	switch {
	case jsonOutput:
		output, err := json.MarshalIndent(bookmarks, "", "  ")
		if err != nil {
			log.Fatal(fmt.Errorf("%s: unable to get the JSON output: %w", browserName, err))
		}
		if _, err = os.Stdout.Write(output); err != nil {
			log.Fatal(err)
		}
	default:
		choices, err := ui.FindBookmarksMulti(browserName, bookmarks)
		if err != nil {
			log.Fatal(err)
		}

		for _, idx := range choices {
			if err := openURL(bookmarks[idx].Url); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// openURL opens the given URL in the default browser.
func openURL(url string) error {
	if err := webbrowser.OpenURL(url); err != nil {
		return fmt.Errorf("failed to open the URL (%s): %w", url, err)
	}
	return nil
}

const usage = `Usage: bookmarks [options]

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
`
