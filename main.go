package main

import (
	"github.com/dory-engine/dorycli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var err error
	rootCmd := cmd.NewCmdRoot()
	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
