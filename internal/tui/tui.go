package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choice    choiceState
	generator generatorState
	isQuit    bool
}

// InitialModel to start the program with the choices
func InitialModel() model {
	return model{
		isQuit:    false,
		choice:    initializeChoices(),
		generator: initializeGenerator(),
	}
}

// InitialGeneratorModel to allow skipping choices and go directly to generator
func InitialGeneratorModel() model {
	return model{
		isQuit:    false,
		choice:    choiceState{cursor: 1, chosen: generatorLabel},
		generator: initializeGenerator(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		m.generatorInit,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.isQuit = true
			return m, tea.Quit
		}
	}

	switch m.choice.chosen {
	case "":
		return updateChoices(msg, m)
	case generatorLabel:
		return updateGenerator(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	if m.isQuit {
		return "\n  See you later!\n\n"
	}

	var s string
	if m.choice.chosen == "" {
		s = choicesView(m)
	} else if m.choice.chosen == generatorLabel {
		s = generatorView(m)
	} else {
		s = choices[m.choice.cursor-1].View(m)
	}

	return mainStyle.Render("\n" + s + "\n\n")
}
