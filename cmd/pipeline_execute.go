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

type OptionsPipelineExecute struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Batch          string   `yaml:"batch" json:"batch" bson:"batch" validate:""`
	Params         []string `yaml:"params" json:"params" bson:"params" validate:""`
	Param          struct {
		PipelineName string   `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
		Batches      []string `yaml:"batches" json:"batches" bson:"batches" validate:""`
		QueryParam   string   `yaml:"queryParam" json:"queryParam" bson:"queryParam" validate:""`
	}
}

func NewOptionsPipelineExecute() *OptionsPipelineExecute {
	var o OptionsPipelineExecute
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdPipelineExecute() *cobra.Command {
	o := NewOptionsPipelineExecute()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("execute [pipelineName]")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_pipeline_execute_short")
	msgLong := OptCommon.TransLang("cmd_pipeline_execute_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_pipeline_execute_example", baseName, baseName, baseName))

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
	cmd.Flags().StringVarP(&o.Batch, "batch", "b", "", OptCommon.TransLang("param_pipeline_execute_batch"))
	cmd.Flags().StringSliceVar(&o.Params, "params", []string{}, OptCommon.TransLang("param_pipeline_execute_params"))
	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsPipelineExecute) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			pipelineNames, err := o.GetPipelineNames()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return pipelineNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return err
}

func (o *OptionsPipelineExecute) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		err = fmt.Errorf("args error: only accept one args")
		return err
	}

	s := args[0]
	s = strings.Trim(s, " ")
	err = pkg.ValidatePipelineName(s)
	if err != nil {
		err = fmt.Errorf("pipelineName error: %s", err.Error())
		return err
	}
	o.Param.PipelineName = s

	o.Batch = strings.Trim(o.Batch, " ")
	arr := strings.Split(o.Batch, "::")
	for _, val := range arr {
		val = strings.Trim(val, " ")
		if val != "" {
			o.Param.Batches = append(o.Param.Batches, val)
		}
	}

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

func (o *OptionsPipelineExecute) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{
		"batches": o.Param.Batches,
	}
	urlQuery := fmt.Sprintf("api/cicd/pipeline/%s?%s", o.Param.PipelineName, o.Param.QueryParam)
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
