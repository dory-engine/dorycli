package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

type OptionsInstallCheck struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Mode           string `yaml:"mode" json:"mode" bson:"mode" validate:""`
	Runtime        string `yaml:"runtime" json:"runtime" bson:"runtime" validate:""`
}

func NewOptionsInstallCheck() *OptionsInstallCheck {
	var o OptionsInstallCheck
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallCheck() *cobra.Command {
	o := NewOptionsInstallCheck()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("check")
	msgShort := fmt.Sprintf("check install prerequisite")
	msgLong := fmt.Sprintf(`check installing dory in docker or kubernetes prerequisite`)
	msgExample := fmt.Sprintf(`  # check installing dory in kubernetes prerequisite, with dory managed kubernetes cluster over containerd runtime
  %s install check --mode kubernetes --runtime containerd

  # check installing dory in docker prerequisite, with dory managed kubernetes cluster over docker runtime
  %s install check --mode docker --runtime docker`, baseName, baseName)

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

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallCheck) Complete(cmd *cobra.Command) error {
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

	err = cmd.MarkFlagRequired("runtime")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallCheck) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("--mode must be docker or kubernetes")
		return err
	}

	if o.Runtime != "docker" && o.Runtime != "containerd" && o.Runtime != "crio" {
		err = fmt.Errorf("--runtime must be docker, containerd or crio")
		return err
	}

	return err
}

// Run executes the appropriate steps to check a model's documentation
func (o *OptionsInstallCheck) Run(args []string) error {
	var err error

	defer color.Unset()

	log.Info("check kubernetes installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl get pods -A -o wide"), ".")
	if err != nil {
		err = fmt.Errorf("check kubernetes installed error: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	log.Success("check kubernetes installed success")

	log.Info("check openssl installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("openssl version"), ".")
	if err != nil {
		err = fmt.Errorf("check openssl installed error: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	log.Success("check openssl installed success")

	if o.Mode == "docker" {
		log.Info("check dockerd installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("dockerd --version"), ".")
		if err != nil {
			err = fmt.Errorf("check dockerd installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check dockerd installed success")

		log.Info("check docker installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check docker installed success")

		log.Info("check docker-compose installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker-compose installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check docker-compose installed success")
	} else if o.Mode == "kubernetes" {
		log.Info("check helm v3.x installed")
		outputHelm, _, err := pkg.CommandExec(fmt.Sprintf("helm version --template='{{.Version}}'"), ".")
		if err != nil {
			err = fmt.Errorf("check helm v3.x installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		if !strings.HasPrefix(outputHelm, "v3.") {
			err = fmt.Errorf("check helm v3.x installed error: helm version must be v3.x")
			log.Error(err.Error())
			return err
		}
		log.Success("check helm v3.x installed success")
	}

	var cmdImagePull, cmdImageTag string
	switch o.Runtime {
	case "docker":
		cmdImagePull = "docker pull"
		cmdImageTag = "docker tag"

		log.Info("check kubernetes cluster dockerd installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("dockerd --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster dockerd installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster dockerd installed success")

		log.Info("check kubernetes cluster docker installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster docker installed success")
	case "containerd":
		cmdImagePull = "nerdctl -n k8s.io pull"
		cmdImageTag = "nerdctl -n k8s.io tag"

		log.Info("check kubernetes cluster containerd installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("containerd --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster containerd installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster containerd installed success")

		log.Info("check kubernetes cluster nerdctl installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("nerdctl --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster nerdctl installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster nerdctl installed success")

		log.Info("check kubernetes cluster buildkitd installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("buildkitd --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster buildkitd installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster buildkitd installed success")

		log.Info("check kubernetes cluster buildctl installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("buildctl --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster buildctl installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster buildctl installed success")
	case "crio":
		cmdImagePull = "podman pull"
		cmdImageTag = "podman tag"

		log.Info("check kubernetes cluster crio installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("crio --version"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster crio installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster crio installed success")

		log.Info("check kubernetes cluster podman installed")
		strOut, _, err := pkg.CommandExec(fmt.Sprintf("podman version -f json"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster podman installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal([]byte(strOut), &m)
		if err != nil {
			err = fmt.Errorf("check kubernetes cluster podman installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		_, ok := m["Client"]
		if !ok {
			err = fmt.Errorf("podman version is too old")
			log.Error(err.Error())
			return err
		}
		log.Success("check kubernetes cluster podman installed success")
	default:
		err = fmt.Errorf("--runtime must be docker, containerd or crio")
		return err
	}

	bs, err := pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-README-check.md", pkg.DirInstallScripts, o.Language))
	if err != nil {
		err = fmt.Errorf("get readme error: %s", err.Error())
		return err
	}

	vals := map[string]interface{}{
		"mode":         o.Mode,
		"runtime":      o.Runtime,
		"cmdImagePull": cmdImagePull,
		"cmdImageTag":  cmdImageTag,
	}
	strReadme, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create readme error: %s", err.Error())
		return err
	}

	log.Warning(fmt.Sprintf("########################################################"))
	log.Warning(fmt.Sprintf("KUBERNETES PREREQUISITE README INFO"))
	log.Warning(fmt.Sprintf("########################################################"))
	log.Warning(fmt.Sprintf("\n%s", strReadme))

	return err
}
