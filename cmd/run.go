package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdRun() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("run")
	msgShort := fmt.Sprintf("manage pipeline run resources")
	msgLong := fmt.Sprintf(`manage pipeline run resources in dory-engine server`)
	msgExample := fmt.Sprintf(`  # get pipeline run resources
  %s run get
  
  # show pipeline run logs
  %s run logs test-project1-develop-1
  
  # delete run, project maintainer permission required
  %s run abort test-project1-develop-1`, baseName, baseName, baseName)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(NewCmdRunGet())
	cmd.AddCommand(NewCmdRunLog())
	cmd.AddCommand(NewCmdRunAbort())
	return cmd
}
