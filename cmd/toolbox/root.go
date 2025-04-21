package toolbox

import (
	"fmt"
	"os"

	tui "github.com/chaninlaw/toolbox/internal/tui"
	"github.com/chaninlaw/toolbox/pkgs/logs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var version = "v0.1.2"

var cmd = &cobra.Command{
	Use:   "toolbox",
	Short: "toolbox - code generator CLI",
	Long: `toolbox - code generator CLI

A tool for generating project boilerplate and utilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.InitialModel())
		if _, err := p.Run(); err != nil {
			logs.Error.Fatalf("Oops, Could not start the program: %v\n", err)
		}
	},
}

func Execute() {
	cmd.Version = version
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
