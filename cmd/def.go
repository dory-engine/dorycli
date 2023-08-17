package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdDef() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("def")
	msgShort := fmt.Sprintf("manage project definitions")
	msgLong := fmt.Sprintf(`manage project definitions in dory-engine server`)
	msgExample := fmt.Sprintf(`  # get project all definitions
  %s def get test-project1 all

  # apply project definitions from file or directory
  %s def apply -f def1.yaml -f def2.json

  # clone project definitions deploy modules to another environments
  %s def clone test-project1 deploy --from-env=test --modules=tp1-gin-demo,tp1-node-demo --to-envs=uat,prod

  # delete modules from project build definitions
  %s def delete test-project1 build --modules=tp1-gin-demo,tp1-node-demo

  # patch project build modules definitions, update tp1-gin-demo,tp1-go-demo buildChecks commands
  %s def patch test-project1 build --modules=tp1-go-demo,tp1-gin-demo --patch='[{"action": "update", "path": "buildChecks", "value": ["ls -alh"]}]'`, baseName, baseName, baseName, baseName, baseName)

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

	cmd.AddCommand(NewCmdDefGet())
	cmd.AddCommand(NewCmdDefApply())
	cmd.AddCommand(NewCmdDefDelete())
	cmd.AddCommand(NewCmdDefClone())
	cmd.AddCommand(NewCmdDefPatch())
	return cmd
}
