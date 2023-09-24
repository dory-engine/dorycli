package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsDefPatch struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	BranchNames    []string `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	StepName       string   `yaml:"stepName" json:"stepName" bson:"stepName" validate:""`
	Patch          string   `yaml:"patch" json:"patch" bson:"patch" validate:""`
	FileName       string   `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Runs           []string `yaml:"runs" json:"runs" bson:"runs" validate:""`
	NoRuns         []string `yaml:"noRuns" json:"noRuns" bson:"noRuns" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind         string            `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName  string            `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
		PatchActions []pkg.PatchAction `yaml:"patchActions" json:"patchActions" bson:"patchActions" validate:""`
	}
}

func NewOptionsDefPatch() *OptionsDefPatch {
	var o OptionsDefPatch
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefPatch() *cobra.Command {
	o := NewOptionsDefPatch()

	defCmdKinds := []string{
		pkg.DefKindBuild,
		pkg.DefKindPackage,
		pkg.DefKindArtifact,
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindIstioGateway,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
		pkg.DefKindCustomStep,
		pkg.DefKindPipeline,
	}

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`patch [projectName] [kind] [--output=json|yaml] [--patch=patchAction] [--file=patchFile]... [--modules=moduleName1,moduleName2] [--envs=envName1,envName2] [--branches=branchName1,branchName2] [--step=stepName1,stepName2]
  # kind options: %s`, strings.Join(defCmdKinds, " / "))
	msgShort := fmt.Sprintf("patch project definitions")
	msgLong := fmt.Sprintf(`patch project definitions in dory-engine server`)
	msgExample := fmt.Sprintf(`  # print current project build modules definitions for patched
  %s def patch test-project1 %s --modules=tp1-go-demo,tp1-gin-demo -o yaml

  # patch project build modules definitions, update tp1-gin-demo,tp1-go-demo buildChecks commands
  %s def patch test-project1 %s --modules=tp1-go-demo,tp1-gin-demo --patch='[{"action": "update", "path": "buildChecks", "value": ["ls -alh"]}]'

  # patch project deploy modules definitions, delete test environment tp1-go-demo,tp1-gin-demo deployResources settings
  %s def patch test-project1 %s --modules=tp1-go-demo,tp1-gin-demo --envs=test --patch='[{"action": "delete", "path": "deployResources"}]'

  # patch project deploy modules definitions, delete test environment tp1-gin-demo deployNodePorts.0.nodePort to 30109
  %s def patch test-project1 %s --modules=tp1-gin-demo --envs=test --patch='[{"action": "update", "path": "deployNodePorts.0.nodePort", "value": 30109}]'

  # patch project pipeline definitions, update builds dp1-gin-demo run setting to true 
  %s def patch test-project1 %s --branches=develop,release --patch='[{"action": "update", "path": "builds.#(name==\"dp1-gin-demo\").run", "value": true}]'

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to true 
  %s def patch test-project1 %s --branches=develop,release --runs=dp1-gin-demo,dp1-go-demo

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to false 
  %s def patch test-project1 %s --branches=develop,release --no-runs=dp1-gin-demo,dp1-go-demo

  # patch project custom step modules definitions, update customStepName2 step in test environment tp1-gin-demo paramInputYaml
  %s def patch test-project1 %s --envs=test --step=customStepName2 --modules=tp1-gin-demo --patch='[{"action": "update", "path": "paramInputYaml", "value": "path: Tests"}]'

  # patch project pipeline definitions from stdin, support JSON and YAML
  cat << EOF | %s def patch test-project1 %s --branches=develop,release -f -
  - action: update
    path: builds
    value:
      - name: dp1-go-demo
        run: true
      - name: dp1-vue-demo
        run: true
  - action: update
    path: pipelineStep.deploy.enable
    value: false
  - action: delete
    value: customStepInsertDefs.build
  EOF

  # patch project pipeline definitions from file, support JSON and YAML
  %s def patch test-project1 %s --branches=develop,release -f patch.yaml`, baseName, pkg.DefKindBuild, baseName, pkg.DefKindBuild, baseName, pkg.DefKindDeployContainer, baseName, pkg.DefKindDeployContainer, baseName, pkg.DefKindPipeline, baseName, pkg.DefKindPipeline, baseName, pkg.DefKindPipeline, baseName, pkg.DefKindCustomStep, baseName, pkg.DefKindPipeline, baseName, pkg.DefKindPipeline)

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
	cmd.Flags().StringSliceVar(&o.ModuleNames, "modules", []string{}, "filter moduleNames to patch")
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, fmt.Sprintf("filter envNames to patch, required if kind is %s / %s / %s / %s", pkg.DefKindDeployContainer, pkg.DefKindDeployArtifact, pkg.DefKindIstio, pkg.DefKindIstioGateway))
	cmd.Flags().StringSliceVar(&o.BranchNames, "branches", []string{}, fmt.Sprintf("filter branchNames to patch, required if kind is %s", pkg.DefKindPipeline))
	cmd.Flags().StringVar(&o.StepName, "step", "", fmt.Sprintf("filter stepName to patch, required if kind is %s", pkg.DefKindCustomStep))
	cmd.Flags().StringVarP(&o.Patch, "patch", "p", "", "patch actions in JSON format")
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "project definitions file name or directory, support *.json and *.yaml and *.yml file")
	cmd.Flags().StringSliceVar(&o.Runs, "runs", []string{}, fmt.Sprintf("set pipeline which build modules enable run, only uses with kind is %s", pkg.DefKindPipeline))
	cmd.Flags().StringSliceVar(&o.NoRuns, "no-runs", []string{}, fmt.Sprintf("set pipeline which build modules disable run, only uses with kind is %s", pkg.DefKindPipeline))
	cmd.Flags().BoolVar(&o.Try, "try", false, "try to check input project definitions only, not apply to dory-engine server, use with --output option")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project definitions in full version, use with --output option")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsDefPatch) Complete(cmd *cobra.Command) error {
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
		pkg.DefKindIstioGateway,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
		pkg.DefKindCustomStep,
		pkg.DefKindPipeline,
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

	err = cmd.RegisterFlagCompletionFunc("patch", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		patchActions := []string{
			`'[{"action":"update","path":"xxx","value":"xxx"}]'`,
		}
		return patchActions, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
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

	err = cmd.RegisterFlagCompletionFunc("branches", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		branchNames := []string{}
		for _, pp := range project.ProjectPipelines {
			branchNames = append(branchNames, pp.BranchName)
		}
		return branchNames, cobra.ShellCompDirectiveNoFileComp
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

	err = cmd.RegisterFlagCompletionFunc("runs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		moduleNames := []string{}
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		for _, def := range project.ProjectDef.BuildDefs {
			moduleNames = append(moduleNames, def.BuildName)
		}
		return moduleNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("no-runs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		moduleNames := []string{}
		projectName := args[0]
		project, err := o.GetProjectDef(projectName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		for _, def := range project.ProjectDef.BuildDefs {
			moduleNames = append(moduleNames, def.BuildName)
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

	return err
}

func (o *OptionsDefPatch) Validate(args []string) error {
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

	defCmdKinds := []string{
		pkg.DefKindBuild,
		pkg.DefKindPackage,
		pkg.DefKindArtifact,
		pkg.DefKindDeployContainer,
		pkg.DefKindDeployArtifact,
		pkg.DefKindIstio,
		pkg.DefKindIstioGateway,
		pkg.DefKindCustomOps,
		pkg.DefKindOpsBatch,
		pkg.DefKindCustomStep,
		pkg.DefKindPipeline,
	}

	var found bool
	for _, cmdKind := range defCmdKinds {
		if cmdKind == kind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct, options: %s", kind, strings.Join(defCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	err = pkg.ValidateMinusNameID(projectName)
	if err != nil {
		err = fmt.Errorf("projectName %s format error: %s", projectName, err.Error())
		return err
	}
	o.Param.ProjectName = projectName

	if kind != pkg.DefKindPipeline && kind != pkg.DefKindIstioGateway && len(o.ModuleNames) == 0 {
		err = fmt.Errorf("--modules required")
		return err
	}
	if kind == pkg.DefKindPipeline && len(o.BranchNames) == 0 {
		err = fmt.Errorf("kind is %s, --branches required", pkg.DefKindPipeline)
		return err
	}
	if kind == pkg.DefKindDeployContainer && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindDeployContainer)
		return err
	}
	if kind == pkg.DefKindDeployArtifact && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindDeployArtifact)
		return err
	}
	if kind == pkg.DefKindIstio && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindIstio)
		return err
	}
	if kind == pkg.DefKindIstioGateway && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is %s, --envs required", pkg.DefKindIstioGateway)
		return err
	}
	if kind == pkg.DefKindCustomStep && o.StepName == "" {
		err = fmt.Errorf("kind is %s, --step required", pkg.DefKindCustomStep)
		return err
	}

	for _, runName := range o.Runs {
		err = pkg.ValidateRunName(runName)
		if err != nil {
			err = fmt.Errorf("run runName %s format error: %s", runName, err.Error())
			return err
		}
	}

	for _, runName := range o.NoRuns {
		err = pkg.ValidateRunName(runName)
		if err != nil {
			err = fmt.Errorf("no-run runName %s format error: %s", runName, err.Error())
			return err
		}
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}

	patchActions := []pkg.PatchAction{}
	pas := []pkg.PatchAction{}
	baseName := pkg.GetCmdBaseName()
	if o.FileName == "-" {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			err = fmt.Errorf("--file read stdin error: %s", err.Error())
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s def patch test-project1 %s --modules=tp1-gin-demo -f -", baseName, pkg.DefKindBuild)
			return err
		}
		err = json.Unmarshal(bs, &pas)
		if err != nil {
			err = yaml.Unmarshal(bs, &pas)
			if err != nil {
				err = fmt.Errorf("--file parse error: %s", err.Error())
				return err
			}
		}
	} else if o.FileName != "" {
		if o.FileName != "-" {
			ext := filepath.Ext(o.FileName)
			if ext != ".json" && ext != ".yaml" && ext != ".yml" {
				err = fmt.Errorf("--file %s read error: file extension must be json or yaml or yml", o.FileName)
				return err
			}
			bs, err := os.ReadFile(o.FileName)
			if err != nil {
				err = fmt.Errorf("--file %s read error: %s", o.FileName, err.Error())
				return err
			}
			switch ext {
			case ".json":
				err = json.Unmarshal(bs, &pas)
				if err != nil {
					err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
					return err
				}
			case ".yaml", ".yml":
				err = yaml.Unmarshal(bs, &pas)
				if err != nil {
					err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
					return err
				}
			}
		} else {
			bs, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			if len(bs) == 0 {
				err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s def patch [...] -f -", baseName)
				return err
			}
			if json.Unmarshal(bs, &pas) != nil {
				err = yaml.Unmarshal(bs, &pas)
				if err != nil {
					err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
					return err
				}
			}
		}
	}
	for _, pa := range pas {
		patchActions = append(patchActions, pa)
	}

	if o.Patch != "" {
		pas = []pkg.PatchAction{}
		err = json.Unmarshal([]byte(o.Patch), &pas)
		if err != nil {
			err = fmt.Errorf("--patch %s parse error: %s", o.Patch, err.Error())
			return err
		}
		for _, pa := range pas {
			patchActions = append(patchActions, pa)
		}
	}

	for _, patchAction := range patchActions {
		b, _ := json.Marshal(patchAction.Value)
		patchAction.Str = string(b)
		bs, _ := json.Marshal(patchAction)
		if patchAction.Action != "update" && patchAction.Action != "delete" {
			err = fmt.Errorf("--patch %s parse error: action must be update or delete", string(bs))
			return err
		}
		if patchAction.Path == "" {
			err = fmt.Errorf("--patch %s parse error: path can not be empty", string(bs))
			return err
		}
		o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
	}

	if kind == pkg.DefKindPipeline && len(o.Runs) > 0 {
		for _, name := range o.Runs {
			patchAction := pkg.PatchAction{
				Action: "update",
				Path:   fmt.Sprintf(`builds.#(name=="%s").run`, name),
				Value:  true,
				Str:    "true",
			}
			o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
		}
	}
	if kind == pkg.DefKindPipeline && len(o.NoRuns) > 0 {
		for _, name := range o.NoRuns {
			patchAction := pkg.PatchAction{
				Action: "update",
				Path:   fmt.Sprintf(`builds.#(name=="%s").run`, name),
				Value:  false,
				Str:    "false",
			}
			o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
		}
	}

	return err
}

