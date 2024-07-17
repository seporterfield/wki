package main

import (
	"fmt"
	"testing"
)

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
	"File": {
		input: `

[[File:IBM360-67AtUmichWithMikeAlexander.jpg|thumb|right|An [[IBM System/360]] in use at the [[University of Michigan]] {{Circa|1969}}]]
[[File:Saturn_IB_and_V_Instrument_Unit.jpg|thumb|IBM guidance computer hardware for the [[Saturn V Instrument Unit]]]]
`,
		result: "\n\n",
	},
	"Brackets": {
		input:  "here are [[brackets]] wow",
		result: fmt.Sprintf("here are %s wow", linkStyle("brackets")),
	},
	"double curly brace with newlines": {
		input:  "{{Multiple |\n things\n}}",
		result: "",
	},
	"Infobox": {
		input: `{{Infobox company
| fetchwikidata = no
| name = OpenAI, Inc.
| logo = OpenAI Logo.svg
| logo_size = 250px
| image = 
| image_size = 
| image_caption = 
| type = [[Privately held company|Private]]
| industry = [[Artificial intelligence]]
| founders = <!-- listing any or all the "founders" here lacks context. Please discuss before adding particular individuals here. -->
| founded = {{Start date and age|2015|12|11}}
| hq_location = [[San Francisco]], [[California]] U.S.<ref>{{Cite web |date=December 20, 2022 |title=I Tried To Visit OpenAI's Office. Hilarity Ensued |url=https://sfstandard.com/technology/i-tried-to-visit-openais-office-hilarity-ensued/ |access-date=June 3, 2023 |website=The San Francisco Standard |language=en-US |archive-date=June 3, 2023 |archive-url=https://web.archive.org/web/20230603194312/https://sfstandard.com/technology/i-tried-to-visit-openais-office-hilarity-ensued/ |url-status=live }}</ref>
| key_people = {{Unbulleted list
| [[Bret Taylor]] ([[chairman]])
| [[Sam Altman]] ([[Chief executive officer|CEO]])
| [[Greg Brockman]] ([[President (corporate title)|president]])
| [[Mira Murati]] ([[Chief technology officer|CTO]])
}}
| area_served = 
| products = [[OpenAI Five]]
{{flatlist|
* [[GPT-1]]
* [[GPT-2|2]]
* [[GPT-3|3]]
* [[GPT-4|4]]
* [[GPT-4o|4o]]
}}
{{Unbulleted list
| [[DALL·E]]
| [[OpenAI Codex]]
| [[ChatGPT]]
| [[Sora (text-to-video model)|Sora]]
}}
| services = 
| revenue = {{increase}} [[US$]]28 million<ref name=2022-fin>{{cite web |last1=Woo |first1=Erin |last2=Efrati |first2=Amir |date=May 4, 2023 |title=OpenAI's Losses Doubled to $540 Million as It Developed ChatGPT |url=https://www.theinformation.com/articles/openais-losses-doubled-to-540-million-as-it-developed-chatgpt |website=[[The Information (website)|The Information]] |url-access=subscription |quote=In 2022, by comparison, revenue was just $28 million, mainly from selling access to its AI software... OpenAI's losses roughly doubled to around $540 million last year as it developed ChatGPT... |access-date=June 19, 2023 |archive-date=June 19, 2023 |archive-url=https://web.archive.org/web/20230619191257/https://www.theinformation.com/articles/openais-losses-doubled-to-540-million-as-it-developed-chatgpt |url-status=live }}</ref>
| revenue_year = 2022
| net_income = {{decrease}} US${{color|red|&minus;540}} million<ref name=2022-fin />
| net_income_year = 2022
| equity = 
| equity_year = 
| num_employees = {{circa|1,200}} (2024)<ref>{{cite web|title=OpenAI Sees 'Tremendous Growth' in Corporate Version of ChatGPT|url=https://www.bloomberg.com/news/articles/2024-04-04/openai-sees-tremendous-growth-in-corporate-version-of-chatgpt|website=Bloomberg|language=en|date=April 4, 2024|archive-date=April 4, 2024|archive-url=https://archive.today/20240404174825/https://www.bloomberg.com/news/articles/2024-04-04/openai-sees-tremendous-growth-in-corporate-version-of-chatgpt|url-status=live}}</ref>
| subsid = 
| homepage = {{URL|https://openai.com/}}
| footnotes = 
}}OpenAI is an`,
		result: "OpenAI is an",
	},
}

func TestCleanWikimediaHTML(t *testing.T) {
	for name, test := range tests {
		// test := test // NOTE: uncomment for Go < 1.22, see /doc/faq#closures_and_goroutines
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := CleanWikimediaHTML(test.input), test.result; got != expected {
				t.Fatalf("function CleanWikimediaHTML\n---INPUT\n%q\n---GOT\n%q\n---EXPECTED\n%q\n---", shorten(test.input), shorten(got), expected)
			}
		})
	}
}

func shorten(s string) string {
	if len(s) > 500 {
		return s[:500] + "...(CUT OFF)"
	}
	return s
}
