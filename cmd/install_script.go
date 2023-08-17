package cmd

import (
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

type OptionsInstallScript struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileName       string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	OutputDir      string `yaml:"outputDir" json:"outputDir" bson:"outputDir" validate:""`
	Param          struct {
		Stdin []byte `yaml:"stdin" json:"stdin" bson:"stdin" validate:""`
	}
}

func NewOptionsInstallScript() *OptionsInstallScript {
	var o OptionsInstallScript
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallScript() *cobra.Command {
	o := NewOptionsInstallScript()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("script")
	msgShort := fmt.Sprintf("create dory-engine install scripts and config files")
	msgLong := fmt.Sprintf(`create dory-engine install scripts and config files, run the scripts by manual, for experts`)
	msgExample := fmt.Sprintf(`  # create dory-engine install scripts and config files with docker-compose or kubernetes
  %s install script -o readme-install -f install-config.yaml
  or
  cat install-config.yaml | %s install script -o readme-install -f -
`, baseName, baseName)

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
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "install settings YAML file")
	cmd.Flags().StringVarP(&o.OutputDir, "output", "o", "", "output README, script and config files directory")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallScript) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("file")
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("output")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallScript) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.FileName == "" {
		err = fmt.Errorf("--file required")
		return err
	}
	baseName := pkg.GetCmdBaseName()
	if o.FileName == "-" {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		o.Param.Stdin = bs
		if len(o.Param.Stdin) == 0 {
			err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s install script -o readme-install -f -", baseName)
			return err
		}
	}
	if o.OutputDir == "" {
		err = fmt.Errorf("--output required")
		return err
	}
	if strings.Contains(o.OutputDir, "/") {
		err = fmt.Errorf("--output can not contain /")
		return err
	}
	return err
}

// Run executes the appropriate steps to run a model's documentation
func (o *OptionsInstallScript) Run(args []string) error {
	var err error

	bs := []byte{}

	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	if o.FileName == "-" {
		bs = o.Param.Stdin
	} else {
		bs, err = os.ReadFile(o.FileName)
		if err != nil {
			err = fmt.Errorf("install script error: %s", err.Error())
			return err
		}
	}

	_, err = os.Stat(pkg.NexusInitData)
	if err != nil {
		err = fmt.Errorf("%s not exists in current directory", pkg.NexusInitData)
		return err
	}

	_, err = os.Stat(pkg.TrivyDb)
	if err != nil {
		err = fmt.Errorf("%s not exists in current directory", pkg.TrivyDb)
		return err
	}

	var installConfig pkg.InstallConfig
	err = yaml.Unmarshal(bs, &installConfig)
	if err != nil {
		err = fmt.Errorf("install script error: %s", err.Error())
		return err
	}

	err = installConfig.VerifyInstallConfig()
	if err != nil {
		err = fmt.Errorf("install script error: %s", err.Error())
		return err
	}

	if installConfig.InstallMode == "docker" {
		err = o.ScriptWithDocker(installConfig)
		if err != nil {
			return err
		}
	} else if installConfig.InstallMode == "kubernetes" {
		err = o.ScriptWithKubernetes(installConfig)
		if err != nil {
			return err
		}
	} else {
		err = fmt.Errorf("install script error: installMode not correct, must be docker or kubernetes")
		return err
	}
	return err
}

