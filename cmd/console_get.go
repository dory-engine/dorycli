package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type OptionsConsoleGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Items          []string `yaml:"items" json:"items" bson:"items" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	BranchNames    []string `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kinds       []string `yaml:"kinds" json:"kinds" bson:"kinds" validate:""`
		ProjectName string   `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
		IsAllKind   bool     `yaml:"isAllKind" json:"isAllKind" bson:"isAllKind" validate:""`
	}
}

func NewOptionsConsoleGet() *OptionsConsoleGet {
	var o OptionsConsoleGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdConsoleGet() *cobra.Command {
	o := NewOptionsConsoleGet()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`get [projectName] [kind],[kind] [--envs=envName1,envName2] [--branches=branch1,branch2] [--items=itemName1,itemName2] [--output=json|yaml] ...`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_console_get_short")
	msgLong := OptCommon.TransLang("cmd_console_get_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_console_get_example", strings.Join(pkg.ConsoleKinds, ", "), baseName, pkg.ConsoleKindAll, baseName, pkg.ConsoleKindAll, baseName, pkg.ConsoleKindMember, pkg.ConsoleKindPipeline, baseName, pkg.ConsoleKindComponent, baseName, pkg.ConsoleKindPipelineTrigger))

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
	cmd.Flags().StringSliceVar(&o.Items, "items", []string{}, OptCommon.TransLang("param_console_get_items"))
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, OptCommon.TransLang("param_console_get_envs"))
	cmd.Flags().StringSliceVar(&o.BranchNames, "branches", []string{}, OptCommon.TransLang("param_console_get_branches"))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_console_get_output"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_console_get_full"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsConsoleGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	consoleCmdKinds := []string{}
	for k, _ := range pkg.ConsoleCmdKinds {
		consoleCmdKinds = append(consoleCmdKinds, k)
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			projectNames, err := o.GetConsoleProjectNames()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return projectNames, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			return consoleCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("envs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		projectConsole, err := o.GetConsoleProject(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		envNames := []string{}
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			envNames = append(envNames, pae.EnvName)
		}
		return envNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("branches", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		projectConsole, err := o.GetConsoleProject(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		branchNames := []string{}
		for _, pp := range projectConsole.Pipelines {
			branchNames = append(branchNames, pp.BranchName)
		}
		return branchNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("items", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items := []string{}
		projectName := args[0]
		kindStr := args[1]
		var isAllKind bool
		kinds := strings.Split(kindStr, ",")
		for _, kind := range kinds {
			if kind == pkg.ConsoleKindAll {
				isAllKind = true
			}
		}
		envs, _ := cmd.Flags().GetStringSlice("envs")
		branches, _ := cmd.Flags().GetStringSlice("branches")
		projectConsole, err := o.GetConsoleProject(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		m := map[string]string{}
		for _, kind := range kinds {
			if kind == pkg.ConsoleKindMember || isAllKind {
				for _, item := range projectConsole.ProjectMembers {
					m[item.Username] = ""
				}
			}
			if kind == pkg.ConsoleKindPipeline || isAllKind {
				for _, item := range projectConsole.Pipelines {
					m[item.BranchName] = ""
				}
			}
			if kind == pkg.ConsoleKindPipelineTrigger || isAllKind {
				if len(branches) == 0 {
					for _, pipeline := range projectConsole.Pipelines {
						for _, item := range pipeline.PipelineTriggers {
							m[item.StepAction] = ""
						}
					}
				} else {
					pipelines := []pkg.ProjectPipeline{}
					for _, pipeline := range projectConsole.Pipelines {
						for _, branch := range branches {
							if branch == pipeline.BranchName {
								pipelines = append(pipelines, pipeline)
								break
							}
						}
					}
					for _, pipeline := range pipelines {
						for _, item := range pipeline.PipelineTriggers {
							m[item.StepAction] = ""
						}
					}
				}
			}
			if kind == pkg.ConsoleKindHost || kind == pkg.ConsoleKindDatabase || kind == pkg.ConsoleKindComponent || isAllKind {
				if len(envs) == 0 {
					for _, pae := range projectConsole.ProjectAvailableEnvs {
						if kind == pkg.ConsoleKindHost || isAllKind {
							for _, item := range pae.Hosts {
								m[item.HostName] = ""
							}
						}
						if kind == pkg.ConsoleKindDatabase || isAllKind {
							for _, item := range pae.Databases {
								m[item.DbName] = ""
							}
						}
						if kind == pkg.ConsoleKindComponent || isAllKind {
							for _, item := range pae.Components {
								m[item.ComponentName] = ""
							}
						}
					}
				} else {
					paes := []pkg.ProjectAvailableEnvConsole{}
					for _, pae := range projectConsole.ProjectAvailableEnvs {
						for _, env := range envs {
							if env == pae.EnvName {
								paes = append(paes, pae)
								break
							}
						}
					}
					for _, pae := range paes {
						if kind == pkg.ConsoleKindHost || isAllKind {
							for _, item := range pae.Hosts {
								m[item.HostName] = ""
							}
						}
						if kind == pkg.ConsoleKindDatabase || isAllKind {
							for _, item := range pae.Databases {
								m[item.DbName] = ""
							}
						}
						if kind == pkg.ConsoleKindComponent || isAllKind {
							for _, item := range pae.Components {
								m[item.ComponentName] = ""
							}
						}
					}
				}
			}
			if isAllKind {
				break
			}
		}
		for k, _ := range m {
			items = append(items, k)
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsConsoleGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("projectName required")
		return err
	}
	var projectName string
	var kinds, kindParams []string
	var kindsStr string
	projectName = args[0]
	if len(args) == 1 {
		kindsStr = pkg.ConsoleKindAll
	} else {
		kindsStr = args[1]
	}
	arr := strings.Split(kindsStr, ",")
	for _, s := range arr {
		a := strings.Trim(s, " ")
		if a != "" {
			kinds = append(kinds, a)
		}
	}
	var foundAll bool
	for _, kind := range kinds {
		var found bool
		for cmdKind, _ := range pkg.ConsoleCmdKinds {
			if kind == cmdKind {
				found = true
				break
			}
		}
		if !found {
			defCmdKinds := []string{}
			for k, _ := range pkg.ConsoleCmdKinds {
				defCmdKinds = append(defCmdKinds, k)
			}
			err = fmt.Errorf("kind %s format error: not correct, options: %s", kind, strings.Join(defCmdKinds, " / "))
			return err
		}
		if kind == pkg.ConsoleKindAll {
			foundAll = true
		}
		kindParams = append(kindParams, pkg.ConsoleCmdKinds[kind])
	}
	if foundAll == true {
		o.Param.IsAllKind = true
	}
	o.Param.Kinds = kindParams

	err = pkg.ValidateMinusNameID(projectName)
	if err != nil {
		err = fmt.Errorf("projectName %s format error: %s", projectName, err.Error())
		return err
	}

	o.Param.ProjectName = projectName

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsConsoleGet) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	projectConsole, err := o.GetConsoleProject(o.Param.ProjectName)
	if err != nil {
		return err
	}

	consoleKinds := []pkg.ConsoleKind{}
	consoleKindProject := pkg.ConsoleKind{
		Kind: "",
		Metadata: pkg.ConsoleMetadata{
			ProjectName: projectConsole.ProjectInfo.ProjectName,
			EnvName:     "",
			BranchName:  "",
		},
		Items: []interface{}{},
	}
	if len(projectConsole.ProjectMembers) > 0 {
		ck := consoleKindProject
		ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindMember]
		for _, item := range projectConsole.ProjectMembers {
			var isShow bool
			if len(o.Items) == 0 {
				isShow = true
			} else {
				for _, itemName := range o.Items {
					if itemName == item.Username {
						isShow = true
						break
					}
				}
			}
			if isShow {
				ck.Items = append(ck.Items, item)
			}
		}
		consoleKinds = append(consoleKinds, ck)
	}

	if len(projectConsole.Pipelines) > 0 {
		ck := consoleKindProject
		ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]
		for _, item := range projectConsole.Pipelines {
			item.PipelineTriggers = []pkg.PipelineTrigger{}
			item.PipelineCrons = []pkg.PipelineCron{}
			var isShow bool
			if len(o.Items) == 0 {
				isShow = true
			} else {
				for _, itemName := range o.Items {
					if itemName == item.BranchName {
						isShow = true
						break
					}
				}
			}
			if isShow {
				ck.Items = append(ck.Items, item)
			}
		}
		consoleKinds = append(consoleKinds, ck)
	}

	if len(projectConsole.Pipelines) > 0 {
		pipelines := []pkg.ProjectPipeline{}
		for _, pipeline := range projectConsole.Pipelines {
			if len(o.BranchNames) == 0 {
				pipelines = append(pipelines, pipeline)
			} else {
				for _, branchName := range o.BranchNames {
					if branchName == pipeline.BranchName {
						pipelines = append(pipelines, pipeline)
						break
					}
				}
			}
		}
		for _, pipeline := range pipelines {
			if len(pipeline.PipelineTriggers) > 0 {
				ck := consoleKindProject
				ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindPipelineTrigger]
				ck.Metadata.BranchName = pipeline.BranchName
				for _, item := range pipeline.PipelineTriggers {
					var isShow bool
					if len(o.Items) == 0 {
						isShow = true
					} else {
						for _, itemName := range o.Items {
							if itemName == item.StepAction {
								isShow = true
								break
							}
						}
					}
					if isShow {
						ck.Items = append(ck.Items, item)
					}
				}
				consoleKinds = append(consoleKinds, ck)
			}
		}
	}

	if len(projectConsole.ProjectAvailableEnvs) > 0 {
		paes := []pkg.ProjectAvailableEnvConsole{}
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			if len(o.EnvNames) == 0 {
				paes = append(paes, pae)
			} else {
				for _, envName := range o.EnvNames {
					if envName == pae.EnvName {
						paes = append(paes, pae)
						break
					}
				}
			}
		}

		for _, pae := range paes {
			if len(pae.Hosts) > 0 {
				ck := consoleKindProject
				ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindHost]
				ck.Metadata.EnvName = pae.EnvName
				for _, item := range pae.Hosts {
					var isShow bool
					if len(o.Items) == 0 {
						isShow = true
					} else {
						for _, itemName := range o.Items {
							if itemName == item.HostName {
								isShow = true
								break
							}
						}
					}
					if isShow {
						ck.Items = append(ck.Items, item)
					}
				}
				consoleKinds = append(consoleKinds, ck)
			}
		}

		for _, pae := range paes {
			if len(pae.Databases) > 0 {
				ck := consoleKindProject
				ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindDatabase]
				ck.Metadata.EnvName = pae.EnvName
				for _, item := range pae.Databases {
					var isShow bool
					if len(o.Items) == 0 {
						isShow = true
					} else {
						for _, itemName := range o.Items {
							if itemName == item.DbName {
								isShow = true
								break
							}
						}
					}
					if isShow {
						ck.Items = append(ck.Items, item)
					}
				}
				consoleKinds = append(consoleKinds, ck)
			}
		}

		for _, pae := range paes {
			if len(pae.Components) > 0 {
				ck := consoleKindProject
				ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindComponent]
				ck.Metadata.EnvName = pae.EnvName
				for _, item := range pae.Components {
					var isShow bool
					if len(o.Items) == 0 {
						isShow = true
					} else {
						for _, itemName := range o.Items {
							if itemName == item.ComponentName {
								isShow = true
								break
							}
						}
					}
					if isShow {
						ck.Items = append(ck.Items, item)
					}
				}
				consoleKinds = append(consoleKinds, ck)
			}
		}

		for _, pae := range paes {
			if pae.ComponentDebug.Arch != "" {
				ck := consoleKindProject
				ck.Kind = pkg.ConsoleCmdKinds[pkg.ConsoleKindDebugComponent]
				ck.Metadata.EnvName = pae.EnvName
				ck.Items = append(ck.Items, pae.ComponentDebug)
				consoleKinds = append(consoleKinds, ck)
			}
		}

	}

	consoleKindFilters := []pkg.ConsoleKind{}
	if len(o.Param.Kinds) == 0 || o.Param.IsAllKind {
		consoleKindFilters = consoleKinds
	} else {
		for _, consoleKind := range consoleKinds {
			for _, kind := range o.Param.Kinds {
				if kind == consoleKind.Kind {
					consoleKindFilters = append(consoleKindFilters, consoleKind)
					break
				}
			}
		}
	}

	consoleKindList := pkg.ConsoleKindList{
		Kind:     "list",
		Consoles: consoleKindFilters,
	}

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(consoleKindList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		dataOutput = m
	} else {
		dataOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(dataOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = pkg.YamlIndent(dataOutput)
		fmt.Println(string(bs))
	default:
		for _, consoleKind := range consoleKindList.Consoles {

			dataHeader := []string{}
			dataRows := [][]string{}
			bs, _ := json.Marshal(consoleKind.Items)
			switch consoleKind.Kind {
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindMember]:
				items := []pkg.ProjectMember{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, item.Username), item.AccessLevel, strings.Join(item.DisablePipelines, ","), strings.Join(item.DisableRepoSecrets, ","), strings.Join(item.DisableProjectDefs, ",")}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "AccessLevel", "DisablePipelines", "DisableRepoSecrets", "DisableProjectDefs"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]:
				items := []pkg.ProjectPipeline{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s-%s", consoleKind.Kind, consoleKind.Metadata.ProjectName, item.BranchName), strings.Join(item.Envs, ","), strings.Join(item.EnvProductions, ","), fmt.Sprintf("%v", item.IsDefault), fmt.Sprintf("%v", item.WebhookPushEvent)}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "TestEnvs", "ProdEnvs", "IsDefault", "WebhookPushEvent"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipelineTrigger]:
				items := []pkg.PipelineTrigger{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, item.StepAction), fmt.Sprintf("%s-%s", consoleKind.Metadata.ProjectName, consoleKind.Metadata.BranchName), strings.Join(item.StatusResults, ","), fmt.Sprintf("%v", item.Enable), fmt.Sprintf("%v", item.BeforeExecute)}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "Pipeline", "StatusResults", "Enable", "BeforeExecute"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindHost]:
				items := []pkg.Host{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, item.HostName), consoleKind.Metadata.EnvName, fmt.Sprintf("%s:%d", item.HostAddr, item.HostPort), item.HostUser, strings.Join(item.Groups, ",")}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "Env", "Address", "User", "Groups"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindDatabase]:
				items := []pkg.Database{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, item.DbName), consoleKind.Metadata.EnvName, item.DbUrl, item.DbUser}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "Env", "URL", "User"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindComponent]:
				items := []pkg.Component{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, item.ComponentName), consoleKind.Metadata.EnvName, item.ComponentDesc, item.Arch, item.DeploySpecStatic.DeployImage}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "Env", "Desc", "Arch", "Image"}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindDebugComponent]:
				items := []pkg.ComponentDebug{}
				_ = json.Unmarshal(bs, &items)
				for _, item := range items {
					dataRow := []string{fmt.Sprintf("%s/%s", consoleKind.Kind, consoleKind.Metadata.EnvName), consoleKind.Metadata.EnvName, item.Arch, item.DeploySpecDebug.DebugQuota.CpuLimit, item.DeploySpecDebug.DebugQuota.MemoryLimit}
					dataRows = append(dataRows, dataRow)
				}
				dataHeader = []string{"Name", "Env", "Arch", "CPU", "MEM"}
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(dataHeader)
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
			table.AppendBulk(dataRows)
			table.Render()
			fmt.Println("------------")
			fmt.Println()
		}
	}

	return err
}
