package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dorycli/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type OptionsLogin struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Username       string `yaml:"username" json:"username" bson:"username" validate:""`
	Password       string `yaml:"password" json:"password" bson:"password" validate:""`
	ExpireDays     int    `yaml:"expireDays" json:"expireDays" bson:"expireDays" validate:""`
}

func NewOptionsLogin() *OptionsLogin {
	var o OptionsLogin
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdLogin() *cobra.Command {
	o := NewOptionsLogin()

	baseName := pkg.GetCmdBaseName()
	msgUse := fmt.Sprintf("login")
	msgShort := fmt.Sprintf("login to dory-engine server")
	msgLong := fmt.Sprintf("login first before use %s to control your dory-engine server, it will save dory-engine server settings in %s config file", baseName, baseName)
	msgExample := fmt.Sprintf(`  # login with username and password input prompt
  %s login --serverURL http://dory.example.com:8080

  # login without password input prompt
  %s login --serverURL http://dory.example.com:8080 --username test-user

  # login without input prompt
  %s login --serverURL http://dory.example.com:8080 --username test-user --password xxx

  # login with access token
  %s login --serverURL http://dory.example.com:8080 --token xxx`, baseName, baseName, baseName, baseName)

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
	cmd.Flags().StringVarP(&o.Username, "username", "U", "", "dory-engine server username")
	cmd.Flags().StringVarP(&o.Password, "password", "P", "", "dory-engine server password")
	cmd.Flags().IntVar(&o.ExpireDays, "expireDays", 90, "dory-engine server token expires days")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsLogin) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsLogin) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		err = fmt.Errorf("command args must be empty")
		return err
	}
	if o.ServerURL == "" {
		err = fmt.Errorf("--serverURL required")
		return err
	}
	if !strings.HasPrefix(o.ServerURL, "http://") && strings.HasPrefix(o.ServerURL, "https://") {
		err = fmt.Errorf("--serverURL must start with http:// or https://")
		return err
	}
	if o.ExpireDays < 0 {
		err = fmt.Errorf("--expireDays can not less than 0")
		return err
	}

	return err
}

func (o *OptionsLogin) Run(args []string) error {
	var err error

	var accessToken string
	var accessTokenExpires string
	if o.AccessToken == "" {
		if o.Password != "" {
			log.Warning("set password in command line args is not safe!")
		}
		for {
			if o.Username == "" {
				log.Info("please input username")
				reader := bufio.NewReader(os.Stdin)
				username, _ := reader.ReadString('\n')
				username = strings.Trim(username, "\n")
				username = strings.Trim(username, " ")
				o.Username = username
			} else {
				break
			}
		}
		for {
			if o.Password == "" {
				log.Info("please input password")
				bytePassword, _ := terminal.ReadPassword(0)
				password := string(bytePassword)
				password = strings.Trim(password, " ")
				o.Password = password
			} else {
				break
			}
		}

		bs, _ := pkg.YamlIndent(o)
		log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

		param := map[string]interface{}{
			"username": o.Username,
			"password": o.Password,
		}
		_, xUserToken, err := o.QueryAPI("api/public/login", http.MethodPost, "", param, true)
		if err != nil {
			return err
		}

		param = map[string]interface{}{}
		result, _, err := o.QueryAPI("api/account/accessTokens", http.MethodGet, xUserToken, param, true)
		if err != nil {
			return err
		}
		accessTokens := result.Get("data.accessTokens").Array()
		tokens := []pkg.AccessToken{}
		expireTimeHalf := time.Now().Add(time.Hour * time.Duration(o.ExpireDays*12)).Format("2006-01-02 15:04:05")
		expireTime := time.Now().Add(time.Hour * time.Duration(o.ExpireDays*24)).Format("2006-01-02 15:04:05")
		for _, tk := range accessTokens {
			var account pkg.AccessToken
			err = json.Unmarshal([]byte(tk.Raw), &account)
			if err != nil {
				err = fmt.Errorf("parse accessToken error: %s", err.Error())
				return err
			}
			if !account.Expired {
				if expireTimeHalf < account.ExpireTime && expireTime > account.ExpireTime {
					tokens = append(tokens, account)
				}
			}
		}
		sort.SliceStable(tokens, func(i, j int) bool {
			return tokens[i].ExpireTime > tokens[j].ExpireTime
		})
		if len(tokens) > 0 {
			accessToken = tokens[0].AccessToken
			accessTokenExpires = tokens[0].ExpireTime
		}

		if accessToken == "" {
			baseName := pkg.GetCmdBaseName()
			accessTokenName := fmt.Sprintf("%s-%s", baseName, time.Now().Format("20060102030405"))
			param = map[string]interface{}{
				"accessTokenName": accessTokenName,
				"expireDays":      o.ExpireDays,
			}
			result, _, err = o.QueryAPI("api/account/accessToken", http.MethodPost, xUserToken, param, true)
			if err != nil {
				return err
			}
			accessToken = result.Get("data.accessToken").String()
			accessTokenExpires = time.Now().Add(time.Hour * time.Duration(o.ExpireDays*24)).Format("2006-01-02 15:04:05")
			if accessToken == "" {
				err = fmt.Errorf("get accessToken error: accessToken is empty")
				return err
			}
		}
		if accessToken == "" {
			err = fmt.Errorf("get accessToken error: accessToken is empty")
			return err
		}
	} else {
		param := map[string]interface{}{}
		result, _, err := o.QueryAPI("api/account/accessTokens", http.MethodGet, "", param, true)
		// if get accessTokens failed, clear .dorycli config file
		if err != nil {
			doryConfig := pkg.DoryConfig{
				ServerURL:   "",
				Insecure:    o.Insecure,
				Timeout:     o.Timeout,
				AccessToken: "",
				Language:    o.Language,
			}
			bs, _ := pkg.YamlIndent(doryConfig)
			e := os.WriteFile(o.ConfigFile, bs, 0600)
			if e != nil {
				err = fmt.Errorf("%s\nwrite config file error: %s", err.Error(), e.Error())
			}
			return err
		}
		accessTokens := result.Get("data.accessTokens").Array()
		for _, tk := range accessTokens {
			var account pkg.AccessToken
			err = json.Unmarshal([]byte(tk.Raw), &account)
			if err != nil {
				err = fmt.Errorf("parse accessToken error: %s", err.Error())
				return err
			}
			if account.AccessToken == o.AccessToken && !account.Expired {
				accessToken = account.AccessToken
				accessTokenExpires = account.ExpireTime
				break
			}
		}
		if accessToken == "" {
			err = fmt.Errorf("get accessToken error: accessToken is empty")
			return err
		}
	}

	accessTokenBase64 := base64.StdEncoding.EncodeToString([]byte(accessToken))
	o.AccessToken = accessTokenBase64
	doryConfig := pkg.DoryConfig{
		ServerURL:   o.ServerURL,
		Insecure:    o.Insecure,
		Timeout:     o.Timeout,
		AccessToken: o.AccessToken,
		Language:    o.Language,
	}
	bs, _ := pkg.YamlIndent(doryConfig)
	err = os.WriteFile(o.ConfigFile, bs, 0600)
	if err != nil {
		return err
	}

	log.Success("login success")
	log.Success(fmt.Sprintf("access token expires time: %s", accessTokenExpires))
	log.Debug(fmt.Sprintf("update %s success", o.ConfigFile))

	return err
}
