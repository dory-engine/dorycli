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
)

type OptionsProjectGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Output         string `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		ProjectNames []string `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
	}
}

func NewOptionsProjectGet() *OptionsProjectGet {
	var o OptionsProjectGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdProjectGet() *cobra.Command {
	o := NewOptionsProjectGet()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("get [projectName] ...")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_project_get_short")
	msgLong := OptCommon.TransLang("cmd_project_get_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_project_get_example", baseName, baseName, baseName))

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_project_get_output"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsProjectGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectNames, err := o.GetProjectNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if len(args) >= 0 {
			return projectNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml", "table"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsProjectGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	projectNames := args
	for _, s := range projectNames {
		s = strings.Trim(s, " ")
		err = pkg.ValidateMinusNameID(s)
		if err != nil {
			err = fmt.Errorf("projectNames error: %s", err.Error())
			return err
		}
		o.Param.ProjectNames = append(o.Param.ProjectNames, s)
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" && o.Output != "table" {
			err = fmt.Errorf("--output must be yaml or json or table")
			return err
		}
	}
	return err
}

func (o *OptionsProjectGet) Run(args []string) error {
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
		"projectNames": o.Param.ProjectNames,
		"page":         1,
		"perPage":      1000,
	}
	result, _, err := o.QueryAPI("api/cicd/projects", http.MethodPost, "", param, false)
	if err != nil {
		return err
	}
	rs := result.Get("data.projects").Array()
	projects := []pkg.Project{}
	for _, r := range rs {
		project := pkg.Project{}
		err = json.Unmarshal([]byte(r.Raw), &project)
		if err != nil {
			return err
		}
		projects = append(projects, project)
	}

	if len(projects) > 0 {
		dataOutput := map[string]interface{}{}
		if len(o.Param.ProjectNames) == 1 && len(projects) == 1 && o.Param.ProjectNames[0] == projects[0].ProjectInfo.ProjectName {
			dataOutput["project"] = projects[0]
		} else {
			dataOutput["projects"] = projects
		}
		switch o.Output {
		case "json":
			bs, _ = json.MarshalIndent(dataOutput, "", "  ")
			fmt.Println(string(bs))
		case "yaml":
			bs, _ = pkg.YamlIndent(dataOutput)
			fmt.Println(string(bs))
		default:
			data := [][]string{}
			for _, project := range projects {
				projectName := project.ProjectInfo.ProjectName
				projectShortName := project.ProjectInfo.ProjectShortName
				projectDesc := project.ProjectInfo.ProjectDesc
				projectEnvs := []string{}
				for _, pnp := range project.ProjectNodePorts {
					nodePorts := []string{}
					for _, envNodePort := range pnp.EnvNodePorts {
						nodePorts = append(nodePorts, fmt.Sprintf("%d", envNodePort.NodePortStart))
					}
					envInfo := fmt.Sprintf("%s:%s", pnp.EnvName, strings.Join(nodePorts, ","))
					projectEnvs = append(projectEnvs, envInfo)
				}
				projectEnvNames := strings.Join(projectEnvs, "\n")
				pipelines := []string{}
				for _, pp := range project.Pipelines {
					p := fmt.Sprintf("%s:%d/%d/%d", pp.PipelineName, pp.SuccessCount, pp.FailCount, pp.AbortCount)
					pipelines = append(pipelines, p)
				}
				pipelineNames := strings.Join(pipelines, "\n")

				builds := []string{}
				for mt, mds := range project.Modules {
					if mt == "build" {
						for _, module := range mds {
							builds = append(builds, module.ModuleName)
						}
						break
					}
				}
				buildNames := strings.Join(builds, "\n")

				data = append(data, []string{projectName, projectShortName, projectDesc, projectEnvNames, pipelineNames, buildNames})
			}

			dataHeader := []string{"Name", "ShortName", "Desc", "EnvNames", "Pipelines", "BuildNames"}
			table = tablewriter.NewTable(os.Stdout, tableRender, tableCellConfig)
			table.Header(dataHeader)
			table.Bulk(data)
			table.Render()
		}
	}

	return err
}