func (o *OptionsInstallScript) DoryCreateConfig(installConfig pkg.InstallConfig, rootDir string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	doryengineDir := fmt.Sprintf("%s/%s/dory-engine", rootDir, installConfig.Dory.Namespace)
	doryengineConfigDir := fmt.Sprintf("%s/config", doryengineDir)
	doryengineScriptDir := "dory/dory-engine"
	doryengineConfigName := "config.yaml"
	doryengineEnvK8sName := "env-k8s.yaml"
	doryengineLicenseName := "license.yaml"
	_ = os.RemoveAll(doryengineConfigDir)
	_ = os.MkdirAll(doryengineConfigDir, 0700)

	// create config.yaml
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s-%s", pkg.DirInstallScripts, doryengineScriptDir, o.Language, doryengineConfigName))
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	strDoryengineConfig, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryengineConfigDir, doryengineConfigName), []byte(strDoryengineConfig), 0600)
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	// create env-k8s-test.yaml
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s-%s", pkg.DirInstallScripts, doryengineScriptDir, o.Language, doryengineEnvK8sName))
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	strDoryengineEnvK8s, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryengineConfigDir, fmt.Sprintf("env-k8s-%s.yaml", installConfig.Kubernetes.EnvName)), []byte(strDoryengineEnvK8s), 0600)
	if err != nil {
		err = fmt.Errorf("create dory-engine config files error: %s", err.Error())
		return err
	}
	// create license.yaml
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, doryengineScriptDir, doryengineLicenseName))
	if err != nil {
		err = fmt.Errorf("create dory-engine license files error: %s", err.Error())
		return err
	}
	strDoryengineLicense, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory-engine license files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryengineConfigDir, doryengineLicenseName), []byte(strDoryengineLicense), 0600)
	if err != nil {
		err = fmt.Errorf("create dory-engine license files error: %s", err.Error())
		return err
	}

	return err
}

func (o *OptionsInstallScript) DoryCreateNginxGitlabConfig(installConfig pkg.InstallConfig, rootDir string) error {
	var err error
	var bs []byte

	if installConfig.Dory.GitRepo.Type == "gitlab" {
		vals, err := installConfig.UnmarshalMapValues()
		if err != nil {
			return err
		}

		nginxGitlabDir := fmt.Sprintf("%s/%s/nginx-%s", rootDir, installConfig.Dory.Namespace, installConfig.Dory.GitRepo.Type)
		nginxGitlabScriptDir := "dory/nginx-gitlab"
		nginxGitlabConfigName := "nginx-gitlab.conf"
		_ = os.RemoveAll(nginxGitlabDir)
		_ = os.MkdirAll(nginxGitlabDir, 0700)

		// create nginx-gitlab.conf
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, nginxGitlabScriptDir, nginxGitlabConfigName))
		if err != nil {
			err = fmt.Errorf("create nginx-gitlab.conf file error: %s", err.Error())
			return err
		}
		strNginxGitlabConfig, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create nginx-gitlab.conf file error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", nginxGitlabDir, nginxGitlabConfigName), []byte(strNginxGitlabConfig), 0600)
		if err != nil {
			err = fmt.Errorf("create nginx-gitlab.conf file error: %s", err.Error())
			return err
		}
	}

	return err
}

func (o *OptionsInstallScript) DoryCreateDockerCertsConfig(installConfig pkg.InstallConfig, rootDir string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	dockerDir := fmt.Sprintf("%s/%s/%s", rootDir, installConfig.Dory.Namespace, installConfig.Dory.Docker.DockerName)
	_ = os.RemoveAll(dockerDir)
	_ = os.MkdirAll(dockerDir, 0700)
	dockerScriptDir := "dory/docker"
	dockerScriptName := "docker_certs.sh"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerScriptName))
	if err != nil {
		err = fmt.Errorf("create docker certificates error: %s", err.Error())
		return err
	}
	strDockerCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create docker certificates error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerScriptName), []byte(strDockerCertScript), 0600)
	if err != nil {
		err = fmt.Errorf("create docker certificates error: %s", err.Error())
		return err
	}

	dockerDaemonJsonName := "daemon.json"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerDaemonJsonName))
	if err != nil {
		err = fmt.Errorf("create docker config error: %s", err.Error())
		return err
	}
	strDockerDaemonJson, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create docker config error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerDaemonJsonName), []byte(strDockerDaemonJson), 0600)
	if err != nil {
		err = fmt.Errorf("create docker config error: %s", err.Error())
		return err
	}

	dockerConfigJsonName := "config.json"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerConfigJsonName))
	if err != nil {
		err = fmt.Errorf("create docker config files error: %s", err.Error())
		return err
	}
	strDockerConfigJson, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create docker config files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerConfigJsonName), []byte(strDockerConfigJson), 0600)
	if err != nil {
		err = fmt.Errorf("create docker config files error: %s", err.Error())
		return err
	}

	return err
}

