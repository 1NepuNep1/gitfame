package cmd

import (
	"io"
	"os"
)

func Do(args []string, in io.Reader, out io.Writer, errOut io.Writer) int {
	rootCmd := NewRootCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		return 2
	}

	return 0
}
