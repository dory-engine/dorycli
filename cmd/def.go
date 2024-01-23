package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

func NewCmdDef() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("def")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_def_short")
	msgLong := OptCommon.TransLang("cmd_def_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_def_example", baseName, pkg.DefKindAll, baseName, baseName, pkg.DefKindDeployContainer, baseName, pkg.DefKindBuild, baseName, pkg.DefKindBuild))

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
			} else {
				var found bool
				subcommands := []string{"get", "apply", "delete", "clone", "patch"}
				sort.Strings(subcommands)
				for _, subcommand := range subcommands {
					if args[0] == subcommand {
						found = true
						break
					}
				}
				if !found {
					log.Error(fmt.Sprintf("subcommand options: %s\n", strings.Join(subcommands, " / ")))
					cmd.Help()
					os.Exit(0)
				}
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