func (o *OptionsInstallScript) DoryCreateOpenldapCertsConfig(installConfig pkg.InstallConfig, rootDir string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	doryDir := fmt.Sprintf("%s/%s", rootDir, installConfig.Dory.Namespace)
	openldapCertsDir := fmt.Sprintf("%s/%s/certs", doryDir, installConfig.Dory.Openldap.ServiceName)
	_ = os.RemoveAll(openldapCertsDir)
	_ = os.MkdirAll(openldapCertsDir, 0700)
	openldapScriptDir := "dory/openldap"
	openldapScriptName := "openldap_certs.sh"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, openldapScriptDir, openldapScriptName))
	if err != nil {
		err = fmt.Errorf("create openldap certificates error: %s", err.Error())
		return err
	}
	strDockerCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create openldap certificates error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", openldapCertsDir, openldapScriptName), []byte(strDockerCertScript), 0600)
	if err != nil {
		err = fmt.Errorf("create openldap certificates error: %s", err.Error())
		return err
	}

	return err
}

func (o *OptionsInstallScript) DoryCreateKubernetesDataPod(installConfig pkg.InstallConfig, rootDir string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	kubernetesDir := "kubernetes"
	projectDataPodName := "project-data-pod.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, kubernetesDir, projectDataPodName))
	if err != nil {
		err = fmt.Errorf("create project-data-pod in kubernetes error: %s", err.Error())
		return err
	}
	strProjectDataAlpine, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create project-data-pod in kubernetes error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", rootDir, projectDataPodName), []byte(strProjectDataAlpine), 0600)
	if err != nil {
		err = fmt.Errorf("create project-data-pod in kubernetes error: %s", err.Error())
		return err
	}

	return err
}

func (o *OptionsInstallScript) DoryCreateReadme(installConfig pkg.InstallConfig, readmeInstallDir, readmeName string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	var cmdContainerLogin, cmdContainerTag, cmdContainerPush string
	switch installConfig.Kubernetes.Runtime {
	case "docker":
		cmdContainerLogin = "docker login"
		cmdContainerTag = "docker tag"
		cmdContainerPush = "docker push"
	case "containerd":
		cmdContainerLogin = "nerdctl -n k8s.io login"
		cmdContainerTag = "nerdctl -n k8s.io tag"
		cmdContainerPush = "nerdctl -n k8s.io push"
	case "crio":
		cmdContainerLogin = "podman login"
		cmdContainerTag = "podman tag"
		cmdContainerPush = "podman push"
	}
	vals["cmdLogin"] = cmdContainerLogin
	vals["cmdTag"] = cmdContainerTag
	vals["cmdPush"] = cmdContainerPush

	// get pull docker images
	dockerImages, err := pkg.GetDockerImages(installConfig)
	if err != nil {
		return err
	}
	bs, _ = pkg.YamlIndent(dockerImages)
	m := map[string]interface{}{}
	_ = yaml.Unmarshal(bs, &m)
	for k, v := range m {
		vals[k] = v
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-%s", pkg.DirInstallScripts, o.Language, readmeName))
	if err != nil {
		err = fmt.Errorf("create %s error: %s", readmeName, err.Error())
		return err
	}
	strDoryInstallSettings, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create %s error: %s", readmeName, err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", readmeInstallDir, readmeName), []byte(strDoryInstallSettings), 0600)
	if err != nil {
		err = fmt.Errorf("create %s error: %s", readmeName, err.Error())
		return err
	}

	return err
}

