package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/text"
	"github.com/PuerkitoBio/goquery"
)

func generateScroll(url string) string {
	// Read page
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
	result += tagsLine

	return result
}

func usage() {
	fmt.Println("Usage: alexandria-import-usefulscience.org [URL]...")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		os.Exit(0)
	}
	for _, url := range os.Args[1:] {
		println(generateScroll(url))
	}
}
