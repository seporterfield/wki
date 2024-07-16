package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	search          string
	displayArticles map[int]string
	cursor          int
}

func initialModel(topic string) model {
	return model{
		search:          topic,
		displayArticles: map[int]string{0: "Lions", 1: "India", 2: "Submarines", 3: "Turtles", 4: "Canada", 5: "Go_(programming_language)"}, //make(map[int]string)
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// Navigates the TUI to the selected display article
func navigateDisplayArticle(displayArticleNum int) {

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.displayArticles)-1 {
				m.cursor++
			}

		case "enter":
			navigateDisplayArticle(m.cursor)
		default:
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := "wki - Search Wikipedia\n\n"
	for i := 0; i < len(m.displayArticles); i++ {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] \n", cursor, m.displayArticles[i])
	}

	// The footer
	s += "\nPress esc to quit.\n"

	// Send the UI for rendering
	return s
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
