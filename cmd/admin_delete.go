package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type OptionsAdminDelete struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Param          struct {
		Kind      string   `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ItemNames []string `yaml:"itemNames" json:"itemNames" bson:"itemNames" validate:""`
	}
}

func NewOptionsAdminDelete() *OptionsAdminDelete {
	var o OptionsAdminDelete
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminDelete() *cobra.Command {
	o := NewOptionsAdminDelete()

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`delete [kind] [itemName1] [itemName2]...
# kind options: %s`, strings.Join(adminCmdKinds, " / "))
	msgShort := fmt.Sprintf("delete configurations, admin permission required")
	msgLong := fmt.Sprintf(`delete configurations in dory-engine server, admin permission required`)
	msgExample := fmt.Sprintf(`kind: %s

  # delete users, admin permission required
  %s admin delete %s test-user01 test-user02

  # delete custom step configurations, admin permission required
  %s admin delete %s customStepName1 customStepName2

  # delete kubernetes environment configurations, admin permission required
  %s admin delete %s test uat

  # delete component template configurations, admin permission required
  %s admin delete %s mysql-v8`, strings.Join(pkg.AdminKinds, ", "), baseName, pkg.AdminKindUser, baseName, pkg.AdminKindCustomStep, baseName, pkg.AdminKindEnvK8s, baseName, pkg.AdminKindComponentTemplate)

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

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsAdminDelete) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return adminCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) >= 1 {
			kind := args[0]
			itemNames := []string{}
			switch kind {
			case pkg.AdminKindUser:
				itemNames, _ = o.GetUserNames()
			case pkg.AdminKindCustomStep:
				itemNames, _ = o.GetStepNames()
			case pkg.AdminKindEnvK8s:
				itemNames, _ = o.GetEnvNames()
			case pkg.AdminKindComponentTemplate:
				itemNames, _ = o.GetComponentTemplateNames()
			case pkg.AdminKindDockerBuildEnv:
				itemNames, _ = o.GetBuildEnvNames()
			case pkg.AdminKindGitRepoConfig:
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.GitRepoNames
			case pkg.AdminKindImageRepoConfig:
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.ImageRepoNames
			case pkg.AdminKindArtifactRepoConfig:
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.ArtifactRepoNames
			case pkg.AdminKindScanCodeRepoConfig:
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.ScanCodeRepoNames
			default:
				err = fmt.Errorf("kind not correct")
			}
			if err != nil {
				return itemNames, cobra.ShellCompDirectiveNoFileComp
			}
			return itemNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return err
}

func (o *OptionsAdminDelete) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("kind required")
		return err
	}
	var kind string
	kind = args[0]

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	var found bool
	for _, cmdKind := range adminCmdKinds {
		if kind == cmdKind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct: kind options: %s", kind, strings.Join(adminCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	if len(args) < 2 {
		err = fmt.Errorf("itemName to delete required")
		return err
	}

	o.Param.ItemNames = args[1:]

	return err
}

func (o *OptionsAdminDelete) Run(args []string) error {
	var err error
	for _, itemName := range o.Param.ItemNames {
		logHeader := fmt.Sprintf("delete %s/%s", pkg.AdminCmdKinds[o.Param.Kind], itemName)
		switch o.Param.Kind {
		case pkg.AdminKindUser:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/user/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindCustomStep:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConf/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindEnvK8s:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/env/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindComponentTemplate:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplate/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindDockerBuildEnv:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnv/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindGitRepoConfig:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/gitRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindImageRepoConfig:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/imageRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindArtifactRepoConfig:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/artifactRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case pkg.AdminKindScanCodeRepoConfig:
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/scanCodeRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		}

	}
	return err
}
