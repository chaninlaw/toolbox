package toolbox

import (
	"log"

	tui "github.com/chaninlaw/toolbox/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const version = "v0.1.3"

var cmd = &cobra.Command{
	Use:   "toolbox",
	Short: "toolbox - code generator CLI",
	Long: `toolbox - code generator CLI,

A tool for generating project boilerplate and utilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.InitialModel())
		if _, err := p.Run(); err != nil {
			log.Fatalf("Oops, Could not start the program: %v\n", err)
		}
	},
	Version: version,
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Oops, Could not execute the command: %v\n", err)
	}
}
