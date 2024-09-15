package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/go-playground/validator/v10"
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

type OptionsConsoleApply struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames    []string          `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Consoles     []pkg.ConsoleKind `yaml:"consoles" json:"consoles" bson:"consoles" validate:""`
		ProjectNames []string          `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
	}
}

func NewOptionsConsoleApply() *OptionsConsoleApply {
	var o OptionsConsoleApply
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdConsoleApply() *cobra.Command {
	o := NewOptionsConsoleApply()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`apply -f [filename]`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_console_apply_short")
	msgLong := OptCommon.TransLang("cmd_console_apply_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_console_apply_example", baseName, baseName))

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_console_apply_output"))
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, OptCommon.TransLang("param_console_apply_recursive"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_console_apply_full"))
	cmd.Flags().StringSliceVarP(&o.FileNames, "files", "f", []string{}, OptCommon.TransLang("param_console_apply_files"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_console_apply_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsConsoleApply) CheckConsoleKind(consoleKind pkg.ConsoleKind) error {
	var err error
	if consoleKind.Metadata.ProjectName == "" {
		err = fmt.Errorf("metadata.projectName can not be empty")
		return err
	}
	var foundProject bool
	for _, projectName := range o.Param.ProjectNames {
		if consoleKind.Metadata.ProjectName == projectName {
			foundProject = true
			break
		}
	}
	if !foundProject {
		err = fmt.Errorf("projectName %s not exists", consoleKind.Metadata.ProjectName)
		return err
	}

	validate := validator.New()
	switch consoleKind.Kind {
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindMember]:
		for _, item := range consoleKind.Items {
			var d pkg.ProjectMember
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}

			var foundAccessLevel bool
			for _, s := range pkg.AccessLevels {
				if s == d.AccessLevel {
					foundAccessLevel = true
					break
				}
			}
			if !foundAccessLevel {
				err = fmt.Errorf("%s/%s accessLevel %s not exists", consoleKind.Kind, d.Username, d.AccessLevel)
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]:
		for _, item := range consoleKind.Items {
			var d pkg.ProjectPipeline
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			if d.SourceBranch == "" {
				err = fmt.Errorf("%s/%s sourceBranch required", consoleKind.Kind, d.BranchName)
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipelineTrigger]:
		if consoleKind.Metadata.BranchName == "" {
			err = fmt.Errorf("kind is %s, metadata.branchName can not be empty", consoleKind.Kind)
			return err
		}
		for _, item := range consoleKind.Items {
			var d pkg.PipelineTrigger
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindHost]:
		if consoleKind.Metadata.EnvName == "" {
			err = fmt.Errorf("kind is %s, metadata.envName can not be empty", consoleKind.Kind)
			return err
		}
		for _, item := range consoleKind.Items {
			var d pkg.Host
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindDatabase]:
		if consoleKind.Metadata.EnvName == "" {
			err = fmt.Errorf("kind is %s, metadata.envName can not be empty", consoleKind.Kind)
			return err
		}
		for _, item := range consoleKind.Items {
			var d pkg.Database
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindComponent]:
		if consoleKind.Metadata.EnvName == "" {
			err = fmt.Errorf("kind is %s, metadata.envName can not be empty", consoleKind.Kind)
			return err
		}
		for _, item := range consoleKind.Items {
			var d pkg.Component
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
		}
	case pkg.ConsoleCmdKinds[pkg.ConsoleKindDebugComponent]:
		if consoleKind.Metadata.EnvName == "" {
			err = fmt.Errorf("kind is %s, metadata.envName can not be empty", consoleKind.Kind)
			return err
		}
		for _, item := range consoleKind.Items {
			var d pkg.ComponentDebug
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
			err = validate.Struct(d)
			if err != nil {
				err = fmt.Errorf("kind is %s, but item parse error: %s\n%s", consoleKind.Kind, err.Error(), string(bs))
				return err
			}
		}
	}
	return err
}

func GetConsoleKindsFromJson(fileName string, bs []byte) ([]pkg.ConsoleKind, error) {
	var err error
	consoleKinds := []pkg.ConsoleKind{}
	var list pkg.ConsoleKindList
	err = json.Unmarshal(bs, &list)
	if err == nil {
		if list.Kind == "list" {
			consoleKinds = append(consoleKinds, list.Consoles...)
		} else {
			var console pkg.ConsoleKind
			err = json.Unmarshal(bs, &console)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return consoleKinds, err
			}
			if console.Kind != "" {
				consoleKinds = append(consoleKinds, console)
			}
		}
	} else {
		var console pkg.ConsoleKind
		err = json.Unmarshal(bs, &console)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return consoleKinds, err
		}
		if console.Kind != "" {
			consoleKinds = append(consoleKinds, console)
		}
	}
	return consoleKinds, err
}

func GetConsoleKindsFromYaml(fileName string, bs []byte) ([]pkg.ConsoleKind, error) {
	var err error
	consoleKinds := []pkg.ConsoleKind{}
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
			return consoleKinds, err
		} else {
			ms = append(ms, m)
		}
	}
	for _, m := range ms {
		b, _ := yaml.Marshal(m)
		var list pkg.ConsoleKindList
		err = yaml.Unmarshal(b, &list)
		if err == nil {
			if list.Kind == "list" {
				consoleKinds = append(consoleKinds, list.Consoles...)
			} else {
				var console pkg.ConsoleKind
				err = yaml.Unmarshal(b, &console)
				if err != nil {
					err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
					return consoleKinds, err
				}
				if console.Kind != "" {
					consoleKinds = append(consoleKinds, console)
				}
			}
		} else {
			var console pkg.ConsoleKind
			err = yaml.Unmarshal(b, &console)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return consoleKinds, err
			}
			if console.Kind != "" {
				consoleKinds = append(consoleKinds, console)
			}
		}
	}

	return consoleKinds, err
}

