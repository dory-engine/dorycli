package pkg

import "embed"

const (
	VersionDoryCli      = "v1.5.0"
	VersionDoryEngine   = "v2.5.0"
	VersionDoryFrontend = "v2.5.0"

	TrivyDb       = "trivy-db-20230718.tar.gz"
	NexusInitData = "nexus-data-init.tar.gz"

	ConfigDirDefault  = ".dorycli"
	ConfigFileDefault = "config.yaml"
	EnvVarConfigFile  = "DORYCONFIG"
	DirInstallScripts = "install_scripts"
	DirInstallConfigs = "install_configs"

	TimeoutDefault = 5

	LogTypeInfo    = "INFO"
	LogTypeWarning = "WARNING"
	LogTypeError   = "ERROR"

	StatusSuccess = "SUCCESS"
	StatusFail    = "FAIL"

	InputValueAbort   = "ABORT"
	InputValueConfirm = "CONFIRM"

	LogStatusInput = "INPUT" // special usage for websocket send notice directives
)

var (
	// !!! go embed function will ignore _* and .* file
	//go:embed install_scripts/* install_scripts/kubernetes/harbor/.helmignore install_scripts/kubernetes/harbor/templates/_helpers.tpl
	FsInstallScripts embed.FS
	//go:embed install_configs/*
	FsInstallConfigs embed.FS

	DefCmdKinds = map[string]string{
		"all":      "",
		"build":    "buildDefs",
		"package":  "packageDefs",
		"artifact": "artifactDefs",
		"deploy":   "deployContainerDefs",
		"setup":    "deployArtifactDefs",
		"istio":    "istioDefs",
		"gateway":  "istioGatewayDef",
		"step":     "customStepDef",
		"pipeline": "pipelineDef",
		"ops":      "customOpsDefs",
		"batch":    "opsBatchDefs",
		"ignore":   "dockerIgnoreDefs",
	}

	AdminCmdKinds = map[string]string{
		"all":  "",
		"user": "user",
		"step": "customStepConf",
		"env":  "envK8s",
		"ct":   "componentTemplate",
		"dbe":  "dockerBuildEnv",
		"grc":  "gitRepoConfig",
		"irc":  "imageRepoConfig",
		"arc":  "artifactRepoConfig",
		"scrc": "scanCodeRepoConfig",
	}

	AccessLevelMaintainer = "maintainer"
	AccessLevelDeveloper  = "developer"
	AccessLevelRunner     = "runner"
	AccessLevels          = []string{
		AccessLevelMaintainer,
		AccessLevelDeveloper,
		AccessLevelRunner,
	}
)
