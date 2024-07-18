package main

import (
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

const DefaultWikiUrl = "wikipedia.org/wiki"
const DefaultApiUrl = "wikipedia.org/w/api.php?"

type WikipediaJSON interface {
}

type WikipediaPageQueryJSON struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

type WikipediaExtractPageJSON struct {
	Query struct {
		Pages map[string]struct {
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
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

// Removes infoboxes of any type from a Wikitext string
func removeInfobox(input string) string {
	// This regex takes into account cases like
	// {{Infobox ...
	// {{Taxobox ...
	// {{Automatic taxobox ...
	// https://en.wikipedia.org/wiki/Wikipedia:List_of_infoboxes
	re := regexp.MustCompile(`{{[a-zA-Z0-9-_]+(?:\s[a-zA-Z0-9-_]+)?box`)
	startSlice := re.FindStringIndex(input)
	if startSlice == nil {
		return input // No Infobox found
	}
	start := startSlice[0]

	bracketCount := 0
	end := start

	for i := start; i < len(input); i++ {
		if input[i] == '{' {
			bracketCount++
		} else if input[i] == '}' {
			bracketCount--
			if bracketCount == 0 {
				end = i + 1
				break
			}
		}
	}

	return input[:start] + input[end:]
}

// Character entities that wikitext doesn't include.
// Non-examples: @ and © are allowed by wikitext.
var WikiHTMLCharacterEntities = map[string]string{
	"&nbsp;": " ",
	"&quot;": `"`,
}

// Converts Wikitext into a TUI-friendly string
func CleanWikimediaHTML(dirty string) string {
	m := regexp.MustCompile(`<ref[^>]*>.*?</ref>`)
	clean := m.ReplaceAllString(dirty, "")

	clean = removeInfobox(clean)

	m = regexp.MustCompile(`(?s)\{\{(.*?)\}\}`)
	replace := func(match string) string {
		// Format based on content what's inside the {{brackets}}
		bracketContent := match[2 : len(match)-2]
		startWord, rest, _ := strings.Cut(bracketContent, " ")
		switch startWord {

		// Short description
		// On the "Fork" article: {{Short description|Eating utensil}}
		case "Short", "short":
			_, description, _ := strings.Cut(rest, "|")
			return articleDescriptionStyle(description)
		}

		// Stuff like {{... | ...}}
		startWord, rest, found := strings.Cut(bracketContent, "|")
		if !found {
			return ""
		}

		// https://new.wikipedia.org/wiki/Template:Zh
		// Generalized just in case.
		if _, ok := WikipediaLangs[startWord]; ok {
			return rest
		}

		switch startWord {
		// https://en.wikipedia.org/wiki/Template:IPA
		// Only four exceptions this time, not bad
		case "IPA", "IPAc-cmn", "IPAc-yue", "IPAc-hu", "IPAc-pl":
			return rest
		// https://en.wikipedia.org/wiki/Template:Transliteration
		case "transliteration":
			parts := strings.Split(bracketContent, "|")
			if len(parts) == 0 {
				return ""
			}
			return parts[len(parts)-1]
		}
		if len(startWord) < 4 {
			return ""
		}

		// https://en.wikipedia.org/wiki/Template:Lang
		// It's much worse than it seems
		startWord = strings.ToLower(startWord)[:4]
		switch startWord {
		case "lang":
			_, phrase, found := strings.Cut(rest, "|")
			if found {
				return phrase
			}
			parts := strings.Split(bracketContent, "|")

			// Check if there are at least two parts
			if len(parts) >= 2 {
				return parts[1]
			}
		}
		return ""
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Files and images
	lines := strings.Split(clean, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, "[[File:") {
			cleanedLines = append(cleanedLines, line)
		}
	}
	clean = strings.Join(cleanedLines, "\n")

	// Hyperlinks
	m = regexp.MustCompile(`(?s)\[\[(.*?)\]\]`)
	replace = func(match string) string {
		bracketContent := match[2 : len(match)-2]
		parts := strings.Split(bracketContent, "|")
		if len(parts) != 2 {
			return parts[0]
		}
		return linkStyle(parts[1])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Other brackets
	replace = func(match string) string {
		return linkStyle(match[2 : len(match)-2])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Strip HTML tags only after removing
	// Wikimedia/XML tags
	clean = strip.StripTags(clean)

	// Replace HTML character entities
	for entity, entityStr := range WikiHTMLCharacterEntities {
		m = regexp.MustCompile(entity)
		clean = m.ReplaceAllString(clean, entityStr)
	}

	// Bold italic
	m = regexp.MustCompile(`'''''(.*?)'''''`)
	replace = func(match string) string {
		return articleBoldedItalicStyle(match[5 : len(match)-5])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Bold
	m = regexp.MustCompile(`'''(.*?)'''`)
	replace = func(match string) string {
		return articleBoldedStyle(match[3 : len(match)-3])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Italic
	m = regexp.MustCompile(`''(.*?)''`)
	replace = func(match string) string {
		return articleItalicStyle(match[2 : len(match)-2])
	}
	clean = m.ReplaceAllStringFunc(clean, replace)

	// Anything more than three consecutive newlines is excessive
	m = regexp.MustCompile(`\n{4,}`)
	clean = m.ReplaceAllString(clean, "\n\n\n")
	return clean
}
