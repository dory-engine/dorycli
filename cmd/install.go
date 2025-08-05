package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

func NewCmdInstall() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("install")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_install_short")
	msgLong := OptCommon.TransLang("cmd_install_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_install_example", baseName, baseName, baseName, baseName, baseName))

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
				subcommands := []string{"check", "print", "pull", "run", "script", "ha"}
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

	cmd.AddCommand(NewCmdInstallCheck())
	cmd.AddCommand(NewCmdInstallPrint())
	cmd.AddCommand(NewCmdInstallPull())
	cmd.AddCommand(NewCmdInstallScript())
	cmd.AddCommand(NewCmdInstallHa())
	return cmd
}
