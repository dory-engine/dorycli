package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdInstall() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("install")
	msgShort := fmt.Sprintf("install dory-engine with docker or kubernetes")
	msgLong := fmt.Sprintf(`install dory-engine and relative components with docker-compose or kubernetes`)
	msgExample := fmt.Sprintf(`  %s install should run on a managed kubernetes cluster node
  
  ##############################
  # please follow these steps to install dory-engine with kubernetes (for production recommended):
  
  # 1. check installing dory in kubernetes prerequisite, with dory managed kubernetes cluster over docker runtime
  %s install check --mode kubernetes --runtime docker
  
  # 2. print installing dory in kubernetes settings YAML file, with dory managed kubernetes cluster over docker runtime
  %s install print --mode kubernetes --runtime docker > install-config-kubernetes.yaml
  
  # 3. modify installation config file by manual
  vi install-config-kubernetes.yaml
  
  # 4. pull and build all container images required for installing dory in kubernetes over docker runtime
  %s install pull --mode kubernetes --runtime docker -f install-config-kubernetes.yaml
  
  # 5. (option 1) install dory with kubernetes automatically
  %s install run -o readme-install-kubernetes -f install-config-kubernetes.yaml
  
  # 5. (option 2) install dory with kubernetes by manual, it will output readme files, deploy files and config files, follow the readme files to customize install dory
  %s install script -o readme-install-kubernetes -f install-config-kubernetes.yaml

  ##############################
  please follow these steps to install dory-engine with docker (for test only):
  
  # 1. check installing dory in docker prerequisite, with dory managed cluster kubernetes over docker runtime
  %s install check --mode docker --runtime docker
  
  # 2. print installing dory in docker settings YAML file, with dory managed kubernetes cluster over docker runtime
  %s install print --mode docker --runtime docker > install-config-docker.yaml
  
  # 3. modify installation config file by manual
  vi install-config-docker.yaml
  
  # 4. pull and build all container images required for installing dory in docker
  %s install pull --mode docker -f install-config-docker.yaml
  
  # 5. (option 1) install dory with docker automatically
  %s install run -o readme-install-docker -f install-config-docker.yaml
  
  # 5. (option 2) install dory with docker by manual, it will output readme files, deploy files and config files, follow the readme files to customize install dory
  %s install script -o readme-install-docker -f install-config-docker.yaml
`, baseName, baseName, baseName, baseName, baseName, baseName, baseName, baseName, baseName, baseName, baseName)

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

	cmd.AddCommand(NewCmdInstallCheck())
	cmd.AddCommand(NewCmdInstallPrint())
	cmd.AddCommand(NewCmdInstallPull())
	cmd.AddCommand(NewCmdInstallRun())
	cmd.AddCommand(NewCmdInstallScript())
	cmd.AddCommand(NewCmdInstallHa())
	return cmd
}
