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
	// states
	idle = iota
	confirming
	generating
	done

	// focus index
	inputFocus      = 0
	liveReloadFocus = 1
)

// Generator state represents the state of the generator
type generatorState struct {
	projectName textinput.Model
	spinner     spinner.Model
	errMsg      string
	state       int
	progressMsg string
	focusIndex  int  // 0 = text input, 1 = live reload
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
	ti.Placeholder = "Enter project name..."
	ti.Prompt = "./"
	ti.PromptStyle = keywordStyle
	ti.Cursor.Style = keywordStyle
	ti.Width = 30
	ti.CharLimit = 30
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return generatorState{
		projectName: ti,
		state:       idle,
		spinner:     s,
		liveReload:  true, // default true
		focusIndex:  0,    // start with text input focused
	}
}

// --- Update ---

func updateGenerator(msg tea.Msg, m model) (model, tea.Cmd) {
	switch m.generator.state {
	case idle:
		return updateGeneratorIdle(msg, m)
	case confirming:
		return updateGeneratorConfirming(msg, m)
	case generating:
		return updateGeneratorGenerating(msg, m)
	case done:
		return updateGeneratorDone(msg, m)
	default:
		return m, nil
	}
}

// --- State Handlers ---

func updateGeneratorIdle(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.generator.focusIndex = (m.generator.focusIndex + 1) % 2
			if m.generator.focusIndex == 0 {
				m.generator.projectName.Focus()
			} else {
				m.generator.projectName.Blur()
			}
			return m, nil
		case "shift+tab", "up":
			m.generator.focusIndex = (m.generator.focusIndex + 1) % 2
			if m.generator.focusIndex == 0 {
				m.generator.projectName.Focus()
			} else {
				m.generator.projectName.Blur()
			}
			return m, nil
		case "enter":
			if m.generator.projectName.Err != nil {
				m.generator.projectName.Err = nil
			}
			if m.generator.projectName.Value() == "" {
				m.generator.projectName.Placeholder = "Project name cannot be empty"
				m.generator.projectName.PlaceholderStyle = destructiveStyle
				return m, nil
			}
			m.generator.state = confirming
			return m, nil
		case " ":
			if m.generator.focusIndex == inputFocus {
				return m, nil
			}
			m.generator.liveReload = !m.generator.liveReload
			return m, nil
		}
	case projectGeneratedMsg:
		// Should not happen in idle, but handle gracefully
		if msg.err != nil {
			m.generator.projectName.Err = fmt.Errorf("%s", destructiveStyle.Render("Failed to generate project: "+msg.err.Error()))
			return m, nil
		}
		m.generator.state = done
		return m, tea.Quit
	}
	var cmd tea.Cmd
	if m.generator.focusIndex == inputFocus {
		m.generator.projectName, cmd = m.generator.projectName.Update(msg)
	}
	return m, cmd
}

func updateGeneratorConfirming(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			m.generator.state = generating
			m.generator.projectName.Blur()
			m.generator.progressMsg = fmt.Sprintf("Generating %s project...", m.generator.projectName.Value())
			return m, tea.Batch(
				m.generator.spinner.Tick,
				generateProjectCmd(generator.Options{
					ProjectName: m.generator.projectName.Value(),
					LiveReload:  m.generator.liveReload,
				}),
			)
		case "n", "esc":
			m.generator.state = idle
			return m, nil
		}
	}
	return m, nil
}

func updateGeneratorGenerating(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case projectGeneratedMsg:
		if msg.err != nil {
			m.generator.errMsg = destructiveStyle.Render("Failed to generate project: " + msg.err.Error())
			m.generator.state = idle
			return m, nil
		}
		m.generator.state = done
		return m, nil
	}
	var cmd tea.Cmd
	m.generator.spinner, cmd = m.generator.spinner.Update(msg)
	return m, cmd
}

func updateGeneratorDone(_ tea.Msg, m model) (model, tea.Cmd) {
	m.generator.state = done
	return m, tea.Quit
}

// --- View ---

func generatorView(m model) string {
	switch m.generator.state {
	case done:
		return viewDone(m)
	case confirming:
		return viewConfirming(m)
	case generating:
		return viewGenerating(m)
	case idle:
		fallthrough
	default:
		return viewIdle(m)
	}
}

func viewDone(m model) string {
	var s string
	s += "Project generation complete!\n\n"
	s += "You can now navigate to your project directory.\n"
	s += keywordStyle.Render("cd "+m.generator.projectName.Value()) + "\n"
	s += keywordStyle.Render("make dev") + "\n"
	s += "Have fun...\n"
	return s
}

func viewConfirming(m model) string {
	var s string
	s += "Please confirm your project setup:\n\n"
	s += "Project name: " + keywordStyle.Render(m.generator.projectName.Value()) + "\n"
	s += "Live reload: " + keywordStyle.Render(fmt.Sprintf("%v", m.generator.liveReload)) + "\n\n"
	s += subtleStyle.Render("y/enter: confirm  n/esc: edit") + "\n"
	return s
}

func viewGenerating(m model) string {
	var s string
	if m.generator.errMsg != "" {
		return m.generator.errMsg
	}
	s += fmt.Sprintf("\n%s %s", m.generator.spinner.View(), subtleStyle.Render(m.generator.progressMsg))
	return s
}

func viewIdle(m model) string {
	var s string
	s += "You have selected: " + m.choice.chosen + "\n\n"
	s += "Please enter the project name:\n"
	// Show text input with focus indicator
	if m.generator.focusIndex == 0 {
		s += "> " + m.generator.projectName.View() + "\n\n"
	} else {
		s += "  " + m.generator.projectName.View() + "\n\n"
	}
	// Show live reload with focus indicator
	if m.generator.focusIndex == 1 {
		s += "> " + checkbox("Enable live reload (Air)", m.generator.liveReload) + "\n"
	} else {
		s += "  " + checkbox("Enable live reload (Air)", m.generator.liveReload) + "\n"
	}
	if m.generator.projectName.Err != nil {
		s += m.generator.projectName.Err.Error()
	}
	s += "\n"
	s += subtleStyle.Render("tab / ↑↓: move") + dotStyle +
		subtleStyle.Render("space: toggle") + dotStyle +
		subtleStyle.Render("enter: execute") + dotStyle +
		subtleStyle.Render("ctrl+c: quit") + "\n"
	return s
}
