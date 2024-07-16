package pkg

const DefaultWikiUrl = "wikipedia.org/wiki"
const DefaultApiUrl = "wikipedia.org/w/api.php?"

var WikipediaLangs = map[string]bool{
	"en": true,
	"de": true,
	"fr": true,
}

func CleanWikimediaHTML(dirty string) string {
	return dirty
}
