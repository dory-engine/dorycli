package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

type OptionsInstallHaScript struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileName       string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	OutputDir      string `yaml:"outputDir" json:"outputDir" bson:"outputDir" validate:""`
	Param          struct {
		Stdin []byte `yaml:"stdin" json:"stdin" bson:"stdin" validate:""`
	}
}

func NewOptionsInstallHaScript() *OptionsInstallHaScript {
	var o OptionsInstallHaScript
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallHaScript() *cobra.Command {
	o := NewOptionsInstallHaScript()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("script")

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_install_ha_script_short")
	msgLong := OptCommon.TransLang("cmd_install_ha_script_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_install_ha_script_example", baseName, baseName))

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
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", OptCommon.TransLang("param_install_ha_script_file"))
	cmd.Flags().StringVarP(&o.OutputDir, "output", "o", "", OptCommon.TransLang("param_install_ha_script_output"))

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallHaScript) Complete(cmd *cobra.Command) error {
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

func (o *OptionsInstallHaScript) Validate(args []string) error {
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
func (o *OptionsInstallHaScript) Run(args []string) error {
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
			err = fmt.Errorf("install ha script error: %s", err.Error())
			return err
		}
	}

	var kubernetesHaCluster pkg.KubernetesHaCluster
	err = yaml.Unmarshal(bs, &kubernetesHaCluster)
	if err != nil {
		err = fmt.Errorf("install ha script error: %s", err.Error())
		return err
	}

	err = kubernetesHaCluster.VerifyKubernetesHaCluster()
	if err != nil {
		err = fmt.Errorf("install ha script error: %s", err.Error())
		return err
	}

	bs, _ = json.Marshal(kubernetesHaCluster)
	vals := map[string]interface{}{}
	_ = json.Unmarshal(bs, &vals)

	outputDir := o.OutputDir
	_ = os.MkdirAll(outputDir, 0700)

	fileName := "kubeadm-config.yaml"
	err = o.CreateFile(fmt.Sprintf("%s/kubernetes-ha/%s", pkg.DirInstallScripts, fileName), fmt.Sprintf("%s/%s", outputDir, fileName), vals)
	if err != nil {
		err = fmt.Errorf("create kubernetes ha install file %s error: %s", fileName, err.Error())
		return err
	}

	fileName = "README.md"
	err = o.CreateFile(fmt.Sprintf("%s/kubernetes-ha/%s-%s", pkg.DirInstallScripts, o.Language, fileName), fmt.Sprintf("%s/%s", outputDir, fileName), vals)
	if err != nil {
		err = fmt.Errorf("create kubernetes ha install file %s error: %s", fileName, err.Error())
		return err
	}
	readmeFile := fmt.Sprintf("%s/%s", outputDir, fileName)

	for _, host := range kubernetesHaCluster.MasterHosts {
		keepalivedDir := fmt.Sprintf("%s/%s/keepalived", outputDir, host.Hostname)
		nginxDir := fmt.Sprintf("%s/%s/nginx-lb", outputDir, host.Hostname)
		_ = os.MkdirAll(keepalivedDir, 0700)
		_ = os.MkdirAll(nginxDir, 0700)

		fileNames := []string{
			"keepalived/check_apiserver.sh",
			"keepalived/docker-compose.yaml",
			"nginx-lb/nginx-lb.conf",
			"nginx-lb/docker-compose.yaml",
		}
		for _, name := range fileNames {
			err = o.CreateFile(fmt.Sprintf("%s/kubernetes-ha/%s", pkg.DirInstallScripts, name), fmt.Sprintf("%s/%s/%s", outputDir, host.Hostname, name), vals)
			if err != nil {
				err = fmt.Errorf("create kubernetes ha install file %s error: %s", name, err.Error())
				return err
			}
		}
		bs, _ = json.Marshal(host)
		v := map[string]interface{}{}
		_ = json.Unmarshal(bs, &v)
		vals["host"] = v

		fileName = "keepalived/keepalived.conf"
		err = o.CreateFile(fmt.Sprintf("%s/kubernetes-ha/%s", pkg.DirInstallScripts, fileName), fmt.Sprintf("%s/%s/%s", outputDir, host.Hostname, fileName), vals)
		if err != nil {
			err = fmt.Errorf("create kubernetes ha install file %s error: %s", fileName, err.Error())
			return err
		}
	}

	log.Warning(fmt.Sprintf("all scripts and config files located at: %s", outputDir))
	log.Warning(fmt.Sprintf("change your work directory to %s", outputDir))
	log.Warning(fmt.Sprintf("please check README.md"))

	bs, err = os.ReadFile(readmeFile)
	if err != nil {
		err = fmt.Errorf("get %s error: %s", readmeFile, err.Error())
		return err
	}
	log.Warning(fmt.Sprintf("\n\n%s", string(bs)))

	return err
}

func (o *OptionsInstallHaScript) CreateFile(tplFileName, fileName string, vals map[string]interface{}) error {
	bs, err := pkg.FsInstallScripts.ReadFile(tplFileName)
	if err != nil {
		err = fmt.Errorf("read %s error: %s", fileName, err.Error())
		return err
	}
	str, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("parse %s error: %s", fileName, err.Error())
		return err
	}
	err = os.WriteFile(fileName, []byte(str), 0600)
	if err != nil {
		err = fmt.Errorf("write %s error: %s", fileName, err.Error())
		return err
	}

	return err
}
