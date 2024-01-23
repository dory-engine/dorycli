package pkg

import (
	"bytes"
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
)

func GetCmdBaseName() string {
	args := os.Args
	arr := strings.Split(args[0], "/")
	baseName := arr[len(arr)-1]
	return baseName
}

func ExtractEmbedFile(f embed.FS, rootDir string, targetDir string) error {
	return fs.WalkDir(f, rootDir, func(path string, d fs.DirEntry, err error) error {
		rootDir = strings.TrimSuffix(rootDir, "/")
		if err != nil {
			return err
		}
		if path != "." && path != rootDir {
			pathTarget := fmt.Sprintf("%s/%s", targetDir, strings.TrimPrefix(path, fmt.Sprintf("%s/", rootDir)))
			if d.IsDir() {
				_ = os.MkdirAll(pathTarget, 0700)
			} else {
				bs, err := f.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR:", err.Error())
					return err
				}
				_ = os.WriteFile(pathTarget, bs, 0600)
			}
		}
		return nil
	})
}

func CommandExec(command, workDir string) (string, string, error) {
	var err error
	errInfo := fmt.Sprintf("exec %s error", command)
	var strOut, strErr string

	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = workDir

	prOut, pwOut := io.Pipe()
	prErr, pwErr := io.Pipe()
	execCmd.Stdout = pwOut
	execCmd.Stderr = pwErr

	rOut := io.TeeReader(prOut, os.Stdout)
	rErr := io.TeeReader(prErr, os.Stderr)

	err = execCmd.Start()
	if err != nil {
		err = fmt.Errorf("%s: exec start error: %s", errInfo, err.Error())
		return strOut, strErr, err
	}

	var bOut, bErr bytes.Buffer

	go func() {
		_, _ = io.Copy(&bOut, rOut)
	}()

	go func() {
		_, _ = io.Copy(&bErr, rErr)
	}()

	err = execCmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("%s: exit status: %d", errInfo, status.ExitStatus())
			}
		} else {
			err = fmt.Errorf("%s: exec run error: %s", errInfo, err.Error())
			return strOut, strErr, err
		}
	}

	strOut = bOut.String()
	strErr = bErr.String()

	return strOut, strErr, err
}

func CheckRandomStringStrength(password string, length int, enableSpecialChar bool) error {
	var err error

	if len(password) < length {
		err = fmt.Errorf("password must at least %d charactors", length)
		return err
	}
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars := "0123456789"
	specialChars := `~!@#$%^&*()_+-={}[]\|:";'<>?,./`
	lowerOK := strings.ContainsAny(password, lowerChars)
	upperOK := strings.ContainsAny(password, upperChars)
	numberOK := strings.ContainsAny(password, numberChars)
	specialOK := strings.ContainsAny(password, specialChars)
	for _, c := range password {
		if c > unicode.MaxASCII {
			err = fmt.Errorf("password can not include unicode charactors")
			return err
		}
	}
	if enableSpecialChar && !(lowerOK && upperOK && numberOK && specialOK) {
		err = fmt.Errorf("password must include lower upper case charactors and number and special charactors")
		return err
	} else if !enableSpecialChar && !(lowerOK && upperOK && numberOK) {
		err = fmt.Errorf("password must include lower upper case charactors and number")
		return err
	}

	return err
}

func RandomString(n int, enableSpecialChar bool, suffix string) string {
	var letter []rune
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars := "0123456789"
	specialChars := `~!@#$%^&*()_+-={}[]\|:";'<>?,./`
	if enableSpecialChar {
		chars := fmt.Sprintf("%s%s%s%s", lowerChars, upperChars, numberChars, specialChars)
		letter = []rune(chars)
	} else {
		chars := fmt.Sprintf("%s%s%s", lowerChars, upperChars, numberChars)
		letter = []rune(chars)
	}
	var pwd string
	for {
		b := make([]rune, n)
		seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = letter[seededRand.Intn(len(letter))]
		}
		pwd = string(b)
		err := CheckRandomStringStrength(pwd, n, enableSpecialChar)
		if err == nil {
			break
		}
	}
	return fmt.Sprintf("%s%s", pwd, suffix)
}

func ValidateIpAddress(s string) error {
	var err error
	if net.ParseIP(s).To4() == nil {
		err = fmt.Errorf(`not ipv4 address`)
		return err
	}
	return err
}

func ValidateMinusNameID(s string) error {
	var err error
	RegExp := regexp.MustCompile(`^(([a-z])[a-z0-9]+(-[a-z0-9]+)+)$`)
	match := RegExp.MatchString(s)
	if !match {
		err = fmt.Errorf(`should include lower case and number, format should like "hello-world-no3"`)
		return err
	}
	return err
}

func ValidatePipelineName(s string) error {
	var err error
	RegExp := regexp.MustCompile(`^(([a-z])[a-z0-9]+(-[a-z0-9_]+)+)$`)
	match := RegExp.MatchString(s)
	if !match {
		err = fmt.Errorf(`pipeline name format should like "test-project1-develop"`)
		return err
	}
	arr := strings.Split(s, "-")
	for i, s2 := range arr {
		if i < len(arr)-1 && strings.Contains(s2, "_") {
			err = fmt.Errorf(`project name can not contains "_"`)
			return err
		} else if i == len(arr)-1 {
			if strings.HasPrefix(s2, "_") {
				err = fmt.Errorf(`branch name can not starts with "_"`)
				return err
			}
			if strings.HasSuffix(s2, "_") {
				err = fmt.Errorf(`branch name can not ends with "_"`)
				return err
			}
		}
	}
	return err
}

