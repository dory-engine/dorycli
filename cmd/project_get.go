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
	ProjectTeam    string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:""`
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
	msgShort := fmt.Sprintf("get project resources")
	msgLong := fmt.Sprintf(`get project resources in dory-engine server`)
	msgExample := fmt.Sprintf(`  # get all project resources
  %s project get

  # get single project resoure
  %s project get test-project1

  # get multiple project resources
  %s project get test-project1 test-project2`, baseName, baseName, baseName)

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
	cmd.Flags().StringVar(&o.ProjectTeam, "team", "", "filters by projectTeam")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")

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
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
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
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsProjectGet) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{
		"projectNames": o.Param.ProjectNames,
		"projectTeam":  o.ProjectTeam,
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
				projectEnvs := []string{}
				for _, pae := range project.ProjectAvailableEnvs {
					projectEnvs = append(projectEnvs, pae.EnvName)
				}
				projectEnvNames := strings.Join(projectEnvs, ",")
				projectNodePorts := []string{}
				for _, pnp := range project.ProjectNodePorts {
					np := fmt.Sprintf("%d-%d", pnp.NodePortStart, pnp.NodePortEnd)
					projectNodePorts = append(projectNodePorts, np)
				}
				projectNodePortNames := strings.Join(projectNodePorts, ",")
				pipelines := []string{}
				for _, pp := range project.Pipelines {
					pipelines = append(pipelines, pp.PipelineName)
				}
				pipelineNames := strings.Join(pipelines, ",")

				data = append(data, []string{projectName, projectShortName, projectEnvNames, projectNodePortNames, pipelineNames})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ShortName", "EnvNames", "NodePorts", "Pipelines"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)
			table.AppendBulk(data)
			table.Render()
		}
	}

	return err
}
