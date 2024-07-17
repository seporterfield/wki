package main

import "github.com/charmbracelet/lipgloss"

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

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Render

	listArticleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Render
	articleDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Render
	articleBoldedStyle = lipgloss.NewStyle().
				Bold(true).
				Render
	noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Render
)
