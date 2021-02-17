package main

import (
	"os"

	"github.com/aberestyak/kubectl-get-yaml/pkg/cmd"
)

func main() {
	command := cmd.NewCmdGetYaml()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