func (o *OptionsInstallScript) ScriptWithDocker(installConfig pkg.InstallConfig) error {
	var err error
	bs := []byte{}

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}
	vals["versionDoryEngine"] = pkg.VersionDoryEngine
	vals["versionDoryFrontend"] = pkg.VersionDoryFrontend

	outputDir := o.OutputDir
	_ = os.MkdirAll(outputDir, 0700)

	readmeDockerResetName := "README-2-docker-reset.md"
	defer o.DoryCreateReadme(installConfig, outputDir, readmeDockerResetName)

	if installConfig.Dory.ImageRepo.Internal.Hostname != "" {
		// create harbor certificates
		harborDir := fmt.Sprintf("%s/%s", outputDir, installConfig.Dory.ImageRepo.Internal.Namespace)
		_ = os.RemoveAll(harborDir)
		_ = os.MkdirAll(harborDir, 0700)
		harborScriptDir := "harbor"
		harborScriptName := "harbor_certs.sh"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborScriptDir, harborScriptName))
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}
		strHarborCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborScriptName), []byte(strHarborCertScript), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}

		// extract harbor install files
		err = pkg.ExtractEmbedFile(pkg.FsInstallScripts, fmt.Sprintf("%s/harbor/harbor", pkg.DirInstallScripts), harborDir)
		if err != nil {
			err = fmt.Errorf("extract harbor install files error: %s", err.Error())
			return err
		}

		harborInstallerDir := "harbor/harbor"
		harborYamlName := "harbor.yml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborYamlName))
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}
		strHarborYaml, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborYamlName), []byte(strHarborYaml), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}

		harborPrepareName := "prepare"
		_ = os.Rename(fmt.Sprintf("%s/harbor", installConfig.RootDir), harborDir)
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborPrepareName))
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}
		strHarborPrepare, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborPrepareName), []byte(strHarborPrepare), 0700)
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}

		_ = os.Chmod(fmt.Sprintf("%s/common.sh", harborDir), 0700)
		_ = os.Chmod(fmt.Sprintf("%s/install.sh", harborDir), 0700)
		_ = os.Chmod(fmt.Sprintf("%s/prepare", harborDir), 0700)
	}

	////////////////////////////////////////////////////

	// create dory docker-compose.yaml
	doryDir := fmt.Sprintf("%s/%s", outputDir, installConfig.Dory.Namespace)
	_ = os.RemoveAll(doryDir)
	_ = os.MkdirAll(doryDir, 0700)
	dockerComposeDir := "dory"
	dockerComposeName := "docker-compose.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerComposeDir, dockerComposeName))
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}
	strCompose, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryDir, dockerComposeName), []byte(strCompose), 0600)
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}

	// create dory-engine config files
	err = o.DoryCreateConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create nginx-gitlab config files
	err = o.DoryCreateNginxGitlabConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create docker certificates and config
	err = o.DoryCreateDockerCertsConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create openldap certificates
	err = o.DoryCreateOpenldapCertsConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	// create project-data-pod in kubernetes
	err = o.DoryCreateKubernetesDataPod(installConfig, outputDir)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	readmeDockerConfigName := "README-1-docker-config.md"
	err = o.DoryCreateReadme(installConfig, outputDir, readmeDockerConfigName)
	if err != nil {
		return err
	}

	readmeDockerInstallName := "README-0-docker-install.md"
	err = o.DoryCreateReadme(installConfig, outputDir, readmeDockerInstallName)
	if err != nil {
		return err
	}

	log.Warning(fmt.Sprintf("all scripts and config files located at: %s", outputDir))
	log.Warning(fmt.Sprintf("change your work directory to %s", outputDir))
	log.Warning(fmt.Sprintf("1. please follow %s to install dory by manual", readmeDockerInstallName))
	log.Warning(fmt.Sprintf("2. please follow %s to config dory by manual after install", readmeDockerConfigName))
	log.Warning(fmt.Sprintf("3. if install fail, please follow %s to stop all dory services and install again", readmeDockerResetName))

	return err
}

