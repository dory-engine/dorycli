package pkg

import "time"

type DoryConfig struct {
	ServerURL   string `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure    bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout     int    `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken string `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	Language    string `yaml:"language" json:"language" bson:"language" validate:""`
}

type InstallDockerImage struct {
	Source     string `yaml:"source" json:"source" bson:"source" validate:"required"`
	Target     string `yaml:"target" json:"target" bson:"target" validate:"required"`
	DockerFile string `yaml:"dockerFile" json:"dockerFile" bson:"dockerFile" validate:""`
	Built      string `yaml:"built" json:"built" bson:"built" validate:""`
	Arm64      string `yaml:"arm64" json:"arm64" bson:"arm64" validate:""`
}

type InstallDockerImages struct {
	InstallDockerImages []InstallDockerImage `yaml:"dockerImages" json:"dockerImages" bson:"dockerImages" validate:""`
}

type ProxyRepo struct {
	Maven  string            `yaml:"maven" json:"maven" bson:"maven" validate:""`
	Npm    string            `yaml:"npm" json:"npm" bson:"npm" validate:""`
	Pip    string            `yaml:"pip" json:"pip" bson:"pip" validate:""`
	Gradle string            `yaml:"gradle" json:"gradle" bson:"gradle" validate:""`
	Go     string            `yaml:"go" json:"go" bson:"go" validate:""`
	Apt    map[string]string `yaml:"apt" json:"apt" bson:"apt" validate:""`
}

type GitRepo struct {
	Type     string `yaml:"type" json:"type" bson:"type" validate:""`
	Internal struct {
		Image      string `yaml:"image" json:"image" bson:"image" validate:"required_with=Image Port"`
		ImageDB    string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:""`
		ImageNginx string `yaml:"imageNginx" json:"imageNginx" bson:"imageNginx" validate:""`
		Port       int    `yaml:"port" json:"port" bson:"port" validate:"required_with=Image Port"`
	} `yaml:"internal" json:"internal" bson:"internal" validate:""`
	External struct {
		ViewUrl       string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Url           string `yaml:"url" json:"url" bson:"url" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Username      string `yaml:"username" json:"username" bson:"username" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Name          string `yaml:"name" json:"name" bson:"name" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Mail          string `yaml:"mail" json:"mail" bson:"mail" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Password      string `yaml:"password" json:"password" bson:"password" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		Token         string `yaml:"token" json:"token" bson:"token" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
		GitWebhookUrl string `yaml:"gitWebhookUrl" json:"gitWebhookUrl" bson:"gitWebhookUrl" validate:"required_with=ViewUrl Url Username Name Mail Password Token GitWebhookUrl"`
	} `yaml:"external" json:"external" bson:"external" validate:""`
}

type ImageRepo struct {
	Type     string `yaml:"type" json:"type" bson:"type" validate:""`
	Internal struct {
		Hostname         string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required_with=Hostname Namespace Version"`
		Namespace        string `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required_with=Hostname Namespace Version"`
		Version          string `yaml:"version" json:"version" bson:"version" validate:"required_with=Hostname Namespace Version"`
		Password         string `yaml:"password" json:"password" bson:"password" validate:""`
		CertsDir         string `yaml:"certsDir" json:"certsDir" bson:"certsDir" validate:"required_with=CertsDir DataDir"`
		DataDir          string `yaml:"dataDir" json:"dataDir" bson:"dataDir" validate:"required_with=CertsDir DataDir"`
		RegistryPassword string `yaml:"registryPassword" json:"registryPassword" bson:"registryPassword" validate:""`
		RegistryHtpasswd string `yaml:"registryHtpasswd" json:"registryHtpasswd" bson:"registryHtpasswd" validate:""`
		VersionBig       string `yaml:"versionBig" json:"versionBig" bson:"versionBig" validate:""`
	} `yaml:"internal" json:"internal" bson:"internal" validate:""`
	External struct {
		Ip       string `yaml:"ip" json:"ip" bson:"ip" validate:"required_with=Ip Hostname Username Password Email"`
		Hostname string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required_with=Ip Hostname Username Password Email"`
		Username string `yaml:"username" json:"username" bson:"username" validate:"required_with=Ip Hostname Username Password Email"`
		Password string `yaml:"password" json:"password" bson:"password" validate:"required_with=Ip Hostname Username Password Email"`
		Email    string `yaml:"email" json:"email" bson:"email" validate:"required_with=Ip Hostname Username Password Email"`
	} `yaml:"external" json:"external" bson:"external" validate:""`
}

type ArtifactRepo struct {
	Type     string `yaml:"type" json:"type" bson:"type" validate:""`
	Internal struct {
		Image    string `yaml:"image" json:"image" bson:"image" validate:"required_with=Image Port PortHub PortGcr PortQuay"`
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required_with=Image Port PortHub PortGcr PortQuay"`
		PortHub  int    `yaml:"portHub" json:"portHub" bson:"portHub" validate:"required_with=Image Port PortHub PortGcr PortQuay"`
		PortGcr  int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required_with=Image Port PortHub PortGcr PortQuay"`
		PortQuay int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required_with=Image Port PortHub PortGcr PortQuay"`
	} `yaml:"internal" json:"internal" bson:"internal" validate:""`
	External struct {
		ViewUrl        string    `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		Schema         string    `yaml:"schema" json:"schema" bson:"schema" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		Hostname       string    `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		Username       string    `yaml:"username" json:"username" bson:"username" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		Password       string    `yaml:"password" json:"password" bson:"password" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PublicRole     string    `yaml:"publicRole" json:"publicRole" bson:"publicRole" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PublicUser     string    `yaml:"publicUser" json:"publicUser" bson:"publicUser" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PublicPassword string    `yaml:"publicPassword" json:"publicPassword" bson:"publicPassword" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PublicEmail    string    `yaml:"publicEmail" json:"publicEmail" bson:"publicEmail" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		Port           int       `yaml:"port" json:"port" bson:"port" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PortHub        int       `yaml:"portHub" json:"portHub" bson:"portHub" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PortGcr        int       `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		PortQuay       int       `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
		ProxyRepo      ProxyRepo `yaml:"proxyRepo" json:"proxyRepo" bson:"proxyRepo" validate:"required_with=ViewUrl Schema Hostname Username Password PublicRole PublicUser PublicPassword PublicEmail Port PortHub PortGcr PortQuay"`
	} `yaml:"external" json:"external" bson:"external" validate:""`
}

type ScanCodeRepo struct {
	Type     string `yaml:"type" json:"type" bson:"type" validate:""`
	Internal struct {
		Image   string `yaml:"image" json:"image" bson:"image" validate:"required_with=Image ImageDB Port"`
		ImageDB string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:"required_with=Image ImageDB Port"`
		Port    int    `yaml:"port" json:"port" bson:"port" validate:"required_with=Image ImageDB Port"`
	} `yaml:"internal" json:"internal" bson:"internal" validate:""`
	External struct {
		ViewUrl string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required_with=ViewUrl Url Token"`
		Url     string `yaml:"url" json:"url" bson:"url" validate:"required_with=ViewUrl Url Token"`
		Token   string `yaml:"token" json:"token" bson:"token" validate:"required_with=ViewUrl Url Token"`
	} `yaml:"external" json:"external" bson:"external" validate:""`
}

