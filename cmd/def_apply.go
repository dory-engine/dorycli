package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsDefApply struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames []string      `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Defs      []pkg.DefKind `yaml:"defs" json:"defs" bson:"defs" validate:""`
	}
}

func NewOptionsDefApply() *OptionsDefApply {
	var o OptionsDefApply
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefApply() *cobra.Command {
	o := NewOptionsDefApply()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`apply -f [filename]`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_def_apply_short")
	msgLong := OptCommon.TransLang("cmd_def_apply_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_def_apply_example", baseName, baseName))

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_def_apply_output"))
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, OptCommon.TransLang("param_def_apply_recursive"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_def_apply_full"))
	cmd.Flags().StringSliceVarP(&o.FileNames, "files", "f", []string{}, OptCommon.TransLang("param_def_apply_files"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_def_apply_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func CheckDefKind(def pkg.DefKind) error {
	var err error
	switch def.Kind {
	case "buildDefs":
		for _, item := range def.Items {
			var d pkg.BuildDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is buildDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "packageDefs":
		for _, item := range def.Items {
			var d pkg.PackageDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is packageDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "artifactDefs":
		for _, item := range def.Items {
			var d pkg.ArtifactDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is artifactDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "deployContainerDefs":
		var envName string
		for k, v := range def.Metadata.Labels {
			if k == "envName" {
				envName = v
				break
			}
		}
		if envName == "" {
			err = fmt.Errorf("kind is deployContainerDefs, but projectName %s metadata.Labels.envName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.DeployContainerDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "deployArtifactDefs":
		var envName string
		for k, v := range def.Metadata.Labels {
			if k == "envName" {
				envName = v
				break
			}
		}
		if envName == "" {
			err = fmt.Errorf("kind is deployArtifactDefs, but projectName %s metadata.Labels.envName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.DeployArtifactDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is deployArtifactDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "istioDefs":
		var envName string
		for k, v := range def.Metadata.Labels {
			if k == "envName" {
				envName = v
				break
			}
		}
		if envName == "" {
			err = fmt.Errorf("kind is istioDefs, but projectName %s metadata.Labels.envName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.IstioDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is istioDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "istioGatewayDef":
		var envName string
		for k, v := range def.Metadata.Labels {
			if k == "envName" {
				envName = v
				break
			}
		}
		if envName == "" {
			err = fmt.Errorf("kind is istioGatewayDef, but projectName %s metadata.Labels.envName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.IstioGatewayDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is istioGatewayDef, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "pipelineDef":
		var branchName string
		for k, v := range def.Metadata.Labels {
			if k == "branchName" {
				branchName = v
				break
			}
		}
		if branchName == "" {
			err = fmt.Errorf("kind is pipelineDef, but projectName %s metadata.Labels.branchName is empty", def.Metadata.ProjectName)
			return err
		}
		if len(def.Items) != 1 {
			err = fmt.Errorf("kind is pipelineDef, but projectName %s items size is not 1", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.PipelineDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is pipelineDef, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "dockerIgnoreDefs":
		for _, item := range def.Items {
			switch item.(type) {
			case string:
			default:
				err = fmt.Errorf("kind is dockerIgnoreDefs, but item parse error: items must be string array")
				return err
			}
		}
	case "customOpsDefs":
		for _, item := range def.Items {
			var d pkg.CustomOpsDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "opsBatchDefs":
		for _, item := range def.Items {
			var d pkg.OpsBatchDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is opsBatchDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	case "customStepDef":
		var stepName string
		for k, v := range def.Metadata.Labels {
			if k == "stepName" {
				stepName = v
				break
			}
		}
		if stepName == "" {
			err = fmt.Errorf("kind is customStepDef, but projectName %s metadata.Labels.stepName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.CustomStepModuleDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is customStepDef, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
		}
	}
	return err
}

func GetDefKindsFromJson(fileName string, bs []byte) ([]pkg.DefKind, error) {
	var err error
	defKinds := []pkg.DefKind{}
	var list pkg.DefKindList
	err = json.Unmarshal(bs, &list)
	if err == nil {
		if list.Kind == "list" {
			defKinds = append(defKinds, list.Defs...)
		} else {
			var def pkg.DefKind
			err = json.Unmarshal(bs, &def)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return defKinds, err
			}
			if def.Kind != "" {
				defKinds = append(defKinds, def)
			}
		}
	} else {
		var def pkg.DefKind
		err = json.Unmarshal(bs, &def)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return defKinds, err
		}
		if def.Kind != "" {
			defKinds = append(defKinds, def)
		}
	}
	return defKinds, err
}

func GetDefKindsFromYaml(fileName string, bs []byte) ([]pkg.DefKind, error) {
	var err error
	defKinds := []pkg.DefKind{}
	dec := yaml.NewDecoder(bytes.NewReader(bs))
	ms := []map[string]interface{}{}
	for {
		var m map[string]interface{}
		err = dec.Decode(&m)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return defKinds, err
		} else {
			ms = append(ms, m)
		}
	}
	for _, m := range ms {
		b, _ := yaml.Marshal(m)
		var list pkg.DefKindList
		err = yaml.Unmarshal(b, &list)
		if err == nil {
			if list.Kind == "list" {
				defKinds = append(defKinds, list.Defs...)
			} else {
				var def pkg.DefKind
				err = yaml.Unmarshal(b, &def)
				if err != nil {
					err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
					return defKinds, err
				}
				if def.Kind != "" {
					defKinds = append(defKinds, def)
				}
			}
		} else {
			var def pkg.DefKind
			err = yaml.Unmarshal(b, &def)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return defKinds, err
			}
			if def.Kind != "" {
				defKinds = append(defKinds, def)
			}
		}
	}

	return defKinds, err
}

func GetDefKinds(fileName string, bs []byte) ([]pkg.DefKind, error) {
	var err error
	defKinds := []pkg.DefKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		defKinds, err = GetDefKindsFromJson(fileName, bs)
		if err != nil {
			return defKinds, err
		}
	} else if ext == ".yaml" || ext == ".yml" {
		defKinds, err = GetDefKindsFromYaml(fileName, bs)
		if err != nil {
			return defKinds, err
		}
	} else if fileName == "" {
		defKinds, err = GetDefKindsFromJson(fileName, bs)
		if err != nil {
			defKinds, err = GetDefKindsFromYaml(fileName, bs)
			if err != nil {
				return defKinds, err
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return defKinds, err
	}

	for _, def := range defKinds {
		if def.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return defKinds, err
		}
		if def.Metadata.ProjectName == "" {
			err = fmt.Errorf("parse file %s error: metadata.projectName is empty", fileName)
			return defKinds, err
		}
		err = pkg.ValidateMinusNameID(def.Metadata.ProjectName)
		if err != nil {
			err = fmt.Errorf("parse file %s error: metadata.projectName %s format error: %s", fileName, def.Metadata.ProjectName, err.Error())
			return defKinds, err
		}

		var found bool

		var kinds []string
		for _, v := range pkg.DefCmdKinds {
			if v != "" {
				kinds = append(kinds, v)
			}
		}
		for _, d := range kinds {
			if def.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, def.Kind)
			return defKinds, err
		}
		err = CheckDefKind(def)
		if err != nil {
			return defKinds, err
		}
	}
	return defKinds, err
}

func (o *OptionsDefApply) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("files")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsDefApply) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(o.FileNames) == 0 {
		err = fmt.Errorf("--files required")
		return err
	}
	var fileNames []string
	for _, name := range o.FileNames {
		fileNames = append(fileNames, strings.Trim(name, " "))
	}
	var isStdin bool
	for _, name := range fileNames {
		if name == "-" {
			isStdin = true
			break
		}
	}
	if isStdin && len(fileNames) > 1 {
		err = fmt.Errorf(`"--files -" found, can not use multiple --files options`)
		return err
	}

	baseName := pkg.GetCmdBaseName()
	if isStdin {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--files - required os.stdin\n example: echo 'xxx' | %s def apply -f -", baseName)
			return err
		}
		defs, err := GetDefKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Defs = append(o.Param.Defs, defs...)
	} else {
		for _, fileName := range fileNames {
			fi, err := os.Stat(fileName)
			if err != nil {
				return err
			}
			if fi.IsDir() {
				if o.Recursive {
					err = filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						ext := filepath.Ext(path)
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							o.Param.FileNames = append(o.Param.FileNames, path)
						}
						return nil
					})
				} else {
					infos, err := ioutil.ReadDir(fileName)
					if err != nil {
						return err
					}
					for _, info := range infos {
						ext := filepath.Ext(info.Name())
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							if strings.HasSuffix(fileName, "/") {
								fileName = strings.TrimSuffix(fileName, "/")
							}
							o.Param.FileNames = append(o.Param.FileNames, fmt.Sprintf("%s/%s", fileName, info.Name()))
						}
					}
				}
			} else {
				ext := filepath.Ext(fileName)
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					err = fmt.Errorf("file %s error: file extension name not json, yaml or yml", fileName)
					return err
				}
				o.Param.FileNames = append(o.Param.FileNames, fileName)
			}
		}

		fileNames = []string{}
		m := map[string]bool{}
		for _, fileName := range o.Param.FileNames {
			m[fileName] = true
		}
		for fileName, _ := range m {
			fileNames = append(fileNames, fileName)
		}
		sort.Strings(fileNames)
		o.Param.FileNames = fileNames

		for _, fileName := range o.Param.FileNames {
			bs, err := os.ReadFile(fileName)
			if err != nil {
				err = fmt.Errorf("read file %s error: %s", fileName, err.Error())
				return err
			}

			defs, err := GetDefKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Defs = append(o.Param.Defs, defs...)
		}
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsDefApply) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	mapDefProjects := map[string][]pkg.DefKind{}
	projects := []pkg.ProjectOutput{}
	for _, def := range o.Param.Defs {
		mapDefProjects[def.Metadata.ProjectName] = append(mapDefProjects[def.Metadata.ProjectName], def)
	}
	for projectName, defs := range mapDefProjects {
		param := map[string]interface{}{}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", projectName), http.MethodGet, "", param, false)
		if err != nil {
			return err
		}
		project := pkg.ProjectOutput{}
		err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
		if err != nil {
			return err
		}

		for _, def := range defs {
			switch def.Kind {
			case "buildDefs":
				for _, item := range def.Items {
					var d pkg.BuildDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, buildDef := range project.ProjectDef.BuildDefs {
						if buildDef.BuildName == d.BuildName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.BuildDefs[idx] = d
					} else {
						project.ProjectDef.BuildDefs = append(project.ProjectDef.BuildDefs, d)
					}
					project.ProjectDef.UpdateBuildDefs = true
				}
			case "packageDefs":
				for _, item := range def.Items {
					var d pkg.PackageDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, packageDef := range project.ProjectDef.PackageDefs {
						if packageDef.PackageName == d.PackageName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.PackageDefs[idx] = d
					} else {
						project.ProjectDef.PackageDefs = append(project.ProjectDef.PackageDefs, d)
					}
					project.ProjectDef.UpdatePackageDefs = true
				}
			case "artifactDefs":
				for _, item := range def.Items {
					var d pkg.ArtifactDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, artifactDef := range project.ProjectDef.ArtifactDefs {
						if artifactDef.ArtifactName == d.ArtifactName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.ArtifactDefs[idx] = d
					} else {
						project.ProjectDef.ArtifactDefs = append(project.ProjectDef.ArtifactDefs, d)
					}
					project.ProjectDef.UpdateArtifactDefs = true
				}
			case "deployContainerDefs":
				var envName string
				for k, v := range def.Metadata.Labels {
					if k == "envName" {
						envName = v
						break
					}
				}
				var projectAvailableEnv pkg.ProjectAvailableEnv
				index := -1
				for i, pae := range project.ProjectAvailableEnvs {
					if pae.EnvName == envName {
						projectAvailableEnv = pae
						index = i
						break
					}
				}
				if projectAvailableEnv.EnvName == "" {
					err = fmt.Errorf("kind is deployContainerDefs, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.DeployContainerDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, deployContainerDef := range projectAvailableEnv.DeployContainerDefs {
						if deployContainerDef.DeployName == d.DeployName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						projectAvailableEnv.DeployContainerDefs[idx] = d
					} else {
						projectAvailableEnv.DeployContainerDefs = append(projectAvailableEnv.DeployContainerDefs, d)
					}
					projectAvailableEnv.UpdateDeployContainerDefs = true
				}
				project.ProjectAvailableEnvs[index] = projectAvailableEnv
			case "deployArtifactDefs":
				var envName string
				for k, v := range def.Metadata.Labels {
					if k == "envName" {
						envName = v
						break
					}
				}
				var projectAvailableEnv pkg.ProjectAvailableEnv
				index := -1
				for i, pae := range project.ProjectAvailableEnvs {
					if pae.EnvName == envName {
						projectAvailableEnv = pae
						index = i
						break
					}
				}
				if projectAvailableEnv.EnvName == "" {
					err = fmt.Errorf("kind is deployArtifactDefs, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.DeployArtifactDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, deployArtifactDef := range projectAvailableEnv.DeployArtifactDefs {
						if deployArtifactDef.DeployArtifactName == d.DeployArtifactName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						projectAvailableEnv.DeployArtifactDefs[idx] = d
					} else {
						projectAvailableEnv.DeployArtifactDefs = append(projectAvailableEnv.DeployArtifactDefs, d)
					}
					projectAvailableEnv.UpdateDeployArtifactDefs = true
				}
				project.ProjectAvailableEnvs[index] = projectAvailableEnv
			case "istioDefs":
				var envName string
				for k, v := range def.Metadata.Labels {
					if k == "envName" {
						envName = v
						break
					}
				}
				var projectAvailableEnv pkg.ProjectAvailableEnv
				index := -1
				for i, pae := range project.ProjectAvailableEnvs {
					if pae.EnvName == envName {
						projectAvailableEnv = pae
						index = i
						break
					}
				}
				if projectAvailableEnv.EnvName == "" {
					err = fmt.Errorf("kind is istioDefs, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.IstioDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, istioDef := range projectAvailableEnv.IstioDefs {
						if istioDef.DeployName == d.DeployName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						projectAvailableEnv.IstioDefs[idx] = d
					} else {
						projectAvailableEnv.IstioDefs = append(projectAvailableEnv.IstioDefs, d)
					}
					projectAvailableEnv.UpdateIstioDefs = true
				}
				project.ProjectAvailableEnvs[index] = projectAvailableEnv
			case "istioGatewayDef":
				var envName string
				for k, v := range def.Metadata.Labels {
					if k == "envName" {
						envName = v
						break
					}
				}
				var projectAvailableEnv pkg.ProjectAvailableEnv
				index := -1
				for i, pae := range project.ProjectAvailableEnvs {
					if pae.EnvName == envName {
						projectAvailableEnv = pae
						index = i
						break
					}
				}
				if projectAvailableEnv.EnvName == "" {
					err = fmt.Errorf("kind is istioGatewayDef, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.IstioGatewayDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					projectAvailableEnv.IstioGatewayDef = d
					projectAvailableEnv.UpdateIstioGatewayDef = true
				}
				project.ProjectAvailableEnvs[index] = projectAvailableEnv
			case "pipelineDef":
				var branchName string
				for k, v := range def.Metadata.Labels {
					if k == "branchName" {
						branchName = v
						break
					}
				}
				var projectPipeline pkg.ProjectPipeline
				index := -1
				for i, pp := range project.ProjectPipelines {
					if pp.BranchName == branchName {
						projectPipeline = pp
						index = i
						break
					}
				}
				if projectPipeline.BranchName == "" {
					err = fmt.Errorf("kind is pipelineDef, but projectName %s metadata.Labels.branchName %s not exists", def.Metadata.ProjectName, branchName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.PipelineDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					projectPipeline.PipelineDef = d
					projectPipeline.UpdatePipelineDef = true
				}
				project.ProjectPipelines[index] = projectPipeline
			case "dockerIgnoreDefs":
				dockerIgnoreDefs := []string{}
				for _, item := range def.Items {
					switch v := item.(type) {
					case string:
						dockerIgnoreDefs = append(dockerIgnoreDefs, v)
					}
				}
				project.ProjectDef.DockerIgnoreDefs = dockerIgnoreDefs
				project.ProjectDef.UpdateDockerIgnoreDefs = true
			case "customOpsDefs":
				for _, item := range def.Items {
					var d pkg.CustomOpsDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, customOpsDef := range project.ProjectDef.CustomOpsDefs {
						if customOpsDef.CustomOpsName == d.CustomOpsName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.CustomOpsDefs[idx] = d
					} else {
						project.ProjectDef.CustomOpsDefs = append(project.ProjectDef.CustomOpsDefs, d)
					}
					project.ProjectDef.UpdateCustomOpsDefs = true
				}
			case "opsBatchDefs":
				for _, item := range def.Items {
					var d pkg.OpsBatchDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, opsBatchDef := range project.ProjectDef.OpsBatchDefs {
						if opsBatchDef.OpsBatchName == d.OpsBatchName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.OpsBatchDefs[idx] = d
					} else {
						project.ProjectDef.OpsBatchDefs = append(project.ProjectDef.OpsBatchDefs, d)
					}
					project.ProjectDef.UpdateOpsBatchDefs = true
				}
			case "customStepDef":
				var stepName string
				var envName string
				var enableMode string
				for k, v := range def.Metadata.Labels {
					if k == "stepName" {
						stepName = v
					}
					if k == "envName" {
						envName = v
					}
					if k == "enableMode" {
						enableMode = v
					}
				}
				if envName != "" {
					var projectAvailableEnv pkg.ProjectAvailableEnv
					index := -1
					for i, pae := range project.ProjectAvailableEnvs {
						if pae.EnvName == envName {
							projectAvailableEnv = pae
							index = i
							break
						}
					}
					if projectAvailableEnv.EnvName == "" {
						err = fmt.Errorf("kind is customStepDef, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
						return err
					}
					var found bool
					var customStepDef pkg.CustomStepDef
					for name, csd := range projectAvailableEnv.CustomStepDefs {
						if name == stepName {
							customStepDef = csd
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("kind is customStepDef, but projectName %s metadata.Labels.stepName %s not exists", def.Metadata.ProjectName, stepName)
						return err
					}
					for _, item := range def.Items {
						var d pkg.CustomStepModuleDef
						bs, _ := pkg.YamlIndent(item)
						_ = yaml.Unmarshal(bs, &d)
						idx := -1
						for i, moduleDef := range customStepDef.CustomStepModuleDefs {
							if d.ModuleName == moduleDef.ModuleName {
								idx = i
								break
							}
						}
						if idx >= 0 {
							customStepDef.CustomStepModuleDefs[idx] = d
						} else {
							customStepDef.CustomStepModuleDefs = append(customStepDef.CustomStepModuleDefs, d)
						}
						customStepDef.UpdateCustomStepModuleDefs = true
					}
					customStepDef.EnableMode = enableMode
					projectAvailableEnv.CustomStepDefs[stepName] = customStepDef
					project.ProjectAvailableEnvs[index] = projectAvailableEnv
				} else {
					var found bool
					var customStepDef pkg.CustomStepDef
					for name, csd := range project.ProjectDef.CustomStepDefs {
						if name == stepName {
							customStepDef = csd
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("kind is customStepDef, but projectName %s metadata.Labels.stepName %s not exists", def.Metadata.ProjectName, stepName)
						return err
					}
					for _, item := range def.Items {
						var d pkg.CustomStepModuleDef
						bs, _ := pkg.YamlIndent(item)
						_ = yaml.Unmarshal(bs, &d)
						idx := -1
						for i, moduleDef := range customStepDef.CustomStepModuleDefs {
							if d.ModuleName == moduleDef.ModuleName {
								idx = i
								break
							}
						}
						if idx >= 0 {
							customStepDef.CustomStepModuleDefs[idx] = d
						} else {
							customStepDef.CustomStepModuleDefs = append(customStepDef.CustomStepModuleDefs, d)
						}
						customStepDef.UpdateCustomStepModuleDefs = true
					}
					customStepDef.EnableMode = enableMode
					project.ProjectDef.CustomStepDefs[stepName] = customStepDef
				}
			}
		}
		projects = append(projects, project)
	}

	defUpdates := []pkg.DefUpdate{}

	for _, project := range projects {
		if project.ProjectDef.UpdateBuildDefs {
			sort.SliceStable(project.ProjectDef.BuildDefs, func(i, j int) bool {
				return project.ProjectDef.BuildDefs[i].BuildName < project.ProjectDef.BuildDefs[j].BuildName
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "buildDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.BuildDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}

		if project.ProjectDef.UpdatePackageDefs {
			sort.SliceStable(project.ProjectDef.PackageDefs, func(i, j int) bool {
				return project.ProjectDef.PackageDefs[i].PackageName < project.ProjectDef.PackageDefs[j].PackageName
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "packageDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.PackageDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}

		if project.ProjectDef.UpdateArtifactDefs {
			sort.SliceStable(project.ProjectDef.ArtifactDefs, func(i, j int) bool {
				return project.ProjectDef.ArtifactDefs[i].ArtifactName < project.ProjectDef.ArtifactDefs[j].ArtifactName
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "artifactDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.ArtifactDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}

		for _, pae := range project.ProjectAvailableEnvs {
			if pae.UpdateDeployContainerDefs {
				sort.SliceStable(pae.DeployContainerDefs, func(i, j int) bool {
					return pae.DeployContainerDefs[i].DeployName < pae.DeployContainerDefs[j].DeployName
				})
				defUpdate := pkg.DefUpdate{
					Kind:        "deployContainerDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.DeployContainerDefs,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}

			if pae.UpdateDeployArtifactDefs {
				sort.SliceStable(pae.DeployArtifactDefs, func(i, j int) bool {
					return pae.DeployArtifactDefs[i].DeployArtifactName < pae.DeployArtifactDefs[j].DeployArtifactName
				})
				defUpdate := pkg.DefUpdate{
					Kind:        "deployArtifactDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.DeployArtifactDefs,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}

			if pae.UpdateIstioDefs {
				sort.SliceStable(pae.IstioDefs, func(i, j int) bool {
					return pae.IstioDefs[i].DeployName < pae.IstioDefs[j].DeployName
				})
				defUpdate := pkg.DefUpdate{
					Kind:        "istioDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.IstioDefs,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}

			if pae.UpdateIstioGatewayDef {
				defUpdate := pkg.DefUpdate{
					Kind:        "istioGatewayDef",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.IstioGatewayDef,
					EnvName:     pae.EnvName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}

			for stepName, csd := range pae.CustomStepDefs {
				if csd.UpdateCustomStepModuleDefs {
					sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
						return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
					})
					defUpdate := pkg.DefUpdate{
						Kind:           "customStepDef",
						ProjectName:    project.ProjectInfo.ProjectName,
						Def:            csd,
						EnvName:        pae.EnvName,
						CustomStepName: stepName,
					}
					defUpdates = append(defUpdates, defUpdate)
				}
			}
		}

		for stepName, csd := range project.ProjectDef.CustomStepDefs {
			if csd.UpdateCustomStepModuleDefs {
				sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
					return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
				})
				defUpdate := pkg.DefUpdate{
					Kind:           "customStepDef",
					ProjectName:    project.ProjectInfo.ProjectName,
					Def:            csd,
					CustomStepName: stepName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}

		for _, pp := range project.ProjectPipelines {
			if pp.UpdatePipelineDef {
				defUpdate := pkg.DefUpdate{
					Kind:        "pipelineDef",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pp.PipelineDef,
					BranchName:  pp.BranchName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}

		if project.ProjectDef.UpdateCustomOpsDefs {
			sort.SliceStable(project.ProjectDef.CustomOpsDefs, func(i, j int) bool {
				return project.ProjectDef.CustomOpsDefs[i].CustomOpsName < project.ProjectDef.CustomOpsDefs[j].CustomOpsName
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "customOpsDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.CustomOpsDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}

		if project.ProjectDef.UpdateOpsBatchDefs {
			sort.SliceStable(project.ProjectDef.OpsBatchDefs, func(i, j int) bool {
				return project.ProjectDef.OpsBatchDefs[i].OpsBatchName < project.ProjectDef.OpsBatchDefs[j].OpsBatchName
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "opsBatchDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.OpsBatchDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}

		if project.ProjectDef.UpdateDockerIgnoreDefs {
			sort.SliceStable(project.ProjectDef.DockerIgnoreDefs, func(i, j int) bool {
				return project.ProjectDef.DockerIgnoreDefs[i] < project.ProjectDef.DockerIgnoreDefs[j]
			})
			defUpdate := pkg.DefUpdate{
				Kind:        "dockerIgnoreDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.DockerIgnoreDefs,
			}
			defUpdates = append(defUpdates, defUpdate)
		}
	}

	outputs := []map[string]interface{}{}
	for _, defUpdate := range defUpdates {
		out := map[string]interface{}{}
		m := map[string]interface{}{}
		bs, _ := json.Marshal(defUpdate)
		_ = json.Unmarshal(bs, &m)
		if o.Full {
			out = m
		} else {
			out = pkg.RemoveMapEmptyItems(m)
		}
		outputs = append(outputs, out)
	}

	bs = []byte{}
	if o.Output == "json" {
		bs, _ = json.MarshalIndent(outputs, "", "  ")
	} else if o.Output == "yaml" {
		bs, _ = pkg.YamlIndent(outputs)
	}
	if len(bs) > 0 {
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
			case "dockerIgnoreDefs":
				param["dockerIgnoreDefsYaml"] = string(bs)
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
