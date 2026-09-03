// package commands cmd/calvin/commands/root.go
package commands

import (
	"bufio"
	"fmt"
	"github.com/0magnet/calvin"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// rootCmd represents the base command for the application
var RootCmd = &cobra.Command{
	Use:   "calvin",
	Short: "generate calvin ascii font from text",
	Long:  calvin.AsciiFont("calvin") + "\ngenerate calvin ascii font from text",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string

		// Arguments win over stdin. Checking stdin first meant that whenever
		// it was not a terminal — a script, a Makefile, CI — calvin ignored
		// its arguments and blocked waiting for an EOF that never came.
		// A failure here means stdin cannot be described, which is treated the
		// same as it not being a pipe: fall through to the argument.
		stat, statErr := os.Stdin.Stat()
		switch {
		case len(args) > 0:
			input = strings.Join(args, " ")
		case statErr == nil && (stat.Mode()&os.ModeCharDevice) == 0:
			scanner := bufio.NewScanner(os.Stdin)
			var sb strings.Builder
			for scanner.Scan() {
				sb.WriteString(scanner.Text())
				sb.WriteString("\n")
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading stdin: %w", err)
			}
			input = sb.String()
		default:
			return fmt.Errorf("no input provided; pipe text or pass as arguments")
		}

		// Generate and print the ASCII font
		output := calvin.AsciiFont(input)
		fmt.Println(output)
		return nil
	},
}
