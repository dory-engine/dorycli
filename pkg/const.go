package pkg

import (
	"embed"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

const (
	VersionDoryCli      = "v1.7.0"
	VersionDoryEngine   = "v2.7.0"
	VersionDoryFrontend = "v2.7.0"

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
	AdminKindAdminWebhook       = "admin-webhook"

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
		AdminKindAdminWebhook,
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
		AdminKindAdminWebhook:       "adminWebhook",
	}

	ConsoleKindAll             = "all"
	ConsoleKindMember          = "member"
	ConsoleKindPipeline        = "pipeline"
	ConsoleKindPipelineTrigger = "pipeline-trigger"
	ConsoleKindHost            = "host"
	ConsoleKindDatabase        = "database"
	ConsoleKindDebugComponent  = "debug-component"
	ConsoleKindComponent       = "component"

	ConsoleKinds = []string{
		ConsoleKindMember,
		ConsoleKindPipeline,
		ConsoleKindPipelineTrigger,
		ConsoleKindHost,
		ConsoleKindDatabase,
		ConsoleKindDebugComponent,
		ConsoleKindComponent,
	}

	ConsoleCmdKinds = map[string]string{
		ConsoleKindAll:             "",
		ConsoleKindMember:          "member",
		ConsoleKindPipeline:        "pipeline",
		ConsoleKindPipelineTrigger: "pipelineTrigger",
		ConsoleKindHost:            "host",
		ConsoleKindDatabase:        "database",
		ConsoleKindDebugComponent:  "debugComponent",
		ConsoleKindComponent:       "component",
	}

	AccessLevelMaintainer = "maintainer"
	AccessLevelDeveloper  = "developer"
	AccessLevelRunner     = "runner"
	AccessLevels          = []string{
		AccessLevelMaintainer,
		AccessLevelDeveloper,
		AccessLevelRunner,
	}

	TableRenderBorder = tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
	}))

	TableRenderBorderNone = tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Borders:  tw.BorderNone,
		Settings: tw.Settings{Separators: tw.SeparatorsNone, Lines: tw.LinesNone},
	}))

	TableCellConfig = tablewriter.WithConfig(tablewriter.Config{
		Header: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		},
		Row: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		},
	})
)
