package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ChoicesState represents the state of the choices
type choiceState struct {
	cursor int
	chosen string
}

type choice struct {
	Label string
	View  func(m model) string
}

var choices []choice

func initializeChoices() choiceState {
	return choiceState{
		cursor: 1,
		chosen: "",
	}
}

func registerChoice(label string, view func(m model) string) {
	choices = append(choices, choice{
		Label: label,
		View:  view,
	})
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice.cursor++
			maxChoices := len(choices)
			if m.choice.cursor > maxChoices {
				m.choice.cursor = maxChoices
			}
		case "k", "up":
			m.choice.cursor--
			if m.choice.cursor <= 0 {
				m.choice.cursor = 1
			}
		case "enter", " ":
			m.choice.chosen = choices[m.choice.cursor-1].Label
			return m, nil
		}
	}

	return m, nil
}

func choicesView(m model) string {
	c := m.choice.cursor

	s := "What you would like to do today?\n\n"

	for i, choice := range choices {
		s += fmt.Sprintf("%s\n", checkbox(choice.Label, c == i+1))
	}
	s += "\nMore options coming soon...\n"

	s += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter, space: choose") + dotStyle +
		subtleStyle.Render("ctrl+c: quit")

	return s
}
