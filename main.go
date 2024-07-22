package main

// A Wikipedia TUI

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ExtendedUsage = `
wki - Wikipedia at your fingertips

Type into the search bar to search for articles.
- Move back and forth:         left and right arrow keys
- Move cursor ` + "`*`" + `:             up and down arrow keys
- Open the selected article:   enter
- Navigate the article reader: arrow keys or vim/less controls
- Return to search page:       left arrow key
- Quit:                        escape or Ctrl+C`

// Helper struct enabling multiple TUI pages
// along with the pages map and model.pageName
type Page struct {
	update func(model, tea.Msg) (tea.Model, tea.Cmd)
	view   func(model) string
}

// New Update/View methods go here
var pages = map[string]Page{
	"search":  {update: SearchUpdate, view: SearchView},
	"article": {update: ArticleUpdate, view: ArticleView},
}

// ---------------------------------------
// Bubbletea model, Update, View, and Init
// ---------------------------------------

type model struct {
	pageName string
	client   *Client
	// Used in search view
	textInput textinput.Model
	Articles  map[int]Article
	cursor    int
	info      string
	// Article view
	shownArticle string
	viewport     viewport.Model
	ready        bool
	content      string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		// Wait for window dimensions before initializing viewport
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}
	// Use Update method of current page
	if page, ok := pages[m.pageName]; ok {
		return page.update(m, msg)
	}
	return m, tea.Quit
}

func (m model) View() string {
	if page, ok := pages[m.pageName]; ok {
		return page.view(m)
	}
	return "I don't know how you ended up here.."
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// --------------------
// Initial model & main
// --------------------

func initialModel(topic string) model {
	ti := textinput.New()
	ti.Placeholder = "Giraffe"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.SetValue(topic)

	client, err := NewClient("en", DefaultWikiUrl, DefaultApiUrl)
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	var vp viewport.Model
	vp.Style = lipgloss.NewStyle()

	return model{
		pageName:  "search",
		client:    client,
		textInput: ti,
		Articles:  DefaultArticleMap,
		content:   "Waiting for content...",
		ready:     false,
		viewport:  vp,
	}
}

func main() {
	topic := flag.String("t", "", "Optional starting topic to search\nExample: wki -t Lions")
	help := flag.Bool("help", false, "Show this help menu")
	flag.Parse()
	if *help {
		fmt.Println(ExtendedUsage)
		flag.Usage()
		os.Exit(0)
	}
	if flag.NArg() > 0 {
		fmt.Println(ExtendedUsage)
		flag.Usage()
		os.Exit(1)
	}

	p := tea.NewProgram(
		initialModel(*topic),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
