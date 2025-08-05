package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
)

type OptionsInstallPrint struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Runtime        string `yaml:"runtime" json:"runtime" bson:"runtime" validate:""`
	Full           bool   `yaml:"full" json:"full" bson:"full" validate:""`
}

func NewOptionsInstallPrint() *OptionsInstallPrint {
	var o OptionsInstallPrint
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallPrint() *cobra.Command {
	o := NewOptionsInstallPrint()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("print")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_install_print_short")
	msgLong := OptCommon.TransLang("cmd_install_print_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_install_print_example", baseName))

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
	cmd.Flags().StringVar(&o.Runtime, "runtime", "", OptCommon.TransLang("param_install_print_runtime"))
	cmd.Flags().BoolVarP(&o.Full, "full", "", false, OptCommon.TransLang("param_install_print_full"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallPrint) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("runtime", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"docker", "containerd", "crio"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("runtime")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallPrint) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.Runtime != "docker" && o.Runtime != "containerd" && o.Runtime != "crio" {
		err = fmt.Errorf("--runtime must be docker, containerd or crio")
		return err
	}

	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsInstallPrint) Run(args []string) error {
	var err error

	bs, err := pkg.FsInstallConfigs.ReadFile(fmt.Sprintf("%s/%s-install-config.yaml", pkg.DirInstallConfigs, o.Language))
	if err != nil {
		return err
	}
	vals := map[string]interface{}{
		"runtime":  o.Runtime,
		"full":     o.Full,
		"language": o.Language,
		"baseName": pkg.GetCmdBaseName(),
	}
	strInstallConfig, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("parse install config error: %s", err.Error())
		return err
	}
	fmt.Println(strInstallConfig)
	return err
}
