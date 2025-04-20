package tui

import (
	"fmt"

	"github.com/chaninlaw/toolbox/internal/generator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	generatorLabel = "Initialize Project"
	idle           = iota
	generating
	done
)

// Generator state represents the state of the generator
type generatorState struct {
	projectName textinput.Model
	state       int
	spinner     spinner.Model
	progressMsg string
	liveReload  bool // Live reload option
}

// Message to signal project generation is done
type projectGeneratedMsg struct {
	err error
}

// Update generateProjectCmd to accept a reporter callback
func generateProjectCmd(options generator.Options) tea.Cmd {
	return func() tea.Msg {
		err := generator.Generate(options)
		return projectGeneratedMsg{err: err}
	}
}

func (m model) generatorInit() tea.Msg {
	registerChoice(generatorLabel, generatorView)
	return nil
}

func initializeGenerator() generatorState {
	ti := textinput.New()
	ti.Placeholder = "Project Name"
	ti.Prompt = "./"
	ti.PromptStyle = keywordStyle
	ti.Cursor.Style = keywordStyle
	ti.Width = 20
	ti.CharLimit = 20
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return generatorState{
		projectName: ti,
		state:       idle,
		spinner:     s,
		liveReload:  true, // default true
	}
}

func updateGenerator(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			// Toggle live reload
			m.generator.liveReload = !m.generator.liveReload
			return m, nil
		case " ":
			return m, nil
		case "enter":
			if m.generator.projectName.Err != nil {
				m.generator.projectName.Err = nil
			}

			if m.generator.projectName.Value() == "" {
				m.generator.projectName.Err = fmt.Errorf("%s", destructiveStyle.Render("Project name cannot be empty"))
				return m, nil
			}

			m.generator.state = generating
			m.generator.projectName.TextStyle = keywordStyle
			m.generator.projectName.Blur()
			// Start spinner and project generation
			return m, tea.Batch(
				m.generator.spinner.Tick,
				generateProjectCmd(generator.Options{
					ProjectName: m.generator.projectName.Value(),
					LiveReload:  m.generator.liveReload,
				}),
			)
		}
	case projectGeneratedMsg:
		if msg.err != nil {
			m.generator.projectName.Err = fmt.Errorf("%s", destructiveStyle.Render("Failed to generate project: "+msg.err.Error()))
			return m, nil
		}
		m.generator.state = done
		return m, tea.Quit
	}

	var cmd tea.Cmd

	// If generating, update spinner and return its Tick
	if m.generator.state == generating {
		m.generator.spinner, cmd = m.generator.spinner.Update(msg)
		return m, cmd
	}

	m.generator.projectName, cmd = m.generator.projectName.Update(msg)

	return m, cmd
}

func generatorView(m model) string {
	var s string

	if m.generator.state == done {
		s += "Project generation complete!\n\n"
		s += "You can now navigate to your project directory.\n"
		s += keywordStyle.Render("cd "+m.generator.projectName.Value()) + "\n"
		s += keywordStyle.Render("make dev") + "\n"
		s += "Have fun...\n"
		return s
	}

	if m.generator.state == idle {
		s += "You have selected: " + m.choice.chosen + "\n\n"
		s += "Please enter the project name:\n"
	}

	s += fmt.Sprintf("%s\n\n", m.generator.projectName.View())
	if m.generator.state == idle {
		s += checkbox("Enable live reload (Air)", m.generator.liveReload) + "\n"
	}
	if m.generator.state == generating {
		s += fmt.Sprintf("\n%s %s", m.generator.spinner.View(), subtleStyle.Render(m.generator.progressMsg))
	}
	if m.generator.projectName.Err != nil {
		s += m.generator.projectName.Err.Error()
	}

	s += "\n"
	if m.generator.state == idle {
		s += subtleStyle.Render("enter: to next") + "\n"
		s += subtleStyle.Render("l: toggle live reload") + "\n"
	}
	s += subtleStyle.Render("q: quit")

	return s
}
