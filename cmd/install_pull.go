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
	FileName       string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Build          bool   `yaml:"build" json:"build" bson:"build" validate:""`
	ForceDownload  bool   `yaml:"forceDownload" json:"forceDownload" bson:"forceDownload" validate:""`
	Yes            bool   `yaml:"yes" json:"yes" bson:"yes" validate:""`
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

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_install_pull_short")
	msgLong := OptCommon.TransLang("cmd_install_pull_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_install_pull_example", baseName))

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

	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", OptCommon.TransLang("param_install_pull_file"))
	cmd.Flags().BoolVar(&o.Build, "build", false, OptCommon.TransLang("param_install_pull_build"))
	cmd.Flags().BoolVar(&o.ForceDownload, "force-download", false, OptCommon.TransLang("param_install_pull_force_download"))
	cmd.Flags().BoolVar(&o.Yes, "yes", false, OptCommon.TransLang("param_install_pull_yes"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallPull) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
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

	var cmdImagePull, cmdImageBuild, cmdImagePullArm64, cmdImageTag string
	switch installConfig.Kubernetes.Runtime {
	case "docker":
		cmdImagePull = "docker pull"
		cmdImageBuild = "docker build"
		cmdImagePullArm64 = "--platform=arm64"
		cmdImageTag = "docker tag"
	case "containerd":
		cmdImagePull = "nerdctl -n k8s.io pull -q"
		cmdImageBuild = "nerdctl -n k8s.io build"
		cmdImagePullArm64 = "--platform=arm64"
		cmdImageTag = "nerdctl -n k8s.io tag"
	case "crio":
		cmdImagePull = "podman pull"
		cmdImageBuild = "podman build"
		cmdImagePullArm64 = "--arch=arm64"
		cmdImageTag = "podman tag"
	}

	var dockerImages pkg.InstallDockerImages
	dockerImages, err = pkg.GetDockerImages(installConfig)
	if err != nil {
		return err
	}

	dockerFileDir := fmt.Sprintf("dory-docker-files")
	dockerFileTplDir := "docker-files"
	if installConfig.Dory.ImageRepo.Type == "harbor" {
		if o.Build {
			_ = os.RemoveAll(dockerFileDir)
			_ = os.MkdirAll(dockerFileDir, 0700)

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
		}
	}

	if installConfig.Dory.ArtifactRepo.Type == "nexus" {
		log.Info("nexus initial data need to download")
		fmt.Println(fmt.Sprintf("curl -O -L https://doryengine.com/attachments/%s", pkg.NexusInitData))
	}

	if installConfig.Dory.ImageRepo.Type == "harbor" {
		log.Info("container images need to pull")
		for _, idi := range dockerImages.InstallDockerImages {
			if o.Build {
				fmt.Println(fmt.Sprintf("%s %s", cmdImagePull, idi.Source))
				if idi.Arm64 != "" {
					fmt.Println(fmt.Sprintf("%s %s %s", cmdImagePull, cmdImagePullArm64, idi.Arm64))
				}
			} else {
				if idi.Built != "" {
					fmt.Println(fmt.Sprintf("%s %s", cmdImagePull, idi.Built))
				} else {
					fmt.Println(fmt.Sprintf("%s %s", cmdImagePull, idi.Source))
				}
				if idi.Arm64 != "" {
					if idi.Built != "" {
						fmt.Println(fmt.Sprintf("%s %s %s-arm64v8", cmdImagePull, cmdImagePullArm64, idi.Built))
					} else {
						fmt.Println(fmt.Sprintf("%s %s %s", cmdImagePull, cmdImagePullArm64, idi.Arm64))
					}
				}
			}
		}

		if o.Build {
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
		} else {
			for _, idi := range dockerImages.InstallDockerImages {
				if idi.Built != "" {
					fmt.Println(fmt.Sprintf("%s %s %s", cmdImageTag, idi.Built, idi.Target))
					if idi.Arm64 != "" {
						fmt.Println(fmt.Sprintf("%s %s-arm64v8 %s-arm64v8", cmdImageTag, idi.Built, idi.Target))
					}
				}
			}
		}
	}

	if installConfig.Dory.ImageRepo.Type == "harbor" || installConfig.Dory.ArtifactRepo.Type == "nexus" {
		arr := []string{}
		if installConfig.Dory.ImageRepo.Type == "harbor" {
			arr = append(arr, "pull container images")
		}
		if installConfig.Dory.ArtifactRepo.Type == "nexus" {
			arr = append(arr, "download nexus initial data")
		}
		if !o.Yes {
			log.Warning(fmt.Sprintf("make sure current host can connect internet, are you sure %s now? [YES/NO]", strings.Join(arr, " and ")))
			reader := bufio.NewReader(os.Stdin)
			userInput, _ := reader.ReadString('\n')
			userInput = strings.Trim(userInput, "\n")
			if userInput != "YES" {
				err = fmt.Errorf("user cancelled")
				return err
			}
		}

		var isDownloadNexus bool
		if installConfig.Dory.ArtifactRepo.Type == "nexus" {
			fi, err := os.Stat(pkg.NexusInitData)
			if err == nil {
				if !fi.IsDir() {
					if o.ForceDownload {
						isDownloadNexus = true
					}
				} else {
					_ = os.RemoveAll(pkg.NexusInitData)
					isDownloadNexus = true
				}
			} else {
				isDownloadNexus = true
			}
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

		if installConfig.Dory.ImageRepo.Type == "harbor" {
			log.Info("pull container images begin")
			for i, idi := range dockerImages.InstallDockerImages {
				imageName := idi.Source
				imageNameArm64 := idi.Arm64
				if !o.Build && idi.Built != "" {
					imageName = idi.Built
					imageNameArm64 = fmt.Sprintf("%s-arm64v8", idi.Built)
				}
				_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s", cmdImagePull, imageName), ".")
				if err != nil {
					err = fmt.Errorf("pull container image %s error: %s", imageName, err.Error())
					return err
				}
				if idi.Arm64 != "" {
					_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s %s", cmdImagePull, cmdImagePullArm64, imageNameArm64), ".")
					if err != nil {
						err = fmt.Errorf("pull container image %s error: %s", imageNameArm64, err.Error())
						return err
					}
				}
				if o.Build {
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
				} else {
					if idi.Built != "" {
						_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s %s", cmdImageTag, idi.Built, idi.Target), ".")
						if err != nil {
							err = fmt.Errorf("tag container image %s error: %s", idi.Source, err.Error())
							return err
						}
						if idi.Arm64 != "" {
							_, _, err = pkg.CommandExec(fmt.Sprintf("%s %s-arm64v8 %s-arm64v8", cmdImageTag, idi.Built, idi.Target), ".")
							if err != nil {
								err = fmt.Errorf("tag container image %s error: %s", idi.Source, err.Error())
								return err
							}
						}
					}
				}
				log.Success(fmt.Sprintf("# progress: %d/%d %s", i+1, len(dockerImages.InstallDockerImages), idi.Target))
			}
			log.Success(fmt.Sprintf("pull and build container images success"))
		}
	} else {
		log.Success(fmt.Sprintf("make sure current host can connect internet, no container images need to pull now"))
	}

	defer color.Unset()
	return err
}
