package pkg

import (
	"fmt"
)

type Client struct {
	Lang    string
	BaseUrl string
}

var wikipediaLangs = map[string]bool{
	"en": true,
	"de": true,
	"fr": true,
}

func NewClient(lang string, unformattedBaseUrl string) (*Client, error) {
	if _, ok := wikipediaLangs[lang]; ok {
		return nil, fmt.Errorf("wikipedia language %s does not exist", lang)
	}
	client := &Client{
		Lang:    lang,
		BaseUrl: fmt.Sprintf("https://%s.%s", lang, unformattedBaseUrl),
	}
	return client, nil
}

func (c *Client) LoadArticle(article Article) {
}
