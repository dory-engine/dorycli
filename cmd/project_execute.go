package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strings"
)

type OptionsProjectExecute struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Params         []string `yaml:"params" json:"params" bson:"params" validate:""`
	Param          struct {
		ProjectName  string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
		OpsBatchName string `yaml:"opsBatchName" json:"opsBatchName" bson:"opsBatchName" validate:""`
		QueryParam   string `yaml:"queryParam" json:"queryParam" bson:"queryParam" validate:""`
	}
}

func NewOptionsProjectExecute() *OptionsProjectExecute {
	var o OptionsProjectExecute
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdProjectExecute() *cobra.Command {
	o := NewOptionsProjectExecute()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("execute [projectName] [opsBatchName]")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_project_execute_short")
	msgLong := OptCommon.TransLang("cmd_project_execute_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_project_execute_example", baseName, baseName))

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
	cmd.Flags().StringSliceVar(&o.Params, "params", []string{}, OptCommon.TransLang("param_project_execute_params"))
	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsProjectExecute) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			projectNames, err := o.GetProjectNames()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return projectNames, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			projectNames, err := o.GetOpsBatchNames(args[0])
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return projectNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return err
}

func (o *OptionsProjectExecute) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) != 2 {
		err = fmt.Errorf("args error: only accept two args")
		return err
	}

	s := args[0]
	s = strings.Trim(s, " ")
	err = pkg.ValidateMinusNameID(s)
	if err != nil {
		err = fmt.Errorf("projectName error: %s", err.Error())
		return err
	}
	o.Param.ProjectName = s

	s = args[1]
	s = strings.Trim(s, " ")
	err = pkg.ValidateMinusNameID(s)
	if err != nil {
		err = fmt.Errorf("opsBatchName error: %s", err.Error())
		return err
	}
	o.Param.OpsBatchName = s

	params := []string{}
	for _, param := range o.Params {
		ss := strings.Split(param, "=")
		if len(ss) < 2 {
			err = fmt.Errorf("param %s format error: example: varName=varValue", param)
			return err
		}
		k := url.QueryEscape(ss[0])
		v := url.QueryEscape(strings.Join(ss[1:], "="))
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	o.Param.QueryParam = strings.Join(params, "&")
	return err
}

func (o *OptionsProjectExecute) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{}
	urlQuery := fmt.Sprintf("api/cicd/batch/%s/%s?%s", o.Param.ProjectName, o.Param.OpsBatchName, o.Param.QueryParam)
	result, _, err := o.QueryAPI(urlQuery, http.MethodPost, "", param, false)
	if err != nil {
		return err
	}
	runName := result.Get("data.runName").String()
	if runName == "" {
		err = fmt.Errorf("runName is empty")
		return err
	}

	urlQuery = fmt.Sprintf("api/cicd/run/%s", runName)
	result, _, err = o.QueryAPI(urlQuery, http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	run := pkg.Run{}
	err = json.Unmarshal([]byte(result.Get("data.run").Raw), &run)
	if err != nil {
		return err
	}

	if run.RunName == "" {
		err = fmt.Errorf("runName %s not exists", runName)
		return err
	}

	urlQuery = fmt.Sprintf("api/ws/log/run/%s", runName)
	err = o.QueryWebsocket(urlQuery, runName)
	if err != nil {
		return err
	}

	return err
}
