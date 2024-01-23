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

type OptionsDefClone struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FromEnvName    string   `yaml:"fromEnvName" json:"fromEnvName" bson:"fromEnvName" validate:""`
	StepName       string   `yaml:"stepName" json:"stepName" bson:"stepName" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	ToEnvNames     []string `yaml:"toEnvNames" json:"toEnvNames" bson:"toEnvNames" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind        string `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	}
}

func NewOptionsDefClone() *OptionsDefClone {
	var o OptionsDefClone
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefClone() *cobra.Command {
	o := NewOptionsDefClone()

	defCmdKinds := []string{
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomStep,
	}
	sort.Strings(defCmdKinds)

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`clone [projectName] [kind] [--from-env=envName] [--step=stepName] [--modules=moduleName1,moduleName2] [--to-envs=envName1,envName2] [--output=json|yaml]`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_def_clone_short")
	msgLong := OptCommon.TransLang("cmd_def_clone_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_def_clone_example", strings.Join(defCmdKinds, " / "), baseName, pkg.DefKindDeployContainer, baseName, pkg.DefKindCustomStep))

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
	cmd.Flags().StringVar(&o.FromEnvName, "from-env", "", OptCommon.TransLang("param_def_clone_from_env"))
	cmd.Flags().StringVar(&o.StepName, "step", "", OptCommon.TransLang("param_def_clone_step", pkg.DefKindCustomStep))
	cmd.Flags().StringSliceVar(&o.ModuleNames, "modules", []string{}, OptCommon.TransLang("param_def_clone_modules"))
	cmd.Flags().StringSliceVar(&o.ToEnvNames, "to-envs", []string{}, OptCommon.TransLang("param_def_clone_to_envs"))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_def_clone_output"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_def_clone_full"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_def_clone_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsDefClone) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	defCmdKinds := []string{
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomStep,
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectNames, err := o.GetProjectNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if len(args) == 0 {
			return projectNames, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			return defCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("from-env", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		envNames := []string{}
		for _, pae := range project.ProjectAvailableEnvs {
			envNames = append(envNames, pae.EnvName)
		}
		return envNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("to-envs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		envNames := []string{}
		for _, pae := range project.ProjectAvailableEnvs {
			envNames = append(envNames, pae.EnvName)
		}
		return envNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("step", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		stepNames := []string{}
		for _, conf := range project.CustomStepConfs {
			if conf.IsEnvDiff {
				stepNames = append(stepNames, conf.CustomStepName)
			}
		}
		return stepNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("modules", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		moduleNames := []string{}
		projectName := args[0]
		kind := args[1]
		step, _ := cmd.Flags().GetString("step")
		envName, _ := cmd.Flags().GetString("from-env")
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		switch kind {
		case pkg.DefKindDeployContainer:
			m := map[string]string{}
			if envName == "" {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.DeployContainerDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				pae := pkg.ProjectAvailableEnv{}
				for _, p := range project.ProjectAvailableEnvs {
					if envName == p.EnvName {
						pae = p
						break
					}
				}
				for _, def := range pae.DeployContainerDefs {
					m[def.DeployName] = def.DeployName
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindDeployArtifact:
			m := map[string]string{}
			if envName == "" {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.DeployArtifactDefs {
						m[def.DeployArtifactName] = def.DeployArtifactName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				pae := pkg.ProjectAvailableEnv{}
				for _, p := range project.ProjectAvailableEnvs {
					if envName == p.EnvName {
						pae = p
						break
					}
				}
				for _, def := range pae.DeployArtifactDefs {
					m[def.DeployArtifactName] = def.DeployArtifactName
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindIstio:
			m := map[string]string{}
			if envName == "" {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.IstioDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				pae := pkg.ProjectAvailableEnv{}
				for _, p := range project.ProjectAvailableEnvs {
					if envName == p.EnvName {
						pae = p
						break
					}
				}
				for _, def := range pae.IstioDefs {
					m[def.DeployName] = def.DeployName
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindCustomStep:
			if step != "" {
				if envName == "" {
					for stepName, csd := range project.ProjectDef.CustomStepDefs {
						if stepName == step {
							for _, def := range csd.CustomStepModuleDefs {
								moduleNames = append(moduleNames, def.ModuleName)
							}
							break
						}
					}
				} else {
					m := map[string]string{}
					pae := pkg.ProjectAvailableEnv{}
					for _, p := range project.ProjectAvailableEnvs {
						if envName == p.EnvName {
							pae = p
							break
						}
					}
					for stepName, csd := range pae.CustomStepDefs {
						if stepName == step {
							for _, def := range csd.CustomStepModuleDefs {
								m[def.ModuleName] = def.ModuleName
							}
							break
						}
					}
					for k, _ := range m {
						moduleNames = append(moduleNames, k)
					}
				}
			}
		}
		return moduleNames, cobra.ShellCompDirectiveNoFileComp
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

	err = cmd.MarkFlagRequired("from-env")
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("to-envs")
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("modules")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsDefClone) Validate(args []string) error {
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

	defCmdKinds := []string{
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomStep,
	}
	var found bool
	for _, cmdKind := range defCmdKinds {
		if kind == cmdKind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct: kind options: %s", kind, strings.Join(defCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	if len(o.ModuleNames) == 0 {
		err = fmt.Errorf("--modules required")
		return err
	}

	if o.FromEnvName == "" {
		err = fmt.Errorf("--from-env required")
		return err
	}

	if len(o.ToEnvNames) == 0 {
		err = fmt.Errorf("--to-envs required")
		return err
	}

	if o.Param.Kind == pkg.DefKindCustomStep && o.StepName == "" {
		err = fmt.Errorf("kind is %s, --step required", pkg.DefKindCustomStep)
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

func (o *OptionsDefClone) Run(args []string) error {
	var err error

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", o.Param.ProjectName), http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	project := pkg.ProjectOutput{}
	err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
	if err != nil {
		return err
	}

	for _, envName := range o.ToEnvNames {
		var found bool
		for _, pae := range project.ProjectAvailableEnvs {
			if envName == pae.EnvName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("to envName %s not exists", envName)
			return err
		}
	}

	var defClone pkg.DefClone
	switch o.Param.Kind {
	case pkg.DefKindDeployContainer:
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("from envName %s not exists", o.FromEnvName)
			return err
		}
		defs := []pkg.DeployContainerDef{}
		for _, def := range pae.DeployContainerDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.DeployName == moduleName {
					found = true
					break
				}
			}
			if found {
				defs = append(defs, def)
			}
		}
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = defs
	case pkg.DefKindDeployArtifact:
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("from envName %s not exists", o.FromEnvName)
			return err
		}
		defs := []pkg.DeployArtifactDef{}
		for _, def := range pae.DeployArtifactDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.DeployArtifactName == moduleName {
					found = true
					break
				}
			}
			if found {
				defs = append(defs, def)
			}
		}
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = defs
	case pkg.DefKindIstio:
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("from envName %s not exists", o.FromEnvName)
			return err
		}
		defs := []pkg.IstioDef{}
		for _, def := range pae.IstioDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.DeployName == moduleName {
					found = true
					break
				}
			}
			if found {
				defs = append(defs, def)
			}
		}
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = defs
	case pkg.DefKindCustomStep:
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("from envName %s not exists", o.FromEnvName)
			return err
		}

		csd := pkg.CustomStepDef{}
		for stepName, c := range pae.CustomStepDefs {
			if o.StepName == stepName {
				csd = c
				break
			}
		}
		defs := []pkg.CustomStepModuleDef{}
		for _, def := range csd.CustomStepModuleDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.ModuleName == moduleName {
					found = true
					break
				}
			}
			if found {
				defs = append(defs, def)
			}
		}
		csd.CustomStepModuleDefs = defs
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = csd
	}

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ := json.Marshal(defClone)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		dataOutput = m
	} else {
		dataOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ := json.MarshalIndent(dataOutput["def"], "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ := pkg.YamlIndent(dataOutput["def"])
		fmt.Println(string(bs))
	}

	if !o.Try {
		bs, _ = pkg.YamlIndent(dataOutput["def"])
		urlKind := defClone.Kind
		param["envNames"] = o.ToEnvNames
		switch defClone.Kind {
		case "deployContainerDefs":
			param["deployContainerDefsYaml"] = string(bs)
		case "deployArtifactDefs":
			param["deployArtifactDefsYaml"] = string(bs)
		case "istioDefs":
			param["istioDefsYaml"] = string(bs)
		case "customStepDef":
			urlKind = fmt.Sprintf("%s/env", urlKind)
			param["customStepName"] = o.StepName
			param["customStepDefYaml"] = string(bs)
		}
		logHeader := fmt.Sprintf("[%s/%s]", defClone.ProjectName, defClone.Kind)
		result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defClone.ProjectName, urlKind), http.MethodPut, "", param, false)
		if err != nil {
			err = fmt.Errorf("%s: %s", logHeader, err.Error())
			return err
		}
		msg := result.Get("msg").String()
		log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
	}

	return err
}