func (o *OptionsDefPatch) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	project, err := o.GetProjectDef(o.Param.ProjectName)
	if err != nil {
		return err
	}

	for _, envName := range o.EnvNames {
		var found bool
		for _, pae := range project.ProjectAvailableEnvs {
			if envName == pae.EnvName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("envName %s not exists", envName)
			return err
		}
	}

	for _, branchName := range o.BranchNames {
		var found bool
		for _, pp := range project.ProjectPipelines {
			if branchName == pp.BranchName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("branchName %s not exists", branchName)
			return err
		}
	}

	if o.StepName != "" {
		var found bool
		for _, conf := range project.CustomStepConfs {
			if conf.CustomStepName == o.StepName {
				if len(o.EnvNames) == 0 && !conf.IsEnvDiff {
					found = true
					break
				} else if len(o.EnvNames) > 0 && conf.IsEnvDiff {
					found = true
					break
				}
			}
		}
		if !found {
			err = fmt.Errorf("stepName %s not exists", o.StepName)
			return err
		}
	}

	for _, run := range o.Runs {
		var found bool
		for _, def := range project.ProjectDef.BuildDefs {
			if run == def.BuildName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("run %s not exists", run)
			return err
		}
	}

	for _, noRun := range o.NoRuns {
		var found bool
		for _, def := range project.ProjectDef.BuildDefs {
			if noRun == def.BuildName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("no-run %s not exists", noRun)
			return err
		}
	}

	defUpdates := []pkg.DefUpdate{}
	defUpdateFilters := []pkg.DefUpdate{}

	switch o.Param.Kind {
	case pkg.DefKindBuild:
		sort.SliceStable(project.ProjectDef.BuildDefs, func(i, j int) bool {
			return project.ProjectDef.BuildDefs[i].BuildName < project.ProjectDef.BuildDefs[j].BuildName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.BuildDefs {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}

		defs := []pkg.BuildDef{}
		ds := []pkg.BuildDef{}
		for _, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if found == true {
				ds = append(ds, def)
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
		defUpdateFilter := defUpdate
		defUpdateFilter.Def = ds
		defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
	case pkg.DefKindPackage:
		sort.SliceStable(project.ProjectDef.PackageDefs, func(i, j int) bool {
			return project.ProjectDef.PackageDefs[i].PackageName < project.ProjectDef.PackageDefs[j].PackageName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.PackageDefs {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.PackageDef{}
		ds := []pkg.PackageDef{}
		for _, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			if found == true {
				ds = append(ds, def)
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
		defUpdateFilter := defUpdate
		defUpdateFilter.Def = ds
		defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
	case pkg.DefKindArtifact:
		sort.SliceStable(project.ProjectDef.ArtifactDefs, func(i, j int) bool {
			return project.ProjectDef.ArtifactDefs[i].ArtifactName < project.ProjectDef.ArtifactDefs[j].ArtifactName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.ArtifactDefs {
				if def.ArtifactName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.ArtifactDef{}
		ds := []pkg.ArtifactDef{}
		for _, def := range project.ProjectDef.ArtifactDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.ArtifactName == moduleName {
					found = true
					break
				}
			}
			if found == true {
				ds = append(ds, def)
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
		defUpdateFilter := defUpdate
		defUpdateFilter.Def = ds
		defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
	case pkg.DefKindDeployContainer:
		for _, pae := range project.ProjectAvailableEnvs {
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				sort.SliceStable(pae.DeployContainerDefs, func(i, j int) bool {
					return pae.DeployContainerDefs[i].DeployName < pae.DeployContainerDefs[j].DeployName
				})
				for _, moduleName := range o.ModuleNames {
					var found bool
					for _, def := range pae.DeployContainerDefs {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("%s module %s in envName %s not exists", o.Param.Kind, moduleName, pae.EnvName)
						return err
					}
				}
				defs := []pkg.DeployContainerDef{}
				ds := []pkg.DeployContainerDef{}
				for _, def := range pae.DeployContainerDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if found == true {
						ds = append(ds, def)
					}
					def.IsPatch = found
					defs = append(defs, def)
				}
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					EnvName:     pae.EnvName,
					Def:         defs,
				}
				defUpdates = append(defUpdates, defUpdate)
				defUpdateFilter := defUpdate
				defUpdateFilter.Def = ds
				defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
			}
		}
	case pkg.DefKindDeployArtifact:
		for _, pae := range project.ProjectAvailableEnvs {
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				sort.SliceStable(pae.DeployArtifactDefs, func(i, j int) bool {
					return pae.DeployArtifactDefs[i].DeployArtifactName < pae.DeployArtifactDefs[j].DeployArtifactName
				})
				for _, moduleName := range o.ModuleNames {
					var found bool
					for _, def := range pae.DeployArtifactDefs {
						if def.DeployArtifactName == moduleName {
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("%s module %s in envName %s not exists", o.Param.Kind, moduleName, pae.EnvName)
						return err
					}
				}
				defs := []pkg.DeployArtifactDef{}
				ds := []pkg.DeployArtifactDef{}
				for _, def := range pae.DeployArtifactDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployArtifactName == moduleName {
							found = true
							break
						}
					}
					if found == true {
						ds = append(ds, def)
					}
					def.IsPatch = found
					defs = append(defs, def)
				}
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					EnvName:     pae.EnvName,
					Def:         defs,
				}
				defUpdates = append(defUpdates, defUpdate)
				defUpdateFilter := defUpdate
				defUpdateFilter.Def = ds
				defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
			}
		}
	case pkg.DefKindIstio:
		for _, pae := range project.ProjectAvailableEnvs {
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				sort.SliceStable(pae.IstioDefs, func(i, j int) bool {
					return pae.IstioDefs[i].DeployName < pae.IstioDefs[j].DeployName
				})
				for _, moduleName := range o.ModuleNames {
					var found bool
					for _, def := range pae.IstioDefs {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("%s module %s in envName %s not exists", o.Param.Kind, moduleName, pae.EnvName)
						return err
					}
				}
				defs := []pkg.IstioDef{}
				ds := []pkg.IstioDef{}
				for _, def := range pae.IstioDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if found == true {
						ds = append(ds, def)
					}
					def.IsPatch = found
					defs = append(defs, def)
				}
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					EnvName:     pae.EnvName,
					Def:         defs,
				}
				defUpdates = append(defUpdates, defUpdate)
				defUpdateFilter := defUpdate
				defUpdateFilter.Def = ds
				defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
			}
		}
	case pkg.DefKindIstioGateway:
		for _, pae := range project.ProjectAvailableEnvs {
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.IstioGatewayDef,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
				defUpdateFilter := defUpdate
				defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
			}
		}
	case pkg.DefKindCustomStep:
		if len(o.EnvNames) == 0 {
			for stepName, csd := range project.ProjectDef.CustomStepDefs {
				if stepName == o.StepName {
					sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
						return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
					})
					for _, moduleName := range o.ModuleNames {
						var found bool
						for _, def := range csd.CustomStepModuleDefs {
							if def.ModuleName == moduleName {
								found = true
								break
							}
						}
						if !found {
							err = fmt.Errorf("%s module %s step %s not exists", o.Param.Kind, moduleName, stepName)
							return err
						}
					}
					defs := []pkg.CustomStepModuleDef{}
					ds := []pkg.CustomStepModuleDef{}
					for _, def := range csd.CustomStepModuleDefs {
						var found bool
						for _, moduleName := range o.ModuleNames {
							if def.ModuleName == moduleName {
								found = true
								break
							}
						}
						if found == true {
							ds = append(ds, def)
						}
						def.IsPatch = found
						defs = append(defs, def)
					}
					csd.CustomStepModuleDefs = defs
					defUpdate := pkg.DefUpdate{
						Kind:           pkg.DefCmdKinds[o.Param.Kind],
						ProjectName:    project.ProjectInfo.ProjectName,
						Def:            csd,
						CustomStepName: stepName,
					}
					defUpdates = append(defUpdates, defUpdate)
					defUpdateFilter := defUpdate
					csd.CustomStepModuleDefs = ds
					defUpdateFilter.Def = csd
					defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
					break
				}
			}
		} else {
			for _, pae := range project.ProjectAvailableEnvs {
				for stepName, csd := range pae.CustomStepDefs {
					var found bool
					for _, envName := range o.EnvNames {
						if pae.EnvName == envName {
							found = true
							break
						}
					}
					if found {
						sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
							return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
						})
						for _, moduleName := range o.ModuleNames {
							var found bool
							for _, def := range csd.CustomStepModuleDefs {
								if def.ModuleName == moduleName {
									found = true
									break
								}
							}
							if !found {
								err = fmt.Errorf("%s module %s step %s in envName %s not exists", o.Param.Kind, moduleName, stepName, pae.EnvName)
								return err
							}
						}
						defs := []pkg.CustomStepModuleDef{}
						ds := []pkg.CustomStepModuleDef{}
						for _, def := range csd.CustomStepModuleDefs {
							var found bool
							for _, moduleName := range o.ModuleNames {
								if def.ModuleName == moduleName {
									found = true
									break
								}
							}
							if found == true {
								ds = append(ds, def)
							}
							def.IsPatch = found
							defs = append(defs, def)
						}
						csd.CustomStepModuleDefs = defs
						defUpdate := pkg.DefUpdate{
							Kind:           pkg.DefCmdKinds[o.Param.Kind],
							ProjectName:    project.ProjectInfo.ProjectName,
							Def:            csd,
							EnvName:        pae.EnvName,
							CustomStepName: stepName,
						}
						defUpdates = append(defUpdates, defUpdate)
						defUpdateFilter := defUpdate
						csd.CustomStepModuleDefs = ds
						defUpdateFilter.Def = csd
						defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
					}
				}
			}
		}
	case pkg.DefKindPipeline:
		for _, pp := range project.ProjectPipelines {
			var found bool
			for _, branchName := range o.BranchNames {
				if pp.BranchName == branchName {
					found = true
					break
				}
			}
			if found {
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pp.PipelineDef,
					BranchName:  pp.BranchName,
				}
				defUpdates = append(defUpdates, defUpdate)
				defUpdateFilter := defUpdate
				defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
			}
		}
	case pkg.DefKindCustomOps:
		sort.SliceStable(project.ProjectDef.CustomOpsDefs, func(i, j int) bool {
			return project.ProjectDef.CustomOpsDefs[i].CustomOpsName < project.ProjectDef.CustomOpsDefs[j].CustomOpsName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.CustomOpsDefs {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.CustomOpsDef{}
		ds := []pkg.CustomOpsDef{}
		for _, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			if found == true {
				ds = append(ds, def)
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
		defUpdateFilter := defUpdate
		defUpdateFilter.Def = ds
		defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
	case pkg.DefKindOpsBatch:
		sort.SliceStable(project.ProjectDef.OpsBatchDefs, func(i, j int) bool {
			return project.ProjectDef.OpsBatchDefs[i].OpsBatchName < project.ProjectDef.OpsBatchDefs[j].OpsBatchName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.OpsBatchDefs {
				if def.OpsBatchName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.OpsBatchDef{}
		ds := []pkg.OpsBatchDef{}
		for _, def := range project.ProjectDef.OpsBatchDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.OpsBatchName == moduleName {
					found = true
					break
				}
			}
			if found == true {
				ds = append(ds, def)
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
		defUpdateFilter := defUpdate
		defUpdateFilter.Def = ds
		defUpdateFilters = append(defUpdateFilters, defUpdateFilter)
	}

	if len(defUpdates) == 0 {
		err = fmt.Errorf("nothing to patch")
		return err
	}

	defPatches := []pkg.DefUpdate{}
	if len(o.Param.PatchActions) > 0 {
		for idx, defUpdate := range defUpdates {
			bs, _ := json.Marshal(defUpdate.Def)
			switch defUpdate.Kind {
			case "buildDefs":
				defs := []pkg.BuildDef{}
				dps := []pkg.BuildDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.BuildDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.BuildDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "packageDefs":
				defs := []pkg.PackageDef{}
				dps := []pkg.PackageDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.PackageDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.PackageDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "artifactDefs":
				defs := []pkg.ArtifactDef{}
				dps := []pkg.ArtifactDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.ArtifactDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.ArtifactDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "deployContainerDefs":
				defs := []pkg.DeployContainerDef{}
				dps := []pkg.DeployContainerDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.DeployContainerDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.DeployContainerDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "deployArtifactDefs":
				defs := []pkg.DeployArtifactDef{}
				dps := []pkg.DeployArtifactDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.DeployArtifactDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.DeployArtifactDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "istioDefs":
				defs := []pkg.IstioDef{}
				dps := []pkg.IstioDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.IstioDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.IstioDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "istioGatewayDef":
				def := pkg.IstioGatewayDef{}
				_ = json.Unmarshal(bs, &def)
				var dp pkg.IstioGatewayDef
				var s string
				for _, patchAction := range o.Param.PatchActions {
					switch patchAction.Action {
					case "update":
						s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
							return err
						}
					case "delete":
						s, err = sjson.Delete(string(bs), patchAction.Path)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
							return err
						}
					}
					var dd pkg.IstioGatewayDef
					err = json.Unmarshal([]byte(s), &dd)
					if err != nil {
						err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
						return err
					}
					bs = []byte(s)
					dp = dd
				}
				defUpdate.Def = dp
				defUpdates[idx] = defUpdate
				defPatches = append(defPatches, defUpdate)
			case "customStepDef":
				defs := pkg.CustomStepDef{}
				dps := []pkg.CustomStepModuleDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs.CustomStepModuleDefs {
					if d.IsPatch {
						var dp pkg.CustomStepModuleDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.CustomStepModuleDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs.CustomStepModuleDefs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defs.CustomStepModuleDefs = dps
				defUpdate.Def = defs
				defPatch := defUpdate
				defPatches = append(defPatches, defPatch)
			case "pipelineDef":
				def := pkg.PipelineDef{}
				_ = json.Unmarshal(bs, &def)
				var dp pkg.PipelineDef
				var s string
				for _, patchAction := range o.Param.PatchActions {
					switch patchAction.Action {
					case "update":
						s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
							return err
						}
					case "delete":
						s, err = sjson.Delete(string(bs), patchAction.Path)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
							return err
						}
					}
					var dd pkg.PipelineDef
					err = json.Unmarshal([]byte(s), &dd)
					if err != nil {
						err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
						return err
					}
					bs = []byte(s)
					dp = dd
				}
				defUpdate.Def = dp
				defUpdates[idx] = defUpdate
				defPatches = append(defPatches, defUpdate)
			case "customOpsDefs":
				defs := []pkg.CustomOpsDef{}
				dps := []pkg.CustomOpsDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.CustomOpsDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.CustomOpsDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "opsBatchDefs":
				defs := []pkg.OpsBatchDef{}
				dps := []pkg.OpsBatchDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.OpsBatchDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.OpsBatchDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			}
		}
	}

	defUpdateList := pkg.DefUpdateList{
		Kind: "list",
	}
	if len(defPatches) == 0 {
		defUpdateList.Defs = defUpdateFilters
	} else {
		defUpdateList.Defs = defPatches
	}

	mapOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(defUpdateList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		mapOutput = m
	} else {
		mapOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(mapOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = pkg.YamlIndent(mapOutput)
		fmt.Println(string(bs))
	}

	if !o.Try && len(defPatches) > 0 {
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
			case "buildDefs":
				param["buildDefsYaml"] = string(bs)
			case "packageDefs":
				param["packageDefsYaml"] = string(bs)
			case "artifactDefs":
				param["artifactDefsYaml"] = string(bs)
			case "deployContainerDefs":
				param["deployContainerDefsYaml"] = string(bs)
			case "deployArtifactDefs":
				param["deployArtifactDefsYaml"] = string(bs)
			case "istioDefs":
				param["istioDefsYaml"] = string(bs)
			case "istioGatewayDef":
				param["istioGatewayDefYaml"] = string(bs)
			case "customStepDef":
				param["customStepDefYaml"] = string(bs)
				if defUpdate.EnvName != "" {
					urlKind = fmt.Sprintf("%s/env", urlKind)
				}
			case "customOpsDefs":
				param["customOpsDefsYaml"] = string(bs)
			case "opsBatchDefs":
				param["opsBatchDefsYaml"] = string(bs)
			case "pipelineDef":
				param["pipelineDefYaml"] = string(bs)
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
