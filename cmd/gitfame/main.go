//go:build !solution

package main

import (
	"os"

	"github.com/1NepuNep1/gitfame/internal/cmd"
)

func main() {
	os.Exit(cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
