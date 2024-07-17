package pkg

import "testing"

var tests = map[string]struct {
	input  string
	result string
}{
	"empty string": {
		input:  "",
		result: "",
	},
	"one character": {
		input:  "x",
		result: "x",
	},
	"one multi byte glyph": {
		input:  "🎉",
		result: "🎉",
	},
	"string with multiple multi-byte glyphs": {
		input:  "🥳🎉🐶",
		result: "🥳🎉🐶",
	},
	"string with double brackets": {
		input:  "[[x]]",
		result: "x",
	},
	"single a tag": {
		input:  "<a>link</a>",
		result: "link",
	},
	"IBM wiki": {
		input:  "IBM was founded in 1911 as the [[Computing-Tabulating-Recording Company]] (CTR), a [[holding company]] of manufacturers of record-keeping and measuring systems.",
		result: "IBM was founded in 1911 as the Computing-Tabulating-Recording Company (CTR), a holding company of manufacturers of record-keeping and measuring systems.",
	},
}

func TestCleanWikimediaHTML(t *testing.T) {
	for name, test := range tests {
		// test := test // NOTE: uncomment for Go < 1.22, see /doc/faq#closures_and_goroutines
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := CleanWikimediaHTML(test.input), test.result; got != expected {
				t.Fatalf("reverse(%q) returned %q; expected %q", test.input, got, expected)
			}
		})
	}
}
