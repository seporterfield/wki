package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Article struct {
	title       string
	description string
	content     string
	url         string
}

type model struct {
	pageName string
	// Used in search view
	textInput textinput.Model
	Articles  map[int]Article
	cursor    int
	// Article view
	shownArticle string
	viewport     viewport.Model
	ready        bool
	content      string
}

func initialModel(topic string) model {
	ti := textinput.New()
	ti.Placeholder = topic
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		pageName:  "search",
		textInput: ti,
		Articles:  map[int]Article{0: {"...", "... waiting", "", ""}, 1: {"...", "... waiting", "", ""}},
		content:   "Waiting for content...",
		ready:     false,
	}
}

func (m model) headerView() string {
	title := titleStyle.Render(fmt.Sprintf("wki - %s", m.shownArticle))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func ArticleView(m model) string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func ArticleUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "left", "h":
			m.pageName = "search"
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func SearchView(m model) string {
	s := "wki - Search Wikipedia\n\n"
	s += m.textInput.View()
	s += "\n\n"
	for i := 0; i < len(m.Articles); i++ {

		cursor := " "
		if m.cursor == i {
			cursor = "*"
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] — %s \n", cursor, m.Articles[i].title, m.Articles[i].description)
	}

	// The footer
	s += "\nPress esc to quit. Arrow keys to navigate.\n"

	// Send the UI for rendering
	return s
}

func SearchUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down":
			if m.cursor < len(m.Articles)-1 {
				m.cursor++
			}

		case "enter", "right":
			m.pageName = "article"
			m.shownArticle = m.Articles[m.cursor].title
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

type Page struct {
	update func(model, tea.Msg) (tea.Model, tea.Cmd)
	view   func(model) string
}

var pages = map[string]Page{
	"search":  {update: SearchUpdate, view: SearchView},
	"article": {update: ArticleUpdate, view: ArticleView},
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
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

func main() {
	topic := flag.String("topic", "", "Topic to search")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
	}

	p := tea.NewProgram(initialModel(*topic))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
