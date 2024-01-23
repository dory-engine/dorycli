package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type OptionsRunLog struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Param          struct {
		RunName string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	}
}

func NewOptionsRunLog() *OptionsRunLog {
	var o OptionsRunLog
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdRunLog() *cobra.Command {
	o := NewOptionsRunLog()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("logs [runName]")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_run_logs_short")
	msgLong := OptCommon.TransLang("cmd_run_logs_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_run_logs_example", baseName))

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

func (o *OptionsRunLog) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			runNames, err := o.GetRunNames()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return runNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return err
}

func (o *OptionsRunLog) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		err = fmt.Errorf("runName error: only accept one runName")
		return err
	}

	s := args[0]
	s = strings.Trim(s, " ")
	err = pkg.ValidateRunName(s)
	if err != nil {
		err = fmt.Errorf("runName error: %s", err.Error())
		return err
	}
	o.Param.RunName = s
	return err
}

func (o *OptionsRunLog) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", o.Param.RunName), http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	run := pkg.Run{}
	err = json.Unmarshal([]byte(result.Get("data.run").Raw), &run)
	if err != nil {
		return err
	}

	if run.RunName == "" {
		err = fmt.Errorf("runName %s not exists", o.Param.RunName)
		return err
	}

	url := fmt.Sprintf("api/ws/log/run/%s", o.Param.RunName)
	err = o.QueryWebsocket(url, o.Param.RunName)
	if err != nil {
		return err
	}

	return err
}
