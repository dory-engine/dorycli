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

type OptionsDefDelete struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	StepNames      []string `yaml:"stepNames" json:"stepNames" bson:"stepNames" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind        string `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	}
}

func NewOptionsDefDelete() *OptionsDefDelete {
	var o OptionsDefDelete
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefDelete() *cobra.Command {
	o := NewOptionsDefDelete()

	defCmdKinds := []string{
		pkg.DefKindBuild,
		pkg.DefKindPackage,
		pkg.DefKindArtifact,
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
		pkg.DefKindCustomStep,
	}
	sort.Strings(defCmdKinds)

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`delete [projectName] [kind] [--modules=moduleName1,moduleName2] [--envs=envName1,envName2] [--steps=stepName1,stepName2] [--output=json|yaml]`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_def_delete_short")
	msgLong := OptCommon.TransLang("cmd_def_delete_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_def_delete_example", strings.Join(defCmdKinds, " / "), baseName, pkg.DefKindBuild, baseName, pkg.DefKindDeployContainer, baseName, pkg.DefKindCustomStep, baseName, pkg.DefKindCustomStep))

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
	cmd.Flags().StringSliceVar(&o.ModuleNames, "modules", []string{}, OptCommon.TransLang("param_def_delete_modules"))
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, OptCommon.TransLang("param_def_delete_envs", pkg.DefKindDeployContainer, pkg.DefKindDeployArtifact, pkg.DefKindIstio))
	cmd.Flags().StringSliceVar(&o.StepNames, "steps", []string{}, OptCommon.TransLang("param_def_delete_steps", pkg.DefKindCustomStep))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_def_delete_output"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_def_delete_full"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_def_delete_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsDefDelete) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	defCmdKinds := []string{
		pkg.DefKindBuild,
		pkg.DefKindPackage,
		pkg.DefKindArtifact,
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
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

	err = cmd.RegisterFlagCompletionFunc("envs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

	err = cmd.RegisterFlagCompletionFunc("steps", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		stepNames := []string{}
		for _, conf := range project.CustomStepConfs {
			stepNames = append(stepNames, conf.CustomStepName)
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
		envs, _ := cmd.Flags().GetStringSlice("envs")
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		switch kind {
		case pkg.DefKindBuild:
			for _, def := range project.ProjectDef.BuildDefs {
				moduleNames = append(moduleNames, def.BuildName)
			}
		case pkg.DefKindPackage:
			for _, def := range project.ProjectDef.PackageDefs {
				moduleNames = append(moduleNames, def.PackageName)
			}
		case pkg.DefKindArtifact:
			for _, def := range project.ProjectDef.ArtifactDefs {
				moduleNames = append(moduleNames, def.ArtifactName)
			}
		case pkg.DefKindDeployContainer:
			m := map[string]string{}
			if len(envs) == 0 {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.DeployContainerDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				paes := []pkg.ProjectAvailableEnv{}
				for _, pae := range project.ProjectAvailableEnvs {
					for _, env := range envs {
						if env == pae.EnvName {
							paes = append(paes, pae)
							break
						}
					}
				}
				for _, pae := range paes {
					for _, def := range pae.DeployContainerDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindDeployArtifact:
			m := map[string]string{}
			if len(envs) == 0 {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.DeployArtifactDefs {
						m[def.DeployArtifactName] = def.DeployArtifactName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				paes := []pkg.ProjectAvailableEnv{}
				for _, pae := range project.ProjectAvailableEnvs {
					for _, env := range envs {
						if env == pae.EnvName {
							paes = append(paes, pae)
							break
						}
					}
				}
				for _, pae := range paes {
					for _, def := range pae.DeployArtifactDefs {
						m[def.DeployArtifactName] = def.DeployArtifactName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindIstio:
			m := map[string]string{}
			if len(envs) == 0 {
				for _, pae := range project.ProjectAvailableEnvs {
					for _, def := range pae.IstioDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			} else {
				paes := []pkg.ProjectAvailableEnv{}
				for _, pae := range project.ProjectAvailableEnvs {
					for _, env := range envs {
						if env == pae.EnvName {
							paes = append(paes, pae)
							break
						}
					}
				}
				for _, pae := range paes {
					for _, def := range pae.IstioDefs {
						m[def.DeployName] = def.DeployName
					}
				}
				for k, _ := range m {
					moduleNames = append(moduleNames, k)
				}
			}
		case pkg.DefKindCustomOps:
			for _, def := range project.ProjectDef.CustomOpsDefs {
				moduleNames = append(moduleNames, def.CustomOpsName)
			}
		case pkg.DefKindOpsBatch:
			for _, def := range project.ProjectDef.OpsBatchDefs {
				moduleNames = append(moduleNames, def.OpsBatchName)
			}
		case pkg.DefKindCustomStep:
			if step != "" {
				if len(envs) == 0 {
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
					paes := []pkg.ProjectAvailableEnv{}
					for _, pae := range project.ProjectAvailableEnvs {
						for _, env := range envs {
							if env == pae.EnvName {
								paes = append(paes, pae)
								break
							}
						}
					}
					for _, pae := range paes {
						for stepName, csd := range pae.CustomStepDefs {
							if stepName == step {
								for _, def := range csd.CustomStepModuleDefs {
									m[def.ModuleName] = def.ModuleName
								}
								break
							}
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

	err = cmd.MarkFlagRequired("modules")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsDefDelete) Validate(args []string) error {
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
		pkg.DefKindBuild,
		pkg.DefKindPackage,
		pkg.DefKindArtifact,
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
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

	if o.Param.Kind == pkg.DefKindDeployContainer && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindDeployContainer)
		return err
	}
	if o.Param.Kind == pkg.DefKindDeployArtifact && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindDeployArtifact)
		return err
	}
	if o.Param.Kind == pkg.DefKindIstio && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindIstio)
		return err
	}
	if o.Param.Kind == pkg.DefKindCustomStep && len(o.StepNames) == 0 {
		err = fmt.Errorf("kind is %s, --steps required", pkg.DefKindCustomStep)
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

func (o *OptionsDefDelete) Run(args []string) error {
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

	defKinds := []pkg.DefKind{}
	defUpdates := []pkg.DefUpdate{}
	defKindProject := pkg.DefKind{
		Kind: "",
		Metadata: pkg.DefMetadata{
			ProjectName: project.ProjectInfo.ProjectName,
			Labels:      map[string]string{},
		},
		Items: []interface{}{},
	}

	switch o.Param.Kind {
	case pkg.DefKindBuild:
		defKind := defKindProject
		defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
		ids := []int{}
		for i, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
		}
		defUpdates = append(defUpdates, defUpdate)
	case pkg.DefKindPackage:
		defKind := defKindProject
		defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defKind.Status.ErrMsg = project.ProjectDef.ErrMsgPackageDefs
		ids := []int{}
		for i, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
		}
		defUpdates = append(defUpdates, defUpdate)
	case pkg.DefKindArtifact:
		defKind := defKindProject
		defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defKind.Status.ErrMsg = project.ProjectDef.ErrMsgArtifactDefs
		ids := []int{}
		for i, def := range project.ProjectDef.ArtifactDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.ArtifactName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.ArtifactDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
		}
		defUpdates = append(defUpdates, defUpdate)
	case pkg.DefKindDeployContainer:
		paes := []pkg.ProjectAvailableEnv{}
		for _, pae := range project.ProjectAvailableEnvs {
			for _, envName := range o.EnvNames {
				if envName == pae.EnvName {
					paes = append(paes, pae)
					break
				}
			}
		}
		for _, pae := range paes {
			if len(pae.DeployContainerDefs) > 0 {
				defKind := defKindProject
				defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
				defKind.Status.ErrMsg = pae.ErrMsgDeployContainerDefs
				defKind.Metadata.Labels = map[string]string{
					"envName": pae.EnvName,
				}
				ids := []int{}
				for i, def := range pae.DeployContainerDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if found {
						ids = append(ids, i)
					}
				}
				for i, def := range pae.DeployContainerDefs {
					var found bool
					for _, id := range ids {
						if i == id {
							found = true
							break
						}
					}
					if !found {
						defKind.Items = append(defKind.Items, def)
					}
				}

				defKinds = append(defKinds, defKind)

				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         defKind.Items,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	case pkg.DefKindDeployArtifact:
		paes := []pkg.ProjectAvailableEnv{}
		for _, pae := range project.ProjectAvailableEnvs {
			for _, envName := range o.EnvNames {
				if envName == pae.EnvName {
					paes = append(paes, pae)
					break
				}
			}
		}
		for _, pae := range paes {
			if len(pae.DeployArtifactDefs) > 0 {
				defKind := defKindProject
				defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
				defKind.Status.ErrMsg = pae.ErrMsgDeployArtifactDefs
				defKind.Metadata.Labels = map[string]string{
					"envName": pae.EnvName,
				}
				ids := []int{}
				for i, def := range pae.DeployArtifactDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployArtifactName == moduleName {
							found = true
							break
						}
					}
					if found {
						ids = append(ids, i)
					}
				}
				for i, def := range pae.DeployArtifactDefs {
					var found bool
					for _, id := range ids {
						if i == id {
							found = true
							break
						}
					}
					if !found {
						defKind.Items = append(defKind.Items, def)
					}
				}

				defKinds = append(defKinds, defKind)

				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         defKind.Items,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	case pkg.DefKindIstio:
		paes := []pkg.ProjectAvailableEnv{}
		for _, pae := range project.ProjectAvailableEnvs {
			for _, envName := range o.EnvNames {
				if envName == pae.EnvName {
					paes = append(paes, pae)
					break
				}
			}
		}
		for _, pae := range paes {
			if len(pae.IstioDefs) > 0 {
				defKind := defKindProject
				defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
				defKind.Status.ErrMsg = pae.ErrMsgIstioDefs
				defKind.Metadata.Labels = map[string]string{
					"envName": pae.EnvName,
				}
				ids := []int{}
				for i, def := range pae.IstioDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if found {
						ids = append(ids, i)
					}
				}
				for i, def := range pae.IstioDefs {
					var found bool
					for _, id := range ids {
						if i == id {
							found = true
							break
						}
					}
					if !found {
						defKind.Items = append(defKind.Items, def)
					}
				}

				defKinds = append(defKinds, defKind)

				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         defKind.Items,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	case pkg.DefKindCustomOps:
		defKind := defKindProject
		defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
		ids := []int{}
		for i, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
		}
		defUpdates = append(defUpdates, defUpdate)
	case pkg.DefKindOpsBatch:
		defKind := defKindProject
		defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
		ids := []int{}
		for i, def := range project.ProjectDef.OpsBatchDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.OpsBatchName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.OpsBatchDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
		}
		defUpdates = append(defUpdates, defUpdate)
	case pkg.DefKindCustomStep:
		if len(o.EnvNames) > 0 {
			paes := []pkg.ProjectAvailableEnv{}
			for _, pae := range project.ProjectAvailableEnvs {
				for _, envName := range o.EnvNames {
					if envName == pae.EnvName {
						paes = append(paes, pae)
						break
					}
				}
			}
			for _, pae := range paes {
				if len(pae.CustomStepDefs) > 0 {
					csds := pkg.CustomStepDefs{}
					for stepName, csd := range pae.CustomStepDefs {
						if len(o.StepNames) == 0 {
							csds[stepName] = csd
						} else {
							for _, name := range o.StepNames {
								if name == stepName {
									csds[stepName] = csd
									break
								}
							}
						}
					}
					for stepName, csd := range csds {
						defKind := defKindProject
						defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
						var errMsg string
						for name, msg := range pae.ErrMsgCustomStepDefs {
							if name == stepName {
								errMsg = msg
							}
						}
						defKind.Status.ErrMsg = errMsg
						defKind.Metadata.Labels = map[string]string{
							"envName":    pae.EnvName,
							"stepName":   stepName,
							"enableMode": csd.EnableMode,
						}

						ids := []int{}
						for i, csmd := range csd.CustomStepModuleDefs {
							var found bool
							for _, moduleName := range o.ModuleNames {
								if csmd.ModuleName == moduleName {
									found = true
									break
								}
							}
							if found {
								ids = append(ids, i)
							}
						}

						csmds := []pkg.CustomStepModuleDef{}
						for i, csmd := range csd.CustomStepModuleDefs {
							var found bool
							for _, id := range ids {
								if i == id {
									found = true
									break
								}
							}
							if !found {
								defKind.Items = append(defKind.Items, csmd)
								csmds = append(csmds, csmd)
							}
						}

						defKinds = append(defKinds, defKind)

						defUpdate := pkg.DefUpdate{
							Kind:        pkg.DefCmdKinds[o.Param.Kind],
							ProjectName: project.ProjectInfo.ProjectName,
							Def: pkg.CustomStepDef{
								EnableMode:                 csd.EnableMode,
								CustomStepModuleDefs:       csmds,
								UpdateCustomStepModuleDefs: false,
							},
							CustomStepName: stepName,
							EnvName:        pae.EnvName,
						}
						defUpdates = append(defUpdates, defUpdate)
					}
				}
			}
		} else {
			csds := pkg.CustomStepDefs{}
			for stepName, csd := range project.ProjectDef.CustomStepDefs {
				if len(o.StepNames) == 0 {
					csds[stepName] = csd
				} else {
					for _, name := range o.StepNames {
						if name == stepName {
							csds[stepName] = csd
							break
						}
					}
				}
			}
			for stepName, csd := range csds {
				defKind := defKindProject
				defKind.Kind = pkg.DefCmdKinds[o.Param.Kind]
				var errMsg string
				for name, msg := range project.ProjectDef.ErrMsgCustomStepDefs {
					if name == stepName {
						errMsg = msg
					}
				}
				defKind.Status.ErrMsg = errMsg
				defKind.Metadata.Labels = map[string]string{
					"stepName":   stepName,
					"enableMode": csd.EnableMode,
				}

				ids := []int{}
				for i, csmd := range csd.CustomStepModuleDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if csmd.ModuleName == moduleName {
							found = true
							break
						}
					}
					if found {
						ids = append(ids, i)
					}
				}

				csmds := []pkg.CustomStepModuleDef{}
				for i, csmd := range csd.CustomStepModuleDefs {
					var found bool
					for _, id := range ids {
						if i == id {
							found = true
							break
						}
					}
					if !found {
						defKind.Items = append(defKind.Items, csmd)
						csmds = append(csmds, csmd)
					}
				}

				defKinds = append(defKinds, defKind)

				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def: pkg.CustomStepDef{
						EnableMode:                 csd.EnableMode,
						CustomStepModuleDefs:       csmds,
						UpdateCustomStepModuleDefs: false,
					},
					CustomStepName: stepName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	}

	defKindList := pkg.DefKindList{
		Kind: "list",
		Defs: defKinds,
	}

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ := json.Marshal(defKindList)
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
		for _, defUpdate := range defUpdates {
			bs, _ = pkg.YamlIndent(defUpdate.Def)

			param := map[string]interface{}{
				"envName":        defUpdate.EnvName,
				"customStepName": defUpdate.CustomStepName,
				"branchName":     defUpdate.BranchName,
			}
			paramOutput := map[string]interface{}{}
			for k, v := range param {
				paramOutput[k] = v
			}

			urlKind := defUpdate.Kind
			switch defUpdate.Kind {
			case pkg.DefCmdKinds[pkg.DefKindBuild]:
				param["buildDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindPackage]:
				param["packageDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindArtifact]:
				param["artifactDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindDeployContainer]:
				param["deployContainerDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindDeployArtifact]:
				param["deployArtifactDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindIstio]:
				param["istioDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindCustomStep]:
				param["customStepDefYaml"] = string(bs)
				if defUpdate.EnvName != "" {
					urlKind = fmt.Sprintf("%s/env", urlKind)
				}
			case pkg.DefCmdKinds[pkg.DefKindCustomOps]:
				param["customOpsDefsYaml"] = string(bs)
			case pkg.DefCmdKinds[pkg.DefKindOpsBatch]:
				param["opsBatchDefsYaml"] = string(bs)
			}
			paramOutput = pkg.RemoveMapEmptyItems(paramOutput)
			bs, _ = json.Marshal(paramOutput)
			logHeader := fmt.Sprintf("[%s/%s] %s", defUpdate.ProjectName, defUpdate.Kind, string(bs))
			result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defUpdate.ProjectName, urlKind), http.MethodPost, "", param, false)
			if err != nil {
				err = fmt.Errorf("%s: %s", logHeader, err.Error())
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		}
	}

	return err
}