func ValidateRunName(s string) error {
	var err error
	arr := strings.Split(s, "-")
	if len(arr) < 3 {
		err = fmt.Errorf(`run name format should like "test-project1-develop-1"`)
		return err
	}
	err = ValidatePipelineName(strings.Join(arr[:len(arr)-1], "-"))
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(arr[len(arr)-1])
	if err != nil {
		err = fmt.Errorf(`run number must be integer`)
		return err
	}
	return err
}

func RemoveMapEmptyItems(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			if vv.IsZero() {
				delete(m, k)
			}
		case reflect.Bool:
			if vv.Bool() == false {
				delete(m, k)
			}
		case reflect.String:
			if vv.String() == "" {
				delete(m, k)
			}
		case reflect.Slice, reflect.Array:
			if vv.Len() == 0 {
				delete(m, k)
			} else {
				var isMap bool
				var x []map[string]interface{}
				for i := 0; i < vv.Len(); i++ {
					vvv := reflect.ValueOf(vv.Index(i))
					if vvv.Kind() == reflect.Map {
						vm, ok := vv.Index(i).Interface().(map[string]interface{})
						if ok {
							isMap = true
							v3 := RemoveMapEmptyItems(vm)
							x = append(x, v3)
						}
					} else if vvv.Kind() == reflect.Struct {
						vm, ok := vv.Index(i).Interface().(map[string]interface{})
						if ok {
							isMap = true
							v3 := RemoveMapEmptyItems(vm)
							x = append(x, v3)
						}
					}
				}
				if isMap {
					m[k] = x
				}
			}
		case reflect.Struct, reflect.Map:
			v2 := RemoveMapEmptyItems(v.(map[string]interface{}))
			if len(v2) == 0 {
				delete(m, k)
			} else {
				m[k] = v2
			}
		default:
			if !vv.IsValid() {
				delete(m, k)
			}
		}
	}
	m2 := m
	return m2
}

func YamlIndent(obj interface{}) ([]byte, error) {
	var err error
	var bs []byte
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(&obj)
	if err != nil {
		return bs, err
	}
	bs = b.Bytes()

	return bs, err
}

func GetDockerImages(installConfig InstallConfig) (InstallDockerImages, error) {
	var err error
	var bs []byte
	var dockerImages InstallDockerImages
	// get pull container images
	bs, err = FsInstallScripts.ReadFile(fmt.Sprintf("%s/harbor/docker-images.yaml", DirInstallScripts))
	if err != nil {
		err = fmt.Errorf("get pull container images error: %s", err.Error())
		return dockerImages, err
	}
	err = yaml.Unmarshal(bs, &dockerImages)
	if err != nil {
		err = fmt.Errorf("get pull container images error: %s", err.Error())
		return dockerImages, err
	}

	var image string
	var dockerImage InstallDockerImage
	image = fmt.Sprintf("doryengine/dorycli:%s-alpine", VersionDoryCli)
	dockerImage = InstallDockerImage{
		Source: image,
		Target: fmt.Sprintf("hub/%s", image),
	}
	dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	if installConfig.Dory.GitRepo.Internal.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.GitRepo.Internal.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.GitRepo.Internal.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.GitRepo.Internal.ImageDB != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.GitRepo.Internal.ImageDB,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.GitRepo.Internal.ImageDB),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.GitRepo.Internal.ImageNginx != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.GitRepo.Internal.ImageNginx,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.GitRepo.Internal.ImageNginx),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.ArtifactRepo.Internal.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.ArtifactRepo.Internal.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.ArtifactRepo.Internal.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Openldap.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.Openldap.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.Openldap.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Openldap.ImageAdmin != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.Openldap.ImageAdmin,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.Openldap.ImageAdmin),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.ScanCodeRepo.Internal.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.ScanCodeRepo.Internal.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.ScanCodeRepo.Internal.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.ScanCodeRepo.Internal.ImageDB != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.ScanCodeRepo.Internal.ImageDB,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.ScanCodeRepo.Internal.ImageDB),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Redis.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.Redis.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.Redis.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Mongo.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.Mongo.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.Mongo.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.DemoDatabase.Internal.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.DemoDatabase.Internal.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.DemoDatabase.Internal.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.DemoHost.Internal.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.DemoHost.Internal.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.DemoHost.Internal.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Docker.Image != "" {
		dockerImage = InstallDockerImage{
			Source: installConfig.Dory.Docker.Image,
			Target: fmt.Sprintf("hub/%s", installConfig.Dory.Docker.Image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Doryengine.Port != 0 {
		image = fmt.Sprintf("doryengine/dory-engine:%s-alpine", VersionDoryEngine)
		dockerImage = InstallDockerImage{
			Source: image,
			Target: fmt.Sprintf("hub/%s", image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	if installConfig.Dory.Doryengine.Port != 0 {
		image = fmt.Sprintf("doryengine/dory-console:%s", VersionDoryFrontend)
		dockerImage = InstallDockerImage{
			Source: image,
			Target: fmt.Sprintf("hub/%s", image),
		}
		dockerImages.InstallDockerImages = append(dockerImages.InstallDockerImages, dockerImage)
	}
	return dockerImages, err
}

func Indent(s string) string {
	var str string
	arr := strings.Split(s, "\n")
	for i, a := range arr {
		arr[i] = fmt.Sprintf("%s%s", "  ", a)
	}
	str = strings.Join(arr, "\n")
	return str
}
