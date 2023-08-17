package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"net/http"
)

type OptionsVersionRun struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
}

func NewOptionsVersionRun() *OptionsVersionRun {
	var o OptionsVersionRun
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdVersion() *cobra.Command {
	o := NewOptionsVersionRun()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("version")
	msgShort := fmt.Sprintf("show %s version info", baseName)
	msgLong := fmt.Sprintf(`show %s and isntall dory-engine, dory-console version info`, baseName)
	msgExample := fmt.Sprintf(`  # show %s and dory-engine, dory-console version info:
  %s version`, baseName, baseName)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Complete(cmd))
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}

	return cmd
}

func (o *OptionsVersionRun) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsVersionRun) Validate(args []string) error {
	var err error
	return err
}

func (o *OptionsVersionRun) Run(args []string) error {
	var err error
	baseName := pkg.GetCmdBaseName()
	fmt.Println(fmt.Sprintf("%s version: %s", baseName, pkg.VersionDoryCli))
	fmt.Println(fmt.Sprintf("install dory-engine version: %s", pkg.VersionDoryEngine))
	fmt.Println(fmt.Sprintf("install dory-console version: %s", pkg.VersionDoryFrontend))
	if o.ServerURL != "" {
		fmt.Println(fmt.Sprintf("serverURL: %s", o.ServerURL))
		if o.AccessToken != "" {
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/public/about"), http.MethodGet, "", param, false)
			if err != nil {
				return err
			}
			appInfo := result.Get("data.app").String()
			versionInfo := result.Get("data.version").String()
			fmt.Println(fmt.Sprintf("versionInfo: %s/%s", appInfo, versionInfo))
		}
	}

	return err
}
