package cmd

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OptionsCommon struct {
	ServerURL    string `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure     bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout      int    `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken  string `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	Language     string `yaml:"language" json:"language" bson:"language" validate:""`
	ConfigFile   string `yaml:"configFile" json:"configFile" bson:"configFile" validate:""`
	Verbose      bool   `yaml:"verbose" json:"verbose" bson:"verbose" validate:""`
	ConfigExists bool   `yaml:"configExists" json:"configExists" bson:"configExists" validate:""`
	LangBundle   *i18n.Bundle
}

type Log struct {
	Verbose bool `yaml:"verbose" json:"verbose" bson:"verbose" validate:""`
}

func (log *Log) SetVerbose(verbose bool) {
	log.Verbose = verbose
}

func (log *Log) Debug(msg string) {
	if log.Verbose {
		defer color.Unset()
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
	}
}

func (log *Log) Success(msg string) {
	defer color.Unset()
	color.Set(color.FgGreen)
	fmt.Println(fmt.Sprintf("[SUCC] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Info(msg string) {
	defer color.Unset()
	color.Set(color.FgBlue)
	fmt.Println(fmt.Sprintf("[INFO] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Warning(msg string) {
	defer color.Unset()
	color.Set(color.FgMagenta)
	fmt.Println(fmt.Sprintf("[WARN] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Error(msg string) {
	defer color.Unset()
	color.Set(color.FgRed)
	fmt.Println(fmt.Sprintf("[ERRO] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) RunLog(msg pkg.WsRunLog) {
	defer color.Unset()
	bs, _ := json.Marshal(msg)
	strJson := string(bs)
	switch msg.LogType {
	case pkg.LogTypeInfo:
		color.Set(color.FgBlue)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	case pkg.LogTypeWarning:
		color.Set(color.FgMagenta)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	case pkg.LogTypeError:
		color.Set(color.FgRed)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	}
	if log.Verbose {
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), strJson))
	}
}

func (log *Log) AdminLog(msg pkg.WsAdminLog) {
	defer color.Unset()
	bs, _ := json.Marshal(msg)
	strJson := string(bs)
	switch msg.LogType {
	case pkg.LogTypeInfo:
		color.Set(color.FgBlue)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	case pkg.StatusFail:
		color.Set(color.FgRed)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	case pkg.StatusSuccess:
		color.Set(color.FgGreen)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	}
	if log.Verbose {
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), strJson))
	}
}

func CheckError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func PrettyJson(strJson string) (string, error) {
	var err error
	var strPretty string
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(strJson), "", "  ")
	if err != nil {
		return strPretty, err
	}
	strPretty = buf.String()
	return strPretty, err
}

func NewOptionsCommon() *OptionsCommon {
	var o OptionsCommon
	bsEN, _ := pkg.FsLanguage.ReadFile(fmt.Sprintf("language/lang.en.yaml"))

	bsZH, _ := pkg.FsLanguage.ReadFile(fmt.Sprintf("language/lang.zh.yaml"))

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustParseMessageFileBytes(bsEN, "lang.en.yaml")
	bundle.MustParseMessageFileBytes(bsZH, "lang.zh.yaml")

	o.LangBundle = bundle

	return &o
}

var OptCommon = NewOptionsCommon()
var log Log

func NewCmdRoot() *cobra.Command {
	o := OptCommon
	baseName := pkg.GetCmdBaseName()
	msgUse := baseName

	_ = OptCommon.GetOptionsCommon()
	msgShort := OptCommon.TransLang("cmd_short")
	msgLong := OptCommon.TransLang("cmd_long")
	msgExample := pkg.Indent(OptCommon.TransLang("cmd_example", baseName))

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&o.ConfigFile, "config", "c", "", OptCommon.TransLang("param_config", baseName, pkg.EnvVarConfigFile, pkg.ConfigDirDefault, pkg.ConfigFileDefault))
	cmd.PersistentFlags().StringVarP(&o.ServerURL, "server-url", "s", "", OptCommon.TransLang("param_server_url"))
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", false, OptCommon.TransLang("param_insecure"))
	cmd.PersistentFlags().IntVar(&o.Timeout, "timeout", pkg.TimeoutDefault, OptCommon.TransLang("param_timeout"))
	cmd.PersistentFlags().StringVar(&o.AccessToken, "token", "", OptCommon.TransLang("param_token"))
	cmd.PersistentFlags().StringVar(&o.Language, "language", "", OptCommon.TransLang("param_language"))
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, fmt.Sprintf("show logs in verbose mode"))

	cmd.AddCommand(NewCmdLogin())
	cmd.AddCommand(NewCmdLogout())
	cmd.AddCommand(NewCmdProject())
	cmd.AddCommand(NewCmdPipeline())
	cmd.AddCommand(NewCmdRun())
	cmd.AddCommand(NewCmdDef())
	cmd.AddCommand(NewCmdAdmin())
	cmd.AddCommand(NewCmdInstall())
	cmd.AddCommand(NewCmdVersion())

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsCommon) TransLang(msg string, args ...interface{}) string {
	var err error
	var s string
	m := map[string]interface{}{}
	if o.Language != "en" && o.Language != "zh" {
		o.Language = "en"
	}
	for i, arg := range args {
		m[fmt.Sprintf("_%d", i)] = arg
	}
	loc := i18n.NewLocalizer(o.LangBundle, o.Language)
	s, err = loc.Localize(&i18n.LocalizeConfig{
		MessageID:    msg,
		TemplateData: m,
	})
	if err != nil {
		loc = i18n.NewLocalizer(o.LangBundle, "en")
		s, err = loc.Localize(&i18n.LocalizeConfig{
			MessageID:    msg,
			TemplateData: m,
		})
		if err != nil {
			s = err.Error()
		}
	}
	return s
}

func (o *OptionsCommon) CheckConfigFile() error {
	errInfo := fmt.Sprintf("check config file error")
	var err error

	if o.ConfigFile == "" {
		v, exists := os.LookupEnv(pkg.EnvVarConfigFile)
		if exists {
			o.ConfigFile = v
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
			defaultConfigFile := fmt.Sprintf("%s/%s/%s", homeDir, pkg.ConfigDirDefault, pkg.ConfigFileDefault)
			o.ConfigFile = defaultConfigFile
		}
	}
	fi, err := os.Stat(o.ConfigFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			configDir := filepath.Dir(o.ConfigFile)
			err = os.MkdirAll(configDir, 0700)
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
			err = os.WriteFile(o.ConfigFile, []byte{}, 0600)
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
		} else {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
	} else {
		if fi.IsDir() {
			err = fmt.Errorf("%s: %s must be a file", errInfo, o.ConfigFile)
			return err
		}
	}
	bs, err := os.ReadFile(o.ConfigFile)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}
	var doryConfig pkg.DoryConfig
	err = yaml.Unmarshal(bs, &doryConfig)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}

	if doryConfig.AccessToken == "" {
		bs, err = pkg.YamlIndent(doryConfig)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}

		err = os.WriteFile(o.ConfigFile, bs, 0600)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
	}

	return err
}

func (o *OptionsCommon) GetOptionsCommon() error {
	errInfo := fmt.Sprintf("get common option error")
	var err error

	err = o.CheckConfigFile()
	if err != nil {
		return err
	}

	bs, err := os.ReadFile(o.ConfigFile)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}
	var doryConfig pkg.DoryConfig
	err = yaml.Unmarshal(bs, &doryConfig)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}

	if o.ServerURL == "" && doryConfig.ServerURL != "" {
		o.ServerURL = doryConfig.ServerURL
	}

	if o.AccessToken == "" && doryConfig.AccessToken != "" {
		bs, err = base64.StdEncoding.DecodeString(doryConfig.AccessToken)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
		o.AccessToken = string(bs)
	}

	if o.Language == "" {
		lang := "en"
		l, err := locale.Detect()
		if err == nil {
			b, _ := l.Base()
			if strings.ToLower(b.String()) == "zh" {
				lang = "zh"
			}
		}
		o.Language = lang
	}
	if o.Language == "" && doryConfig.Language != "" {
		o.Language = doryConfig.Language
	}

	if o.Timeout == 0 && doryConfig.Timeout != 0 && doryConfig.Timeout != pkg.TimeoutDefault {
		o.Timeout = doryConfig.Timeout
	}

	if o.Verbose {
		log.SetVerbose(o.Verbose)
	}

	return err
}

func (o *OptionsCommon) QueryAPI(url, method, userToken string, param map[string]interface{}, showSuccess bool) (gjson.Result, string, error) {
	var err error
	var result gjson.Result
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	var xUserToken string
	client := &http.Client{
		Timeout: time.Second * time.Duration(o.Timeout),
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if !strings.HasPrefix(url, "api/public/") && url != "api/account/accessToken" && url != "api/account/accessTokens" && (o.AccessToken == "" || o.ServerURL == "") {
		err = fmt.Errorf("please login first")
		return result, xUserToken, err
	}
	if o.ServerURL == "" {
		err = fmt.Errorf("--server-url required")
		return result, xUserToken, err
	}
	url = fmt.Sprintf("%s/%s", o.ServerURL, url)

	var strReqBody string
	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return result, xUserToken, err
		}
		strReqBody = string(bs)
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return result, xUserToken, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return result, xUserToken, err
		}
	}
	headerMap := map[string]string{}
	req.Header.Set("Language", o.Language)
	headerMap["Language"] = o.Language
	req.Header.Set("Content-Type", "application/json")
	headerMap["Content-Type"] = "application/json"
	if userToken != "" {
		req.Header.Set("X-User-Token", userToken)
		headerMap["X-User-Token"] = "******"
	} else {
		req.Header.Set("X-Access-Token", o.AccessToken)
		headerMap["X-Access-Token"] = "******"
	}

	headers := []string{}
	for key, val := range headerMap {
		header := fmt.Sprintf(`-H "%s: %s"`, key, val)
		headers = append(headers, header)
	}
	msgCurlParam := strings.Join(headers, " ")
	if strReqBody != "" {
		msgCurlParam = fmt.Sprintf("%s -d '%s'", msgCurlParam, strReqBody)
	}
	msgCurl := fmt.Sprintf(`curl -v -X%s %s '%s'`, method, msgCurlParam, url)
	log.Debug(msgCurl)

	resp, err = client.Do(req)
	if err != nil {
		return result, xUserToken, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, xUserToken, err
	}

	strJson = string(bs)
	result = gjson.Parse(strJson)

	strPrettyJson, err := PrettyJson(strJson)
	if err != nil {
		return result, xUserToken, err
	}

	log.Debug(fmt.Sprintf("%s %s %s in %s", method, url, resp.Status, result.Get("duration").String()))
	log.Debug(fmt.Sprintf("Response Header:"))
	for key, val := range resp.Header {
		log.Debug(fmt.Sprintf("  %s: %s", key, strings.Join(val, ",")))
	}
	log.Debug(fmt.Sprintf("Response Body:\n%s", strPrettyJson))

	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf("%s %s [%s] %s", method, url, result.Get("status").String(), result.Get("msg").String())
		return result, xUserToken, err
	}
	xUserToken = resp.Header.Get("X-User-Token")

	msg := fmt.Sprintf("%s %s [%s] %s", method, url, result.Get("status").String(), result.Get("msg").String())
	if showSuccess {
		log.Success(msg)
	} else {
		log.Debug(msg)
	}

	return result, xUserToken, err
}

func (o *OptionsCommon) QueryWebsocket(url, runName string) error {
	var err error
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var serverURL string
	if strings.HasPrefix(o.ServerURL, "http://") {
		serverURL = strings.Replace(o.ServerURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(o.ServerURL, "https://") {
		serverURL = strings.Replace(o.ServerURL, "https://", "wss://", 1)
	}
	if serverURL == "" {
		return err
	}

	if o.AccessToken == "" || o.ServerURL == "" {
		err = fmt.Errorf("please login first")
		return err
	}
	url = fmt.Sprintf("%s/%s", serverURL, url)

	header := http.Header{}
	header.Add("X-Access-Token", o.AccessToken)
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	conn, resp, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Debug(fmt.Sprintf("WEBSOCKET %s %s", url, resp.Status))

	go func(conn *websocket.Conn) {
		for {
			err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}(conn)

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			break
		}
		switch msgType {
		case websocket.TextMessage:
			if runName != "" {
				var msg pkg.WsRunLog
				err = json.Unmarshal(msgData, &msg)
				if err != nil {
					err = fmt.Errorf("parse msg error: %s", err.Error())
					return err
				}
				log.RunLog(msg)
				if msg.LogType == pkg.LogStatusInput {
					param := map[string]interface{}{}
					var r gjson.Result

					r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", runName), http.MethodGet, "", param, false)
					if err != nil {
						return err
					}
					run := pkg.Run{}
					err = json.Unmarshal([]byte(r.Get("data.run").Raw), &run)
					if err != nil {
						return err
					}
					if run.RunName == "" {
						err = fmt.Errorf("runName %s not exists", runName)
						return err
					}
					if run.Status.Duration == "" {
						r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s/input", runName), http.MethodGet, "", param, false)
						if err != nil {
							return err
						}
						var runInput pkg.RunInput
						err = json.Unmarshal([]byte(r.Get("data").Raw), &runInput)
						if err != nil {
							err = fmt.Errorf("parse run input error: %s", err.Error())
							return err
						}
						if runInput.PhaseID == msg.PhaseID {
							if runInput.IsApiOnly {
								log.Warning("# waiting for system callback to continue")
							} else {
								opts := []string{}
								for _, opt := range runInput.Options {
									opts = append(opts, opt.Value)
								}
								if len(opts) == 0 {
									opts = append(opts, pkg.InputValueConfirm, pkg.InputValueAbort)
								} else {
									opts = append(opts, pkg.InputValueAbort)
								}
								strOptions := strings.Join(opts, ",")
								log.Warning(fmt.Sprintf("# %s, %s", runInput.Title, runInput.Desc))
								log.Warning(fmt.Sprintf("# options: %s", strOptions))

								var inputValue string

								for {
									if inputValue == "" {
										if runInput.IsMultiple {
											log.Warning("# please input options (support multiple options, example: opt1,opt2)")
										} else {
											log.Warning("# please input option")
										}
										reader := bufio.NewReader(os.Stdin)
										inputValue, _ = reader.ReadString('\n')
										inputValue = strings.Trim(inputValue, "\n")
										inputValue = strings.Trim(inputValue, " ")
									} else {
										break
									}
								}

								param = map[string]interface{}{
									"phaseID":    runInput.PhaseID,
									"inputValue": inputValue,
								}
								r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s/input", runName), http.MethodPost, "", param, false)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			} else {
				var msg pkg.WsAdminLog
				err = json.Unmarshal(msgData, &msg)
				if err != nil {
					err = fmt.Errorf("parse msg error: %s", err.Error())
					return err
				}
				log.AdminLog(msg)
			}
		case websocket.CloseMessage:
			break
		default:
			break
		}
	}

	return err
}

func (o *OptionsCommon) GetProjectNames() ([]string, error) {
	var err error
	projectNames := []string{}
	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectNames"), http.MethodGet, "", param, false)
	if err != nil {
		return projectNames, err
	}
	err = json.Unmarshal([]byte(result.Get("data.projectNames").Raw), &projectNames)
	if err != nil {
		return projectNames, err
	}
	return projectNames, err
}

func (o *OptionsCommon) GetProjectDef(projectName string) (pkg.ProjectOutput, error) {
	var err error
	var project pkg.ProjectOutput

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", projectName), http.MethodGet, "", param, false)
	if err != nil {
		return project, err
	}
	err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
	if err != nil {
		return project, err
	}

	return project, err
}

func (o *OptionsCommon) GetPipelineNames() ([]string, error) {
	var err error
	var pipelineNames []string

	param := map[string]interface{}{
		"projectNames": []string{},
		"projectTeam":  "",
		"page":         1,
		"perPage":      1000,
	}
	result, _, err := o.QueryAPI("api/cicd/projects", http.MethodPost, "", param, false)
	if err != nil {
		return pipelineNames, err
	}
	rs := result.Get("data.projects").Array()
	projects := []pkg.Project{}
	for _, r := range rs {
		project := pkg.Project{}
		err = json.Unmarshal([]byte(r.Raw), &project)
		if err != nil {
			return pipelineNames, err
		}
		projects = append(projects, project)
	}
	for _, project := range projects {
		for _, pipeline := range project.Pipelines {
			pipelineNames = append(pipelineNames, pipeline.PipelineName)
		}
	}

	return pipelineNames, err
}

func (o *OptionsCommon) GetOpsBatchNames(projectName string) ([]string, error) {
	var err error
	var opsBatchNames []string

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/project/%s", projectName), http.MethodGet, "", param, false)
	if err != nil {
		return opsBatchNames, err
	}
	rs := result.Get("data.project.opsBatchDefs").Array()
	opsBatchDefs := []pkg.OpsBatchDef{}
	for _, r := range rs {
		opsBatchDef := pkg.OpsBatchDef{}
		err = json.Unmarshal([]byte(r.Raw), &opsBatchDef)
		if err != nil {
			return opsBatchNames, err
		}
		opsBatchDefs = append(opsBatchDefs, opsBatchDef)
	}
	for _, opsBatchDef := range opsBatchDefs {
		opsBatchNames = append(opsBatchNames, opsBatchDef.OpsBatchName)
	}

	return opsBatchNames, err
}

func (o *OptionsCommon) GetRunNames() ([]string, error) {
	var err error
	var runNames []string

	param := map[string]interface{}{
		"page":    1,
		"perPage": 200,
	}
	result, _, err := o.QueryAPI("api/cicd/runs", http.MethodPost, "", param, false)
	if err != nil {
		return runNames, err
	}
	rs := result.Get("data.runs").Array()
	runs := []pkg.Run{}
	for _, r := range rs {
		run := pkg.Run{}
		err = json.Unmarshal([]byte(r.Raw), &run)
		if err != nil {
			return runNames, err
		}
		runs = append(runs, run)
	}
	for _, run := range runs {
		runNames = append(runNames, run.RunName)
	}

	return runNames, err
}

func (o *OptionsCommon) GetUserNames() ([]string, error) {
	var err error
	var userNames []string

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI("api/admin/userNames", http.MethodGet, "", param, false)
	if err != nil {
		return userNames, err
	}
	rs := result.Get("data.userNames").Array()
	for _, r := range rs {
		userNames = append(userNames, r.String())
	}

	return userNames, err
}

func (o *OptionsCommon) GetStepNames() ([]string, error) {
	var err error
	var stepNames []string

	param := map[string]interface{}{
		"page":    1,
		"perPage": 1,
	}
	result, _, err := o.QueryAPI("api/admin/customStepConfs", http.MethodPost, "", param, false)
	if err != nil {
		return stepNames, err
	}
	rs := result.Get("data.customStepNames").Array()
	for _, r := range rs {
		stepNames = append(stepNames, r.String())
	}

	return stepNames, err
}

func (o *OptionsCommon) GetEnvNames() ([]string, error) {
	var err error
	var envNames []string

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI("api/admin/envNames", http.MethodGet, "", param, false)
	if err != nil {
		return envNames, err
	}
	rs := result.Get("data.envNames").Array()
	for _, r := range rs {
		envNames = append(envNames, r.String())
	}

	return envNames, err
}

func (o *OptionsCommon) GetComponentTemplateNames() ([]string, error) {
	var err error
	var componentTemplateNames []string

	param := map[string]interface{}{
		"page":    1,
		"perPage": 1000,
	}
	result, _, err := o.QueryAPI("api/admin/componentTemplates", http.MethodPost, "", param, false)
	if err != nil {
		return componentTemplateNames, err
	}
	rs := result.Get("data.componentTemplates").Array()
	componentTemplates := []pkg.ComponentTemplate{}
	for _, r := range rs {
		componentTemplate := pkg.ComponentTemplate{}
		err = json.Unmarshal([]byte(r.Raw), &componentTemplate)
		if err != nil {
			return componentTemplateNames, err
		}
		componentTemplates = append(componentTemplates, componentTemplate)
	}
	for _, componentTemplate := range componentTemplates {
		componentTemplateNames = append(componentTemplateNames, componentTemplate.ComponentTemplateName)
	}

	return componentTemplateNames, err
}

func (o *OptionsCommon) GetRepoNames() (pkg.RepoNameList, error) {
	var err error
	var repoNameList pkg.RepoNameList

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI("api/admin/repoNames", http.MethodGet, "", param, false)
	if err != nil {
		return repoNameList, err
	}

	err = json.Unmarshal([]byte(result.Get("data").Raw), &repoNameList)
	if err != nil {
		return repoNameList, err
	}

	return repoNameList, err
}

func (o *OptionsCommon) GetBuildEnvNames() ([]string, error) {
	var err error
	var envNames []string

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI("api/admin/buildEnvNames", http.MethodGet, "", param, false)
	if err != nil {
		return envNames, err
	}
	rs := result.Get("data.buildEnvNames").Array()
	for _, r := range rs {
		envNames = append(envNames, r.String())
	}

	return envNames, err
}

func (o *OptionsCommon) Complete(cmd *cobra.Command) error {
	var err error

	err = cmd.RegisterFlagCompletionFunc("language", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"zh", "en"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}
