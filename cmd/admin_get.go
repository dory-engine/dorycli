package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

type OptionsAdminGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Full           bool   `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kinds     []string `yaml:"kinds" json:"kinds" bson:"kinds" validate:""`
		ItemNames []string `yaml:"itemNames" json:"itemNames" bson:"itemNames" validate:""`
		IsAllKind bool     `yaml:"isAllKind" json:"isAllKind" bson:"isAllKind" validate:""`
	}
}

func NewOptionsAdminGet() *OptionsAdminGet {
	var o OptionsAdminGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminGet() *cobra.Command {
	o := NewOptionsAdminGet()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf(`get [kind],[kind]... [itemName1] [itemName2]... [--output=json|yaml]`)
	msgShort := fmt.Sprintf("get configurations, admin permission required")
	msgLong := fmt.Sprintf(`get users, custom steps, kubernetes environments and component templates configurations in dory-engine server, admin permission required`)
	msgExample := fmt.Sprintf(`kind: all, %s

  # get all configurations, admin permission required
  %s admin get %s --output=yaml

  # get all configurations, and show in full version, admin permission required
  %s admin get %s --output=yaml --full

  # get custom steps and component templates configurations, admin permission required
  %s admin get %s,%s

  # get users configurations, and filter by userNames, admin permission required
  %s admin get %s test-user1 test-user2

  # get kubernetes environments configurations, and filter by envNames, admin permission required
  %s admin get %s test uat prod`, strings.Join(pkg.AdminKinds, ", "), baseName, pkg.AdminKindAll, baseName, pkg.AdminKindAll, baseName, pkg.AdminKindCustomStep, pkg.AdminKindComponentTemplate, baseName, pkg.AdminKindUser, baseName, pkg.AdminKindEnvK8s)

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project configurations in full version, use with --output option")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsAdminGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	adminCmdKinds := []string{}
	for k, _ := range pkg.AdminCmdKinds {
		adminCmdKinds = append(adminCmdKinds, k)
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return adminCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) >= 1 {
			kindStr := args[0]
			var isAllKind bool
			kinds := strings.Split(kindStr, ",")
			for _, kind := range kinds {
				if kind == pkg.AdminKindAll {
					isAllKind = true
				}
			}
			if len(kinds) == 1 && !isAllKind {
				kind := kinds[0]
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
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return itemNames, cobra.ShellCompDirectiveNoFileComp
			} else {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsAdminGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("kind required")
		return err
	}
	var kinds, kindParams []string
	kindsStr := args[0]
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
		for cmdKind, _ := range pkg.AdminCmdKinds {
			if kind == cmdKind {
				found = true
				break
			}
		}
		if !found {
			adminCmdKinds := []string{}
			for k, _ := range pkg.AdminCmdKinds {
				adminCmdKinds = append(adminCmdKinds, k)
			}
			err = fmt.Errorf("kind %s format error: not correct, options: %s", kind, strings.Join(adminCmdKinds, " / "))
			return err
		}
		if kind == pkg.AdminKindAll {
			foundAll = true
		}
		kindParams = append(kindParams, pkg.AdminCmdKinds[kind])
	}
	if foundAll == true {
		o.Param.IsAllKind = true
	}
	o.Param.Kinds = kindParams

	o.Param.ItemNames = []string{}
	if len(args) > 1 {
		o.Param.ItemNames = args[1:]
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsAdminGet) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	var foundKindUser bool
	foundKindUser = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindUser {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindUser] {
			foundKindUser = true
			break
		}
	}

	var foundKindStep bool
	foundKindStep = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindStep {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindCustomStep] {
			foundKindStep = true
			break
		}
	}

	var foundKindEnv bool
	foundKindEnv = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindEnv {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindEnvK8s] {
			foundKindEnv = true
			break
		}
	}

	var foundKindCt bool
	foundKindCt = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindCt {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindComponentTemplate] {
			foundKindCt = true
			break
		}
	}

	var foundKindDbe bool
	foundKindDbe = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindDbe {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindDockerBuildEnv] {
			foundKindDbe = true
			break
		}
	}

	var foundKindGrc bool
	foundKindGrc = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindGrc {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindGitRepoConfig] {
			foundKindGrc = true
			break
		}
	}

	var foundKindIrc bool
	foundKindIrc = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindIrc {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindImageRepoConfig] {
			foundKindIrc = true
			break
		}
	}

	var foundKindArc bool
	foundKindArc = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindArc {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindArtifactRepoConfig] {
			foundKindArc = true
			break
		}
	}

	var foundKindScrc bool
	foundKindScrc = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindScrc {
			break
		} else if kind == pkg.AdminCmdKinds[pkg.AdminKindScanCodeRepoConfig] {
			foundKindScrc = true
			break
		}
	}

	adminKindList := pkg.AdminKindList{
		Kind: "list",
	}
	adminKinds := []pkg.AdminKind{}

	userFilters := []pkg.UserDetail{}
	if foundKindUser {
		param := map[string]interface{}{
			"sortMode": "username",
			"page":     1,
			"perPage":  1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/users"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		users := []pkg.UserDetail{}
		err = json.Unmarshal([]byte(result.Get("data.users").Raw), &users)
		if err != nil {
			return err
		}

		for _, user := range users {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			} else {
				for _, name := range o.Param.ItemNames {
					if name == user.Username {
						found = true
						break
					}
				}
			}
			if found {
				userFilters = append(userFilters, user)
			}
		}
		for _, user := range userFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = pkg.AdminKindUser
			adminKind.Metadata.Name = user.Username
			var userProjects []string
			for _, up := range user.UserProjects {
				userProjects = append(userProjects, fmt.Sprintf("%s:%s", up.ProjectName, up.AccessLevel))
			}
			adminKind.Metadata.Annotations = map[string]string{
				"avatarUrl":    user.AvatarUrl,
				"createTime":   user.CreateTime,
				"lastLogin":    user.LastLogin,
				"userProjects": strings.Join(userProjects, ","),
			}
			spec := pkg.User{
				Username:     user.Username,
				TenantCode:   user.TenantCode,
				TenantAdmins: user.TenantAdmins,
				Name:         user.Name,
				Mail:         user.Mail,
				Mobile:       user.Mobile,
				IsAdmin:      user.IsAdmin,
				IsActive:     user.IsActive,
			}
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	stepFilters := []pkg.CustomStepConfDetail{}
	if foundKindStep {
		param := map[string]interface{}{
			"customStepNames": o.Param.ItemNames,
			"page":            1,
			"perPage":         1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConfs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(result.Get("data.customStepConfs").Raw), &stepFilters)
		if err != nil {
			return err
		}

		for _, csc := range stepFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "customStepConf"
			adminKind.Metadata.Name = csc.CustomStepName
			adminKind.Metadata.Annotations = map[string]string{
				"projectNames": strings.Join(csc.ProjectNames, ","),
			}
			var spec pkg.CustomStepConf
			bs, _ := json.Marshal(csc)
			_ = json.Unmarshal(bs, &spec)
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	envFilters := []pkg.EnvK8sDetail{}
	if foundKindEnv {
		param := map[string]interface{}{
			"envNames": o.Param.ItemNames,
			"page":     1,
			"perPage":  1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/envs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(result.Get("data.envK8ss").Raw), &envFilters)
		if err != nil {
			return err
		}

		for _, envK8s := range envFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "envK8s"
			adminKind.Metadata.Name = envK8s.EnvName
			adminKind.Metadata.Annotations = map[string]string{
				"ingressVersion": envK8s.ResourceVersion.IngressVersion,
				"hpaVersion":     envK8s.ResourceVersion.HpaVersion,
				"istioVersion":   envK8s.ResourceVersion.IstioVersion,
			}
			var spec pkg.EnvK8s
			bs, _ := json.Marshal(envK8s)
			_ = json.Unmarshal(bs, &spec)
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	ctFilters := []pkg.ComponentTemplate{}
	if foundKindCt {
		param := map[string]interface{}{
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplates"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		cts := []pkg.ComponentTemplate{}
		err = json.Unmarshal([]byte(result.Get("data.componentTemplates").Raw), &cts)
		if err != nil {
			return err
		}

		for _, ct := range cts {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == ct.ComponentTemplateName {
					found = true
					break
				}
			}
			if found {
				ctFilters = append(ctFilters, ct)
			}
		}

		for _, ct := range ctFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "componentTemplate"
			adminKind.Metadata.Name = ct.ComponentTemplateName
			adminKind.Spec = ct
			adminKinds = append(adminKinds, adminKind)
		}
	}

	dbeFilters := []pkg.DockerBuildEnv{}
	if foundKindDbe {
		param := map[string]interface{}{
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/dockerBuildEnvs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		dbes := []pkg.DockerBuildEnv{}
		err = json.Unmarshal([]byte(result.Get("data.dockerBuildEnvs").Raw), &dbes)
		if err != nil {
			return err
		}

		for _, dbe := range dbes {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == dbe.BuildEnvName {
					found = true
					break
				}
			}
			if found {
				dbeFilters = append(dbeFilters, dbe)
			}
		}

		for _, dbe := range dbeFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "dockerBuildEnv"
			adminKind.Metadata.Name = dbe.BuildEnvName
			adminKind.Spec = dbe
			adminKinds = append(adminKinds, adminKind)
		}
	}

	grcFilters := []pkg.GitRepoConfig{}
	if foundKindGrc {
		param := map[string]interface{}{
			"types":   []string{"gitRepoConfig"},
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		grcs := []pkg.GitRepoConfig{}
		err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &grcs)
		if err != nil {
			return err
		}

		for _, grc := range grcs {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == grc.RepoName {
					found = true
					break
				}
			}
			if found {
				grcFilters = append(grcFilters, grc)
			}
		}

		for _, grc := range grcFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "gitRepoConfig"
			adminKind.Metadata.Name = grc.RepoName
			adminKind.Spec = grc
			adminKinds = append(adminKinds, adminKind)
		}
	}

	ircFilters := []pkg.ImageRepoConfig{}
	if foundKindIrc {
		param := map[string]interface{}{
			"types":   []string{"imageRepoConfig"},
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		ircs := []pkg.ImageRepoConfig{}
		err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &ircs)
		if err != nil {
			return err
		}

		for _, irc := range ircs {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == irc.RepoName {
					found = true
					break
				}
			}
			if found {
				ircFilters = append(ircFilters, irc)
			}
		}

		for _, irc := range ircFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "imageRepoConfig"
			adminKind.Metadata.Name = irc.RepoName
			adminKind.Spec = irc
			adminKinds = append(adminKinds, adminKind)
		}
	}

	arcFilters := []pkg.ArtifactRepoConfig{}
	if foundKindArc {
		param := map[string]interface{}{
			"types":   []string{"artifactRepoConfig"},
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		arcs := []pkg.ArtifactRepoConfig{}
		err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &arcs)
		if err != nil {
			return err
		}

		for _, arc := range arcs {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == arc.RepoName {
					found = true
					break
				}
			}
			if found {
				arcFilters = append(arcFilters, arc)
			}
		}

		for _, arc := range arcFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "artifactRepoConfig"
			adminKind.Metadata.Name = arc.RepoName
			adminKind.Spec = arc
			adminKinds = append(adminKinds, adminKind)
		}
	}

	scrcFilters := []pkg.ScanCodeRepoConfig{}
	if foundKindScrc {
		param := map[string]interface{}{
			"types":   []string{"scanCodeRepoConfig"},
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/repoConfigs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		scrcs := []pkg.ScanCodeRepoConfig{}
		err = json.Unmarshal([]byte(result.Get("data.repoConfigs").Raw), &scrcs)
		if err != nil {
			return err
		}

		for _, scrc := range scrcs {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == scrc.RepoName {
					found = true
					break
				}
			}
			if found {
				scrcFilters = append(scrcFilters, scrc)
			}
		}

		for _, scrc := range scrcFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "scanCodeRepoConfig"
			adminKind.Metadata.Name = scrc.RepoName
			adminKind.Spec = scrc
			adminKinds = append(adminKinds, adminKind)
		}
	}

	adminKindList.Items = adminKinds

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(adminKindList)
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
		if len(userFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range userFilters {
				ups := []string{}
				for _, up := range item.UserProjects {
					ups = append(ups, fmt.Sprintf("%s:%s", up.ProjectName, up.AccessLevel))
				}
				dataRow := []string{fmt.Sprintf("user/%s", item.Username), item.TenantCode, strings.Join(item.TenantAdmins, ","), item.Name, item.Mail, fmt.Sprintf("%v", item.IsAdmin), fmt.Sprintf("%v", item.IsActive), strings.Join(ups, "\n")}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Username", "TenantCode", "TenantAdmins", "Name", "Mail", "Admin", "Active", "Projects"}
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

		if len(stepFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range stepFilters {
				dataRow := []string{fmt.Sprintf("customStepConf/%s", item.CustomStepName), item.TenantCode, item.CustomStepActionDesc, fmt.Sprintf("%v", item.IsEnvDiff), strings.Join(item.ProjectNames, ","), item.ParamInputYamlDef}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Desc", "EnvDiff", "Projects", "Input"}
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

		if len(envFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range envFilters {
				arches := strings.Join(item.Arches, ",")
				dataRow := []string{fmt.Sprintf("envK8s/%s", item.EnvName), item.TenantCode, item.EnvDesc, fmt.Sprintf("https://%s:%d", item.Host, item.Port), arches, item.ResourceVersion.IngressVersion, item.ResourceVersion.HpaVersion, item.ResourceVersion.IstioVersion}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Desc", "Host", "Arches", "Ingress", "Hpa", "Istio"}
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

		if len(ctFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range ctFilters {
				dataRow := []string{fmt.Sprintf("componentTemplate/%s", item.ComponentTemplateName), item.TenantCode, item.ComponentTemplateDesc, item.DeploySpecStatic.DeployImage, fmt.Sprintf("%d", item.DeploySpecStatic.DeployReplicas)}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Desc", "Image", "Replicas"}
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

		if len(dbeFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range dbeFilters {
				dataRow := []string{fmt.Sprintf("dockerBuildEnv/%s", item.BuildEnvName), item.TenantCode, item.Image, strings.Join(item.BuildArches, ",")}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Image", "Arches"}
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

		if len(grcFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range grcFilters {
				dataRow := []string{fmt.Sprintf("gitRepoConfig/%s", item.RepoName), item.TenantCode, item.Kind, item.ViewUrl}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Kind", "ViewUrl"}
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

		if len(ircFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range ircFilters {
				dataRow := []string{fmt.Sprintf("imageRepoConfig/%s", item.RepoName), item.TenantCode, item.Kind, item.Hostname}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Kind", "Hostname"}
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

		if len(arcFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range arcFilters {
				dataRow := []string{fmt.Sprintf("artifactRepoConfig/%s", item.RepoName), item.TenantCode, item.Kind, item.ViewUrl}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Kind", "ViewUrl"}
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

		if len(scrcFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range scrcFilters {
				dataRow := []string{fmt.Sprintf("scanCodeRepoConfig/%s", item.RepoName), item.TenantCode, item.Kind, item.ViewUrl}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "TenantCode", "Kind", "ViewUrl"}
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
