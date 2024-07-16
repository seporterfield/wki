package pkg

type Article struct {
	Title       string
	Description string
	Content     string
	Url         string
}

var DefaultArticleMap = map[int]Article{
	0: {Title: "...", Description: "type something!", Content: "", Url: ""},
}
