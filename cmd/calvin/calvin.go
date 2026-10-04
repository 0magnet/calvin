// package main cmd/calvin/calvin.go
/*
convert text to ascii art
*/

package main

import (
	"os"

	"github.com/0magnet/calvin/cmd/calvin/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
