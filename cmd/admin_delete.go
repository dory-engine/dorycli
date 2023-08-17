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
	msgExample := fmt.Sprintf(`kind: ct(component templates), env(kubernetes environments), step(custom steps), user(users), dbe(docker build environments), grc(git repository configs), irc(image repository configs), arc(artifact repository configs), scrc(scan code repository configs)

  # delete users, admin permission required
  %s admin delete user test-user01 test-user02

  # delete custom step configurations, admin permission required
  %s admin delete step customStepName1 customStepName2

  # delete kubernetes environment configurations, admin permission required
  %s admin delete env test uat

  # delete component template configurations, admin permission required
  %s admin delete ct mysql-v8`, baseName, baseName, baseName, baseName)

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
			case "user":
				itemNames, _ = o.GetUserNames()
			case "step":
				itemNames, _ = o.GetStepNames()
			case "env":
				itemNames, _ = o.GetEnvNames()
			case "ct":
				itemNames, _ = o.GetComponentTemplateNames()
			case "dbe":
				itemNames, _ = o.GetBuildEnvNames()
			case "grc":
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.GitRepoNames
			case "irc":
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.ImageRepoNames
			case "arc":
				repoNames, _ := o.GetRepoNames()
				itemNames = repoNames.ArtifactRepoNames
			case "scrc":
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
		case "user":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/user/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "step":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConf/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "env":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/env/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "ct":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplate/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "dbe":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnv/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "grc":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/gitRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "irc":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/imageRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "arc":
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/artifactRepoConfig/%s", itemName), http.MethodDelete, "", param, false)
			if err != nil {
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		case "scrc":
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
