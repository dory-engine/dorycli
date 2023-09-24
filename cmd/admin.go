package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdAdmin() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("admin")
	msgShort := fmt.Sprintf("manage configurations, admin permission required")
	msgLong := fmt.Sprintf(`manage users, custom steps, kubernetes environments, component templates, docker build environments, repository configurations in dory-engine server, admin permission required`)
	msgExample := fmt.Sprintf(`  # get all users, custom steps, kubernetes environments and component templates, docker build environments, repository configurations, admin permission required
  %s admin get %s

  # apply multiple configurations from file or directory, admin permission required
  %s admin apply -f users.yaml -f custom-steps.json

  # delete configuration items, admin permission required
  %s admin delete %s customStepName1`, baseName, pkg.AdminKindAll, baseName, baseName, pkg.AdminKindCustomStep)

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

	cmd.AddCommand(NewCmdAdminGet())
	cmd.AddCommand(NewCmdAdminApply())
	cmd.AddCommand(NewCmdAdminDelete())
	return cmd
}
