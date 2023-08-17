package cmd

import (
	"bufio"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

type OptionsInstallPull struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Mode           string `yaml:"mode" json:"mode" bson:"mode" validate:""`
	Runtime        string `yaml:"runtime" json:"runtime" bson:"runtime" validate:""`
	FileName       string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	DownloadAgain  bool   `yaml:"downloadAgain" json:"downloadAgain" bson:"downloadAgain" validate:""`
}

func NewOptionsInstallPull() *OptionsInstallPull {
	var o OptionsInstallPull
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallPull() *cobra.Command {
	o := NewOptionsInstallPull()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("pull")
	msgShort := fmt.Sprintf("pull and build all container images")
	msgLong := fmt.Sprintf(`pull and build all container images required for installation`)
	msgExample := fmt.Sprintf(`  # pull and build all container images required for installing dory in kubernetes over containerd runtime
  %s install pull --mode kubernetes --runtime containerd -f install-config.yaml

  # pull and build all container images required for installing dory in docker
  %s install pull --mode docker -f install-config.yaml`, baseName, baseName)

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

	cmd.Flags().StringVar(&o.Mode, "mode", "", "install dory in docker or kubernetes, options: docker, kubernetes")
	cmd.Flags().StringVar(&o.Runtime, "runtime", "", "dory managed kubernetes cluster's container runtime, options: docker, containerd, crio")
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "install settings YAML file")
	cmd.Flags().BoolVarP(&o.DownloadAgain, "download-again", "d", false, "download nexus init data and trivy db again even file exists")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallPull) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"kubernetes", "docker"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("runtime", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"docker", "containerd", "crio"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("mode")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallPull) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("--mode must be docker or kubernetes")
		return err
	}

	if o.Mode == "kubernetes" {
		if o.Runtime != "docker" && o.Runtime != "containerd" && o.Runtime != "crio" {
			err = fmt.Errorf("--runtime must be docker, containerd or crio")
			return err
		}
	}

	if o.FileName == "" {
		err = fmt.Errorf("--file required")
		return err
	}

	return err
}