func (o *OptionsConsoleApply) GetConsoleKinds(fileName string, bs []byte) ([]pkg.ConsoleKind, error) {
	var err error
	consoleKinds := []pkg.ConsoleKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		consoleKinds, err = GetConsoleKindsFromJson(fileName, bs)
		if err != nil {
			return consoleKinds, err
		}
	} else if ext == ".yaml" || ext == ".yml" {
		consoleKinds, err = GetConsoleKindsFromYaml(fileName, bs)
		if err != nil {
			return consoleKinds, err
		}
	} else if fileName == "" {
		consoleKinds, err = GetConsoleKindsFromJson(fileName, bs)
		if err != nil {
			consoleKinds, err = GetConsoleKindsFromYaml(fileName, bs)
			if err != nil {
				return consoleKinds, err
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return consoleKinds, err
	}

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/console/projectNames"), http.MethodGet, "", param, false)
	if err != nil {
		return consoleKinds, err
	}
	projectNames := []string{}
	rs := result.Get("data.projectNames").Array()
	for _, r := range rs {
		projectNames = append(projectNames, r.String())
	}
	o.Param.ProjectNames = projectNames

	for _, console := range consoleKinds {
		if console.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return consoleKinds, err
		}
		if console.Metadata.ProjectName == "" {
			err = fmt.Errorf("parse file %s error: metadata.projectName is empty", fileName)
			return consoleKinds, err
		}
		err = pkg.ValidateMinusNameID(console.Metadata.ProjectName)
		if err != nil {
			err = fmt.Errorf("parse file %s error: metadata.projectName %s format error: %s", fileName, console.Metadata.ProjectName, err.Error())
			return consoleKinds, err
		}

		var found bool

		var kinds []string
		for _, v := range pkg.ConsoleCmdKinds {
			if v != "" {
				kinds = append(kinds, v)
			}
		}
		for _, d := range kinds {
			if console.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, console.Kind)
			return consoleKinds, err
		}
		err = o.CheckConsoleKind(console)
		if err != nil {
			return consoleKinds, err
		}
	}
	return consoleKinds, err
}

