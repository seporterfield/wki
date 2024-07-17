package main

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
	"single a tag": {
		input:  "<a>link</a>",
		result: "link",
	},
	"IBM wiki refs": {
		input:  "and present in over 175 countries.<ref>{{Cite web |date=June 27, 2019 |title=Trust and responsibility. Earned and practiced daily. |url=https://www.ibm.com/blogs/corporate-social-responsibility/2019/06/trust-and-responsibility-earned-and-practiced-daily/ |access-date=December 30, 2022 |website=IBM Impact |language=en-US}}</ref><ref name=\"auto\">{{cite web|website=10-K|url=https://www.sec.gov/Archives/edgar/data/51143/104746919000712/0001047469-19-000712-index.htm|title=10-K|access-date=June 1, 2019|ref={{harvid|10-K|2018}}|archive-date=December 5, 2019|archive-url=https://web.archive.org/web/20191205181213/https://www.sec.gov/Archives/edgar/data/51143/104746919000712/0001047469-19-000712-index.htm|url-status=live}}</ref> IBM is the largest industrial research",
		result: "and present in over 175 countries. IBM is the largest industrial research",
	},
}

func TestCleanWikimediaHTML(t *testing.T) {
	for name, test := range tests {
		// test := test // NOTE: uncomment for Go < 1.22, see /doc/faq#closures_and_goroutines
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := CleanWikimediaHTML(test.input), test.result; got != expected {
				t.Fatalf("CleanWikimediaHTML(%q) returned %q; expected %q", test.input, got, expected)
			}
		})
	}
}
