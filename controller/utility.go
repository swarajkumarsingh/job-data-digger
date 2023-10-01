package controller

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func SplitStringByDotAndInsertIntoArray(arr string) []string {
	var words []string
	for _, word := range strings.SplitN(arr, ".", 1) {
		if word == "" {
			continue
		}
		word = strings.TrimSpace(word)
		words = append(words, word)
	}
	return words
}

func GetFullUrl(e *colly.HTMLElement) string {
	return fmt.Sprintf("https://www.google.com/about/careers/applications/%s", e.ChildAttr("a.WpHeLc", "href"))
}