type InstallConfig struct {
	InstallMode string `yaml:"installMode" json:"installMode" bson:"installMode" validate:"required"`
	RootDir     string `yaml:"rootDir" json:"rootDir" bson:"rootDir" validate:"required"`
	HostIP      string `yaml:"hostIP" json:"hostIP" bson:"hostIP" validate:"required"`
	ViewURL     string `yaml:"viewURL" json:"viewURL" bson:"viewURL" validate:"required"`
	Dory        struct {
		Namespace    string            `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required"`
		NodeSelector map[string]string `yaml:"nodeSelector" json:"nodeSelector" bson:"nodeSelector" validate:""`
		LicenseKey   string            `yaml:"licenseKey" json:"licenseKey" bson:"licenseKey" validate:""`
		GitRepo      GitRepo           `yaml:"gitRepo" json:"gitRepo" bson:"gitRepo" validate:""`
		ImageRepo    ImageRepo         `yaml:"imageRepo" json:"imageRepo" bson:"imageRepo" validate:""`
		ArtifactRepo ArtifactRepo      `yaml:"artifactRepo" json:"artifactRepo" bson:"artifactRepo" validate:""`
		ScanCodeRepo ScanCodeRepo      `yaml:"scanCodeRepo" json:"scanCodeRepo" bson:"scanCodeRepo" validate:""`
		Openldap     struct {
			Image       string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageAdmin  string `yaml:"imageAdmin" json:"imageAdmin" bson:"imageAdmin" validate:"required"`
			Port        int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			Password    string `yaml:"password" json:"password" bson:"password" validate:""`
			Domain      string `yaml:"domain" json:"domain" bson:"domain" validate:"required"`
			BaseDN      string `yaml:"baseDN" json:"baseDN" bson:"baseDN" validate:"required"`
			ServiceName string `yaml:"serviceName" json:"serviceName" bson:"serviceName" validate:"required"`
		} `yaml:"openldap" json:"openldap" bson:"openldap" validate:"required"`
		Redis struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:""`
		} `yaml:"redis" json:"redis" bson:"redis" validate:"required"`
		Mongo struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:""`
		} `yaml:"mongo" json:"mongo" bson:"mongo" validate:"required"`
		DemoDatabase struct {
			Internal struct {
				DeployName   string `yaml:"deployName" json:"deployName" bson:"deployName" validate:"required_with=DeployName Image User Database Port"`
				Image        string `yaml:"image" json:"image" bson:"image" validate:"required_with=DeployName Image User Database Port"`
				Password     string `yaml:"password" json:"password" bson:"password" validate:""`
				User         string `yaml:"user" json:"user" bson:"user" validate:"required_with=DeployName Image User Database Port"`
				Database     string `yaml:"database" json:"database" bson:"database" validate:"required_with=DeployName Image User Database Port"`
				UserPassword string `yaml:"userPassword" json:"userPassword" bson:"userPassword" validate:""`
				Port         int    `yaml:"port" json:"port" bson:"port" validate:"required_with=DeployName Image User Database Port"`
			} `yaml:"internal" json:"internal" bson:"internal" validate:""`
			External struct {
				DbUrl      string `yaml:"dbUrl" json:"dbUrl" bson:"dbUrl" validate:"required_with=DbUrl DbUser DbPassword"`
				DbUser     string `yaml:"dbUser" json:"dbUser" bson:"dbUser" validate:"required_with=DbUrl DbUser DbPassword"`
				DbPassword string `yaml:"dbPassword" json:"dbPassword" bson:"dbPassword" validate:"required_with=DbUrl DbUser DbPassword"`
			} `yaml:"external" json:"external" bson:"external" validate:""`
		} `yaml:"demoDatabase" json:"demoDatabase" bson:"demoDatabase" validate:"required"`
		DemoHost struct {
			Internal struct {
				DeployName  string `yaml:"deployName" json:"deployName" bson:"deployName" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
				Image       string `yaml:"image" json:"image" bson:"image" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
				Password    string `yaml:"password" json:"password" bson:"password" validate:""`
				PortSsh     int    `yaml:"portSsh" json:"portSsh" bson:"portSsh" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
				NodePortSsh int    `yaml:"nodePortSsh" json:"nodePortSsh" bson:"nodePortSsh" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
				PortWeb     int    `yaml:"portWeb" json:"portWeb" bson:"portWeb" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
				NodePortWeb int    `yaml:"nodePortWeb" json:"nodePortWeb" bson:"nodePortWeb" validate:"required_with=DeployName Image PortSsh NodePortSsh PortWeb NodePortWeb"`
			} `yaml:"internal" json:"internal" bson:"internal" validate:""`
			External struct {
				HostAddr           string `yaml:"hostAddr" json:"hostAddr" bson:"hostAddr" validate:"required_with=HostAddr HostPort HostUser HostPassword NodePortWeb"`
				HostPort           int    `yaml:"hostPort" json:"hostPort" bson:"hostPort" validate:"required_with=HostAddr HostPort HostUser HostPassword NodePortWeb"`
				HostUser           string `yaml:"hostUser" json:"hostUser" bson:"hostUser" validate:"required_with=HostAddr HostPort HostUser HostPassword NodePortWeb"`
				HostPassword       string `yaml:"hostPassword" json:"hostPassword" bson:"hostPassword" validate:"required_with=HostAddr HostPort HostUser HostPassword NodePortWeb"`
				HostBecome         bool   `yaml:"hostBecome" json:"hostBecome" bson:"hostBecome" validate:""`
				HostBecomeUser     string `yaml:"hostBecomeUser" json:"hostBecomeUser" bson:"hostBecomeUser" validate:""`
				HostBecomePassword string `yaml:"hostBecomePassword" json:"hostBecomePassword" bson:"hostBecomePassword" validate:""`
				NodePortWeb        int    `yaml:"nodePortWeb" json:"nodePortWeb" bson:"nodePortWeb" validate:"required_with=HostAddr HostPort HostUser HostPassword NodePortWeb"`
			} `yaml:"external" json:"external" bson:"external" validate:""`
		} `yaml:"demoHost" json:"demoHost" bson:"demoHost" validate:"required"`
		Docker struct {
			Image        string `yaml:"image" json:"image" bson:"image" validate:"required"`
			DockerName   string `yaml:"dockerName" json:"dockerName" bson:"dockerName" validate:"required"`
			DockerNumber int    `yaml:"dockerNumber" json:"dockerNumber" bson:"dockerNumber" validate:"required"`
		} `yaml:"docker" json:"docker" bson:"docker" validate:"required"`
		Doryengine struct {
			Port int `yaml:"port" json:"port" bson:"port" validate:"required"`
		} `yaml:"doryengine" json:"doryengine" bson:"doryengine" validate:"required"`
		NexusInitData string `yaml:"nexusInitData" json:"nexusInitData" bson:"nexusInitData" validate:""`
	} `yaml:"dory" json:"dory" bson:"dory" validate:"required"`
	Account struct {
		AdminUser struct {
			Username string `yaml:"username" json:"username" bson:"username" validate:"required"`
			Name     string `yaml:"name" json:"name" bson:"name" validate:"required"`
			Mail     string `yaml:"mail" json:"mail" bson:"mail" validate:"required"`
			Mobile   string `yaml:"mobile" json:"mobile" bson:"mobile" validate:"required"`
		} `yaml:"adminUser" json:"adminUser" bson:"adminUser" validate:"required"`
		Mail struct {
			Host     string `yaml:"host" json:"host" bson:"host" validate:"required"`
			Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			Username string `yaml:"username" json:"username" bson:"username" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:"required"`
			Mode     string `yaml:"mode" json:"mode" bson:"mode" validate:""`
			From     string `yaml:"from" json:"from" bson:"from" validate:"required"`
		} `yaml:"mail" json:"mail" bson:"mail" validate:"required"`
	} `yaml:"account" json:"account" bson:"account" validate:"required"`
	Kubernetes struct {
		EnvName                    string `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
		EnvDesc                    string `yaml:"envDesc" json:"envDesc" bson:"envDesc" validate:"required"`
		Timezone                   string `yaml:"timezone" json:"timezone" bson:"timezone" validate:"required"`
		Runtime                    string `yaml:"runtime" json:"runtime" bson:"runtime" validate:"required,oneof=docker containerd crio"`
		Host                       string `yaml:"host" json:"host" bson:"host" validate:"required"`
		ViewHost                   string `yaml:"viewHost" json:"viewHost" bson:"viewHost" validate:"required"`
		Port                       int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Token                      string `yaml:"token" json:"token" bson:"token" validate:"required"`
		CaCrtBase64                string `yaml:"caCrtBase64" json:"caCrtBase64" bson:"caCrtBase64" validate:"required"`
		DashboardUrl               string `yaml:"dashboardUrl" json:"dashboardUrl" bson:"dashboardUrl" validate:"required"`
		DashboardUrlNetworkPolicy  string `yaml:"dashboardUrlNetworkPolicy" json:"dashboardUrlNetworkPolicy" bson:"dashboardUrlNetworkPolicy" validate:"required"`
		IstioNamespace             string `yaml:"istioNamespace" json:"istioNamespace" bson:"istioNamespace" validate:"required"`
		IngressControllerNamespace string `yaml:"ingressControllerNamespace" json:"ingressControllerNamespace" bson:"ingressControllerNamespace" validate:"required"`
		PvConfigLocal              struct {
			LocalPath string `yaml:"localPath" json:"localPath" bson:"localPath" validate:""`
		} `yaml:"pvConfigLocal" json:"pvConfigLocal" bson:"pvConfigLocal" validate:""`
		PvConfigNfs struct {
			NfsPath   string `yaml:"nfsPath" json:"nfsPath" bson:"nfsPath" validate:"required_with=NfsPath NfsServer"`
			NfsServer string `yaml:"nfsServer" json:"nfsServer" bson:"nfsServer" validate:"required_with=NfsPath NfsServer"`
		} `yaml:"pvConfigNfs" json:"pvConfigNfs" bson:"pvConfigNfs" validate:""`
		PvConfigCephfs struct {
			CephPath     string   `yaml:"cephPath" json:"cephPath" bson:"cephPath" validate:"required_with=CephPath CephUser CephSecret CephMonitors"`
			CephUser     string   `yaml:"cephUser" json:"cephUser" bson:"cephUser" validate:"required_with=CephPath CephUser CephSecret CephMonitors"`
			CephSecret   string   `yaml:"cephSecret" json:"cephSecret" bson:"cephSecret" validate:"required_with=CephPath CephUser CephSecret CephMonitors"`
			CephMonitors []string `yaml:"cephMonitors" json:"cephMonitors" bson:"cephMonitors" validate:"required_with=CephPath CephUser CephSecret CephMonitors"`
		} `yaml:"pvConfigCephfs" json:"pvConfigCephfs" bson:"pvConfigCephfs" validate:""`
	} `yaml:"kubernetes" json:"kubernetes" bson:"kubernetes" validate:"required"`
}

type QuotaPod struct {
	MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:"required"`
	CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:"required"`
	MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:"required"`
	CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:"required"`
}

type QuotaResource struct {
	MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:"required"`
	CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:"required"`
	MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:"required"`
	CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:"required"`
	PodsLimit     int    `yaml:"podsLimit" json:"podsLimit" bson:"podsLimit" validate:"required"`
}

type KubePodState struct {
	Waiting struct {
		Reason string `yaml:"reason" json:"reason" bson:"reason" validate:""`
	} `yaml:"waiting" json:"waiting" bson:"waiting" validate:""`
	Running struct {
		StartedAt string `yaml:"startedAt" json:"startedAt" bson:"startedAt" validate:""`
	} `yaml:"running" json:"running" bson:"running" validate:""`
	Terminated struct {
		Reason   string `yaml:"reason" json:"reason" bson:"reason" validate:""`
		ExitCode int    `yaml:"exitCode" json:"exitCode" bson:"exitCode" validate:""`
		Signal   int    `yaml:"signal" json:"signal" bson:"signal" validate:""`
	} `yaml:"terminated" json:"terminated" bson:"terminated" validate:""`
}

type KubePod struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion" bson:"apiVersion" validate:"required"`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	MetaData   struct {
		Name              string            `yaml:"name" json:"name" bson:"name" validate:"required"`
		NameSpace         string            `yaml:"namespace" json:"namespace" bson:"namespace" validate:""`
		Labels            map[string]string `yaml:"labels" json:"labels" bson:"labels" validate:""`
		Annotations       map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
		CreationTimestamp string            `yaml:"creationTimestamp" json:"creationTimestamp" bson:"creationTimestamp" validate:""`
		DeletionTimestamp string            `yaml:"deletionTimestamp" json:"deletionTimestamp" bson:"deletionTimestamp" validate:""`
		OwnerReferences   []struct {
			ApiVersion         string `yaml:"apiVersion" json:"apiVersion" bson:"apiVersion" validate:"required"`
			BlockOwnerDeletion bool   `yaml:"blockOwnerDeletion" json:"blockOwnerDeletion" bson:"blockOwnerDeletion" validate:""`
			Controller         bool   `yaml:"controller" json:"controller" bson:"controller" validate:""`
			Kind               string `yaml:"kind" json:"kind" bson:"kind" validate:""`
			Name               string `yaml:"name" json:"name" bson:"name" validate:""`
			Uid                string `yaml:"uid" json:"uid" bson:"uid" validate:""`
		} `yaml:"ownerReferences" json:"ownerReferences" bson:"ownerReferences" validate:""`
	} `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Spec struct {
		Containers []struct {
			Name  string `yaml:"name" json:"name" bson:"name" validate:""`
			Image string `yaml:"image" json:"image" bson:"image" validate:""`
		} `yaml:"containers" json:"containers" bson:"containers" validate:""`
	} `yaml:"spec" json:"spec" bson:"spec" validate:""`
	Status struct {
		Reason     string `yaml:"reason" json:"reason" bson:"reason" validate:""`
		Conditions []struct {
			Type   string `yaml:"type" json:"type" bson:"type" validate:""`
			Status string `yaml:"status" json:"status" bson:"status" validate:""`
		} `yaml:"conditions" json:"conditions" bson:"conditions" validate:""`
		ContainerStatuses []struct {
			Name         string       `yaml:"name" json:"name" bson:"name" validate:""`
			Ready        bool         `yaml:"ready" json:"ready" bson:"ready" validate:""`
			Started      bool         `yaml:"started" json:"started" bson:"started" validate:""`
			RestartCount int          `yaml:"restartCount" json:"restartCount" bson:"restartCount" validate:""`
			State        KubePodState `yaml:"state" json:"state" bson:"state" validate:""`
		} `yaml:"containerStatuses" json:"containerStatuses" bson:"containerStatuses" validate:""`
		InitContainerStatuses []struct {
			Name         string       `yaml:"name" json:"name" bson:"name" validate:""`
			Ready        bool         `yaml:"ready" json:"ready" bson:"ready" validate:""`
			Started      bool         `yaml:"started" json:"started" bson:"started" validate:""`
			RestartCount int          `yaml:"restartCount" json:"restartCount" bson:"restartCount" validate:""`
			State        KubePodState `yaml:"state" json:"state" bson:"state" validate:""`
		} `yaml:"initContainerStatuses" json:"initContainerStatuses" bson:"initContainerStatuses" validate:""`
		Phase     string    `yaml:"phase" json:"phase" bson:"phase" validate:""`
		PodIP     string    `yaml:"podIP" json:"podIP" bson:"podIP" validate:""`
		StartTime time.Time `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type KubePodList struct {
	Items []KubePod `yaml:"items" json:"items" bson:"items" validate:""`
}

type EnvNodePort struct {
	NodePortStart int  `yaml:"nodePortStart" json:"nodePortStart" bson:"nodePortStart" validate:""`
	NodePortEnd   int  `yaml:"nodePortEnd" json:"nodePortEnd" bson:"nodePortEnd" validate:""`
	IsDefault     bool `yaml:"isDefault" json:"isDefault" bson:"isDefault" validate:""`
}

type EnvDebugSpec struct {
	DebugSSHNodePort int `yaml:"debugSSHNodePort" json:"debugSSHNodePort" bson:"debugSSHNodePort" validate:""`
	DebugVNCNodePort int `yaml:"debugVNCNodePort" json:"debugVNCNodePort" bson:"debugVNCNodePort" validate:""`
}

type ProjectAvailableEnv struct {
	EnvName                   string               `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
	DeployContainerDefs       []DeployContainerDef `yaml:"deployContainerDefs" json:"deployContainerDefs" bson:"deployContainerDefs" validate:""`
	UpdateDeployContainerDefs bool                 `yaml:"updateDeployContainerDefs" json:"updateDeployContainerDefs" bson:"updateDeployContainerDefs" validate:""`
	DeployArtifactDefs        []DeployArtifactDef  `yaml:"deployArtifactDefs" json:"deployArtifactDefs" bson:"deployArtifactDefs" validate:""`
	UpdateDeployArtifactDefs  bool                 `yaml:"updateDeployArtifactDefs" json:"updateDeployArtifactDefs" bson:"updateDeployArtifactDefs" validate:""`
	IstioDefs                 []IstioDef           `yaml:"istioDefs" json:"istioDefs" bson:"istioDefs" validate:""`
	UpdateIstioDefs           bool                 `yaml:"updateIstioDefs" json:"updateIstioDefs" bson:"updateIstioDefs" validate:""`
	IstioGatewayDef           IstioGatewayDef      `yaml:"istioGatewayDef" json:"istioGatewayDef" bson:"istioGatewayDef" validate:""`
	UpdateIstioGatewayDef     bool                 `yaml:"updateIstioGatewayDef" json:"updateIstioGatewayDef" bson:"updateIstioGatewayDef" validate:""`
	CustomStepDefs            CustomStepDefs       `yaml:"customStepDefs" json:"customStepDefs" bson:"customStepDefs" validate:""`
	EnvDebugSpec              EnvDebugSpec         `yaml:"envDebugSpec" json:"envDebugSpec" bson:"envDebugSpec" validate:""`
	EnvNodePorts              []EnvNodePort        `yaml:"envNodePorts" json:"envNodePorts" bson:"envNodePorts" validate:""`
	ErrMsgDeployContainerDefs string               `yaml:"errMsgDeployContainerDefs" json:"errMsgDeployContainerDefs" bson:"errMsgDeployContainerDefs" validate:""`
	ErrMsgDeployArtifactDefs  string               `yaml:"errMsgDeployArtifactDefs" json:"errMsgDeployArtifactDefs" bson:"errMsgDeployArtifactDefs" validate:""`
	ErrMsgIstioDefs           string               `yaml:"errMsgIstioDefs" json:"errMsgIstioDefs" bson:"errMsgIstioDefs" validate:""`
	ErrMsgIstioGatewayDef     string               `yaml:"errMsgIstioGatewayDef" json:"errMsgIstioGatewayDef" bson:"errMsgIstioGatewayDef" validate:""`
	ErrMsgCustomStepDefs      map[string]string    `yaml:"errMsgCustomStepDefs" json:"errMsgCustomStepDefs" bson:"errMsgCustomStepDefs" validate:""`
}

type Module struct {
	ModuleName string `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:""`
	IsLatest   bool   `yaml:"isLatest" json:"isLatest" bson:"isLatest" validate:""`
	Hidden     bool   `yaml:"hidden" json:"hidden" bson:"hidden" validate:""`
}

type PipelineBuild struct {
	Name string `yaml:"name" json:"name" bson:"name" validate:""`
	Run  bool   `yaml:"run" json:"run" bson:"run" validate:""`
}

type Pipeline struct {
	PipelineName   string   `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
	BranchName     string   `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Envs           []string `yaml:"envs" json:"envs" bson:"envs" validate:""`
	EnvProductions []string `yaml:"envProductions" json:"envProductions" bson:"envProductions" validate:""`
	SuccessCount   int      `yaml:"successCount" json:"successCount" bson:"successCount" validate:""`
	FailCount      int      `yaml:"failCount" json:"failCount" bson:"failCount" validate:""`
	AbortCount     int      `yaml:"abortCount" json:"abortCount" bson:"abortCount" validate:""`
	Status         struct {
		Result    string `yaml:"result" json:"result" bson:"result" validate:""`
		StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
		Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
	ErrMsgPipelineDef string `yaml:"errMsgPipelineDef" json:"errMsgPipelineDef" bson:"errMsgPipelineDef" validate:""`
	PipelineDef       struct {
		Builds       []PipelineBuild `yaml:"builds" json:"builds" bson:"builds" validate:""`
		PipelineStep PipelineStepDef `yaml:"pipelineStep" json:"pipelineStep" bson:"pipelineStep" validate:""`
	} `yaml:"pipelineDef" json:"pipelineDef" bson:"pipelineDef" validate:""`
}

type GitRepoDir struct {
	BuildSettingsDir   string `yaml:"buildSettingsDir" json:"buildSettingsDir" bson:"buildSettingsDir" validate:""`
	DatabaseScriptsDir string `yaml:"databaseScriptsDir" json:"databaseScriptsDir" bson:"databaseScriptsDir" validate:""`
	DemoCodesDir       string `yaml:"demoCodesDir" json:"demoCodesDir" bson:"demoCodesDir" validate:""`
	DeployScriptsDir   string `yaml:"deployScriptsDir" json:"deployScriptsDir" bson:"deployScriptsDir" validate:""`
	DocumentsDir       string `yaml:"documentsDir" json:"documentsDir" bson:"documentsDir" validate:""`
	TestScriptsDir     string `yaml:"testScriptsDir" json:"testScriptsDir" bson:"testScriptsDir" validate:""`
}

type Project struct {
	TenantCode  string      `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	ProjectInfo ProjectInfo `yaml:"projectInfo" json:"projectInfo" bson:"projectInfo" validate:""`
	ProjectRepo struct {
		GitRepo          string     `yaml:"gitRepo" json:"gitRepo" bson:"gitRepo" validate:""`
		GitRepoDir       GitRepoDir `yaml:"gitRepoDir" json:"gitRepoDir" bson:"gitRepoDir" validate:""`
		ImageRepo        string     `yaml:"imageRepo" json:"imageRepo" bson:"imageRepo" validate:""`
		ArtifactRepo     string     `yaml:"artifactRepo" json:"artifactRepo" bson:"artifactRepo" validate:""`
		ArtifactRepoType string     `yaml:"artifactRepoType" json:"artifactRepoType" bson:"artifactRepoType" validate:""`
		ScanCodeRepo     string     `yaml:"scanCodeRepo" json:"scanCodeRepo" bson:"scanCodeRepo" validate:""`
	} `yaml:"projectRepo" json:"projectRepo" bson:"projectRepo" validate:""`
	ProjectNodePorts []struct {
		EnvName      string        `yaml:"envName" json:"envName" bson:"envName" validate:""`
		EnvNodePorts []EnvNodePort `yaml:"envNodePorts" json:"envNodePorts" bson:"envNodePorts" validate:""`
		NodePorts    []int         `yaml:"nodePorts" json:"nodePorts" bson:"nodePorts" validate:""`
	} `yaml:"projectNodePorts" json:"projectNodePorts" bson:"projectNodePorts" validate:""`
	Modules      map[string][]Module `yaml:"modules" json:"modules" bson:"modules" validate:""`
	OpsBatchDefs []struct {
		OpsBatchDesc string `yaml:"opsBatchDesc" json:"opsBatchDesc" bson:"opsBatchDesc" validate:""`
		OpsBatchName string `yaml:"opsBatchName" json:"opsBatchName" bson:"opsBatchName" validate:""`
	} `yaml:"opsBatchDefs" json:"opsBatchDefs" bson:"opsBatchDefs" validate:""`
	Pipelines []Pipeline `yaml:"pipelines" json:"pipelines" bson:"pipelines" validate:""`
}

