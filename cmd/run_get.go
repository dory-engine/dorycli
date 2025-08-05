package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
	"time"
)

type OptionsRunGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ProjectNames   []string `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
	PipelineNames  []string `yaml:"pipelineNames" json:"pipelineNames" bson:"pipelineNames" validate:""`
	StatusResults  []string `yaml:"statusResults" json:"statusResults" bson:"statusResults" validate:""`
	StartDate      string   `yaml:"startDate" json:"startDate" bson:"startDate" validate:""`
	EndDate        string   `yaml:"endDate" json:"endDate" bson:"endDate" validate:""`
	Page           int      `yaml:"page" json:"page" bson:"page" validate:""`
	Number         int      `yaml:"number" json:"number" bson:"number" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		StartDate time.Time `yaml:"startDate" json:"startDate" bson:"startDate" validate:""`
		EndDate   time.Time `yaml:"endDate" json:"endDate" bson:"endDate" validate:""`
		RunNames  []string  `yaml:"runNames" json:"runNames" bson:"runNames" validate:""`
	}
}

func NewOptionsRunGet() *OptionsRunGet {
	var o OptionsRunGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdRunGet() *cobra.Command {
	o := NewOptionsRunGet()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("get [runName]")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_run_get_short")
	msgLong := OptCommon.TransLang("cmd_run_get_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_run_get_example", baseName, baseName))

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
	cmd.Flags().StringSliceVar(&o.ProjectNames, "projects", []string{}, OptCommon.TransLang("param_run_get_projects"))
	cmd.Flags().StringSliceVar(&o.PipelineNames, "pipelines", []string{}, OptCommon.TransLang("param_run_get_pipelines"))
	cmd.Flags().StringSliceVar(&o.StatusResults, "statuses", []string{}, OptCommon.TransLang("param_run_get_statuses"))
	cmd.Flags().StringVar(&o.StartDate, "start", "", OptCommon.TransLang("param_run_get_start"))
	cmd.Flags().StringVar(&o.EndDate, "end", "", OptCommon.TransLang("param_run_get_end"))
	cmd.Flags().IntVar(&o.Page, "page", 1, OptCommon.TransLang("param_run_get_page"))
	cmd.Flags().IntVarP(&o.Number, "number", "n", 50, OptCommon.TransLang("param_run_get_number"))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_run_get_output"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsRunGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	statuses := []string{
		"SUCCESS",
		"FAIL",
		"ABORT",
		"RUNNING",
		"INPUT",
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) >= 0 {
			runNames, err := o.GetRunNames()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return runNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("projects", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectNames, err := o.GetProjectNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return projectNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("pipelines", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		pipelineNames, err := o.GetPipelineNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return pipelineNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("statuses", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return statuses, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("start", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		s := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		return []string{s}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("end", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		s := time.Now().Format("2006-01-02")
		return []string{s}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml", "table"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsRunGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	for _, name := range o.ProjectNames {
		err = pkg.ValidateMinusNameID(name)
		if err != nil {
			err = fmt.Errorf("--projects %s error: %s", name, err.Error())
			return err
		}
	}

	for _, name := range o.PipelineNames {
		err = pkg.ValidatePipelineName(name)
		if err != nil {
			err = fmt.Errorf("--pipeilnes %s error: %s", name, err.Error())
			return err
		}
	}

	statuses := []string{
		"SUCCESS",
		"FAIL",
		"ABORT",
		"RUNNING",
		"INPUT",
	}

	for _, name := range o.StatusResults {
		var found bool
		for _, s := range statuses {
			if name == s {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("--statuses %s error: must be %s", name, strings.Join(statuses, " / "))
			return err
		}
	}

	runNames := args
	m := map[string]string{}
	for _, runName := range runNames {
		err = pkg.ValidateRunName(runName)
		if err != nil {
			err = fmt.Errorf("runName %s error: %s", runName, err.Error())
			return err
		}
		arr := strings.Split(runName, "-")
		pipelineName := strings.Join(arr[:len(arr)-1], "-")
		m[pipelineName] = ""
	}
	if len(m) > 0 {
		pipelineNames := []string{}
		for k, _ := range m {
			pipelineNames = append(pipelineNames, k)
		}
		o.PipelineNames = pipelineNames
	}
	o.Param.RunNames = runNames

	if o.EndDate == "" {
		o.EndDate = time.Now().Format("2006-01-02")
	}
	if o.StartDate != "" {
		o.Param.StartDate, err = time.Parse("2006-01-02", o.StartDate)
		if err != nil {
			err = fmt.Errorf("--startDate error: %s", err.Error())
			return err
		}
	}
	if o.EndDate != "" {
		o.Param.EndDate, err = time.Parse("2006-01-02", o.EndDate)
		if err != nil {
			err = fmt.Errorf("--endDate error: %s", err.Error())
			return err
		}
	}
	if o.Param.StartDate.After(o.Param.EndDate) {
		err = fmt.Errorf("--startDate must after --endDate")
		return err
	}

	if o.Page < 1 {
		err = fmt.Errorf("--page must greater than 1")
		return err
	}

	if o.Number < 1 {
		err = fmt.Errorf("--number must greater than 1")
		return err
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" && o.Output != "table" {
			err = fmt.Errorf("--output must be yaml or json or table")
			return err
		}
	}
	return err
}

func (o *OptionsRunGet) Run(args []string) error {
	var err error

	var table *tablewriter.Table
	var tableRender, tableCellConfig tablewriter.Option
	if o.Output == "table" {
		tableRender = pkg.TableRenderBorder
	} else {
		tableRender = pkg.TableRenderBorderNone
	}
	tableCellConfig = pkg.TableCellConfig

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{
		"projectNames":  o.ProjectNames,
		"pipelineNames": o.PipelineNames,
		"runNames":      o.Param.RunNames,
		"statusResults": o.StatusResults,
		"startTimeRage": map[string]string{
			"startDate": o.StartDate,
			"endDate":   o.EndDate,
		},
		"page":    o.Page,
		"perPage": o.Number,
	}
	result, _, err := o.QueryAPI("api/cicd/runs", http.MethodPost, "", param, false)
	if err != nil {
		return err
	}
	rs := result.Get("data.runs").Array()
	runs := []pkg.Run{}
	for _, r := range rs {
		run := pkg.Run{}
		err = json.Unmarshal([]byte(r.Raw), &run)
		if err != nil {
			return err
		}
		runs = append(runs, run)
	}

	if len(runs) > 0 {
		dataOutput := map[string]interface{}{}
		dataOutput["runs"] = runs
		switch o.Output {
		case "json":
			bs, _ = json.MarshalIndent(dataOutput, "", "  ")
			fmt.Println(string(bs))
		case "yaml":
			bs, _ = pkg.YamlIndent(dataOutput)
			fmt.Println(string(bs))
		default:
			data := [][]string{}
			for _, run := range runs {
				data = append(data, []string{run.RunName, run.PipelineArch, run.StartUser, run.Status.StartTime, run.Status.Result, run.Status.Duration})
			}

			dataHeader := []string{"Name", "Arch", "StartUser", "StartTime", "Status", "Duration"}
			table = tablewriter.NewTable(os.Stdout, tableRender, tableCellConfig)
			table.Header(dataHeader)
			table.Bulk(data)
			table.Render()
		}
	}

	return err
}
