package pkg

import "embed"

const (
	VersionDoryCli      = "v1.6.1"
	VersionDoryEngine   = "v2.6.1"
	VersionDoryFrontend = "v2.6.1"

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
	//go:embed language/*
	FsLanguage embed.FS

	DefKindAll             = "all"
	DefKindBuild           = "build"
	DefKindPackage         = "package"
	DefKindArtifact        = "artifact"
	DefKindDeployContainer = "deploy-container"
	DefKindDeployArtifact  = "deploy-artifact"
	DefKindIstio           = "istio"
	DefKindIstioGateway    = "istio-gateway"
	DefKindCustomStep      = "custom-step"
	DefKindPipeline        = "pipeline"
	DefKindCustomOps       = "custom-ops"
	DefKindOpsBatch        = "ops-batch"
	DefKindDockerIgnore    = "docker-ignore"

	DefCmdKinds = map[string]string{
		DefKindAll:             "",
		DefKindBuild:           "buildDefs",
		DefKindPackage:         "packageDefs",
		DefKindArtifact:        "artifactDefs",
		DefKindDeployContainer: "deployContainerDefs",
		DefKindDeployArtifact:  "deployArtifactDefs",
		DefKindIstio:           "istioDefs",
		DefKindIstioGateway:    "istioGatewayDef",
		DefKindCustomStep:      "customStepDef",
		DefKindPipeline:        "pipelineDef",
		DefKindCustomOps:       "customOpsDefs",
		DefKindOpsBatch:        "opsBatchDefs",
		DefKindDockerIgnore:    "dockerIgnoreDefs",
	}

	AdminKindAll                = "all"
	AdminKindComponentTemplate  = "component-template"
	AdminKindEnvK8s             = "env-k8s"
	AdminKindCustomStep         = "custom-step"
	AdminKindUser               = "user"
	AdminKindDockerBuildEnv     = "docker-build-env"
	AdminKindGitRepoConfig      = "git-repo-config"
	AdminKindImageRepoConfig    = "image-repo-config"
	AdminKindArtifactRepoConfig = "artifact-repo-config"
	AdminKindScanCodeRepoConfig = "scan-code-repo-config"

	AdminKinds = []string{
		AdminKindComponentTemplate,
		AdminKindEnvK8s,
		AdminKindCustomStep,
		AdminKindUser,
		AdminKindDockerBuildEnv,
		AdminKindGitRepoConfig,
		AdminKindImageRepoConfig,
		AdminKindArtifactRepoConfig,
		AdminKindScanCodeRepoConfig,
	}

	AdminCmdKinds = map[string]string{
		AdminKindAll:                "",
		AdminKindUser:               "user",
		AdminKindCustomStep:         "customStepConf",
		AdminKindEnvK8s:             "envK8s",
		AdminKindComponentTemplate:  "componentTemplate",
		AdminKindDockerBuildEnv:     "dockerBuildEnv",
		AdminKindGitRepoConfig:      "gitRepoConfig",
		AdminKindImageRepoConfig:    "imageRepoConfig",
		AdminKindArtifactRepoConfig: "artifactRepoConfig",
		AdminKindScanCodeRepoConfig: "scanCodeRepoConfig",
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
