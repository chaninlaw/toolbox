package main

import (
	tui "github.com/chaninlaw/toolbox/internal/tui"
	"github.com/chaninlaw/toolbox/pkgs/logs"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		logs.Fatal("Oops, Could not start the program: %v", err)
	}
}
