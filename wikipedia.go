package main

import (
	"regexp"

	strip "github.com/grokify/html-strip-tags-go"
)

const DefaultWikiUrl = "wikipedia.org/wiki"
const DefaultApiUrl = "wikipedia.org/w/api.php?"

var WikipediaLangs = map[string]bool{
	"en": true,
	"de": true,
	"fr": true,
}

type WikipediaPageQueryJSON struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

type WikipediaPageJSON struct {
	Query struct {
		Pages []struct {
			Title     string `json:"title"`
			Revisions []struct {
				Slots struct {
					Main struct {
						Content string `json:"content"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

func CleanWikimediaHTML(dirty string) string {
	clean := strip.StripTags(dirty)
	m := regexp.MustCompile(`\[\[(.*?)\]\]`)
	replace := func(match string) string {
		return linkStyle(match[2 : len(match)-2])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)
	return clean
}
