package toolbox

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd.AddCommand(&cobra.Command{
		Use:   "greet",
		Short: "Greet the user",
		Long:  `A simple command to greet the user.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, welcome to toolbox!")
		},
	})
}