func (o *OptionsInstallScript) ScriptWithKubernetes(installConfig pkg.InstallConfig) error {
	var err error
	bs := []byte{}

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}
	vals["versionDoryEngine"] = pkg.VersionDoryEngine
	vals["versionDoryFrontend"] = pkg.VersionDoryFrontend

	outputDir := o.OutputDir
	_ = os.MkdirAll(outputDir, 0700)

	readmeKubernetesResetName := "README-2-kubernetes-reset.md"
	defer o.DoryCreateReadme(installConfig, outputDir, readmeKubernetesResetName)

	if installConfig.Dory.ImageRepo.Internal.Hostname != "" {
		harborInstallerDir := "kubernetes/harbor"
		harborInstallYamlDir := fmt.Sprintf("%s/harbor", outputDir)
		_ = os.RemoveAll(harborInstallYamlDir)
		_ = os.MkdirAll(harborInstallYamlDir, 0700)

		// extract harbor helm files
		err = pkg.ExtractEmbedFile(pkg.FsInstallScripts, fmt.Sprintf("%s/%s", pkg.DirInstallScripts, harborInstallerDir), harborInstallYamlDir)
		if err != nil {
			err = fmt.Errorf("extract harbor helm files error: %s", err.Error())
			return err
		}

		harborValuesYamlName := "values.yaml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborValuesYamlName))
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}
		strHarborValuesYaml, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborInstallYamlDir, harborValuesYamlName), []byte(strHarborValuesYaml), 0600)
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}

		// create harbor namespace and pv pvc
		harborInstallDir := fmt.Sprintf("%s/%s", outputDir, installConfig.Dory.ImageRepo.Internal.Namespace)
		_ = os.MkdirAll(harborInstallDir, 0700)
		vals["currentNamespace"] = installConfig.Dory.ImageRepo.Internal.Namespace
		step01NamespacePvName := "step01-namespace-pv.yaml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step01NamespacePvName))
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}
		strStep01NamespacePv, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborInstallDir, step01NamespacePvName), []byte(strStep01NamespacePv), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}

		harborUpdateCertsName := "harbor_update_docker_certs.sh"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, harborUpdateCertsName))
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
		strHarborUpdateCerts, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborInstallDir, harborUpdateCertsName), []byte(strHarborUpdateCerts), 0600)
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
	}

	//////////////////////////////////////////////////

	// create dory namespace and pv pvc
	doryInstallDir := fmt.Sprintf("%s/%s", outputDir, installConfig.Dory.Namespace)
	_ = os.MkdirAll(doryInstallDir, 0700)
	vals["currentNamespace"] = installConfig.Dory.Namespace
	step01NamespacePvName := "step01-namespace-pv.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step01NamespacePvName))
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	strStep01NamespacePv, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryInstallDir, step01NamespacePvName), []byte(strStep01NamespacePv), 0600)
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}

	// create dory install yaml
	doryInstallYamlName := "dory-install.yaml"
	step02StatefulsetName := "step02-statefulset.yaml"
	step03ServiceName := "step03-service.yaml"
	step04NetworkPolicyName := "step04-networkpolicy.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, doryInstallYamlName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strDoryInstallYaml, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	installVals := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(strDoryInstallYaml), &installVals)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	for k, v := range vals {
		installVals[k] = v
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step02StatefulsetName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strStep02Statefulset, err := pkg.ParseTplFromVals(installVals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryInstallDir, step02StatefulsetName), []byte(strStep02Statefulset), 0600)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step03ServiceName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strStep03Service, err := pkg.ParseTplFromVals(installVals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryInstallDir, step03ServiceName), []byte(strStep03Service), 0600)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step04NetworkPolicyName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strStep04NetworkPolicy, err := pkg.ParseTplFromVals(installVals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, step04NetworkPolicyName), []byte(strStep04NetworkPolicy), 0600)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}

	// create dory-engine config files
	err = o.DoryCreateConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create nginx-gitlab config files
	err = o.DoryCreateNginxGitlabConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create docker certificates and config
	err = o.DoryCreateDockerCertsConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	// create openldap certificates
	err = o.DoryCreateOpenldapCertsConfig(installConfig, outputDir)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	// create project-data-pod in kubernetes
	err = o.DoryCreateKubernetesDataPod(installConfig, outputDir)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	readmeKubernetesConfigName := "README-1-kubernetes-config.md"
	err = o.DoryCreateReadme(installConfig, outputDir, readmeKubernetesConfigName)
	if err != nil {
		return err
	}

	readmeKubernetesInstallName := "README-0-kubernetes-install.md"
	err = o.DoryCreateReadme(installConfig, outputDir, readmeKubernetesInstallName)
	if err != nil {
		return err
	}

	log.Warning(fmt.Sprintf("all scripts and config files located at: %s", outputDir))
	log.Warning(fmt.Sprintf("change your work directory to %s", outputDir))
	log.Warning(fmt.Sprintf("1. please follow %s to install dory by manual", readmeKubernetesInstallName))
	log.Warning(fmt.Sprintf("2. please follow %s to config dory by manual after install", readmeKubernetesConfigName))
	log.Warning(fmt.Sprintf("3. if install fail, please follow %s to stop all dory services and install again", readmeKubernetesResetName))
	return err
}
