package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"sort"
	"strings"
)

type OptionsConsoleDelete struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Items          []string `yaml:"items" json:"items" bson:"items" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	BranchNames    []string `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind        string `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	}
}

func NewOptionsConsoleDelete() *OptionsConsoleDelete {
	var o OptionsConsoleDelete
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdConsoleDelete() *cobra.Command {
	o := NewOptionsConsoleDelete()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`delete [projectName] [kind] [--envs=envName1,envName2] [--branches=branch1,branch2] [--items=itemName1,itemName2]...`)

	consoleCmdKinds := []string{}
	for k, _ := range pkg.ConsoleCmdKinds {
		if k != pkg.ConsoleKindAll {
			consoleCmdKinds = append(consoleCmdKinds, k)
		}
	}
	sort.Strings(consoleCmdKinds)
	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_console_delete_short")
	msgLong := OptCommon.TransLang("cmd_console_delete_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_console_delete_example", strings.Join(consoleCmdKinds, ", "), baseName, pkg.ConsoleKindMember, baseName, pkg.ConsoleKindHost, baseName, pkg.ConsoleKindComponent, baseName, pkg.ConsoleKindPipelineTrigger))

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

	cmd.Flags().StringSliceVar(&o.Items, "items", []string{}, OptCommon.TransLang("param_console_delete_items"))
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, OptCommon.TransLang("param_console_delete_envs", pkg.ConsoleKindHost, pkg.ConsoleKindDatabase, pkg.ConsoleKindComponent, pkg.ConsoleKindDebugComponent))
	cmd.Flags().StringSliceVar(&o.BranchNames, "branches", []string{}, OptCommon.TransLang("param_console_delete_branches", pkg.ConsoleKindPipelineTrigger))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_console_delete_output"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_console_delete_full"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_console_delete_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsConsoleDelete) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	consoleCmdKinds := []string{}
	for k, _ := range pkg.ConsoleCmdKinds {
		if k != pkg.ConsoleKindAll {
			consoleCmdKinds = append(consoleCmdKinds, k)
		}
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

func (o *OptionsConsoleDelete) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("projectName required")
		return err
	}
	if len(args) == 1 {
		err = fmt.Errorf("kind required")
		return err
	}
	var projectName string
	var kind string
	projectName = args[0]
	kind = args[1]

	err = pkg.ValidateMinusNameID(projectName)
	if err != nil {
		err = fmt.Errorf("projectName %s format error: %s", projectName, err.Error())
		return err
	}

	o.Param.ProjectName = projectName

	consoleCmdKinds := []string{}
	for k, _ := range pkg.ConsoleCmdKinds {
		if k != pkg.ConsoleKindAll {
			consoleCmdKinds = append(consoleCmdKinds, k)
		}
	}
	sort.Strings(consoleCmdKinds)

	var found bool
	for k, _ := range pkg.ConsoleCmdKinds {
		if kind == k {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct: kind options: %s", kind, strings.Join(consoleCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	if len(o.Items) == 0 && o.Param.Kind != pkg.ConsoleKindDebugComponent {
		err = fmt.Errorf("--items required")
		return err
	}

	if o.Param.Kind == pkg.ConsoleKindHost && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.ConsoleKindHost)
		return err
	}
	if o.Param.Kind == pkg.ConsoleKindDatabase && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.ConsoleKindDatabase)
		return err
	}
	if o.Param.Kind == pkg.ConsoleKindDebugComponent && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.ConsoleKindDebugComponent)
		return err
	}
	if o.Param.Kind == pkg.ConsoleKindComponent && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.ConsoleKindComponent)
		return err
	}
	if o.Param.Kind == pkg.ConsoleKindPipelineTrigger && len(o.BranchNames) == 0 {
		err = fmt.Errorf("kind is %s, --branches required", pkg.ConsoleKindPipelineTrigger)
		return err
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}

	return err
}

func (o *OptionsConsoleDelete) Run(args []string) error {
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

	switch o.Param.Kind {
	case pkg.ConsoleKindMember:
		consoleKind := consoleKindProject
		consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
		for _, item := range projectConsole.ProjectMembers {
			var found bool
			for _, itemName := range o.Items {
				if item.Username == itemName {
					found = true
					break
				}
			}
			if found {
				consoleKind.Items = append(consoleKind.Items, item)
			}
		}
		consoleKinds = append(consoleKinds, consoleKind)
	case pkg.ConsoleKindPipeline:
		consoleKind := consoleKindProject
		consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
		for _, item := range projectConsole.Pipelines {
			var found bool
			for _, itemName := range o.Items {
				if item.BranchName == itemName {
					found = true
					break
				}
			}
			if found {
				consoleKind.Items = append(consoleKind.Items, item)
			}
		}
		consoleKinds = append(consoleKinds, consoleKind)
	case pkg.ConsoleKindPipelineTrigger:
		for _, pipeline := range projectConsole.Pipelines {
			consoleKind := consoleKindProject
			consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
			consoleKind.Metadata.BranchName = pipeline.BranchName
			var found bool
			for _, branchName := range o.BranchNames {
				if pipeline.BranchName == branchName {
					found = true
					break
				}
			}
			if found {
				for _, item := range pipeline.PipelineTriggers {
					for _, itemName := range o.Items {
						if item.StepAction == itemName {
							consoleKind.Items = append(consoleKind.Items, item)
							break
						}
					}
				}
			}
			if len(consoleKind.Items) > 0 {
				consoleKinds = append(consoleKinds, consoleKind)
			}
		}
	case pkg.ConsoleKindHost:
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			consoleKind := consoleKindProject
			consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
			consoleKind.Metadata.EnvName = pae.EnvName
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				for _, item := range pae.Hosts {
					for _, itemName := range o.Items {
						if item.HostName == itemName {
							consoleKind.Items = append(consoleKind.Items, item)
							break
						}
					}
				}
			}
			if len(consoleKind.Items) > 0 {
				consoleKinds = append(consoleKinds, consoleKind)
			}
		}
	case pkg.ConsoleKindDatabase:
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			consoleKind := consoleKindProject
			consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
			consoleKind.Metadata.EnvName = pae.EnvName
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				for _, item := range pae.Databases {
					for _, itemName := range o.Items {
						if item.DbName == itemName {
							consoleKind.Items = append(consoleKind.Items, item)
							break
						}
					}
				}
			}
			if len(consoleKind.Items) > 0 {
				consoleKinds = append(consoleKinds, consoleKind)
			}
		}
	case pkg.ConsoleKindComponent:
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			consoleKind := consoleKindProject
			consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
			consoleKind.Metadata.EnvName = pae.EnvName
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				for _, item := range pae.Components {
					for _, itemName := range o.Items {
						if item.ComponentName == itemName {
							consoleKind.Items = append(consoleKind.Items, item)
							break
						}
					}
				}
			}
			if len(consoleKind.Items) > 0 {
				consoleKinds = append(consoleKinds, consoleKind)
			}
		}
	case pkg.ConsoleKindDebugComponent:
		for _, pae := range projectConsole.ProjectAvailableEnvs {
			consoleKind := consoleKindProject
			consoleKind.Kind = pkg.ConsoleCmdKinds[o.Param.Kind]
			consoleKind.Metadata.EnvName = pae.EnvName
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				consoleKind.Items = append(consoleKind.Items, pae.ComponentDebug)
			}
			if len(consoleKind.Items) > 0 {
				consoleKinds = append(consoleKinds, consoleKind)
			}
		}
	}

	consoleKindList := pkg.ConsoleKindList{
		Kind:     "list",
		Consoles: consoleKinds,
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
		bs, _ := json.MarshalIndent(dataOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ := pkg.YamlIndent(dataOutput)
		fmt.Println(string(bs))
	}

	if !o.Try {
		for _, consoleKind := range consoleKinds {
			for _, item := range consoleKind.Items {
				param := map[string]interface{}{}
				var url string
				var method string
				switch consoleKind.Kind {
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindMember]:
					param["username"] = item.(pkg.ProjectMember).Username
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/memberDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]:
					param["branchName"] = item.(pkg.ProjectPipeline).BranchName
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/pipelineDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipelineTrigger]:
					param["branchName"] = consoleKind.Metadata.BranchName
					param["stepAction"] = item.(pkg.PipelineTrigger).StepAction
					param["beforeExecute"] = item.(pkg.PipelineTrigger).BeforeExecute
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/pipelineTriggerDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindHost]:
					param["envName"] = consoleKind.Metadata.EnvName
					param["hostName"] = item.(pkg.Host).HostName
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/envHostDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindDatabase]:
					param["envName"] = consoleKind.Metadata.EnvName
					param["dbName"] = item.(pkg.Database).DbName
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/envDbDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindComponent]:
					param["envName"] = consoleKind.Metadata.EnvName
					param["componentName"] = item.(pkg.Component).ComponentName
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/envComponentDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindDebugComponent]:
					param["envName"] = consoleKind.Metadata.EnvName
					bs, _ = json.Marshal(param)
					url = fmt.Sprintf("api/console/project/%s/envComponentDebugDelete", consoleKind.Metadata.ProjectName)
					method = http.MethodPost
				}
				logHeader := fmt.Sprintf("[%s/%s] %s", consoleKind.Metadata.ProjectName, consoleKind.Kind, string(bs))
				result, _, err := o.QueryAPI(url, method, "", param, false)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				msg := result.Get("msg").String()
				log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
				switch consoleKind.Kind {
				case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]:
					auditID := result.Get("data.auditID").String()
					if auditID == "" {
						err = fmt.Errorf("can not get auditID")
						return err
					}
					url = fmt.Sprintf("api/ws/log/audit/console/%s", auditID)
					err = o.QueryWebsocket(url, "")
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("##############################"))
					log.Success(fmt.Sprintf("# %s delete finish", logHeader))
				}
			}
		}
	}

	return err
}