func (o *OptionsConsoleApply) Complete(cmd *cobra.Command) error {
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

func (o *OptionsConsoleApply) Validate(args []string) error {
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
			err = fmt.Errorf("--files - required os.stdin\n example: echo 'xxx' | %s console apply -f -", baseName)
			return err
		}
		consoles, err := o.GetConsoleKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Consoles = append(o.Param.Consoles, consoles...)
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

			consoles, err := o.GetConsoleKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Consoles = append(o.Param.Consoles, consoles...)
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

func (o *OptionsConsoleApply) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	consoleKindList := pkg.ConsoleKindList{
		Kind:     "list",
		Consoles: o.Param.Consoles,
	}
	output := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(consoleKindList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		output = m
	} else {
		output = pkg.RemoveMapEmptyItems(m)
	}

	bs = []byte{}
	if o.Output == "json" {
		bs, _ = json.MarshalIndent(output, "", "  ")
	} else if o.Output == "yaml" {
		bs, _ = pkg.YamlIndent(output)
	}
	if len(bs) > 0 {
		fmt.Println(string(bs))
	}

	if !o.Try {
		for _, consoleKind := range o.Param.Consoles {
			logHeader := fmt.Sprintf("%s/%s", consoleKind.Metadata.ProjectName, consoleKind.Kind)
			switch consoleKind.Kind {
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindMember]:
				param := map[string]interface{}{}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/console/project/%s/members", consoleKind.Metadata.ProjectName), http.MethodGet, "", param, false)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				names := []string{}
				rs := result.Get("data.members").Array()
				for _, r := range rs {
					names = append(names, r.Get("username").String())
				}

				for _, item := range consoleKind.Items {
					var obj pkg.ProjectMember
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					var foundItem bool
					for _, s := range names {
						if s == obj.Username {
							foundItem = true
							break
						}
					}
					if !foundItem {
						param = map[string]interface{}{
							"username":           obj.Username,
							"accessLevel":        obj.AccessLevel,
							"disableProjectDefs": obj.DisableProjectDefs,
							"disableRepoSecrets": obj.DisableRepoSecrets,
							"disablePipelines":   obj.DisablePipelines,
						}
						result, _, err = o.QueryAPI(fmt.Sprintf("api/console/project/%s/memberAdd", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
						if err != nil {
							err = fmt.Errorf("%s: %s", logHeader, err.Error())
							return err
						}
						msg := result.Get("msg").String()
						log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					} else {
						param = map[string]interface{}{
							"username":           obj.Username,
							"accessLevel":        obj.AccessLevel,
							"disableProjectDefs": obj.DisableProjectDefs,
							"disableRepoSecrets": obj.DisableRepoSecrets,
							"disablePipelines":   obj.DisablePipelines,
						}
						result, _, err = o.QueryAPI(fmt.Sprintf("api/console/project/%s/memberUpdate", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
						if err != nil {
							err = fmt.Errorf("%s: %s", logHeader, err.Error())
							return err
						}
						msg := result.Get("msg").String()
						log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					}
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipeline]:
				projectConsole, err := o.GetConsoleProject(consoleKind.Metadata.ProjectName)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				names := []string{}
				for _, pipeine := range projectConsole.Pipelines {
					names = append(names, pipeine.BranchName)
				}

				for _, item := range consoleKind.Items {
					var obj pkg.ProjectPipeline
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					var foundItem bool
					for _, s := range names {
						if s == obj.BranchName {
							foundItem = true
							break
						}
					}
					if !foundItem {
						param := map[string]interface{}{
							"branchName":       obj.BranchName,
							"sourceBranch":     obj.SourceBranch,
							"envs":             obj.Envs,
							"envProductions":   obj.EnvProductions,
							"webhookPushEvent": obj.WebhookPushEvent,
						}
						result, _, err := o.QueryAPI(fmt.Sprintf("api/console/project/%s/pipelineAdd", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
						if err != nil {
							err = fmt.Errorf("%s: %s", logHeader, err.Error())
							return err
						}
						msg := result.Get("msg").String()
						log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
						auditID := result.Get("data.auditID").String()
						if auditID == "" {
							err = fmt.Errorf("can not get auditID")
							return err
						}
						url := fmt.Sprintf("api/ws/log/audit/console/%s", auditID)
						err = o.QueryWebsocket(url, "")
						if err != nil {
							return err
						}
						log.Info(fmt.Sprintf("##############################"))
						log.Success(fmt.Sprintf("# %s create finish", logHeader))
					}
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindPipelineTrigger]:
				param := map[string]interface{}{}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/console/project/%s/pipelineTriggerStepActions", consoleKind.Metadata.ProjectName), http.MethodGet, "", param, false)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				stepActions := []string{}
				rs := result.Get("data.stepActions").Array()
				for _, r := range rs {
					stepActions = append(stepActions, r.Get("value").String())
				}

				for _, item := range consoleKind.Items {
					var obj pkg.PipelineTrigger
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					var foundStepAction bool
					for _, s := range stepActions {
						if s == obj.StepAction {
							foundStepAction = true
							break
						}
					}
					if foundStepAction {
						bs, _ = json.Marshal(obj)
						param = map[string]interface{}{}
						_ = json.Unmarshal(bs, &param)
						delete(param, "stepAction")
						param["branchName"] = consoleKind.Metadata.BranchName
						param["stepActions"] = []string{obj.StepAction}
						result, _, err = o.QueryAPI(fmt.Sprintf("api/console/project/%s/pipelineTriggerOp", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
						if err != nil {
							err = fmt.Errorf("%s: %s", logHeader, err.Error())
							return err
						}
						msg := result.Get("msg").String()
						log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					} else {
						err = fmt.Errorf("%s: stepAction %s not correct", logHeader, obj.StepAction)
						return err
					}
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindHost]:
				projectConsole, err := o.GetConsoleProject(consoleKind.Metadata.ProjectName)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				names := []string{}
				for _, pae := range projectConsole.ProjectAvailableEnvs {
					if pae.EnvName == consoleKind.Metadata.EnvName {
						for _, host := range pae.Hosts {
							names = append(names, host.HostName)
						}
					}
				}

				for _, item := range consoleKind.Items {
					var obj pkg.Host
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					var foundItem bool
					for _, s := range names {
						if s == obj.HostName {
							foundItem = true
							break
						}
					}
					bs, _ = json.Marshal(obj)
					param := map[string]interface{}{}
					_ = json.Unmarshal(bs, &param)
					param["envName"] = consoleKind.Metadata.EnvName
					var url string
					if !foundItem {
						url = fmt.Sprintf("api/console/project/%s/envHostAdd", consoleKind.Metadata.ProjectName)
					} else {
						url = fmt.Sprintf("api/console/project/%s/envHostUpdate", consoleKind.Metadata.ProjectName)
					}
					result, _, err := o.QueryAPI(url, http.MethodPost, "", param, false)
					if err != nil {
						err = fmt.Errorf("%s: %s", logHeader, err.Error())
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindDatabase]:
				projectConsole, err := o.GetConsoleProject(consoleKind.Metadata.ProjectName)
				if err != nil {
					err = fmt.Errorf("%s: %s", logHeader, err.Error())
					return err
				}
				names := []string{}
				for _, pae := range projectConsole.ProjectAvailableEnvs {
					if pae.EnvName == consoleKind.Metadata.EnvName {
						for _, database := range pae.Databases {
							names = append(names, database.DbName)
						}
					}
				}

				for _, item := range consoleKind.Items {
					var obj pkg.Database
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					var foundItem bool
					for _, s := range names {
						if s == obj.DbName {
							foundItem = true
							break
						}
					}
					bs, _ = json.Marshal(obj)
					param := map[string]interface{}{}
					_ = json.Unmarshal(bs, &param)
					param["envName"] = consoleKind.Metadata.EnvName
					var url string
					if !foundItem {
						url = fmt.Sprintf("api/console/project/%s/envDbAdd", consoleKind.Metadata.ProjectName)
					} else {
						url = fmt.Sprintf("api/console/project/%s/envDbUpdate", consoleKind.Metadata.ProjectName)
					}
					result, _, err := o.QueryAPI(url, http.MethodPost, "", param, false)
					if err != nil {
						err = fmt.Errorf("%s: %s", logHeader, err.Error())
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindComponent]:
				for _, item := range consoleKind.Items {
					var obj pkg.Component
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					bsYaml, _ := pkg.YamlIndent(obj.DeploySpecStatic)
					param := map[string]interface{}{
						"envName":       consoleKind.Metadata.EnvName,
						"arch":          obj.Arch,
						"componentDesc": obj.ComponentDesc,
						"componentYaml": string(bsYaml),
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/console/project/%s/envComponentUpdate", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
					if err != nil {
						err = fmt.Errorf("%s: %s", logHeader, err.Error())
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					auditID := result.Get("data.auditID").String()
					if auditID == "" {
						err = fmt.Errorf("can not get auditID")
						return err
					}
					url := fmt.Sprintf("api/ws/log/audit/console/%s", auditID)
					err = o.QueryWebsocket(url, "")
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("##############################"))
					log.Success(fmt.Sprintf("# %s finish", logHeader))
				}
			case pkg.ConsoleCmdKinds[pkg.ConsoleKindDebugComponent]:
				for _, item := range consoleKind.Items {
					var obj pkg.ComponentDebug
					bs, _ = json.Marshal(item)
					_ = json.Unmarshal(bs, &obj)
					bsYaml, _ := pkg.YamlIndent(obj.DeploySpecDebug)
					param := map[string]interface{}{
						"envName":            consoleKind.Metadata.EnvName,
						"componentDebugYaml": string(bsYaml),
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/console/project/%s/envComponentDebugUpdate", consoleKind.Metadata.ProjectName), http.MethodPost, "", param, false)
					if err != nil {
						err = fmt.Errorf("%s: %s", logHeader, err.Error())
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					auditID := result.Get("data.auditID").String()
					if auditID == "" {
						err = fmt.Errorf("can not get auditID")
						return err
					}
					url := fmt.Sprintf("api/ws/log/audit/console/%s", auditID)
					err = o.QueryWebsocket(url, "")
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("##############################"))
					log.Success(fmt.Sprintf("# %s finish", logHeader))
				}
			}
		}
	}

	return err
}