// Run executes the appropriate steps to pull a model's documentation
func (o *OptionsInstallPull) Run(args []string) error {
	var err error

	bs := []byte{}

	bs, err = os.ReadFile(o.FileName)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	var installConfig pkg.InstallConfig
	err = yaml.Unmarshal(bs, &installConfig)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	err = installConfig.VerifyInstallConfig()
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	var runtime string
	if o.Mode == "docker" {
		runtime = "docker"
	} else if o.Mode == "kubernetes" {
		runtime = o.Runtime
	}
	if runtime == "" {
		err = fmt.Errorf("--runtime must be docker, containerd or crio")
		return err
	}

	var cmdImagePull, cmdImageBuild, cmdImagePullArm64 string
	switch runtime {
	case "docker":
		cmdImagePull = "docker pull"
		cmdImageBuild = "docker build"
		cmdImagePullArm64 = "--platform=arm64"
	case "containerd":
		cmdImagePull = "nerdctl -n k8s.io pull -q"
		cmdImageBuild = "nerdctl -n k8s.io build"
		cmdImagePullArm64 = "--platform=arm64"
	case "crio":
		cmdImagePull = "podman pull"
		cmdImageBuild = "podman build"
		cmdImagePullArm64 = "--arch=arm64"
	}

	var dockerImages pkg.InstallDockerImages
	dockerImages, err = pkg.GetDockerImages(installConfig)
	if err != nil {
		return err
	}

	dockerFileDir := fmt.Sprintf("dory-docker-files")
	_ = os.RemoveAll(dockerFileDir)
	_ = os.MkdirAll(dockerFileDir, 0700)
	dockerFileTplDir := "docker-files"

	for _, dockerImage := range dockerImages.InstallDockerImages {
		if dockerImage.DockerFile != "" {
			arr := strings.Split(dockerImage.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			dockerFileName := fmt.Sprintf("%s/%s-%s", dockerFileDir, dockerImage.DockerFile, tagName)

			bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerFileTplDir, dockerImage.DockerFile))
			if err != nil {
				err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
				return err
			}
			vals := map[string]interface{}{
				"source":  dockerImage.Source,
				"tagName": tagName,
				"target":  dockerImage.Target,
				"isArm64": false,
			}
			strDockerfile, err := pkg.ParseTplFromVals(vals, string(bs))
			if err != nil {
				err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
				return err
			}
			err = os.WriteFile(fmt.Sprintf("%s", dockerFileName), []byte(strDockerfile), 0600)
			if err != nil {
				err = fmt.Errorf("create values.yaml error: %s", err.Error())
				return err
			}

			if dockerImage.Arm64 != "" {
				arr := strings.Split(dockerImage.Arm64, ":")
				var tagName string
				if len(arr) == 2 {
					tagName = arr[1]
				} else {
					tagName = "latest"
				}
				dockerFileName := fmt.Sprintf("%s/%s-%s-arm64v8", dockerFileDir, dockerImage.DockerFile, tagName)

				bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerFileTplDir, dockerImage.DockerFile))
				if err != nil {
					err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
					return err
				}
				vals := map[string]interface{}{
					"source":  dockerImage.Arm64,
					"tagName": tagName,
					"target":  dockerImage.Target,
					"isArm64": true,
				}
				strDockerfile, err := pkg.ParseTplFromVals(vals, string(bs))
				if err != nil {
					err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
					return err
				}
				err = os.WriteFile(fmt.Sprintf("%s", dockerFileName), []byte(strDockerfile), 0600)
				if err != nil {
					err = fmt.Errorf("create values.yaml error: %s", err.Error())
					return err
				}
			}
		}
	}
	log.Info(fmt.Sprintf("create docker files in %s success", dockerFileDir))
	_, _, err = pkg.CommandExec("ls -alh", dockerFileDir)
	if err != nil {
		err = fmt.Errorf("create docker files %s error: %s", dockerFileDir, err.Error())
		return err
	}
	time.Sleep(time.Second * 1)

	log.Info("nexus initial data need to download")
	fmt.Println(fmt.Sprintf("curl -O -L https://doryengine.com/attachments/%s", pkg.NexusInitData))

	log.Info("trivy vulnerabilities database need to download")
	fmt.Println(fmt.Sprintf("curl -O -L https://doryengine.com/attachments/%s", pkg.TrivyDb))

	log.Info("container images need to pull")
	for _, idi := range dockerImages.InstallDockerImages {
		fmt.Println(fmt.Sprintf("%s %s", cmdImagePull, idi.Source))
		if idi.Arm64 != "" {
			fmt.Println(fmt.Sprintf("%s %s %s", cmdImagePull, cmdImagePullArm64, idi.Arm64))
		}
	}

	log.Info("container images need to build")
	log.Warning(fmt.Sprintf("all docker files in %s folder, if your machine is without internet connection, build container images by manual", dockerFileDir))
	for _, idi := range dockerImages.InstallDockerImages {
		if idi.DockerFile != "" {
			arr := strings.Split(idi.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			fmt.Println(fmt.Sprintf("%s -t %s -f %s/%s-%s %s", cmdImageBuild, idi.Target, dockerFileDir, idi.DockerFile, tagName, dockerFileDir))
			if idi.Arm64 != "" {
				arr := strings.Split(idi.Arm64, ":")
				var tagName string
				if len(arr) == 2 {
					tagName = arr[1]
				} else {
					tagName = "latest"
				}
				fmt.Println(fmt.Sprintf("%s -t %s-arm64v8 -f %s/%s-%s-arm64v8 %s", cmdImageBuild, idi.Target, dockerFileDir, idi.DockerFile, tagName, dockerFileDir))
			}
		}
	}

	log.Warning("Make sure current host can connect internet, are you sure download neuxs initial data, trivy vulnerabilities database and pull container images now? [YES/NO]")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.Trim(userInput, "\n")
	if userInput != "YES" {
		err = fmt.Errorf("user cancelled")
		return err
	}

	var isDownloadNexus bool
	fi, err := os.Stat(pkg.NexusInitData)
	if err == nil {
		if !fi.IsDir() {
			if o.DownloadAgain {
				isDownloadNexus = true
			}
		} else {
			_ = os.RemoveAll(pkg.NexusInitData)
			isDownloadNexus = true
		}
	} else {
		isDownloadNexus = true
	}
	baseName := pkg.GetCmdBaseName()
	if isDownloadNexus {
		log.Info(fmt.Sprintf("start to download nexus initial data"))
		_, _, err = pkg.CommandExec(fmt.Sprintf("curl -O -L https://doryengine.com/attachments/%s", pkg.NexusInitData), ".")
		if err != nil {
			err = fmt.Errorf("download nexus initial data error: %s", err.Error())
			return err
		}
		log.Info(fmt.Sprintf("download nexus initial data %s success", pkg.NexusInitData))
		log.Info(fmt.Sprintf("if run '%s install run [options]' or '%s install script [options]' command, make sure run this command at %s file's directory", baseName, baseName, pkg.NexusInitData))
	}

	var isDownloadTrivy bool
	fi, err = os.Stat(pkg.TrivyDb)
	if err == nil {
		if !fi.IsDir() {
			if o.DownloadAgain {
				isDownloadTrivy = true
			}
		} else {
			_ = os.RemoveAll(pkg.TrivyDb)
			isDownloadTrivy = true
		}
	} else {
		isDownloadTrivy = true
	}
	if isDownloadTrivy {
		log.Info(fmt.Sprintf("start to download trivy vulnerabilities database"))
		_, _, err = pkg.CommandExec(fmt.Sprintf("curl -O -L https://doryengine.com/attachments/%s", pkg.TrivyDb), ".")
		if err != nil {
			err = fmt.Errorf("download trivy vulnerabilities database error: %s", err.Error())
			return err
		}
		log.Info(fmt.Sprintf("download trivy vulnerabilities database %s success", pkg.TrivyDb))
		log.Info(fmt.Sprintf("if run '%s install run [options]' or '%s install script [options]' command, make sure run this command at %s file's directory", baseName, baseName, pkg.TrivyDb))
	}

	log.Info("pull and build container images begin")
	for i, idi := range dockerImages.InstallDockerImages {
		_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s", cmdImagePull, idi.Source), ".")
		if err != nil {
			err = fmt.Errorf("pull container image %s error: %s", idi.Source, err.Error())
			return err
		}
		if idi.Arm64 != "" {
			_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s %s", cmdImagePull, cmdImagePullArm64, idi.Arm64), ".")
			if err != nil {
				err = fmt.Errorf("pull container image %s error: %s", idi.Arm64, err.Error())
				return err
			}
		}
		if idi.DockerFile != "" {
			arr := strings.Split(idi.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			_, _, err = pkg.CommandExec(fmt.Sprintf("%s -t %s -f %s/%s-%s %s", cmdImageBuild, idi.Target, dockerFileDir, idi.DockerFile, tagName, dockerFileDir), ".")
			if err != nil {
				err = fmt.Errorf("build container image %s error: %s", idi.Source, err.Error())
				return err
			}
			if idi.Arm64 != "" {
				arr := strings.Split(idi.Arm64, ":")
				var tagName string
				if len(arr) == 2 {
					tagName = arr[1]
				} else {
					tagName = "latest"
				}
				_, _, err = pkg.CommandExec(fmt.Sprintf("%s -t %s-arm64v8 -f %s/%s-%s-arm64v8 %s", cmdImageBuild, idi.Target, dockerFileDir, idi.DockerFile, tagName, dockerFileDir), ".")
				if err != nil {
					err = fmt.Errorf("build container image %s error: %s", idi.Arm64, err.Error())
					return err
				}
			}
		}
		log.Success(fmt.Sprintf("# progress: %d/%d %s", i+1, len(dockerImages.InstallDockerImages), idi.Target))
	}
	log.Success(fmt.Sprintf("pull and build container images success"))

	defer color.Unset()
	return err
}
