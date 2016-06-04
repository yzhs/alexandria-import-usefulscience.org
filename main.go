package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/kr/text"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/go.uuid"
)

var homeDir string

func usage() {
	fmt.Println("Usage: alexandria-import-usefulscience.org [URL]...")
}

func generateScroll(url string) string {
	// Read page
	url = strings.TrimSpace(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// Extract data
	txt := ""
	var sources []string
	var categories []string
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		txt += s.Find("h1").Text() + "\n\n"
		s.Find("ul.sources a").Each(func(i int, s *goquery.Selection) {
			source := "[" + s.Text() + "](" + s.AttrOr("href", "") + ")"
			sources = append(sources, source)
		})
		categories = s.Find("ul.categories a").Map(func(i int, s *goquery.Selection) string {
			return s.Text()
		})
	})

	// Generate output
	result := text.Wrap(txt, 80) + "\n\n"
	for _, source := range sources {
		result += fmt.Sprintf("%% @source %s\n", source)
	}
	result += fmt.Sprintf("%% @via %s\n", url)
	tagsLine := "% science"
	for _, category := range categories {
		tagsLine += ", " + strings.ToLower(category)
	}
	result += tagsLine + "\n"

	return result
}

func handleURL(url string) {
	if !strings.HasPrefix(url, "http") {
		return
	}
	f, err := os.Create(homeDir + "/.alexandria/import/" + uuid.NewV4().String() + ".tex")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(generateScroll(url))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		usage()
		os.Exit(0)
	}
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = u.HomeDir

	if len(os.Args) == 1 {
		// Read URLs from stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		for _, url := range strings.Split(string(bytes), "\n") {
			handleURL(url)
		}
	} else {
		for _, url := range os.Args[1:] {
			handleURL(url)
		}
	}
}
