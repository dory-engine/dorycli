package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdProject() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("project")
	msgShort := fmt.Sprintf("manage project resources")
	msgLong := fmt.Sprintf(`manage project resources in dory-engine server`)
	msgExample := fmt.Sprintf(`  # get project resources
  %s project get

  # execute project ops batch
  %s project execute test-project1 your-ops-batch-name`, baseName, baseName)

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

	cmd.AddCommand(NewCmdProjectGet())
	cmd.AddCommand(NewCmdProjectExecute())
	return cmd
}
