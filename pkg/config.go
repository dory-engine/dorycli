package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (ic *InstallConfig) VerifyInstallConfig() error {
	var err error
	errInfo := fmt.Sprintf("verify install config error")

	var fieldName, fieldValue string

	validate := validator.New()
	err = validate.Struct(ic)
	if err != nil {
		err = fmt.Errorf("validate install config error: %s", err.Error())
		return err
	}

	fieldName = "dory.namespace"
	fieldValue = ic.Dory.Namespace
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.Dory.GitRepo.Type != "gitlab" && ic.Dory.GitRepo.Type != "gitea" {
		ic.Dory.GitRepo = GitRepo{}
	}

	if ic.Dory.ImageRepo.Type != "harbor" {
		ic.Dory.ImageRepo = ImageRepo{}
	}

	if ic.Dory.ArtifactRepo.Type != "nexus" {
		ic.Dory.ArtifactRepo = ArtifactRepo{}
	}

	if ic.Dory.ScanCodeRepo.Type != "sonarqube" {
		ic.Dory.ScanCodeRepo = ScanCodeRepo{}
	}

	if ic.Dory.GitRepo.Internal.Image != "" && ic.Dory.GitRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: dory.gitRepo.internal and dory.gitRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.GitRepo.Internal.Image != "" {
		fieldName = "dory.gitRepo.internal.imageDB"
		fieldValue = ic.Dory.GitRepo.Internal.ImageDB
		if ic.Dory.GitRepo.Type == "gitea" && ic.Dory.GitRepo.Internal.ImageDB == "" {
			err = fmt.Errorf("%s: %s %s format error: gitea imageDB can not be empty", errInfo, fieldName, fieldValue)
			return err
		}

		fieldName = "dory.gitRepo.internal.imageNginx"
		fieldValue = ic.Dory.GitRepo.Internal.ImageNginx
		if ic.Dory.GitRepo.Type == "gitlab" && ic.Dory.GitRepo.Internal.ImageNginx == "" {
			err = fmt.Errorf("%s: %s %s format error: gitlab imageNginx can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if ic.Dory.ArtifactRepo.Internal.Image != "" && ic.Dory.ArtifactRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: dory.artifactRepo.internal and dory.artifactRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.ArtifactRepo.External.ViewUrl != "" && ic.Dory.ArtifactRepo.External.Schema != "http" && ic.Dory.ArtifactRepo.External.Schema != "https" {
		err = fmt.Errorf("%s: dory.artifactRepo.external.schema must be http or https", errInfo)
		return err
	}

	if ic.Dory.ArtifactRepo.External.ViewUrl != "" && (ic.Dory.ArtifactRepo.External.ProxyRepo.Pip == "" || ic.Dory.ArtifactRepo.External.ProxyRepo.Go == "" || ic.Dory.ArtifactRepo.External.ProxyRepo.Gradle == "" || ic.Dory.ArtifactRepo.External.ProxyRepo.Maven == "" || ic.Dory.ArtifactRepo.External.ProxyRepo.Npm == "") {
		err = fmt.Errorf("%s: dory.artifactRepo.external.proxyRepo required", errInfo)
		return err
	}

	if ic.Dory.ScanCodeRepo.Internal.Image != "" && ic.Dory.ScanCodeRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: dory.scanCodeRepo.internal and dory.scanCodeRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.DemoDatabase.Internal.Image == "" && ic.Dory.DemoDatabase.External.DbUrl == "" {
		err = fmt.Errorf("%s: dory.demoDatabase.internal and dory.demoDatabase.external both empty", errInfo)
		return err
	}

	if ic.Dory.DemoDatabase.Internal.Image != "" && ic.Dory.DemoDatabase.External.DbUrl != "" {
		err = fmt.Errorf("%s: dory.demoDatabase.internal and dory.demoDatabase.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.DemoHost.Internal.Image == "" && ic.Dory.DemoHost.External.HostAddr == "" {
		err = fmt.Errorf("%s: dory.demoHost.internal and dory.demoHost.external both empty", errInfo)
		return err
	}

	if ic.Dory.DemoHost.Internal.Image != "" && ic.Dory.DemoHost.External.HostAddr != "" {
		err = fmt.Errorf("%s: dory.demoHost.internal and dory.demoHost.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.ImageRepo.Internal.Hostname != "" && ic.Dory.ImageRepo.External.Hostname != "" {
		err = fmt.Errorf("%s: imageRepo.internal and imageRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.ImageRepo.Internal.Hostname != "" {
		fieldName = "imageRepo.internal.namespace"
		fieldValue = ic.Dory.ImageRepo.Internal.Namespace
		if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
			err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
			return err
		}

		arr := strings.Split(ic.Dory.ImageRepo.Internal.Version, ".")
		if len(arr) != 3 {
			fieldName = "imageRepo.internal.version"
			fieldValue = ic.Dory.ImageRepo.Internal.Version
			err = fmt.Errorf("%s: %s %s format error: should like v2.4.0", errInfo, fieldName, fieldValue)
			return err
		}
		arr[2] = "0"
		ic.Dory.ImageRepo.Internal.VersionBig = strings.Join(arr, ".")

		if ic.Dory.ImageRepo.Internal.Password == "" {
			ic.Dory.ImageRepo.Internal.Password = RandomString(16, false, "=")
		}
		if ic.Dory.ImageRepo.Internal.RegistryPassword == "" {
			ic.Dory.ImageRepo.Internal.RegistryPassword = RandomString(16, false, "")
		}
		bs, _ := bcrypt.GenerateFromPassword([]byte(ic.Dory.ImageRepo.Internal.RegistryPassword), 10)
		ic.Dory.ImageRepo.Internal.RegistryHtpasswd = string(bs)
	}

	fieldName = "hostIP"
	fieldValue = ic.HostIP
	err = ValidateIpAddress(fieldValue)
	if err != nil {
		err = fmt.Errorf("%s: %s %s format error: %s", errInfo, fieldName, fieldValue, err.Error())
		return err
	}
	if fieldValue == "127.0.0.1" || fieldValue == "localhost" {
		err = fmt.Errorf("%s: %s %s format error: can not be 127.0.0.1 or localhost", errInfo, fieldName, fieldValue)
		return err
	}

	var count int
	if ic.Kubernetes.PvConfigLocal.LocalPath != "" {
		count = count + 1
	}
	if len(ic.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		count = count + 1
	}
	if ic.Kubernetes.PvConfigNfs.NfsServer != "" {
		count = count + 1
	}
	if len(ic.Kubernetes.PvConfigCsiCephfs.CephPath) > 0 {
		count = count + 1
	}
	if ic.Kubernetes.PvConfigCsiNfs.NfsServer != "" {
		count = count + 1
	}
	if count != 1 {
		err = fmt.Errorf("%s: kubernetes.pvConfigLocal/pvConfigNfs/pvConfigCephfs/pvConfigCsiNfs/pvConfigCsiCephfs must set one only", errInfo)
		return err
	}

	if ic.Kubernetes.Runtime != "docker" && ic.Kubernetes.Runtime != "containerd" && ic.Kubernetes.Runtime != "crio" {
		fieldName = "kubernetes.runtime"
		fieldValue = ic.Kubernetes.Runtime
		err = fmt.Errorf("%s: %s %s format error: must be docker or containerd or crio", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.Kubernetes.PvConfigLocal.LocalPath != "" {
		if !strings.HasPrefix(ic.Kubernetes.PvConfigLocal.LocalPath, "/") {
			fieldName = "kubernetes.pvConfigLocal.localPath"
			fieldValue = ic.Kubernetes.PvConfigLocal.LocalPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
		ic.Kubernetes.PvType = "local-path"
		ic.Kubernetes.PvPath = ic.Kubernetes.PvConfigLocal.LocalPath
	}
	if len(ic.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		for _, monitor := range ic.Kubernetes.PvConfigCephfs.CephMonitors {
			fieldName = "kubernetes.pvConfigCephfs.cephMonitors"
			fieldValue = monitor
			arr := strings.Split(monitor, ":")
			if len(arr) != 2 {
				err = fmt.Errorf("%s: %s %s format error: should like 192.168.0.1:6789", errInfo, fieldName, fieldValue)
				return err
			}
			_, err = strconv.Atoi(arr[1])
			if err != nil {
				err = fmt.Errorf("%s: %s %s format error: should like 192.168.0.1:6789", errInfo, fieldName, fieldValue)
				return err
			}
		}
		if ic.Kubernetes.PvConfigCephfs.CephSecret == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephSecret"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephSecret
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if ic.Kubernetes.PvConfigCephfs.CephUser == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephUser"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephUser
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if !strings.HasPrefix(ic.Kubernetes.PvConfigCephfs.CephPath, "/") {
			fieldName = "kubernetes.pvConfigCephfs.cephPath"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}

		ic.Kubernetes.PvType = "cephfs"
		ic.Kubernetes.PvPath = ic.Kubernetes.PvConfigCephfs.CephPath
	}

	if ic.Kubernetes.PvConfigNfs.NfsServer != "" {
		if !strings.HasPrefix(ic.Kubernetes.PvConfigNfs.NfsPath, "/") {
			fieldName = "kubernetes.pvConfigNfs.nfsPath"
			fieldValue = ic.Kubernetes.PvConfigNfs.NfsPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
		ic.Kubernetes.PvType = "nfs"
		ic.Kubernetes.PvPath = ic.Kubernetes.PvConfigNfs.NfsPath
	}

	if ic.Kubernetes.PvConfigCsiNfs.NfsServer != "" {
		if !strings.HasPrefix(ic.Kubernetes.PvConfigCsiNfs.NfsPath, "/") {
			fieldName = "kubernetes.pvConfigCsiNfs.nfsPath"
			fieldValue = ic.Kubernetes.PvConfigCsiNfs.NfsPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
		for _, opt := range ic.Kubernetes.PvConfigCsiNfs.NfsMountOptions {
			arr := strings.Split(opt, "=")
			if len(arr) != 2 {
				fieldName = "kubernetes.pvConfigCsiNfs.nfsMountOptions"
				fieldValue = opt
				err = fmt.Errorf("%s: %s %s format error: should like nfsvers=4.1", errInfo, fieldName, fieldValue)
				return err
			}
			if arr[0] == "" || arr[1] == "" {
				fieldName = "kubernetes.pvConfigCsiNfs.nfsMountOptions"
				fieldValue = opt
				err = fmt.Errorf("%s: %s %s format error: should like nfsvers=4.1", errInfo, fieldName, fieldValue)
				return err
			}
		}
		ic.Kubernetes.PvType = "csi-nfs"
		ic.Kubernetes.PvPath = ic.Kubernetes.PvConfigCsiNfs.NfsPath
	}

	if ic.Kubernetes.PvConfigCsiCephfs.CephPath != "" {
		if ic.Kubernetes.PvConfigCsiCephfs.CephSecret == "" {
			fieldName = "kubernetes.pvConfigCsiCephfs.cephSecret"
			fieldValue = ic.Kubernetes.PvConfigCsiCephfs.CephSecret
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if ic.Kubernetes.PvConfigCsiCephfs.CephUser == "" {
			fieldName = "kubernetes.pvConfigCsiCephfs.cephUser"
			fieldValue = ic.Kubernetes.PvConfigCsiCephfs.CephUser
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if ic.Kubernetes.PvConfigCsiCephfs.CephFsName == "" {
			fieldName = "kubernetes.pvConfigCsiCephfs.cephFsName"
			fieldValue = ic.Kubernetes.PvConfigCsiCephfs.CephFsName
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if ic.Kubernetes.PvConfigCsiCephfs.CephClusterId == "" {
			fieldName = "kubernetes.pvConfigCsiCephfs.cephClusterId"
			fieldValue = ic.Kubernetes.PvConfigCsiCephfs.CephClusterId
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if !strings.HasPrefix(ic.Kubernetes.PvConfigCsiCephfs.CephPath, "/") {
			fieldName = "kubernetes.pvConfigCsiCephfs.cephPath"
			fieldValue = ic.Kubernetes.PvConfigCsiCephfs.CephPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
		ic.Kubernetes.PvType = "csi-cephfs"
		ic.Kubernetes.PvPath = ic.Kubernetes.PvConfigCsiCephfs.CephPath
	}

	if ic.Dory.Openldap.Password == "" {
		ic.Dory.Openldap.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Redis.Password == "" {
		ic.Dory.Redis.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Mongo.Password == "" {
		ic.Dory.Mongo.Password = RandomString(16, false, "=")
	}
	if ic.Dory.DemoDatabase.Internal.Password == "" {
		ic.Dory.DemoDatabase.Internal.Password = RandomString(16, false, "=")
	}
	if ic.Dory.DemoDatabase.Internal.UserPassword == "" {
		ic.Dory.DemoDatabase.Internal.UserPassword = RandomString(16, false, "=")
	}
	if ic.Dory.DemoHost.Internal.Password == "" {
		ic.Dory.DemoHost.Internal.Password = RandomString(16, false, "=")
	}
	return err
}

func (ic *InstallConfig) UnmarshalMapValues() (map[string]interface{}, error) {
	var err error
	errInfo := fmt.Sprintf("unmarshal install config to map error")

	bs, _ := json.Marshal(ic)
	vals := map[string]interface{}{}
	err = json.Unmarshal(bs, &vals)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return vals, err
	}
	gitRepoInternal := true
	gitRepoViewUrl := fmt.Sprintf("%s:%d", ic.ViewURL, ic.Dory.GitRepo.Internal.Port)
	gitRepoUrl := ""
	if ic.Dory.GitRepo.Type == "gitea" {
		gitRepoUrl = fmt.Sprintf("http://%s:3000", ic.Dory.GitRepo.Type)
	} else if ic.Dory.GitRepo.Type == "gitlab" {
		gitRepoUrl = fmt.Sprintf("http://%s", ic.Dory.GitRepo.Type)
	}
	gitRepoUsername := "GIT_REPO_USERNAME"
	gitRepoName := "GIT_REPO_NAME"
	gitRepoMail := "GIT_REPO_MAIL@example.com"
	gitRepoPassword := RandomString(16, false, "=")
	gitRepoToken := "GIT_REPO_TOKEN"
	if ic.Dory.GitRepo.Internal.Image == "" {
		gitRepoInternal = false
		gitRepoViewUrl = ic.Dory.GitRepo.External.ViewUrl
		gitRepoUrl = ic.Dory.GitRepo.External.Url
		gitRepoUsername = ic.Dory.GitRepo.External.Username
		gitRepoName = ic.Dory.GitRepo.External.Name
		gitRepoMail = ic.Dory.GitRepo.External.Mail
		gitRepoPassword = ic.Dory.GitRepo.External.Password
		gitRepoToken = ic.Dory.GitRepo.External.Token
	}
	vals["gitRepoInternal"] = gitRepoInternal
	vals["gitRepoViewUrl"] = gitRepoViewUrl
	vals["gitRepoUrl"] = gitRepoUrl
	vals["gitRepoUsername"] = gitRepoUsername
	vals["gitRepoName"] = gitRepoName
	vals["gitRepoMail"] = gitRepoMail
	vals["gitRepoPassword"] = gitRepoPassword
	vals["gitRepoToken"] = gitRepoToken

	imageRepoInternal := true
	imageRepoDomainName := ic.Dory.ImageRepo.Internal.Hostname
	imageRepoUsername := "admin"
	imageRepoPassword := ic.Dory.ImageRepo.Internal.Password
	imageRepoEmail := "admin@example.com"
	imageRepoIp := ic.HostIP
	if ic.Dory.ImageRepo.Internal.Hostname == "" {
		imageRepoInternal = false
		imageRepoDomainName = ic.Dory.ImageRepo.External.Hostname
		imageRepoUsername = ic.Dory.ImageRepo.External.Username
		imageRepoPassword = ic.Dory.ImageRepo.External.Password
		imageRepoEmail = ic.Dory.ImageRepo.External.Email
		imageRepoIp = ic.Dory.ImageRepo.External.Ip
	}
	vals["imageRepoInternal"] = imageRepoInternal
	vals["imageRepoDomainName"] = imageRepoDomainName
	vals["imageRepoUsername"] = imageRepoUsername
	vals["imageRepoPassword"] = imageRepoPassword
	vals["imageRepoEmail"] = imageRepoEmail
	vals["imageRepoIp"] = imageRepoIp

	artifactRepoInternal := true
	artifactRepoViewUrl := fmt.Sprintf("%s:%d", ic.ViewURL, ic.Dory.ArtifactRepo.Internal.Port)
	artifactRepoSchema := "http"
	artifactRepoPort := ic.Dory.ArtifactRepo.Internal.Port
	artifactRepoPortHub := ic.Dory.ArtifactRepo.Internal.PortHub
	artifactRepoPortGcr := ic.Dory.ArtifactRepo.Internal.PortGcr
	artifactRepoPortQuay := ic.Dory.ArtifactRepo.Internal.PortQuay
	artifactRepoUsername := "admin"
	artifactRepoPassword := RandomString(16, false, "=")
	artifactRepoPublicRole := "public-role"
	artifactRepoPublicUser := "public-user"
	artifactRepoPublicPassword := "public-user"
	artifactRepoPublicEmail := "public-user@example.com"
	artifactRepoIp := ic.HostIP
	if ic.Dory.ArtifactRepo.Internal.Image == "" {
		artifactRepoInternal = false
		artifactRepoViewUrl = ic.Dory.ArtifactRepo.External.ViewUrl
		artifactRepoSchema = ic.Dory.ArtifactRepo.External.Schema
		artifactRepoPort = ic.Dory.ArtifactRepo.External.Port
		artifactRepoPortHub = ic.Dory.ArtifactRepo.External.PortHub
		artifactRepoPortGcr = ic.Dory.ArtifactRepo.External.PortGcr
		artifactRepoPortQuay = ic.Dory.ArtifactRepo.External.PortQuay
		artifactRepoUsername = ic.Dory.ArtifactRepo.External.Username
		artifactRepoPassword = ic.Dory.ArtifactRepo.External.Password
		artifactRepoPublicRole = ic.Dory.ArtifactRepo.External.PublicRole
		artifactRepoPublicUser = ic.Dory.ArtifactRepo.External.PublicUser
		artifactRepoPublicPassword = ic.Dory.ArtifactRepo.External.PublicPassword
		artifactRepoPublicEmail = ic.Dory.ArtifactRepo.External.PublicEmail
		artifactRepoIp = ic.Dory.ArtifactRepo.External.Hostname
	}
	vals["artifactRepoInternal"] = artifactRepoInternal
	vals["artifactRepoViewUrl"] = artifactRepoViewUrl
	vals["artifactRepoSchema"] = artifactRepoSchema
	vals["artifactRepoPort"] = artifactRepoPort
	vals["artifactRepoPortHub"] = artifactRepoPortHub
	vals["artifactRepoPortGcr"] = artifactRepoPortGcr
	vals["artifactRepoPortQuay"] = artifactRepoPortQuay
	vals["artifactRepoUsername"] = artifactRepoUsername
	vals["artifactRepoPassword"] = artifactRepoPassword
	vals["artifactRepoPublicRole"] = artifactRepoPublicRole
	vals["artifactRepoPublicUser"] = artifactRepoPublicUser
	vals["artifactRepoPublicPassword"] = artifactRepoPublicPassword
	vals["artifactRepoPublicEmail"] = artifactRepoPublicEmail
	vals["artifactRepoIp"] = artifactRepoIp

	scanCodeRepoInternal := true
	scanCodeRepoViewUrl := fmt.Sprintf("%s:%d", ic.ViewURL, ic.Dory.ScanCodeRepo.Internal.Port)
	scanCodeRepoUrl := ""
	scanCodeRepoToken := "SCAN_CODE_REPO_TOKEN"
	scanCodeRepoPassword := RandomString(16, false, "=")
	if ic.Dory.ScanCodeRepo.Internal.Image == "" {
		scanCodeRepoInternal = false
		scanCodeRepoViewUrl = ic.Dory.ScanCodeRepo.External.ViewUrl
		scanCodeRepoUrl = ic.Dory.ScanCodeRepo.External.Url
		scanCodeRepoToken = ic.Dory.ScanCodeRepo.External.Token
	} else {
		scanCodeRepoUrl = fmt.Sprintf("http://%s-web:9000", ic.Dory.ScanCodeRepo.Type)
	}
	vals["scanCodeRepoInternal"] = scanCodeRepoInternal
	vals["scanCodeRepoViewUrl"] = scanCodeRepoViewUrl
	vals["scanCodeRepoUrl"] = scanCodeRepoUrl
	vals["scanCodeRepoToken"] = scanCodeRepoToken
	vals["scanCodeRepoPassword"] = scanCodeRepoPassword

	demoDatabaseInternal := true
	demoDatabaseUrl := fmt.Sprintf("jdbc:mysql://%s:%d/%s", ic.HostIP, ic.Dory.DemoDatabase.Internal.Port, ic.Dory.DemoDatabase.Internal.Database)
	demoDatabaseUsername := ic.Dory.DemoDatabase.Internal.User
	demoDatabasePassword := ic.Dory.DemoDatabase.Internal.UserPassword
	if ic.Dory.DemoDatabase.Internal.Image == "" {
		demoDatabaseInternal = false
		demoDatabaseUrl = ic.Dory.DemoDatabase.External.DbUrl
		demoDatabaseUsername = ic.Dory.DemoDatabase.External.DbUser
		demoDatabasePassword = ic.Dory.DemoDatabase.External.DbPassword
	}
	vals["demoDatabaseInternal"] = demoDatabaseInternal
	vals["demoDatabaseUrl"] = demoDatabaseUrl
	vals["demoDatabaseUsername"] = demoDatabaseUsername
	vals["demoDatabasePassword"] = demoDatabasePassword

	demoHostInternal := true
	demoHostAddr := ic.HostIP
	demoHostPort := ic.Dory.DemoHost.Internal.NodePortSsh
	demoHostUser := "root"
	demoHostPassword := ic.Dory.DemoHost.Internal.Password
	demoHostBecome := false
	demoHostBecomeUser := ""
	demoHostBecomePassword := ""
	demoHostNodePortWeb := ic.Dory.DemoHost.Internal.NodePortWeb
	if ic.Dory.DemoHost.Internal.Image == "" {
		demoHostInternal = false
		demoHostAddr = ic.Dory.DemoHost.External.HostAddr
		demoHostPort = ic.Dory.DemoHost.External.HostPort
		demoHostUser = ic.Dory.DemoHost.External.HostUser
		demoHostPassword = ic.Dory.DemoHost.External.HostPassword
		demoHostBecome = ic.Dory.DemoHost.External.HostBecome
		demoHostBecomeUser = ic.Dory.DemoHost.External.HostBecomeUser
		demoHostBecomePassword = ic.Dory.DemoHost.External.HostBecomePassword
		demoHostNodePortWeb = ic.Dory.DemoHost.External.NodePortWeb
	}
	vals["demoHostInternal"] = demoHostInternal
	vals["demoHostAddr"] = demoHostAddr
	vals["demoHostPort"] = demoHostPort
	vals["demoHostUser"] = demoHostUser
	vals["demoHostPassword"] = demoHostPassword
	vals["demoHostBecome"] = demoHostBecome
	vals["demoHostBecomeUser"] = demoHostBecomeUser
	vals["demoHostBecomePassword"] = demoHostBecomePassword
	vals["demoHostNodePortWeb"] = demoHostNodePortWeb

	gitWebhookUrl := "http://dory-engine:9000"
	if (ic.Dory.GitRepo.Type == "gitlab" || ic.Dory.GitRepo.Type == "gitea") && ic.Dory.GitRepo.External.Url != "" && ic.Dory.GitRepo.External.GitWebhookUrl != "" {
		gitWebhookUrl = ic.Dory.GitRepo.External.GitWebhookUrl
	}
	vals["gitWebhookUrl"] = gitWebhookUrl

	return vals, err
}

func (ic *InstallConfig) HarborQuery(url, method string, param map[string]interface{}) (string, int, error) {
	var err error
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	domainName := ic.Dory.ImageRepo.Internal.Hostname
	username := "admin"
	password := ic.Dory.ImageRepo.Internal.Password
	if ic.Dory.ImageRepo.Internal.Hostname == "" {
		domainName = ic.Dory.ImageRepo.External.Hostname
		username = ic.Dory.ImageRepo.External.Username
		password = ic.Dory.ImageRepo.External.Password
	}

	url = fmt.Sprintf("https://%s%s", domainName, url)

	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return strJson, statusCode, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return strJson, statusCode, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return strJson, statusCode, err
		}
	}

	req.SetBasicAuth(username, password)
	resp, err = client.Do(req)
	if err != nil {
		return strJson, statusCode, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return strJson, statusCode, err
	}
	strJson = string(bs)
	return strJson, statusCode, err
}

func (ic *InstallConfig) HarborProjectAdd(projectName string) error {
	var err error
	var statusCode int
	var strJson string

	url := fmt.Sprintf("/api/v2.0/projects")
	param := map[string]interface{}{
		"project_name": projectName,
		"public":       true,
	}
	strJson, statusCode, err = ic.HarborQuery(url, http.MethodPost, param)
	if err != nil {
		return err
	}

	errmsg := fmt.Sprintf("%s %s", gjson.Get(strJson, "errors.0.code").String(), gjson.Get(strJson, "errors.0.message").String())
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf(errmsg)
		return err
	}

	return err
}

func (ic *InstallConfig) KubernetesQuery(url, method string, param map[string]interface{}) (string, int, error) {
	var err error
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url = fmt.Sprintf("https://%s:%d%s", ic.Kubernetes.Host, ic.Kubernetes.Port, url)

	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return strJson, statusCode, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return strJson, statusCode, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return strJson, statusCode, err
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ic.Kubernetes.Token))
	resp, err = client.Do(req)
	if err != nil {
		return strJson, statusCode, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return strJson, statusCode, err
	}
	strJson = string(bs)
	return strJson, statusCode, err
}

func (ic *InstallConfig) KubernetesPodsGet(namespace string) ([]KubePod, error) {
	var err error
	var statusCode int
	var strJson string

	pods := []KubePod{}

	url := fmt.Sprintf("/api/v1/namespaces/%s/pods", namespace)
	param := map[string]interface{}{}
	strJson, statusCode, err = ic.KubernetesQuery(url, http.MethodGet, param)
	if err != nil {
		return pods, err
	}

	errmsg := gjson.Get(strJson, "message").String()
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf(errmsg)
		return pods, err
	}

	var podList KubePodList
	err = json.Unmarshal([]byte(strJson), &podList)
	if err != nil {
		return pods, err
	}

	pods = podList.Items
	return pods, err
}

func (khc *KubernetesHaCluster) VerifyKubernetesHaCluster() error {
	var err error
	errInfo := fmt.Sprintf("verify kubernetes ha config error")

	var fieldName, fieldValue string

	validate := validator.New()
	err = validate.Struct(khc)
	if err != nil {
		err = fmt.Errorf("validate kubernetes ha config error: %s", err.Error())
		return err
	}

	fieldName = "virtualIp"
	fieldValue = khc.VirtualIp
	err = ValidateIpAddress(fieldValue)
	if err != nil {
		err = fmt.Errorf("%s: %s %s format error: %s", errInfo, fieldName, fieldValue, err.Error())
		return err
	}

	m := map[int]int{}
	for _, host := range khc.MasterHosts {
		fieldName = "masterHosts.ipAddress"
		fieldValue = host.IpAddress
		err = ValidateIpAddress(fieldValue)
		if err != nil {
			err = fmt.Errorf("%s: %s %s format error: %s", errInfo, fieldName, fieldValue, err.Error())
			return err
		}

		_, found := m[host.KeepalivedPriority]
		if found {
			fieldName = "masterHosts.keepalivedPriority"
			fieldValue = fmt.Sprintf("%d", host.KeepalivedPriority)
			err = fmt.Errorf("%s: %s %s is duplicated", errInfo, fieldName, fieldValue)
		}

		m[host.KeepalivedPriority] = host.KeepalivedPriority
	}

	if khc.KeepAlivedAuthPass == "" {
		khc.KeepAlivedAuthPass = RandomString(16, false, "=")
	}

	return err
}
