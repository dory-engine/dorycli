package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type OptionsAdminApply struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames []string        `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Items     []pkg.AdminKind `yaml:"items" json:"items" bson:"items" validate:""`
	}
}

func NewOptionsAdminApply() *OptionsAdminApply {
	var o OptionsAdminApply
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminApply() *cobra.Command {
	o := NewOptionsAdminApply()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`apply -f [filename]`)

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_admin_apply_short")
	msgLong := OptCommon.TransLang("cmd_admin_apply_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_admin_apply_example", baseName, baseName))

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", OptCommon.TransLang("param_admin_apply_output"))
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, OptCommon.TransLang("param_admin_apply_recursive"))
	cmd.Flags().BoolVar(&o.Full, "full", false, OptCommon.TransLang("param_admin_apply_full"))
	cmd.Flags().StringSliceVarP(&o.FileNames, "files", "f", []string{}, OptCommon.TransLang("param_admin_apply_files"))
	cmd.Flags().BoolVar(&o.Try, "try", false, OptCommon.TransLang("param_admin_apply_try"))

	CheckError(o.Complete(cmd))
	return cmd
}

func CheckAdminKind(item pkg.AdminKind) error {
	var err error
	switch item.Kind {
	case pkg.AdminCmdKinds[pkg.AdminKindUser]:
		var spec pkg.User
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindCustomStep]:
		var spec pkg.CustomStepConf
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindEnvK8s]:
		var spec pkg.EnvK8s
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindComponentTemplate]:
		var spec pkg.ComponentTemplate
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindDockerBuildEnv]:
		var spec pkg.DockerBuildEnv
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindGitRepoConfig]:
		var spec pkg.GitRepoConfig
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindImageRepoConfig]:
		var spec pkg.ImageRepoConfig
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindArtifactRepoConfig]:
		var spec pkg.ArtifactRepoConfig
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindScanCodeRepoConfig]:
		var spec pkg.ScanCodeRepoConfig
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	case pkg.AdminCmdKinds[pkg.AdminKindAdminWebhook]:
		var spec pkg.AdminWebhook
		bs, _ := json.Marshal(item.Spec)
		err = json.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
			return err
		}
	}
	return err
}

func GetAdminKindsFromJson(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
	var list pkg.AdminKindList
	err = json.Unmarshal(bs, &list)
	if err == nil {
		if list.Kind == "list" {
			items = append(items, list.Items...)
		} else {
			var item pkg.AdminKind
			err = json.Unmarshal(bs, &item)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return items, err
			}
			if item.Kind != "" {
				items = append(items, item)
			}
		}
	} else {
		var item pkg.AdminKind
		err = json.Unmarshal(bs, &item)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return items, err
		}
		if item.Kind != "" {
			items = append(items, item)
		}
	}
	return items, err
}

func GetAdminKindsFromYaml(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
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
			return items, err
		} else {
			ms = append(ms, m)
		}
	}
	for _, m := range ms {
		b, _ := json.Marshal(m)
		var list pkg.AdminKindList
		err = json.Unmarshal(b, &list)
		if err == nil {
			if list.Kind == "list" {
				items = append(items, list.Items...)
			} else {
				var item pkg.AdminKind
				err = json.Unmarshal(b, &item)
				if err != nil {
					err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
					return items, err
				}
				if item.Kind != "" {
					items = append(items, item)
				}
			}
		} else {
			var item pkg.AdminKind
			err = json.Unmarshal(b, &item)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return items, err
			}
			if item.Kind != "" {
				items = append(items, item)
			}
		}
	}

	return items, err
}

func GetAdminKinds(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		items, err = GetAdminKindsFromJson(fileName, bs)
		if err != nil {
			return items, err
		}
	} else if ext == ".yaml" || ext == ".yml" {
		items, err = GetAdminKindsFromYaml(fileName, bs)
		if err != nil {
			return items, err
		}
	} else if fileName == "" {
		items, err = GetAdminKindsFromJson(fileName, bs)
		if err != nil {
			items, err = GetAdminKindsFromYaml(fileName, bs)
			if err != nil {
				return items, err
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return items, err
	}

	for i, item := range items {
		if item.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return items, err
		}
		if item.Kind != pkg.AdminCmdKinds[pkg.AdminKindAdminWebhook] && item.Metadata.Name == "" {
			err = fmt.Errorf("parse file %s error: metadata.name is empty", fileName)
			return items, err
		}

		var found bool

		var kinds []string
		for _, v := range pkg.AdminCmdKinds {
			if v != "" {
				kinds = append(kinds, v)
			}
		}
		for _, d := range kinds {
			if item.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, item.Kind)
			return items, err
		}
		switch item.Kind {
		case pkg.AdminCmdKinds[pkg.AdminKindUser]:
			var spec pkg.User
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindCustomStep]:
			var spec pkg.CustomStepConf
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindEnvK8s]:
			var spec pkg.EnvK8s
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindComponentTemplate]:
			var spec pkg.ComponentTemplate
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindDockerBuildEnv]:
			var spec pkg.DockerBuildEnv
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindGitRepoConfig]:
			var spec pkg.GitRepoConfig
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindImageRepoConfig]:
			var spec pkg.ImageRepoConfig
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindArtifactRepoConfig]:
			var spec pkg.ArtifactRepoConfig
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindScanCodeRepoConfig]:
			var spec pkg.ScanCodeRepoConfig
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		case pkg.AdminCmdKinds[pkg.AdminKindAdminWebhook]:
			var spec pkg.AdminWebhook
			bs, _ := json.Marshal(item.Spec)
			err = json.Unmarshal(bs, &spec)
			if err != nil {
				err = fmt.Errorf("kind is %s, but spec parse error: %s\n%s", item.Kind, err.Error(), string(bs))
				return items, err
			}
			item.Spec = spec
		}
		err = CheckAdminKind(item)
		if err != nil {
			return items, err
		}
		items[i] = item
	}
	return items, err
}

func (o *OptionsAdminApply) Complete(cmd *cobra.Command) error {
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

func (o *OptionsAdminApply) Validate(args []string) error {
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
			err = fmt.Errorf("--files - required os.stdin\n example: echo 'xxx' | %s admin apply -f -", baseName)
			return err
		}
		items, err := GetAdminKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Items = append(o.Param.Items, items...)
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
					infos, err := os.ReadDir(fileName)
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

			items, err := GetAdminKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Items = append(o.Param.Items, items...)
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

func (o *OptionsAdminApply) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	adminKindList := pkg.AdminKindList{
		Kind:  "list",
		Items: o.Param.Items,
	}
	output := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(adminKindList)
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
		param := map[string]interface{}{}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/projectNames"), http.MethodGet, "", param, false)
		if err != nil {
			return err
		}
		projectNames := []string{}
		rs := result.Get("data.projectNames").Array()
		for _, r := range rs {
			projectNames = append(projectNames, r.String())
		}

		for _, item := range o.Param.Items {
			logHeader := fmt.Sprintf("%s/%s", item.Kind, item.Metadata.Name)

			switch item.Kind {
			case pkg.AdminCmdKinds[pkg.AdminKindUser]:
				var user pkg.User
				switch v := item.Spec.(type) {
				case pkg.User:
					user = v
				}
				// do not update user isAdmin flag
				user.IsAdmin = false
				param = map[string]interface{}{}
				bs, _ = json.Marshal(user)
				_ = json.Unmarshal(bs, &param)
				result, _, err = o.QueryAPI(fmt.Sprintf("api/admin/user"), http.MethodPut, "", param, false)
				if err != nil {
					return err
				}
				msg := result.Get("msg").String()
				log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
				val, ok := item.Metadata.Annotations["userProjects"]
				if ok {
					arr := strings.Split(val, ",")
					for _, a := range arr {
						arr2 := strings.Split(a, ":")
						if len(arr2) == 2 {
							projectName := arr2[0]
							var foundProject bool
							for _, name := range projectNames {
								if name == projectName {
									foundProject = true
									break
								}
							}
							roleName := arr2[1]
							var foundAccessLevel bool
							for _, accessLevel := range pkg.AccessLevels {
								if accessLevel == roleName {
									foundAccessLevel = true
									break
								}
							}
							if foundProject && foundAccessLevel {
								result, _, err = o.QueryAPI(fmt.Sprintf("api/admin/user/%s/projects", user.Username), http.MethodGet, "", param, false)
								if err != nil {
									return err
								}
								userProjects := []pkg.UserProject{}
								err = json.Unmarshal([]byte(result.Get("data.projects").Raw), &userProjects)
								if err != nil {
									return err
								}
								var userProject pkg.UserProject
								for _, up := range userProjects {
									if up.ProjectName == projectName {
										userProject = up
										break
									}
								}
								if userProject.ProjectName != "" {
									if userProject.AccessLevel != roleName {
										// update project member
										time.Sleep(time.Second * 2)
										param = map[string]interface{}{
											"projectName": projectName,
											"accessLevel": roleName,
										}
										result, _, err = o.QueryAPI(fmt.Sprintf("api/admin/user/%s/projectUpdate", user.Username), http.MethodPost, "", param, false)
										if err != nil {
											return err
										}
										msg := result.Get("msg").String()
										log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
									}
								} else {
									// add project member
									time.Sleep(time.Second * 2)
									param = map[string]interface{}{
										"projectName": projectName,
										"accessLevel": roleName,
									}
									result, _, err = o.QueryAPI(fmt.Sprintf("api/admin/user/%s/projectAdd", user.Username), http.MethodPost, "", param, false)
									if err != nil {
										return err
									}
									msg := result.Get("msg").String()
									log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
								}
							}
						}
					}
				}
			case pkg.AdminCmdKinds[pkg.AdminKindCustomStep]:
				param := map[string]interface{}{
					"adminWebhooks": []string{item.Metadata.Name},
					"page":          1,
					"perPage":       1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConfs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				customStepNames := []string{}
				err = json.Unmarshal([]byte(result.Get("data.customStepNames").Raw), &customStepNames)
				if err != nil {
					return err
				}
				var found bool
				for _, name := range customStepNames {
					if name == item.Metadata.Name {
						found = true
						break
					}
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(item.Spec)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				customStepConfYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"customStepConfYaml": customStepConfYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConf/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"customStepConfYaml": customStepConfYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConf"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindEnvK8s]:
				param := map[string]interface{}{}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/envNames"), http.MethodGet, "", param, false)
				if err != nil {
					return err
				}
				envNames := []string{}
				err = json.Unmarshal([]byte(result.Get("data.envNames").Raw), &envNames)
				if err != nil {
					return err
				}
				var found bool
				for _, name := range envNames {
					if name == item.Metadata.Name {
						found = true
						break
					}
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(item.Spec)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				envK8sYaml := string(bs)

				var auditID string
				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"envK8sYaml": envK8sYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/env/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					auditID = result.Get("data.auditID").String()
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"envK8sYaml": envK8sYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/env"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
					auditID = result.Get("data.auditID").String()
				}

				if auditID == "" {
					err = fmt.Errorf("can not get auditID")
					return err
				}

				url := fmt.Sprintf("api/ws/log/audit/admin/%s", auditID)
				err = o.QueryWebsocket(url, "")
				if err != nil {
					return err
				}
				log.Info(fmt.Sprintf("##############################"))
				log.Success(fmt.Sprintf("# %s %s finish", logHeader, op))
			case pkg.AdminCmdKinds[pkg.AdminKindComponentTemplate]:
				param := map[string]interface{}{
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplates"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				componentTemplates := []pkg.ComponentTemplate{}
				err = json.Unmarshal([]byte(result.Get("data.componentTemplates").Raw), &componentTemplates)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range componentTemplates {
					if tpl.ComponentTemplateName == item.Metadata.Name {
						found = true
						break
					}
				}

				var componentTemplateDesc string
				var deploySpecStatic pkg.DeploySpecStatic
				switch tpl := item.Spec.(type) {
				case pkg.ComponentTemplate:
					componentTemplateDesc = tpl.ComponentTemplateDesc
					deploySpecStatic = tpl.DeploySpecStatic
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(deploySpecStatic)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				componentTemplateYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"componentTemplateName": item.Metadata.Name,
						"componentTemplateDesc": componentTemplateDesc,
						"componentTemplateYaml": componentTemplateYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplate/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"componentTemplateName": item.Metadata.Name,
						"componentTemplateDesc": componentTemplateDesc,
						"componentTemplateYaml": componentTemplateYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplate"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindDockerBuildEnv]:
				param := map[string]interface{}{
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnvs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				dockerBuildEnvs := []pkg.DockerBuildEnv{}
				err = json.Unmarshal([]byte(result.Get("data.dockerBuildEnvs").Raw), &dockerBuildEnvs)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range dockerBuildEnvs {
					if tpl.BuildEnvName == item.Metadata.Name {
						found = true
						break
					}
				}

				var dockerBuildEnv pkg.DockerBuildEnv
				switch tpl := item.Spec.(type) {
				case pkg.DockerBuildEnv:
					dockerBuildEnv = tpl
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(dockerBuildEnv)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				dockerBuildEnvYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"dockerBuildEnvYaml": dockerBuildEnvYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnv/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"dockerBuildEnvYaml": dockerBuildEnvYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnv"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindGitRepoConfig]:
				param := map[string]interface{}{
					"types":   []string{"gitRepoConfig"},
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				gitRepoConfigs := []pkg.GitRepoConfig{}
				err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &gitRepoConfigs)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range gitRepoConfigs {
					if tpl.RepoName == item.Metadata.Name {
						found = true
						break
					}
				}

				var gitRepoConfig pkg.GitRepoConfig
				switch tpl := item.Spec.(type) {
				case pkg.GitRepoConfig:
					gitRepoConfig = tpl
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(gitRepoConfig)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				gitRepoConfigYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"gitRepoConfigYaml": gitRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/gitRepoConfig/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"gitRepoConfigYaml": gitRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/gitRepoConfig"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindImageRepoConfig]:
				param := map[string]interface{}{
					"types":   []string{"imageRepoConfig"},
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				imageRepoConfigs := []pkg.ImageRepoConfig{}
				err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &imageRepoConfigs)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range imageRepoConfigs {
					if tpl.RepoName == item.Metadata.Name {
						found = true
						break
					}
				}

				var imageRepoConfig pkg.ImageRepoConfig
				switch tpl := item.Spec.(type) {
				case pkg.ImageRepoConfig:
					imageRepoConfig = tpl
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(imageRepoConfig)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				imageRepoConfigYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"imageRepoConfigYaml": imageRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/imageRepoConfig/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"imageRepoConfigYaml": imageRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/imageRepoConfig"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindArtifactRepoConfig]:
				param := map[string]interface{}{
					"types":   []string{"artifactRepoConfig"},
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				artifactRepoConfigs := []pkg.ArtifactRepoConfig{}
				err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &artifactRepoConfigs)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range artifactRepoConfigs {
					if tpl.RepoName == item.Metadata.Name {
						found = true
						break
					}
				}

				var artifactRepoConfig pkg.ArtifactRepoConfig
				switch tpl := item.Spec.(type) {
				case pkg.ArtifactRepoConfig:
					artifactRepoConfig = tpl
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(artifactRepoConfig)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				artifactRepoConfigYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"artifactRepoConfigYaml": artifactRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/artifactRepoConfig/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"artifactRepoConfigYaml": artifactRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/artifactRepoConfig"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindScanCodeRepoConfig]:
				param := map[string]interface{}{
					"types":   []string{"scanCodeRepoConfig"},
					"page":    1,
					"perPage": 1000,
				}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
				if err != nil {
					return err
				}
				scanCodeRepoConfigs := []pkg.ScanCodeRepoConfig{}
				err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &scanCodeRepoConfigs)
				if err != nil {
					return err
				}
				var found bool
				for _, tpl := range scanCodeRepoConfigs {
					if tpl.RepoName == item.Metadata.Name {
						found = true
						break
					}
				}

				var scanCodeRepoConfig pkg.ScanCodeRepoConfig
				switch tpl := item.Spec.(type) {
				case pkg.ScanCodeRepoConfig:
					scanCodeRepoConfig = tpl
				}

				m := map[string]interface{}{}
				bs, _ := json.Marshal(scanCodeRepoConfig)
				_ = json.Unmarshal(bs, &m)
				pm := pkg.RemoveMapEmptyItems(m)
				bs, _ = pkg.YamlIndent(pm)
				scanCodeRepoConfigYaml := string(bs)

				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{
						"scanCodeRepoConfigYaml": scanCodeRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/scanCodeRepoConfig/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{
						"scanCodeRepoConfigYaml": scanCodeRepoConfigYaml,
					}
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/scanCodeRepoConfig"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			case pkg.AdminCmdKinds[pkg.AdminKindAdminWebhook]:
				param := map[string]interface{}{}
				result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/adminWebhooks"), http.MethodGet, "", param, false)
				if err != nil {
					return err
				}
				adminWebhooks := []pkg.AdminWebhook{}
				err = json.Unmarshal([]byte(result.Get("data.adminWebhooks").Raw), &adminWebhooks)
				if err != nil {
					return err
				}
				var found bool
				for _, adminWebhook := range adminWebhooks {
					if adminWebhook.AdminWebhookID == item.Metadata.Name {
						found = true
						break
					}
				}
				var op string
				if found {
					// update
					op = "update"
					param := map[string]interface{}{}
					bs, _ = json.Marshal(item.Spec)
					_ = json.Unmarshal(bs, &param)
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/adminWebhook/%s", item.Metadata.Name), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				} else {
					// add
					op = "add"
					param := map[string]interface{}{}
					bs, _ = json.Marshal(item.Spec)
					_ = json.Unmarshal(bs, &param)
					result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/adminWebhook"), http.MethodPost, "", param, false)
					if err != nil {
						return err
					}
					msg := result.Get("msg").String()
					log.Info(fmt.Sprintf("%s %s: %s", logHeader, op, msg))
				}
			}
		}
	}

	return err
}