type Run struct {
	ProjectName  string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	PipelineName string `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
	RunName      string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	BranchName   string `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	TriggerKind  string `yaml:"triggerKind" json:"triggerKind" bson:"triggerKind" validate:""`
	StartUser    string `yaml:"startUser" json:"startUser" bson:"startUser" validate:""`
	AbortUser    string `yaml:"abortUser" json:"abortUser" bson:"abortUser" validate:""`
	Status       struct {
		Result    string `yaml:"result" json:"result" bson:"result" validate:""`
		StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
		Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type RunInputOption struct {
	Name  string `yaml:"name" json:"name" bson:"name" validate:""`
	Value string `yaml:"value" json:"value" bson:"value" validate:""`
}

type RunInput struct {
	PhaseID    string           `yaml:"phaseID" json:"phaseID" bson:"phaseID" validate:""`
	Title      string           `yaml:"title" json:"title" bson:"title" validate:""`
	Desc       string           `yaml:"desc" json:"desc" bson:"desc" validate:""`
	IsMultiple bool             `yaml:"isMultiple" json:"isMultiple" bson:"isMultiple" validate:""`
	IsApiOnly  bool             `yaml:"isApiOnly" json:"isApiOnly" bson:"isApiOnly" validate:""`
	Options    []RunInputOption `yaml:"options" json:"options" bson:"options" validate:""`
}

type WsRunLog struct {
	ID         string `yaml:"ID" json:"ID" bson:"ID" validate:""`
	LogType    string `yaml:"logType" json:"logType" bson:"logType" validate:""`
	Content    string `yaml:"content" json:"content" bson:"content" validate:""`
	RunName    string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	PhaseID    string `yaml:"phaseID" json:"phaseID" bson:"phaseID" validate:""`
	StageID    string `yaml:"stageID" json:"stageID" bson:"stageID" validate:""`
	StepID     string `yaml:"stepID" json:"stepID" bson:"stepID" validate:""`
	CreateTime string `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
}

type WsAdminLog struct {
	ID        string `yaml:"ID" json:"ID" bson:"ID" validate:""`
	LogType   string `yaml:"logType" json:"logType" bson:"logType" validate:""`
	Content   string `yaml:"content" json:"content" bson:"content" validate:""`
	StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
	EndTime   string `yaml:"endTime" json:"endTime" bson:"endTime" validate:""`
	Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
}

type CustomStepModuleDef struct {
	ModuleName         string   `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:"required"`
	RelatedStepModules []string `yaml:"relatedStepModules" json:"relatedStepModules" bson:"relatedStepModules" validate:""`
	ManualEnable       bool     `yaml:"manualEnable" json:"manualEnable" bson:"manualEnable" validate:""`
	ParamInputYaml     string   `yaml:"paramInputYaml" json:"paramInputYaml" bson:"paramInputYaml" validate:""`
	CheckVarToIgnore   string   `yaml:"checkVarToIgnore" json:"checkVarToIgnore" bson:"checkVarToIgnore" validate:""`
	IsPatch            bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type CustomStepDef struct {
	EnableMode                 string                `yaml:"enableMode" json:"enableMode" bson:"enableMode" validate:""`
	CustomStepModuleDefs       []CustomStepModuleDef `yaml:"customStepModuleDefs" json:"customStepModuleDefs" bson:"customStepModuleDefs" validate:""`
	UpdateCustomStepModuleDefs bool                  `yaml:"updateCustomStepModuleDefs" json:"updateCustomStepModuleDefs" bson:"updateCustomStepModuleDefs" validate:""`
}

type CustomStepDefs map[string]CustomStepDef

type CustomOpsDef struct {
	CustomOpsName  string   `yaml:"customOpsName" json:"customOpsName" bson:"customOpsName" validate:"required"`
	CustomOpsDesc  string   `yaml:"customOpsDesc" json:"customOpsDesc" bson:"customOpsDesc" validate:"required"`
	CustomOpsSteps []string `yaml:"customOpsSteps" json:"customOpsSteps" bson:"customOpsSteps" validate:"required"`
	IsPatch        bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type OpsBatchDef struct {
	OpsBatchName string   `yaml:"opsBatchName" json:"opsBatchName" bson:"opsBatchName" validate:"required"`
	OpsBatchDesc string   `yaml:"opsBatchDesc" json:"opsBatchDesc" bson:"opsBatchDesc" validate:"required"`
	Batches      []string `yaml:"batches" json:"batches" bson:"batches" validate:"required"`
	IsPatch      bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type BuildDef struct {
	BuildName          string   `yaml:"buildName" json:"buildName" bson:"buildName" validate:"required"`
	BuildPhaseID       int      `yaml:"buildPhaseID" json:"buildPhaseID" bson:"buildPhaseID" validate:"required,gt=0"`
	BuildPath          string   `yaml:"buildPath" json:"buildPath" bson:"buildPath" validate:"required"`
	BuildEnv           string   `yaml:"buildEnv" json:"buildEnv" bson:"buildEnv" validate:"required"`
	BuildCmds          []string `yaml:"buildCmds" json:"buildCmds" bson:"buildCmds" validate:"required,dive"`
	BuildChecks        []string `yaml:"buildChecks" json:"buildChecks" bson:"buildChecks" validate:"required,dive"`
	BuildCaches        []string `yaml:"buildCaches" json:"buildCaches" bson:"buildCaches" validate:""`
	SonarExtraSettings []string `yaml:"sonarExtraSettings" json:"sonarExtraSettings" bson:"sonarExtraSettings" validate:""`
	IsPatch            bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type PackageDef struct {
	PackageName   string   `yaml:"packageName" json:"packageName" bson:"packageName" validate:"required"`
	RelatedBuilds []string `yaml:"relatedBuilds" json:"relatedBuilds" bson:"relatedBuilds" validate:"required"`
	DockerFile    string   `yaml:"dockerFile" json:"dockerFile" bson:"dockerFile" validate:"required"`
	IsPatch       bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type ArtifactDef struct {
	ArtifactName  string   `yaml:"artifactName" json:"artifactName" bson:"artifactName" validate:"required"`
	RelatedBuilds []string `yaml:"relatedBuilds" json:"relatedBuilds" bson:"relatedBuilds" validate:"required"`
	Artifacts     []string `yaml:"artifacts" json:"artifacts" bson:"artifacts" validate:"required"`
	IsPatch       bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type DeployHttpGet struct {
	Path        string `yaml:"path" json:"path" bson:"path" validate:""`
	Port        int    `yaml:"port" json:"port" bson:"port" validate:""`
	HttpHeaders []struct {
		Name  string `yaml:"name" json:"name" bson:"name" validate:"required"`
		Value string `yaml:"value" json:"value" bson:"value" validate:"required"`
	} `yaml:"httpHeaders" json:"httpHeaders" bson:"httpHeaders" validate:"dive"`
	Scheme string `yaml:"scheme" json:"scheme" bson:"scheme" validate:""`
}

type NameValue struct {
	Name  string `yaml:"name" json:"name" bson:"name" validate:""`
	Value string `yaml:"value" json:"value" bson:"value" validate:""`
}

type DeployContainerPatch struct {
	ResourceKind string `yaml:"resourceKind" json:"resourceKind" bson:"resourceKind" validate:"required"`
	Path         string `yaml:"path" json:"path" bson:"path" validate:"required"`
	Content      string `yaml:"content" json:"content" bson:"content" validate:"required"`
}

type ConfigPath struct {
	LocalPath string `yaml:"localPath" json:"localPath" bson:"localPath" validate:""`
	PvcName   string `yaml:"pvcName" json:"pvcName" bson:"pvcName" validate:""`
	PodPath   string `yaml:"podPath" json:"podPath" bson:"podPath" validate:""`
}

type DeployConfigMap struct {
	Name         string   `yaml:"name" json:"name" bson:"name" validate:""`
	FromFileType string   `yaml:"fromFileType" json:"fromFileType" bson:"fromFileType" validate:""`
	Paths        []string `yaml:"paths" json:"paths" bson:"paths" validate:""`
}

type DeploySecret struct {
	Name         string   `yaml:"name" json:"name" bson:"name" validate:""`
	SecretType   string   `yaml:"secretType" json:"secretType" bson:"secretType" validate:""`
	FromFileType string   `yaml:"fromFileType" json:"fromFileType" bson:"fromFileType" validate:""`
	Paths        []string `yaml:"paths" json:"paths" bson:"paths" validate:""`
	DockerConfig string   `yaml:"dockerConfig" json:"dockerConfig" bson:"dockerConfig" validate:""`
	Cert         string   `yaml:"cert" json:"cert" bson:"cert" validate:""`
	Key          string   `yaml:"key" json:"key" bson:"key" validate:""`
}

type DeployNodePort struct {
	Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
	NodePort int    `yaml:"nodePort" json:"nodePort" bson:"nodePort" validate:"required"`
	Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=HTTP TCP UDP SCTP"`
}

type DeployLocalPort struct {
	Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
	Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=HTTP TCP UDP SCTP"`
	Ingress  struct {
		DomainName     string `yaml:"domainName" json:"domainName" bson:"domainName" validate:""`
		PathPrefix     string `yaml:"pathPrefix" json:"pathPrefix" bson:"pathPrefix" validate:""`
		CertSelfSigned bool   `yaml:"certSelfSigned" json:"certSelfSigned" bson:"certSelfSigned" validate:""`
		CertPath       string `yaml:"certPath" json:"certPath" bson:"certPath" validate:""`
	} `yaml:"ingress" json:"ingress" bson:"ingress" validate:""`
}

type DeployVolume struct {
	PathInPod string `yaml:"pathInPod" json:"pathInPod" bson:"pathInPod" validate:"required"`
	PathInPv  string `yaml:"pathInPv" json:"pathInPv" bson:"pathInPv" validate:"required"`
	Pvc       string `yaml:"pvc" json:"pvc" bson:"pvc" validate:""`
}

type DeployContainerDef struct {
	DeployName          string            `yaml:"deployName" json:"deployName" bson:"deployName" validate:"required"`
	RelatedPackage      string            `yaml:"relatedPackage" json:"relatedPackage" bson:"relatedPackage" validate:"required"`
	DeployImageTag      string            `yaml:"deployImageTag" json:"deployImageTag" bson:"deployImageTag" validate:""`
	DeployLabels        map[string]string `yaml:"deployLabels" json:"deployLabels" bson:"deployLabels" validate:""`
	DeployType          string            `yaml:"deployType" json:"deployType" bson:"deployType" validate:""`
	DeployHeadless      bool              `yaml:"deployHeadless" json:"deployHeadless" bson:"deployHeadless" validate:""`
	PodManagementPolicy string            `yaml:"podManagementPolicy" json:"podManagementPolicy" bson:"podManagementPolicy" validate:""`
	DeployMeta          struct {
		Labels      []NameValue `yaml:"labels" json:"labels" bson:"labels" validate:""`
		Annotations []NameValue `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
	} `yaml:"deployMeta" json:"deployMeta" bson:"deployMeta" validate:""`
	DeploySessionAffinityTimeoutSeconds int               `yaml:"deploySessionAffinityTimeoutSeconds" json:"deploySessionAffinityTimeoutSeconds" bson:"deploySessionAffinityTimeoutSeconds" validate:""`
	DeployNodePorts                     []DeployNodePort  `yaml:"deployNodePorts" json:"deployNodePorts" bson:"deployNodePorts" validate:"dive"`
	DeployLocalPorts                    []DeployLocalPort `yaml:"deployLocalPorts" json:"deployLocalPorts" bson:"deployLocalPorts" validate:"dive"`
	DeployReplicas                      int               `yaml:"deployReplicas" json:"deployReplicas" bson:"deployReplicas" validate:"required"`
	HpaConfig                           struct {
		MaxReplicas                 int    `yaml:"maxReplicas" json:"maxReplicas" bson:"maxReplicas" validate:""`
		MemoryAverageValue          string `yaml:"memoryAverageValue" json:"memoryAverageValue" bson:"memoryAverageValue" validate:""`
		MemoryAverageRequestPercent int    `yaml:"memoryAverageRequestPercent" json:"memoryAverageRequestPercent" bson:"memoryAverageRequestPercent" validate:""`
		CpuAverageValue             string `yaml:"cpuAverageValue" json:"cpuAverageValue" bson:"cpuAverageValue" validate:""`
		CpuAverageRequestPercent    int    `yaml:"cpuAverageRequestPercent" json:"cpuAverageRequestPercent" bson:"cpuAverageRequestPercent" validate:""`
	} `yaml:"hpaConfig" json:"hpaConfig" bson:"hpaConfig" validate:""`
	DeployEnvs        []string       `yaml:"deployEnvs" json:"deployEnvs" bson:"deployEnvs" validate:""`
	DeployCommand     string         `yaml:"deployCommand" json:"deployCommand" bson:"deployCommand" validate:""`
	DeployCmds        []string       `yaml:"deployCmds" json:"deployCmds" bson:"deployCmds" validate:""`
	DeployArgs        []string       `yaml:"deployArgs" json:"deployArgs" bson:"deployArgs" validate:""`
	DeployResources   QuotaPod       `yaml:"deployResources" json:"deployResources" bson:"deployResources" validate:""`
	DeployVolumes     []DeployVolume `yaml:"deployVolumes" json:"deployVolumes" bson:"deployVolumes" validate:"dive"`
	DeployHealthCheck struct {
		CheckPort              int           `yaml:"checkPort" json:"checkPort" bson:"checkPort" validate:""`
		Exec                   string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
		ExecCmds               []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
		HttpGet                DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		ReadinessDelaySeconds  int           `yaml:"readinessDelaySeconds" json:"readinessDelaySeconds" bson:"readinessDelaySeconds" validate:""`
		ReadinessPeriodSeconds int           `yaml:"readinessPeriodSeconds" json:"readinessPeriodSeconds" bson:"readinessPeriodSeconds" validate:""`
		LivenessDelaySeconds   int           `yaml:"livenessDelaySeconds" json:"livenessDelaySeconds" bson:"livenessDelaySeconds" validate:""`
		LivenessPeriodSeconds  int           `yaml:"livenessPeriodSeconds" json:"livenessPeriodSeconds" bson:"livenessPeriodSeconds" validate:""`
		StartupDelaySeconds    int           `yaml:"startupDelaySeconds" json:"startupDelaySeconds" bson:"startupDelaySeconds" validate:""`
		StartupPeriodSeconds   int           `yaml:"startupPeriodSeconds" json:"startupPeriodSeconds" bson:"startupPeriodSeconds" validate:""`
	} `yaml:"deployHealthCheck" json:"deployHealthCheck" bson:"deployHealthCheck" validate:""`
	DependServices []struct {
		DependName string `yaml:"dependName" json:"dependName" bson:"dependName" validate:"required"`
		DependPort int    `yaml:"dependPort" json:"dependPort" bson:"dependPort" validate:"required"`
		DependType string `yaml:"dependType" json:"dependType" bson:"dependType" validate:"oneof=TCP UDP"`
	} `yaml:"dependServices" json:"dependServices" bson:"dependServices" validate:"dive"`
	HostAliases []struct {
		Ip        string   `yaml:"ip" json:"ip" bson:"ip" validate:"required,ip"`
		Hostnames []string `yaml:"hostnames" json:"hostnames" bson:"hostnames" validate:"required"`
	} `yaml:"hostAliases" json:"hostAliases" bson:"hostAliases" validate:"dive"`
	SecurityContext struct {
		RunAsUser  int `yaml:"runAsUser" json:"runAsUser" bson:"runAsUser" validate:""`
		RunAsGroup int `yaml:"runAsGroup" json:"runAsGroup" bson:"runAsGroup" validate:""`
	} `yaml:"securityContext" json:"securityContext" bson:"securityContext" validate:""`
	DeployConfigSettings []ConfigPath      `yaml:"deployConfigSettings" json:"deployConfigSettings" bson:"deployConfigSettings" validate:""`
	DeployConfigMaps     []DeployConfigMap `yaml:"deployConfigMaps" json:"deployConfigMaps" bson:"deployConfigMaps" validate:""`
	DeploySecrets        []DeploySecret    `yaml:"deploySecrets" json:"deploySecrets" bson:"deploySecrets" validate:""`
	Lifecycle            struct {
		PostStart struct {
			Exec     string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
			ExecCmds []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
			HttpGet  DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		} `yaml:"postStart" json:"postStart" bson:"postStart" validate:""`
		PreStop struct {
			Exec     string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
			ExecCmds []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
			HttpGet  DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		} `yaml:"preStop" json:"preStop" bson:"preStop" validate:""`
	} `yaml:"lifecycle" json:"lifecycle" bson:"lifecycle" validate:""`
	WorkingDir                    string      `yaml:"workingDir" json:"workingDir" bson:"workingDir" validate:""`
	NodeSelector                  []NameValue `yaml:"nodeSelector" json:"nodeSelector" bson:"nodeSelector" validate:""`
	NodeName                      string      `yaml:"nodeName" json:"nodeName" bson:"nodeName" validate:""`
	TerminationGracePeriodSeconds int         `yaml:"terminationGracePeriodSeconds" json:"terminationGracePeriodSeconds" bson:"terminationGracePeriodSeconds" validate:""`
	Subdomain                     string      `yaml:"subdomain" json:"subdomain" bson:"subdomain" validate:""`
	EnableDownwardApi             bool        `yaml:"enableDownwardApi" json:"enableDownwardApi" bson:"enableDownwardApi" validate:""`
	RestartPolicy                 string      `yaml:"restartPolicy" json:"restartPolicy" bson:"restartPolicy" validate:""`
	Job                           struct {
		Completions             int    `yaml:"completions" json:"completions" bson:"completions" validate:""`
		Parallelism             int    `yaml:"parallelism" json:"parallelism" bson:"parallelism" validate:""`
		CompletionMode          string `yaml:"completionMode" json:"completionMode" bson:"completionMode" validate:""`
		BackoffLimit            int    `yaml:"backoffLimit" json:"backoffLimit" bson:"backoffLimit" validate:""`
		ActiveDeadlineSeconds   int    `yaml:"activeDeadlineSeconds" json:"activeDeadlineSeconds" bson:"activeDeadlineSeconds" validate:""`
		TtlSecondsAfterFinished int    `yaml:"ttlSecondsAfterFinished" json:"ttlSecondsAfterFinished" bson:"ttlSecondsAfterFinished" validate:""`
	} `yaml:"job" json:"job" bson:"job" validate:""`
	CronJob struct {
		Schedule                   string `yaml:"schedule" json:"schedule" bson:"schedule" validate:""`
		ConcurrencyPolicy          string `yaml:"concurrencyPolicy" json:"concurrencyPolicy" bson:"concurrencyPolicy" validate:""`
		StartingDeadlineSeconds    int    `yaml:"startingDeadlineSeconds" json:"startingDeadlineSeconds" bson:"startingDeadlineSeconds" validate:""`
		SuccessfulJobsHistoryLimit int    `yaml:"successfulJobsHistoryLimit" json:"successfulJobsHistoryLimit" bson:"successfulJobsHistoryLimit" validate:""`
		FailedJobsHistoryLimit     int    `yaml:"failedJobsHistoryLimit" json:"failedJobsHistoryLimit" bson:"failedJobsHistoryLimit" validate:""`
	} `yaml:"cronJob" json:"cronJob" bson:"cronJob" validate:""`
	Patches []DeployContainerPatch `yaml:"patches" json:"patches" bson:"patches" validate:""`
	IsPatch bool                   `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type DeployArtifactDef struct {
	DeployArtifactName string            `yaml:"deployArtifactName" json:"deployArtifactName" bson:"deployArtifactName" validate:"required"`
	RelatedArtifact    string            `yaml:"relatedArtifact" json:"relatedArtifact" bson:"relatedArtifact" validate:"required"`
	Hosts              string            `yaml:"hosts" json:"hosts" bson:"hosts" validate:"required"`
	Variables          map[string]string `yaml:"variables" json:"variables" bson:"variables" validate:""`
	Tasks              string            `yaml:"tasks" json:"tasks" bson:"tasks" validate:"required"`
	IsPatch            bool              `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type IstioConsistentHash struct {
	ConsistentHashEnable bool   `yaml:"consistentHashEnable" json:"consistentHashEnable" bson:"consistentHashEnable" validate:""`
	HttpHeaderName       string `yaml:"httpHeaderName" json:"httpHeaderName" bson:"httpHeaderName" validate:""`
	HttpCookie           struct {
		Name string `yaml:"name" json:"name" bson:"name" validate:""`
		Path string `yaml:"path" json:"path" bson:"path" validate:""`
		Ttl  string `yaml:"ttl" json:"ttl" bson:"ttl" validate:""`
	} `yaml:"httpCookie" json:"httpCookie" bson:"httpCookie" validate:""`
	UseSourceIp            bool   `yaml:"useSourceIp" json:"useSourceIp" bson:"useSourceIp" validate:""`
	HttpQueryParameterName string `yaml:"httpQueryParameterName" json:"httpQueryParameterName" bson:"httpQueryParameterName" validate:""`
}

type IstioLoadBalancer struct {
	LoadBalancerEnable bool                `yaml:"loadBalancerEnable" json:"loadBalancerEnable" bson:"loadBalancerEnable" validate:""`
	Simple             string              `yaml:"simple" json:"simple" bson:"simple" validate:""`
	ConsistentHash     IstioConsistentHash `yaml:"consistentHash" json:"consistentHash" bson:"consistentHash" validate:""`
}

type IstioConnectionPoolTcp struct {
	TcpEnable      bool   `yaml:"tcpEnable" json:"tcpEnable" bson:"tcpEnable" validate:""`
	MaxConnections int    `yaml:"maxConnections" json:"maxConnections" bson:"maxConnections" validate:""`
	ConnectTimeout string `yaml:"connectTimeout" json:"connectTimeout" bson:"connectTimeout" validate:""`
}

type IstioConnectionPoolHttp struct {
	HttpEnable               bool   `yaml:"httpEnable" json:"httpEnable" bson:"httpEnable" validate:""`
	Http1MaxPendingRequests  int    `yaml:"http1MaxPendingRequests" json:"http1MaxPendingRequests" bson:"http1MaxPendingRequests" validate:""`
	Http2MaxRequests         int    `yaml:"http2MaxRequests" json:"http2MaxRequests" bson:"http2MaxRequests" validate:""`
	MaxRequestsPerConnection int    `yaml:"maxRequestsPerConnection" json:"maxRequestsPerConnection" bson:"maxRequestsPerConnection" validate:""`
	MaxRetries               int    `yaml:"maxRetries" json:"maxRetries" bson:"maxRetries" validate:""`
	IdleTimeout              string `yaml:"idleTimeout" json:"idleTimeout" bson:"idleTimeout" validate:""`
}

type IstioConnectionPool struct {
	ConnectionPoolEnable bool                    `yaml:"connectionPoolEnable" json:"connectionPoolEnable" bson:"connectionPoolEnable" validate:""`
	Tcp                  IstioConnectionPoolTcp  `yaml:"tcp" json:"tcp" bson:"tcp" validate:""`
	Http                 IstioConnectionPoolHttp `yaml:"http" json:"http" bson:"http" validate:""`
}

type IstioOutlierDetection struct {
	OutlierDetectionEnable   bool   `yaml:"outlierDetectionEnable" json:"outlierDetectionEnable" bson:"outlierDetectionEnable" validate:""`
	ConsecutiveGatewayErrors int    `yaml:"consecutiveGatewayErrors" json:"consecutiveGatewayErrors" bson:"consecutiveGatewayErrors" validate:""`
	Consecutive5xxErrors     int    `yaml:"consecutive5xxErrors" json:"consecutive5xxErrors" bson:"consecutive5xxErrors" validate:""`
	Interval                 string `yaml:"interval" json:"interval" bson:"interval" validate:""`
	BaseEjectionTime         string `yaml:"baseEjectionTime" json:"baseEjectionTime" bson:"baseEjectionTime" validate:""`
	MaxEjectionPercent       int    `yaml:"maxEjectionPercent" json:"maxEjectionPercent" bson:"maxEjectionPercent" validate:""`
	MinHealthPercent         int    `yaml:"minHealthPercent" json:"minHealthPercent" bson:"minHealthPercent" validate:""`
}

type IstioDef struct {
	DeployName   string `yaml:"deployName" json:"deployName" bson:"deployName" validate:"required"`
	Port         int    `yaml:"port" json:"port" bson:"port" validate:"required"`
	Protocol     string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"required,oneof=http tcp"`
	HttpSettings struct {
		MatchHeaders []struct {
			Header     string `yaml:"header" json:"header" bson:"header" validate:"required"`
			MatchType  string `yaml:"matchType" json:"matchType" bson:"matchType" validate:"oneof=exact prefix regex"`
			MatchValue string `yaml:"matchValue" json:"matchValue" bson:"matchValue" validate:"required"`
		} `yaml:"matchHeaders" json:"matchHeaders" bson:"matchHeaders" validate:"dive"`
		Gateway struct {
			RewriteUri string `yaml:"rewriteUri" json:"rewriteUri" bson:"rewriteUri" validate:""`
			MatchUris  []struct {
				MatchType  string `yaml:"matchType" json:"matchType" bson:"matchType" validate:"oneof=exact prefix regex"`
				MatchValue string `yaml:"matchValue" json:"matchValue" bson:"matchValue" validate:"required"`
			} `yaml:"matchUris" json:"matchUris" bson:"matchUris" validate:"dive"`
			MatchDefault bool `yaml:"matchDefault" json:"matchDefault" bson:"matchDefault" validate:""`
		} `yaml:"gateway" json:"gateway" bson:"gateway" validate:""`
		Timeout string `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
		Retries struct {
			RetryOn       string `yaml:"retryOn" json:"retryOn" bson:"retryOn" validate:""`
			Attempts      int    `yaml:"attempts" json:"attempts" bson:"attempts" validate:""`
			PerTryTimeout string `yaml:"perTryTimeout" json:"perTryTimeout" bson:"perTryTimeout" validate:""`
		} `yaml:"retries" json:"retries" bson:"retries" validate:""`
		Mirror struct {
			Host          string `yaml:"host" json:"host" bson:"host" validate:""`
			Port          int    `yaml:"port" json:"port" bson:"port" validate:""`
			Subset        string `yaml:"subset" json:"subset" bson:"subset" validate:""`
			MirrorPercent int    `yaml:"mirrorPercent" json:"mirrorPercent" bson:"mirrorPercent" validate:""`
		} `yaml:"mirror" json:"mirror" bson:"mirror" validate:""`
		CorsPolicy struct {
			AllowOrigins     []map[string]string `yaml:"allowOrigins" json:"allowOrigins" bson:"allowOrigins" validate:""`
			AllowMethods     []string            `yaml:"allowMethods" json:"allowMethods" bson:"allowMethods" validate:""`
			AllowCredentials bool                `yaml:"allowCredentials" json:"allowCredentials" bson:"allowCredentials" validate:""`
			AllowHeaders     []string            `yaml:"allowHeaders" json:"allowHeaders" bson:"allowHeaders" validate:""`
			ExposeHeaders    []string            `yaml:"exposeHeaders" json:"exposeHeaders" bson:"exposeHeaders" validate:""`
			MaxAge           string              `yaml:"maxAge" json:"maxAge" bson:"maxAge" validate:""`
		} `yaml:"corsPolicy" json:"corsPolicy" bson:"corsPolicy" validate:""`
		TrafficPolicyEnable bool                  `yaml:"trafficPolicyEnable" json:"trafficPolicyEnable" bson:"trafficPolicyEnable" validate:""`
		LoadBalancer        IstioLoadBalancer     `yaml:"loadBalancer" json:"loadBalancer" bson:"loadBalancer" validate:""`
		ConnectionPool      IstioConnectionPool   `yaml:"connectionPool" json:"connectionPool" bson:"connectionPool" validate:""`
		OutlierDetection    IstioOutlierDetection `yaml:"outlierDetection" json:"outlierDetection" bson:"outlierDetection" validate:""`
	} `yaml:"httpSettings" json:"httpSettings" bson:"httpSettings" validate:""`
	TcpSettings struct {
		SourceServiceNames []string `yaml:"sourceServiceNames" json:"sourceServiceNames" bson:"sourceServiceNames" validate:""`
	} `yaml:"tcpSettings" json:"tcpSettings" bson:"tcpSettings" validate:""`
	LabelName        string `yaml:"labelName" json:"labelName" bson:"labelName" validate:""`
	LocalLabelConfig struct {
		LabelDefault string `yaml:"labelDefault" json:"labelDefault" bson:"labelDefault" validate:""`
		LabelNew     string `yaml:"labelNew" json:"labelNew" bson:"labelNew" validate:""`
	} `yaml:"localLabelConfig" json:"localLabelConfig" bson:"localLabelConfig" validate:""`
	IsPatch bool `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type IstioGatewayDef struct {
	Enable             bool   `yaml:"enable" json:"enable" bson:"enable" validate:""`
	HostDefault        string `yaml:"hostDefault" json:"hostDefault" bson:"hostDefault" validate:""`
	HostNew            string `yaml:"hostNew" json:"hostNew" bson:"hostNew" validate:""`
	SourceSubsetHeader string `yaml:"sourceSubsetHeader" json:"sourceSubsetHeader" bson:"sourceSubsetHeader" validate:""`
	CertSelfSigned     bool   `yaml:"certSelfSigned" json:"certSelfSigned" bson:"certSelfSigned" validate:""`
	CertPath           string `yaml:"certPath" json:"certPath" bson:"certPath" validate:""`
	WeightDefault      int    `yaml:"weightDefault" json:"weightDefault" bson:"weightDefault" validate:""`
	WeightNew          int    `yaml:"weightNew" json:"weightNew" bson:"weightNew" validate:""`
}

type PipelineBuildDef struct {
	Name string `yaml:"name" json:"name" bson:"name" validate:"required"`
	Run  bool   `yaml:"run" json:"run" bson:"run" validate:""`
}

type GitPullStepDef struct {
	Timeout   int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	SelectTag bool `yaml:"selectTag" json:"selectTag" bson:"selectTag" validate:""`
}

type BuildStepDef struct {
	Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type ScanCodeStepDef struct {
	Enable      bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout     int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry       int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type PackageImageStepDef struct {
	Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type ScanImageStepDef struct {
	Enable       bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError  bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout      int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry        int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
	GateMedium   int  `yaml:"gateMedium" json:"gateMedium" bson:"gateMedium" validate:""`
	GateHigh     int  `yaml:"gateHigh" json:"gateHigh" bson:"gateHigh" validate:""`
	GateCritical int  `yaml:"gateCritical" json:"gateCritical" bson:"gateCritical" validate:""`
}

type ArtifactStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type SyncImageStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type DeployContainerStepDef struct {
	Enable                   bool     `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry                    int      `yaml:"retry" json:"retry" bson:"retry" validate:""`
	ForceReplace             bool     `yaml:"forceReplace" json:"forceReplace" bson:"forceReplace" validate:""`
	Archive                  bool     `yaml:"archive" json:"archive" bson:"archive" validate:""`
	Try                      bool     `yaml:"try" json:"try" bson:"try" validate:""`
	IgnoreExecuteModuleNames []string `yaml:"ignoreExecuteModuleNames" json:"ignoreExecuteModuleNames" bson:"ignoreExecuteModuleNames" validate:""`
}

type ApplyIngressStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type ApplyMeshStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type CheckDeployStepDef struct {
	Enable         bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError    bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Retry          int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
	Repeat         int  `yaml:"repeat" json:"repeat" bson:"repeat" validate:""`
	RepeatInterval int  `yaml:"repeatInterval" json:"repeatInterval" bson:"repeatInterval" validate:""`
}

type CheckQuotaStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type DeployArtifactStepDef struct {
	Enable                   bool     `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Timeout                  int      `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry                    int      `yaml:"retry" json:"retry" bson:"retry" validate:""`
	Archive                  bool     `yaml:"archive" json:"archive" bson:"archive" validate:""`
	Try                      bool     `yaml:"try" json:"try" bson:"try" validate:""`
	IgnoreExecuteModuleNames []string `yaml:"ignoreExecuteModuleNames" json:"ignoreExecuteModuleNames" bson:"ignoreExecuteModuleNames" validate:""`
}

type TestApiStepDef struct {
	Enable                   bool     `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError              bool     `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout                  int      `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry                    int      `yaml:"retry" json:"retry" bson:"retry" validate:""`
	PassingRate              int      `yaml:"passingRate" json:"passingRate" bson:"passingRate" validate:""`
	IgnoreExecuteModuleNames []string `yaml:"ignoreExecuteModuleNames" json:"ignoreExecuteModuleNames" bson:"ignoreExecuteModuleNames" validate:""`
}

type TestPerformanceStepDef struct {
	Enable                   bool     `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError              bool     `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout                  int      `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry                    int      `yaml:"retry" json:"retry" bson:"retry" validate:""`
	IgnoreExecuteModuleNames []string `yaml:"ignoreExecuteModuleNames" json:"ignoreExecuteModuleNames" bson:"ignoreExecuteModuleNames" validate:""`
}

type TestWebuiStepDef struct {
	Enable                   bool     `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError              bool     `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout                  int      `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry                    int      `yaml:"retry" json:"retry" bson:"retry" validate:""`
	PassingRate              int      `yaml:"passingRate" json:"passingRate" bson:"passingRate" validate:""`
	IgnoreExecuteModuleNames []string `yaml:"ignoreExecuteModuleNames" json:"ignoreExecuteModuleNames" bson:"ignoreExecuteModuleNames" validate:""`
}

type InputStepDef struct {
	Enable    bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IsApiOnly bool `yaml:"isApiOnly" json:"isApiOnly" bson:"isApiOnly" validate:""`
}

type PipelineStepDef struct {
	GitPullStepDef         GitPullStepDef         `yaml:"gitPull" json:"gitPull" bson:"gitPull" validate:""`
	BuildStepDef           BuildStepDef           `yaml:"build" json:"build" bson:"build" validate:""`
	ScanCodeStepDef        ScanCodeStepDef        `yaml:"scanCode" json:"scanCode" bson:"scanCode" validate:""`
	PackageImageStepDef    PackageImageStepDef    `yaml:"packageImage" json:"packageImage" bson:"packageImage" validate:""`
	ScanImageStepDef       ScanImageStepDef       `yaml:"scanImage" json:"scanImage" bson:"scanImage" validate:""`
	ArtifactStepDef        ArtifactStepDef        `yaml:"artifact" json:"artifact" bson:"artifact" validate:""`
	SyncImageStepDef       SyncImageStepDef       `yaml:"syncImage" json:"syncImage" bson:"syncImage" validate:""`
	DeployContainerStepDef DeployContainerStepDef `yaml:"deploy" json:"deploy" bson:"deploy" validate:""`
	ApplyIngressStepDef    ApplyIngressStepDef    `yaml:"applyIngress" json:"applyIngress" bson:"applyIngress" validate:""`
	ApplyMeshStepDef       ApplyMeshStepDef       `yaml:"applyMesh" json:"applyMesh" bson:"applyMesh" validate:""`
	CheckDeployStepDef     CheckDeployStepDef     `yaml:"checkDeploy" json:"checkDeploy" bson:"checkDeploy" validate:""`
	CheckQuotaStepDef      CheckQuotaStepDef      `yaml:"checkQuota" json:"checkQuota" bson:"checkQuota" validate:""`
	DeployArtifactStepDef  DeployArtifactStepDef  `yaml:"deployArtifact" json:"deployArtifact" bson:"deployArtifact" validate:""`
	TestApiStepDef         TestApiStepDef         `yaml:"testApi" json:"testApi" bson:"testApi" validate:""`
	TestPerformanceStepDef TestPerformanceStepDef `yaml:"testPerformance" json:"testPerformance" bson:"testPerformance" validate:""`
	TestWebuiStepDef       TestWebuiStepDef       `yaml:"testWebui" json:"testWebui" bson:"testWebui" validate:""`
	InputStepDef           InputStepDef           `yaml:"input" json:"input" bson:"input" validate:""`
}

type CustomStepPhaseDef struct {
	CustomStepName string `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:""`
	Enable         bool   `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError    bool   `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout        int    `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry          int    `yaml:"retry" json:"retry" bson:"retry" validate:""`
	EnableInput    bool   `yaml:"enableInput" json:"enableInput" bson:"enableInput" validate:""`
	IsApiOnly      bool   `yaml:"isApiOnly" json:"isApiOnly" bson:"isApiOnly" validate:""`
}

type PipelineDef struct {
	IsAutoDetectBuild    bool                            `yaml:"isAutoDetectBuild" json:"isAutoDetectBuild" bson:"isAutoDetectBuild" validate:""`
	IsQueue              bool                            `yaml:"isQueue" json:"isQueue" bson:"isQueue" validate:""`
	Builds               []PipelineBuildDef              `yaml:"builds" json:"builds" bson:"builds" validate:"dive"`
	PipelineStep         PipelineStepDef                 `yaml:"pipelineStep" json:"pipelineStep" bson:"pipelineStep" validate:"required"`
	CustomStepInsertDefs map[string][]CustomStepPhaseDef `yaml:"customStepInsertDefs" json:"customStepInsertDefs" bson:"customStepInsertDefs" validate:""`
}

type ProjectDef struct {
	BuildDefs              []BuildDef        `yaml:"buildDefs" json:"buildDefs" bson:"buildDefs" validate:""`
	UpdateBuildDefs        bool              `yaml:"updateBuildDefs" json:"updateBuildDefs" bson:"updateBuildDefs" validate:""`
	PackageDefs            []PackageDef      `yaml:"packageDefs" json:"packageDefs" bson:"packageDefs" validate:""`
	UpdatePackageDefs      bool              `yaml:"updatePackageDefs" json:"updatePackageDefs" bson:"updatePackageDefs" validate:""`
	ArtifactDefs           []ArtifactDef     `yaml:"artifactDefs" json:"artifactDefs" bson:"artifactDefs" validate:""`
	UpdateArtifactDefs     bool              `yaml:"updateArtifactDefs" json:"updateArtifactDefs" bson:"updateArtifactDefs" validate:""`
	DockerIgnoreDefs       []string          `yaml:"dockerIgnoreDefs" json:"dockerIgnoreDefs" bson:"dockerIgnoreDefs" validate:""`
	UpdateDockerIgnoreDefs bool              `yaml:"updateDockerIgnoreDefs" json:"updateDockerIgnoreDefs" bson:"updateDockerIgnoreDefs" validate:""`
	CustomStepDefs         CustomStepDefs    `yaml:"customStepDefs" json:"customStepDefs" bson:"customStepDefs" validate:""`
	CustomOpsDefs          []CustomOpsDef    `yaml:"customOpsDefs" json:"customOpsDefs" bson:"customOpsDefs" validate:""`
	UpdateCustomOpsDefs    bool              `yaml:"updateCustomOpsDefs" json:"updateCustomOpsDefs" bson:"updateCustomOpsDefs" validate:""`
	OpsBatchDefs           []OpsBatchDef     `yaml:"opsBatchDefs" json:"opsBatchDefs" bson:"opsBatchDefs" validate:""`
	UpdateOpsBatchDefs     bool              `yaml:"updateOpsBatchDefs" json:"updateOpsBatchDefs" bson:"updateOpsBatchDefs" validate:""`
	ErrMsgPackageDefs      string            `yaml:"errMsgPackageDefs" json:"errMsgPackageDefs" bson:"errMsgPackageDefs" validate:""`
	ErrMsgArtifactDefs     string            `yaml:"errMsgArtifactDefs" json:"errMsgArtifactDefs" bson:"errMsgArtifactDefs" validate:""`
	ErrMsgCustomStepDefs   map[string]string `yaml:"errMsgCustomStepDefs" json:"errMsgCustomStepDefs" bson:"errMsgCustomStepDefs" validate:""`
	ErrMsgCustomOpsDefs    string            `yaml:"errMsgCustomOpsDefs" json:"errMsgCustomOpsDefs" bson:"errMsgCustomOpsDefs" validate:""`
	ErrMsgOpsBatchDefs     string            `yaml:"errMsgOpsBatchDefs" json:"errMsgOpsBatchDefs" bson:"errMsgOpsBatchDefs" validate:""`
}

type ProjectInfo struct {
	ProjectName      string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	ProjectNamespace string `yaml:"projectNamespace" json:"projectNamespace" bson:"projectNamespace" validate:""`
	ProjectShortName string `yaml:"projectShortName" json:"projectShortName" bson:"projectShortName" validate:""`
	ShortName        string `yaml:"shortName" json:"shortName" bson:"shortName" validate:""`
	ProjectArch      string `yaml:"projectArch" json:"projectArch" bson:"projectArch" validate:""`
	Privileged       bool   `yaml:"privileged" json:"privileged" bson:"privileged" validate:""`
	DefaultPv        string `yaml:"defaultPv" json:"defaultPv" bson:"defaultPv" validate:""`
	ProjectDesc      string `yaml:"projectDesc" json:"projectDesc" bson:"projectDesc" validate:""`
	ProjectTeam      string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:""`
}

type ModuleRun struct {
	ModuleName   string `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:""`
	ModuleEnable bool   `yaml:"moduleEnable" json:"moduleEnable" bson:"moduleEnable" validate:""`
}

type ProjectPipeline struct {
	BranchName        string      `yaml:"branchName" json:"branchName" bson:"branchName" validate:"required"`
	IsDefault         bool        `yaml:"isDefault" json:"isDefault" bson:"isDefault" validate:""`
	WebhookPushEvent  bool        `yaml:"webhookPushEvent" json:"webhookPushEvent" bson:"webhookPushEvent" validate:""`
	Envs              []string    `yaml:"envs" json:"envs" bson:"envs" validate:""`
	EnvProductions    []string    `yaml:"envProductions" json:"envProductions" bson:"envProductions" validate:""`
	PipelineDef       PipelineDef `yaml:"pipelineDef" json:"pipelineDef" bson:"pipelineDef" validate:""`
	UpdatePipelineDef bool        `yaml:"updatePipelineDef" json:"updatePipelineDef" bson:"updatePipelineDef" validate:""`
	ErrMsgPipelineDef string      `yaml:"errMsgPipelineDef" json:"errMsgPipelineDef" bson:"errMsgPipelineDef" validate:""`
}

type CustomStepConfOutput struct {
	CustomStepName       string `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:""`
	CustomStepActionDesc string `yaml:"customStepActionDesc" json:"customStepActionDesc" bson:"customStepActionDesc" validate:""`
	CustomStepDesc       string `yaml:"customStepDesc" json:"customStepDesc" bson:"customStepDesc" validate:""`
	CustomStepUsage      string `yaml:"customStepUsage" json:"customStepUsage" bson:"customStepUsage" validate:""`
	IsEnvDiff            bool   `yaml:"isEnvDiff" json:"isEnvDiff" bson:"isEnvDiff" validate:""`
	ParamInputYamlDef    string `yaml:"paramInputYamlDef" json:"paramInputYamlDef" bson:"paramInputYamlDef" validate:""`
	ParamOutputYamlDef   string `yaml:"paramOutputYamlDef" json:"paramOutputYamlDef" bson:"paramOutputYamlDef" validate:""`
}

type ProjectOutput struct {
	ProjectInfo          ProjectInfo           `yaml:"projectInfo" json:"projectInfo" bson:"projectInfo" validate:""`
	ProjectPipelines     []ProjectPipeline     `yaml:"pipelines" json:"pipelines" bson:"pipelines" validate:""`
	ProjectAvailableEnvs []ProjectAvailableEnv `yaml:"projectAvailableEnvs" json:"projectAvailableEnvs" bson:"projectAvailableEnvs" validate:""`
	ProjectNodePorts     []struct {
		EnvName      string        `yaml:"envName" json:"envName" bson:"envName" validate:""`
		EnvNodePorts []EnvNodePort `yaml:"envNodePorts" json:"envNodePorts" bson:"envNodePorts" validate:""`
		NodePorts    []int         `yaml:"nodePorts" json:"nodePorts" bson:"nodePorts" validate:""`
	} `yaml:"projectNodePorts" json:"projectNodePorts" bson:"projectNodePorts" validate:""`
	ProjectDef      ProjectDef             `yaml:"projectDef" json:"projectDef" bson:"projectDef" validate:""`
	BuildEnvs       []string               `yaml:"buildEnvs" json:"buildEnvs" bson:"buildEnvs" validate:""`
	BuildNames      []string               `yaml:"buildNames" json:"buildNames" bson:"buildNames" validate:""`
	PackageNames    []string               `yaml:"packageNames" json:"packageNames" bson:"packageNames" validate:""`
	ArtifactNames   []string               `yaml:"artifactNames" json:"artifactNames" bson:"artifactNames" validate:""`
	CustomStepConfs []CustomStepConfOutput `yaml:"customStepConfs" json:"customStepConfs" bson:"customStepConfs" validate:""`
}

type User struct {
	TenantCode   string   `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	TenantAdmins []string `yaml:"tenantAdmins" json:"tenantAdmins" bson:"tenantAdmins" validate:""`
	Username     string   `yaml:"username" json:"username" bson:"username" validate:""`
	Name         string   `yaml:"name" json:"name" bson:"name" validate:""`
	Mail         string   `yaml:"mail" json:"mail" bson:"mail" validate:""`
	Mobile       string   `yaml:"mobile" json:"mobile" bson:"mobile" validate:""`
	IsAdmin      bool     `yaml:"isAdmin" json:"isAdmin" bson:"isAdmin" validate:""`
	IsActive     bool     `yaml:"isActive" json:"isActive" bson:"isActive" validate:""`
}

type UserProject struct {
	ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	AccessLevel string `yaml:"accessLevel" json:"accessLevel" bson:"accessLevel" validate:""`
	UpdateTime  string `yaml:"updateTime" json:"updateTime" bson:"updateTime" validate:""`
}

type UserDetail struct {
	TenantCode   string        `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	TenantAdmins []string      `yaml:"tenantAdmins" json:"tenantAdmins" bson:"tenantAdmins" validate:""`
	Username     string        `yaml:"username" json:"username" bson:"username" validate:""`
	Name         string        `yaml:"name" json:"name" bson:"name" validate:""`
	Mail         string        `yaml:"mail" json:"mail" bson:"mail" validate:""`
	Mobile       string        `yaml:"mobile" json:"mobile" bson:"mobile" validate:""`
	IsAdmin      bool          `yaml:"isAdmin" json:"isAdmin" bson:"isAdmin" validate:""`
	IsActive     bool          `yaml:"isActive" json:"isActive" bson:"isActive" validate:""`
	AvatarUrl    string        `yaml:"avatarUrl" json:"avatarUrl" bson:"avatarUrl" validate:""`
	UserProjects []UserProject `yaml:"projects" json:"projects" bson:"projects" validate:""`
	CreateTime   string        `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
	LastLogin    string        `yaml:"lastLogin" json:"lastLogin" bson:"lastLogin" validate:""`
}

type CustomStepConfDetail struct {
	CustomStepConf
	ProjectNames []string `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
}

type CustomStepDockerConf struct {
	DockerImage        string   `yaml:"dockerImage" json:"dockerImage" bson:"dockerImage" validate:"required"`
	RegistryUsername   string   `yaml:"registryUsername" json:"registryUsername" bson:"registryUsername" validate:""`
	RegistryPassword   string   `yaml:"registryPassword" json:"registryPassword" bson:"registryPassword" validate:""`
	DockerCommands     []string `yaml:"dockerCommands" json:"dockerCommands" bson:"dockerCommands" validate:"required"`
	DockerShowCommands bool     `yaml:"dockerShowCommands" json:"dockerShowCommands" bson:"dockerShowCommands" validate:""`
	DockerRunAsRoot    bool     `yaml:"dockerRunAsRoot" json:"dockerRunAsRoot" bson:"dockerRunAsRoot" validate:""`
	DockerVolumes      []string `yaml:"dockerVolumes" json:"dockerVolumes" bson:"dockerVolumes" validate:""`
	DockerEnvs         []string `yaml:"dockerEnvs" json:"dockerEnvs" bson:"dockerEnvs" validate:""`
	DockerWorkDir      string   `yaml:"dockerWorkDir" json:"dockerWorkDir" bson:"dockerWorkDir" validate:""`
	ParamInputFormat   string   `yaml:"paramInputFormat" json:"paramInputFormat" bson:"paramInputFormat" validate:"required"`
	ParamOutputFormat  string   `yaml:"paramOutputFormat" json:"paramOutputFormat" bson:"paramOutputFormat" validate:"required"`
}

type CustomStepConf struct {
	TenantCode           string               `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	CustomStepName       string               `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:"required"`
	CustomStepActionDesc string               `yaml:"customStepActionDesc" json:"customStepActionDesc" bson:"customStepActionDesc" validate:"required"`
	CustomStepDesc       string               `yaml:"customStepDesc" json:"customStepDesc" bson:"customStepDesc" validate:"required"`
	CustomStepUsage      string               `yaml:"customStepUsage" json:"customStepUsage" bson:"customStepUsage" validate:"required"`
	GitRepoName          string               `yaml:"gitRepoName" json:"gitRepoName" bson:"gitRepoName" validate:""`
	GitRepoPath          string               `yaml:"gitRepoPath" json:"gitRepoPath" bson:"gitRepoPath" validate:""`
	GitRepoBranch        string               `yaml:"gitRepoBranch" json:"gitRepoBranch" bson:"gitRepoBranch" validate:""`
	CustomStepDockerConf CustomStepDockerConf `yaml:"customStepDockerConf" json:"customStepDockerConf" bson:"customStepDockerConf" validate:"required"`
	ParamInputYamlDef    string               `yaml:"paramInputYamlDef" json:"paramInputYamlDef" bson:"paramInputYamlDef" validate:""`
	ParamOutputYamlDef   string               `yaml:"paramOutputYamlDef" json:"paramOutputYamlDef" bson:"paramOutputYamlDef" validate:""`
	IsEnvDiff            bool                 `yaml:"isEnvDiff" json:"isEnvDiff" bson:"isEnvDiff" validate:""`
}

type QuotaConfig struct {
	DefaultQuota   QuotaPod      `yaml:"defaultQuota" json:"defaultQuota" bson:"defaultQuota" validate:"required"`
	NamespaceQuota QuotaResource `yaml:"namespaceQuota" json:"namespaceQuota" bson:"namespaceQuota" validate:"required"`
}

type ResourceVersion struct {
	IngressVersion string `yaml:"ingressVersion" json:"ingressVersion" bson:"ingressVersion" validate:""`
	HpaVersion     string `yaml:"hpaVersion" json:"hpaVersion" bson:"hpaVersion" validate:""`
	IstioVersion   string `yaml:"istioVersion" json:"istioVersion" bson:"istioVersion" validate:""`
}

type KubeNode struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion" bson:"apiVersion" validate:"required"`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	MetaData   struct {
		Name        string            `yaml:"name" json:"name" bson:"name" validate:"required"`
		Labels      map[string]string `yaml:"labels" json:"labels" bson:"labels" validate:""`
		Annotations map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
	} `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Spec struct {
		PodCIDR       string   `yaml:"podCIDR" json:"podCIDR" bson:"podCIDR" validate:""`
		PodCIDRs      []string `yaml:"podCIDRs" json:"podCIDRs" bson:"podCIDRs" validate:""`
		Unschedulable bool     `yaml:"unschedulable" json:"unschedulable" bson:"unschedulable" validate:""`
	} `yaml:"spec" json:"spec" bson:"spec" validate:""`
	Status struct {
		Addresses []struct {
			Address string `yaml:"address" json:"address" bson:"address" validate:""`
			Type    string `yaml:"type" json:"type" bson:"type" validate:""`
		} `yaml:"addresses" json:"addresses" bson:"addresses" validate:""`
		Capacity struct {
			Cpu    string `yaml:"cpu" json:"cpu" bson:"cpu" validate:""`
			Memory string `yaml:"memory" json:"memory" bson:"memory" validate:""`
			Pods   string `yaml:"pods" json:"pods" bson:"pods" validate:""`
		} `yaml:"capacity" json:"capacity" bson:"capacity" validate:""`
		NodeInfo struct {
			Architecture            string `yaml:"architecture" json:"architecture" bson:"architecture" validate:""`
			ContainerRuntimeVersion string `yaml:"containerRuntimeVersion" json:"containerRuntimeVersion" bson:"containerRuntimeVersion" validate:""`
			KernelVersion           string `yaml:"kernelVersion" json:"kernelVersion" bson:"kernelVersion" validate:""`
			KubeProxyVersion        string `yaml:"kubeProxyVersion" json:"kubeProxyVersion" bson:"kubeProxyVersion" validate:""`
			KubeletVersion          string `yaml:"kubeletVersion" json:"kubeletVersion" bson:"kubeletVersion" validate:""`
			OperatingSystem         string `yaml:"operatingSystem" json:"operatingSystem" bson:"operatingSystem" validate:""`
			OsImage                 string `yaml:"osImage" json:"osImage" bson:"osImage" validate:""`
		} `yaml:"nodeInfo" json:"nodeInfo" bson:"nodeInfo" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type ArchSetting struct {
	Arch         string            `yaml:"arch" json:"arch" bson:"arch" validate:"required"`
	NodeSelector map[string]string `yaml:"nodeSelector" json:"nodeSelector" bson:"nodeSelector" validate:"required"`
}

type EnvK8s struct {
	EnvName                    string `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
	EnvDesc                    string `yaml:"envDesc" json:"envDesc" bson:"envDesc" validate:"required"`
	TenantCode                 string `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile                 bool   `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	Host                       string `yaml:"host" json:"host" bson:"host" validate:"required"`
	ViewHost                   string `yaml:"viewHost" json:"viewHost" bson:"viewHost" validate:"required"`
	Port                       int    `yaml:"port" json:"port" bson:"port" validate:"required"`
	Token                      string `yaml:"token" json:"token" bson:"token" validate:"required"`
	CaCrtBase64                string `yaml:"caCrtBase64" json:"caCrtBase64" bson:"caCrtBase64" validate:"required"`
	DashboardUrl               string `yaml:"dashboardUrl" json:"dashboardUrl" bson:"dashboardUrl" validate:"required"`
	DashboardUrlNetworkPolicy  string `yaml:"dashboardUrlNetworkPolicy" json:"dashboardUrlNetworkPolicy" bson:"dashboardUrlNetworkPolicy" validate:"required"`
	IstioNamespace             string `yaml:"istioNamespace" json:"istioNamespace" bson:"istioNamespace" validate:"required"`
	IngressControllerNamespace string `yaml:"ingressControllerNamespace" json:"ingressControllerNamespace" bson:"ingressControllerNamespace" validate:"required"`
	Timezone                   string `yaml:"timezone" json:"timezone" bson:"timezone" validate:"required"`
	NodePortRange              struct {
		NodePortRangeStart int `yaml:"nodePortRangeStart" json:"nodePortRangeStart" bson:"nodePortRangeStart" validate:"required"`
		NodePortRangeEnd   int `yaml:"nodePortRangeEnd" json:"nodePortRangeEnd" bson:"nodePortRangeEnd" validate:"required"`
	} `yaml:"nodePortRange" json:"nodePortRange" bson:"nodePortRange" validate:"required"`
	ArchSettings   []ArchSetting `yaml:"archSettings" json:"archSettings" bson:"archSettings" validate:""`
	ProjectDataPod struct {
		Namespace       string `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required"`
		StatefulSetName string `yaml:"statefulSetName" json:"statefulSetName" bson:"statefulSetName" validate:"required"`
		Path            string `yaml:"path" json:"path" bson:"path" validate:"required"`
	} `yaml:"projectDataPod" json:"projectDataPod" bson:"projectDataPod" validate:""`
	PodImageSettings struct {
		ProjectDataPodImage string `yaml:"projectDataPodImage" json:"projectDataPodImage" bson:"projectDataPodImage" validate:"required"`
		BusyboxImage        string `yaml:"busyboxImage" json:"busyboxImage" bson:"busyboxImage" validate:"required"`
	} `yaml:"podImageSettings" json:"podImageSettings" bson:"podImageSettings" validate:"required"`
	ImageRepoUseExternal    bool `yaml:"imageRepoUseExternal" json:"imageRepoUseExternal" bson:"imageRepoUseExternal" validate:""`
	ArtifactRepoUseExternal bool `yaml:"artifactRepoUseExternal" json:"artifactRepoUseExternal" bson:"artifactRepoUseExternal" validate:""`
	PvConfigLocal           struct {
		LocalPath string `yaml:"localPath" json:"localPath" bson:"localPath" validate:""`
	} `yaml:"pvConfigLocal" json:"pvConfigLocal" bson:"pvConfigLocal" validate:""`
	PvConfigCephfs struct {
		CephPath     string   `yaml:"cephPath" json:"cephPath" bson:"cephPath" validate:""`
		CephUser     string   `yaml:"cephUser" json:"cephUser" bson:"cephUser" validate:""`
		CephSecret   string   `yaml:"cephSecret" json:"cephSecret" bson:"cephSecret" validate:""`
		CephMonitors []string `yaml:"cephMonitors" json:"cephMonitors" bson:"cephMonitors" validate:""`
	} `yaml:"pvConfigCephfs" json:"pvConfigCephfs" bson:"pvConfigCephfs" validate:""`
	PvConfigNfs struct {
		NfsPath   string `yaml:"nfsPath" json:"nfsPath" bson:"nfsPath" validate:""`
		NfsServer string `yaml:"nfsServer" json:"nfsServer" bson:"nfsServer" validate:""`
	} `yaml:"pvConfigNfs" json:"pvConfigNfs" bson:"pvConfigNfs" validate:""`
	ProjectNodeSelector map[string]string `yaml:"projectNodeSelector" json:"projectNodeSelector" bson:"projectNodeSelector" validate:""`
	QuotaConfig         QuotaConfig       `yaml:"quotaConfig" json:"quotaConfig" bson:"quotaConfig" validate:"required"`
}

type EnvK8sDetail struct {
	EnvK8s
	ResourceVersion ResourceVersion `yaml:"resourceVersion" json:"resourceVersion" bson:"resourceVersion" validate:""`
	Nodes           []KubeNode      `yaml:"nodes" json:"nodes" bson:"nodes" validate:""`
	Arches          []string        `yaml:"arches" json:"arches" bson:"arches" validate:""`
}

type DockerBuildEnv struct {
	TenantCode          string   `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile          bool     `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	BuildEnvName        string   `yaml:"buildEnvName" json:"buildEnvName" bson:"buildEnvName" validate:"required"`
	Image               string   `yaml:"image" json:"image" bson:"image" validate:"required"`
	RegistryUsername    string   `yaml:"registryUsername" json:"registryUsername" bson:"registryUsername" validate:""`
	RegistryPassword    string   `yaml:"registryPassword" json:"registryPassword" bson:"registryPassword" validate:""`
	BuildArches         []string `yaml:"buildArches" json:"buildArches" bson:"buildArches" validate:""`
	MountHomeDir        bool     `yaml:"mountHomeDir" json:"mountHomeDir" bson:"mountHomeDir" validate:""`
	EnableProxy         bool     `yaml:"enableProxy" json:"enableProxy" bson:"enableProxy" validate:""`
	MountExtraCacheDirs []string `yaml:"mountExtraCacheDirs" json:"mountExtraCacheDirs" bson:"mountExtraCacheDirs" validate:""`
	CommandsBeforeBuild []string `yaml:"commandsBeforeBuild" json:"commandsBeforeBuild" bson:"commandsBeforeBuild" validate:""`
	CommandsAfterCheck  []string `yaml:"commandsAfterCheck" json:"commandsAfterCheck" bson:"commandsAfterCheck" validate:""`
}

type GitRepoConfig struct {
	TenantCode string `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile bool   `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	RepoName   string `yaml:"repoName" json:"repoName" bson:"repoName" validate:"required"`
	ViewUrl    string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required"`
	Url        string `yaml:"url" json:"url" bson:"url" validate:"required"`
	Insecure   bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Username   string `yaml:"username" json:"username" bson:"username" validate:"required"`
	Name       string `yaml:"name" json:"name" bson:"name" validate:"required"`
	Mail       string `yaml:"mail" json:"mail" bson:"mail" validate:"required"`
	Password   string `yaml:"password" json:"password" bson:"password" validate:"required"`
	Token      string `yaml:"token" json:"token" bson:"token" validate:"required"`
}

type ImageRepoConfig struct {
	TenantCode string `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile bool   `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	RepoName   string `yaml:"repoName" json:"repoName" bson:"repoName" validate:"required"`
	Hostname   string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required"`
	Insecure   bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Username   string `yaml:"username" json:"username" bson:"username" validate:"required"`
	Password   string `yaml:"password" json:"password" bson:"password" validate:"required"`
	IpInternal string `yaml:"ipInternal" json:"ipInternal" bson:"ipInternal" validate:"required"`
	IpExternal string `yaml:"ipExternal" json:"ipExternal" bson:"ipExternal" validate:"required"`
}

type ArtifactRepoConfig struct {
	TenantCode       string `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile       bool   `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	Kind             string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	RepoName         string `yaml:"repoName" json:"repoName" bson:"repoName" validate:"required"`
	ViewUrl          string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required"`
	SchemaInternal   string `yaml:"schemaInternal" json:"schemaInternal" bson:"schemaInternal" validate:"required"`
	SchemaExternal   string `yaml:"schemaExternal" json:"schemaExternal" bson:"schemaExternal" validate:"required"`
	Schema           string `yaml:"schema" json:"schema" bson:"schema" validate:""`
	HostnameInternal string `yaml:"hostnameInternal" json:"hostnameInternal" bson:"hostnameInternal" validate:"required"`
	HostnameExternal string `yaml:"hostnameExternal" json:"hostnameExternal" bson:"hostnameExternal" validate:"required"`
	Hostname         string `yaml:"hostname" json:"hostname" bson:"hostname" validate:""`
	Port             int    `yaml:"port" json:"port" bson:"port" validate:"required"`
	Url              string `yaml:"url" json:"url" bson:"url" validate:""`
	Insecure         bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Username         string `yaml:"username" json:"username" bson:"username" validate:"required"`
	Password         string `yaml:"password" json:"password" bson:"password" validate:"required"`
	ProxyRepo        struct {
		PublicRole     string `yaml:"publicRole" json:"publicRole" bson:"publicRole" validate:""`
		PublicUser     string `yaml:"publicUser" json:"publicUser" bson:"publicUser" validate:""`
		PublicPassword string `yaml:"publicPassword" json:"publicPassword" bson:"publicPassword" validate:""`
		PortDocker     int    `yaml:"portDocker" json:"portDocker" bson:"portDocker" validate:""`
		PortGcr        int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:""`
		PortQuay       int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:""`
		Maven          string `yaml:"maven" json:"maven" bson:"maven" validate:""`
		Npm            string `yaml:"npm" json:"npm" bson:"npm" validate:""`
		Pip            string `yaml:"pip" json:"pip" bson:"pip" validate:""`
		Gradle         string `yaml:"gradle" json:"gradle" bson:"gradle" validate:""`
		Go             string `yaml:"go" json:"go" bson:"go" validate:""`
		Apt            struct {
			Amd64   string `yaml:"amd64" json:"amd64" bson:"amd64" validate:""`
			Arm64v8 string `yaml:"arm64v8" json:"arm64v8" bson:"arm64v8" validate:""`
		} `yaml:"apt" json:"apt" bson:"apt" validate:""`
	} `yaml:"proxyRepo" json:"proxyRepo" bson:"proxyRepo" validate:""`
}

type ScanCodeRepoConfig struct {
	TenantCode string `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	IsFromFile bool   `yaml:"isFromFile" json:"isFromFile" bson:"isFromFile" validate:""`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	RepoName   string `yaml:"repoName" json:"repoName" bson:"repoName" validate:"required"`
	ViewUrl    string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required"`
	Url        string `yaml:"url" json:"url" bson:"url" validate:"required"`
	Insecure   bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Token      string `yaml:"token" json:"token" bson:"token" validate:"required"`
}

type DeploySpecStatic struct {
	DeployName          string `yaml:"deployName" json:"deployName" bson:"deployName" validate:""`
	DeployImage         string `yaml:"deployImage" json:"deployImage" bson:"deployImage" validate:"required"`
	DeployType          string `yaml:"deployType" json:"deployType" bson:"deployType" validate:""`
	DeployHeadless      bool   `yaml:"deployHeadless" json:"deployHeadless" bson:"deployHeadless" validate:""`
	PodManagementPolicy string `yaml:"podManagementPolicy" json:"podManagementPolicy" bson:"podManagementPolicy" validate:""`
	DeployMeta          struct {
		Labels      []NameValue `yaml:"labels" json:"labels" bson:"labels" validate:""`
		Annotations []NameValue `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
	} `yaml:"deployMeta" json:"deployMeta" bson:"deployMeta" validate:""`
	DeploySessionAffinityTimeoutSeconds int `yaml:"deploySessionAffinityTimeoutSeconds" json:"deploySessionAffinityTimeoutSeconds" bson:"deploySessionAffinityTimeoutSeconds" validate:""`
	DeployNodePorts                     []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		NodePort int    `yaml:"nodePort" json:"nodePort" bson:"nodePort" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=HTTP TCP UDP SCTP"`
	} `yaml:"deployNodePorts" json:"deployNodePorts" bson:"deployNodePorts" validate:"dive"`
	DeployLocalPorts []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=HTTP TCP UDP SCTP"`
		Ingress  struct {
			DomainName     string `yaml:"domainName" json:"domainName" bson:"domainName" validate:""`
			PathPrefix     string `yaml:"pathPrefix" json:"pathPrefix" bson:"pathPrefix" validate:""`
			CertSelfSigned bool   `yaml:"certSelfSigned" json:"certSelfSigned" bson:"certSelfSigned" validate:""`
			CertBranch     string `yaml:"certBranch" json:"certBranch" bson:"certBranch" validate:""`
			CertPath       string `yaml:"certPath" json:"certPath" bson:"certPath" validate:""`
		} `yaml:"ingress" json:"ingress" bson:"ingress" validate:""`
	} `yaml:"deployLocalPorts" json:"deployLocalPorts" bson:"deployLocalPorts" validate:"dive"`
	DeployReplicas int `yaml:"deployReplicas" json:"deployReplicas" bson:"deployReplicas" validate:"required"`
	HpaConfig      struct {
		MaxReplicas                 int    `yaml:"maxReplicas" json:"maxReplicas" bson:"maxReplicas" validate:""`
		MemoryAverageValue          string `yaml:"memoryAverageValue" json:"memoryAverageValue" bson:"memoryAverageValue" validate:""`
		MemoryAverageRequestPercent int    `yaml:"memoryAverageRequestPercent" json:"memoryAverageRequestPercent" bson:"memoryAverageRequestPercent" validate:""`
		CpuAverageValue             string `yaml:"cpuAverageValue" json:"cpuAverageValue" bson:"cpuAverageValue" validate:""`
		CpuAverageRequestPercent    int    `yaml:"cpuAverageRequestPercent" json:"cpuAverageRequestPercent" bson:"cpuAverageRequestPercent" validate:""`
	} `yaml:"hpaConfig" json:"hpaConfig" bson:"hpaConfig" validate:""`
	DeployEnvs      []string `yaml:"deployEnvs" json:"deployEnvs" bson:"deployEnvs" validate:""`
	DeployCommand   string   `yaml:"deployCommand" json:"deployCommand" bson:"deployCommand" validate:""`
	DeployCmds      []string `yaml:"deployCmds" json:"deployCmds" bson:"deployCmds" validate:""`
	DeployArgs      []string `yaml:"deployArgs" json:"deployArgs" bson:"deployArgs" validate:""`
	DeployResources QuotaPod `yaml:"deployResources" json:"deployResources" bson:"deployResources" validate:""`
	DeployVolumes   []struct {
		PathInPod string `yaml:"pathInPod" json:"pathInPod" bson:"pathInPod" validate:"required"`
		PathInPv  string `yaml:"pathInPv" json:"pathInPv" bson:"pathInPv" validate:"required"`
		Pvc       string `yaml:"pvc" json:"pvc" bson:"pvc" validate:""`
	} `yaml:"deployVolumes" json:"deployVolumes" bson:"deployVolumes" validate:"dive"`
	DeployHealthCheck struct {
		CheckPort              int           `yaml:"checkPort" json:"checkPort" bson:"checkPort" validate:""`
		Exec                   string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
		ExecCmds               []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
		HttpGet                DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		ReadinessDelaySeconds  int           `yaml:"readinessDelaySeconds" json:"readinessDelaySeconds" bson:"readinessDelaySeconds" validate:""`
		ReadinessPeriodSeconds int           `yaml:"readinessPeriodSeconds" json:"readinessPeriodSeconds" bson:"readinessPeriodSeconds" validate:""`
		LivenessDelaySeconds   int           `yaml:"livenessDelaySeconds" json:"livenessDelaySeconds" bson:"livenessDelaySeconds" validate:""`
		LivenessPeriodSeconds  int           `yaml:"livenessPeriodSeconds" json:"livenessPeriodSeconds" bson:"livenessPeriodSeconds" validate:""`
		StartupDelaySeconds    int           `yaml:"startupDelaySeconds" json:"startupDelaySeconds" bson:"startupDelaySeconds" validate:""`
		StartupPeriodSeconds   int           `yaml:"startupPeriodSeconds" json:"startupPeriodSeconds" bson:"startupPeriodSeconds" validate:""`
	} `yaml:"deployHealthCheck" json:"deployHealthCheck" bson:"deployHealthCheck" validate:""`
	DependServices []struct {
		DependName string `yaml:"dependName" json:"dependName" bson:"dependName" validate:"required"`
		DependPort int    `yaml:"dependPort" json:"dependPort" bson:"dependPort" validate:"required"`
		DependType string `yaml:"dependType" json:"dependType" bson:"dependType" validate:"oneof=TCP UDP"`
	} `yaml:"dependServices" json:"dependServices" bson:"dependServices" validate:"dive"`
	HostAliases []struct {
		Ip        string   `yaml:"ip" json:"ip" bson:"ip" validate:"required,ip"`
		Hostnames []string `yaml:"hostnames" json:"hostnames" bson:"hostnames" validate:"required"`
	} `yaml:"hostAliases" json:"hostAliases" bson:"hostAliases" validate:"dive"`
	SecurityContext struct {
		RunAsUser  int `yaml:"runAsUser" json:"runAsUser" bson:"runAsUser" validate:""`
		RunAsGroup int `yaml:"runAsGroup" json:"runAsGroup" bson:"runAsGroup" validate:""`
	} `yaml:"securityContext" json:"securityContext" bson:"securityContext" validate:""`
	DeployConfigBranch   string       `yaml:"deployConfigBranch" json:"deployConfigBranch" bson:"deployConfigBranch" validate:""`
	DeployConfigSettings []ConfigPath `yaml:"deployConfigSettings" json:"deployConfigSettings" bson:"deployConfigSettings" validate:""`
	Lifecycle            struct {
		PostStart struct {
			Exec     string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
			ExecCmds []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
			HttpGet  DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		} `yaml:"postStart" json:"postStart" bson:"postStart" validate:""`
		PreStop struct {
			Exec     string        `yaml:"exec" json:"exec" bson:"exec" validate:""`
			ExecCmds []string      `yaml:"execCmds" json:"execCmds" bson:"execCmds" validate:""`
			HttpGet  DeployHttpGet `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		} `yaml:"preStop" json:"preStop" bson:"preStop" validate:""`
	} `yaml:"lifecycle" json:"lifecycle" bson:"lifecycle" validate:""`
	WorkingDir                    string                 `yaml:"workingDir" json:"workingDir" bson:"workingDir" validate:""`
	NodeSelector                  []NameValue            `yaml:"nodeSelector" json:"nodeSelector" bson:"nodeSelector" validate:""`
	NodeName                      string                 `yaml:"nodeName" json:"nodeName" bson:"nodeName" validate:""`
	TerminationGracePeriodSeconds int                    `yaml:"terminationGracePeriodSeconds" json:"terminationGracePeriodSeconds" bson:"terminationGracePeriodSeconds" validate:""`
	Subdomain                     string                 `yaml:"subdomain" json:"subdomain" bson:"subdomain" validate:""`
	EnableDownwardApi             bool                   `yaml:"enableDownwardApi" json:"enableDownwardApi" bson:"enableDownwardApi" validate:""`
	Patches                       []DeployContainerPatch `yaml:"patches" json:"patches" bson:"patches" validate:""`
}

type ComponentTemplate struct {
	TenantCode            string           `yaml:"tenantCode" json:"tenantCode" bson:"tenantCode" validate:""`
	ComponentTemplateName string           `yaml:"componentTemplateName" json:"componentTemplateName" bson:"componentTemplateName" validate:"required"`
	ComponentTemplateDesc string           `yaml:"componentTemplateDesc" json:"componentTemplateDesc" bson:"componentTemplateDesc" validate:"required"`
	DeploySpecStatic      DeploySpecStatic `yaml:"deploySpecStatic" json:"deploySpecStatic" bson:"deploySpecStatic" validate:"required"`
}

type ProjectSummary struct {
	BuildEnvs       []string               `yaml:"buildEnvs" json:"buildEnvs" bson:"buildEnvs" validate:""`
	BuildNames      []string               `yaml:"buildNames" json:"buildNames" bson:"buildNames" validate:""`
	CustomStepConfs []CustomStepConfOutput `yaml:"customStepConfs" json:"customStepConfs" bson:"customStepConfs" validate:""`
	PackageNames    []string               `yaml:"packageNames" json:"packageNames" bson:"packageNames" validate:""`
	ArtifactNames   []string               `yaml:"artifactNames" json:"artifactNames" bson:"artifactNames" validate:""`
	BranchNames     []string               `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	EnvNames        []string               `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
}

type DefMetadata struct {
	ProjectName string            `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	Labels      map[string]string `yaml:"labels" json:"labels" bson:"labels" validate:""`
	Annotations map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
}

type DefKind struct {
	Kind     string        `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Metadata DefMetadata   `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Items    []interface{} `yaml:"items" json:"items" bson:"items" validate:""`
	Status   struct {
		ErrMsg string `yaml:"errMsg" json:"errMsg" bson:"errMsg" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type DefKindList struct {
	Kind   string    `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Defs   []DefKind `yaml:"defs" json:"defs" bson:"defs" validate:""`
	Status struct {
		ErrMsgs []string `yaml:"errMsgs" json:"errMsgs" bson:"errMsgs" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type DefUpdate struct {
	Kind           string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	ProjectName    string      `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	EnvName        string      `yaml:"envName" json:"envName" bson:"envName" validate:""`
	CustomStepName string      `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:""`
	BranchName     string      `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Def            interface{} `yaml:"def" json:"def" bson:"def" validate:""`
}

type DefUpdateList struct {
	Kind string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Defs []DefUpdate `yaml:"defs" json:"defs" bson:"defs" validate:""`
}

type DefClone struct {
	Kind        string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	ProjectName string      `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	Def         interface{} `yaml:"def" json:"def" bson:"def" validate:""`
}

type PatchAction struct {
	Action string      `yaml:"action" json:"action" bson:"action" validate:"required"`
	Path   string      `yaml:"path" json:"path" bson:"path" validate:"required"`
	Value  interface{} `yaml:"value" json:"value" bson:"value" validate:""`
	Str    interface{} `yaml:"str" json:"str" bson:"str" validate:""`
}

type AdminMetadata struct {
	Name        string            `yaml:"name" json:"name" bson:"name" validate:""`
	Annotations map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
}

type AdminKind struct {
	Kind     string        `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Metadata AdminMetadata `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Spec     interface{}   `yaml:"spec" json:"spec" bson:"spec" validate:""`
}

type AdminKindList struct {
	Kind  string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Items []AdminKind `yaml:"items" json:"items" bson:"items" validate:""`
}

type RepoNameList struct {
	ArtifactRepoNames []string `yaml:"artifactRepoNames" json:"artifactRepoNames" bson:"artifactRepoNames" validate:""`
	GitRepoNames      []string `yaml:"gitRepoNames" json:"gitRepoNames" bson:"gitRepoNames" validate:""`
	ImageRepoNames    []string `yaml:"imageRepoNames" json:"imageRepoNames" bson:"imageRepoNames" validate:""`
	ScanCodeRepoNames []string `yaml:"scanCodeRepoNames" json:"scanCodeRepoNames" bson:"scanCodeRepoNames" validate:""`
}

type AccessToken struct {
	AccessToken     string `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	AccessTokenName string `yaml:"accessTokenName" json:"accessTokenName" bson:"accessTokenName" validate:""`
	ExpireTime      string `yaml:"expireTime" json:"expireTime" bson:"expireTime" validate:""`
	Expired         bool   `yaml:"expired" json:"expired" bson:"expired" validate:""`
	Username        string `yaml:"username" json:"username" bson:"username" validate:""`
}

type KubernetesHaCluster struct {
	Version            string `yaml:"version" json:"version" bson:"version" validate:"required"`
	ImageRepository    string `yaml:"imageRepository" json:"imageRepository" bson:"imageRepository" validate:""`
	VirtualIp          string `yaml:"virtualIp" json:"virtualIp" bson:"virtualIp" validate:"required"`
	VirtualPort        int    `yaml:"virtualPort" json:"virtualPort" bson:"virtualPort" validate:"required"`
	VirtualHostname    string `yaml:"virtualHostname" json:"virtualHostname" bson:"virtualHostname" validate:"required"`
	CriSocket          string `yaml:"criSocket" json:"criSocket" bson:"criSocket" validate:"required"`
	PodSubnet          string `yaml:"podSubnet" json:"podSubnet" bson:"podSubnet" validate:""`
	ServiceSubnet      string `yaml:"serviceSubnet" json:"serviceSubnet" bson:"serviceSubnet" validate:""`
	KeepAlivedAuthPass string `yaml:"keepAlivedAuthPass" json:"keepAlivedAuthPass" bson:"keepAlivedAuthPass" validate:""`
	MasterHosts        []struct {
		IpAddress          string `yaml:"ipAddress" json:"ipAddress" bson:"ipAddress" validate:"required"`
		Hostname           string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required"`
		NetworkInterface   string `yaml:"networkInterface" json:"networkInterface" bson:"networkInterface" validate:"required"`
		KeepalivedPriority int    `yaml:"keepalivedPriority" json:"keepalivedPriority" bson:"keepalivedPriority" validate:"required"`
	} `yaml:"masterHosts" json:"masterHosts" bson:"masterHosts" validate:"required"`
}
