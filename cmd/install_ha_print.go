package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
)

type OptionsInstallHaPrint struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
}

func NewOptionsInstallHaPrint() *OptionsInstallHaPrint {
	var o OptionsInstallHaPrint
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallHaPrint() *cobra.Command {
	o := NewOptionsInstallHaPrint()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("print")
	msgShort := fmt.Sprintf("print ha install settings YAML file")
	msgLong := fmt.Sprintf(`print load balancer installation settings YAML file`)
	msgExample := fmt.Sprintf(`  # print installing dory in kubernetes settings YAML file, with dory managed kubernetes cluster over containerd runtime
  %s install print`, baseName)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallHaPrint) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallHaPrint) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsInstallHaPrint) Run(args []string) error {
	var err error

	bs, err := pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes-ha/%s-kubernetes-ha.yaml", pkg.DirInstallScripts, o.Language))
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return err
}
