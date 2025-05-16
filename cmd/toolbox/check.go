package toolbox

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run formatting, vet, and linters",
	Long:  `Run go fmt, go vet, golangci-lint, and staticcheck on the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		commands := [][]string{
			{"go", "fmt", "./..."},
			{"go", "vet", "./..."},
			{"golangci-lint", "run"},
			{"staticcheck", "./..."},
		}
		for _, c := range commands {
			cmdStr := strings.Join(c, " ")
			fmt.Printf("Running %s...\n", cmdStr)
			out, err := exec.Command(c[0], c[1:]...).CombinedOutput()
			fmt.Println(strings.TrimSpace(string(out)))
			if err != nil {
				log.Fatalf("Error running %s: %v\n", cmdStr, err)
			}
		}
		fmt.Println("All checks passed!")
	},
}
