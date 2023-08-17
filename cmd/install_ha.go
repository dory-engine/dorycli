package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdInstallHa() *cobra.Command {
	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("ha")
	msgShort := fmt.Sprintf("create high availability kubernetes cluster load balancer")
	msgLong := fmt.Sprintf(`create high availability kubernetes cluster load balancer with keepalived and nginx, this command will create keepalived and nginx config files and docker-compose files and kuberentes install files`)
	msgExample := fmt.Sprintf(`  high availability kubernetes cluster installation document please check:
  https://github.com/cookeem/kubeadm-ha

  ##############################
  # please follow these steps to create load balancer config files:
  
  # 1. print load balancer installation settings YAML file
  %s install ha print > kubernetes-ha.yaml
  
  # 2. modify load balancer installation settings YAML file by manual
  vi kubernetes-ha.yaml
  
  # 3. create load balancer config files and docker-compose files and kuberentes install files 
  %s install ha script -o readme-kubernetes-ha -f kubernetes-ha.yaml`, baseName, baseName)

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

	cmd.AddCommand(NewCmdInstallHaPrint())
	cmd.AddCommand(NewCmdInstallHaScript())
	return cmd
}
