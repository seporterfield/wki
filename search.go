package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func truncate(str string, length int) (truncated string) {
	if length <= 0 {
		return
	}
	for i, char := range str {
		if i >= length {
			break
		}
		truncated += string(char)
	}
	return
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
		maxW := m.viewport.Width - len(m.Articles[i].Title) - 5
		truncatedDescription := truncate(m.Articles[i].Description, maxW)
		s += fmt.Sprintf("%s %s — %s \n", cursor, listArticleStyle(m.Articles[i].Title), truncatedDescription)
	}

	// The footer
	s += "\nNavigate: ←↑↓→ ↲.\nEnter vim normal mode: ESC.\nPress <CTRL-C> to exit.\n"
	s += m.info

	// Send the UI for rendering
	return s
}

func SearchUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// If a user used the -t option we want to
	// update the article list using their input
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		m.info = ""

		// Cool, what was the actual key pressed?
		msgStr := msg.String()
		switch {
		case msgStr == "esc":
			if !m.normalMode {
				m.normalMode = true
			}
		case msgStr == "ctrl+c":
			return m, tea.Quit
		case msgStr == "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case msgStr == "down":
			if m.cursor < len(m.Articles)-1 {
				m.cursor++
			}
		case msgStr == "j" && m.normalMode:
			if m.cursor < len(m.Articles)-1 {
				m.cursor++
			}
		case msgStr == "k" && m.normalMode:
			if m.cursor > 0 {
				m.cursor--
			}
		case (msgStr == "a" || msgStr == "i") && m.normalMode:
			m.normalMode = !m.normalMode
		case msgStr == "enter":
			// TODO: on right-key press if we're at the last
			// character of the input we should go to the
			// article view.
			article := m.Articles[m.cursor]
			// If the user is fast enough to hit enter before
			// the articles load for the query they provided
			// with -t we want to stop them from getting garbage
			if article.Title == DefaultArticleMap[0].Title {
				break
			}

			m.pageName = "article"
			m.shownArticle = article.Title

			// "Cache" existing content
			if m.Articles[m.cursor].Content != "" {
				break
			}

			newArticle, err := m.client.LoadArticle(article)
			if err != nil {
				m.info = err.Error()
				break
			}

			m.Articles[m.cursor] = newArticle
			m.content = lipgloss.NewStyle().Width(m.viewport.Width).Render(newArticle.Content)
			m.viewport.SetContent(m.content)
		case msgStr == "left" || msgStr == "right":
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		default:
			if !m.normalMode {
				m.textInput, cmd = m.textInput.Update(msg)
				m.cursor = 0
				return m, tea.Batch(cmd, m.queryArticlesCmd())
			}
		}
	case apiResponseMsg:
		if msg.query != m.textInput.Value() {
			break
		}
		m.Articles = msg.articles
	}
	if strings.TrimSpace(m.textInput.Value()) == "" {
		m.Articles = DefaultArticleMap
	}

	// Should be checked towards the end so we don't
	// get stuck in an infinite loop
	if m.Articles[0].Title == DefaultArticleMap[0].Title {
		return m, tea.Batch(cmd, m.queryArticlesCmd())
	}
	return m, cmd
}

func (m model) queryArticlesCmd() tea.Cmd {
	query := m.textInput.Value()
	return func() tea.Msg {
		articles, err := m.client.LoadSearchList(query)
		if err != nil {
			m.info = err.Error()
		}
		return apiResponseMsg{articles: articles, query: query}
	}
}

type apiResponseMsg struct {
	articles map[int]Article
	query    string
}
