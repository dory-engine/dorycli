package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdPipeline() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("pipeline")
	msgShort := fmt.Sprintf("manage pipeline resources")
	msgLong := fmt.Sprintf(`manage pipeline resources in dory-engine server`)
	msgExample := fmt.Sprintf(`  # get all pipeline resources
  %s pipeline get

  # execute pipeline
  %s pipeline execute test-project1-develop`, baseName, baseName)

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

	cmd.AddCommand(NewCmdPipelineGet())
	cmd.AddCommand(NewCmdPipelineExecute())
	return cmd
}
