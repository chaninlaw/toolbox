package toolbox

import (
	tui "github.com/chaninlaw/toolbox/internal/tui"
	"github.com/chaninlaw/toolbox/pkgs/logs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize the project",
		Long:  `Initialize a new Go project with a predefined structure and boilerplate code.`,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(tui.InitialGeneratorModel())
			if _, err := p.Run(); err != nil {
				logs.Fatal("Oops, Could not start the program: %v", err)
			}
		},
	})
}
