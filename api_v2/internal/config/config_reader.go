package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"techytechster.com/roastedoctocats/pkg/iconfig"
)

type jsonConfigReader struct {
}

func (j *jsonConfigReader) GetModelPrompt() string {
	return modelPrompt
}

func getModuleRoot() (string, error) {
	if path, ok := os.LookupEnv("BASEPATH"); ok {
		return path, nil
	}
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

type JsonConfig struct {
	TestPem                 string `json:"testPem"`
	TestPersonalAccessToken string `json:"testPat"`
	GithubAppSecret         string `json:"githubAppSecret"`
}

func (j JsonConfig) GetGithubSecret() string {
	return j.GithubAppSecret
}

func (j JsonConfig) GetPAT() string {
	return j.TestPersonalAccessToken
}

func (j JsonConfig) GetTestPem() string {
	return j.TestPem
}

func (j *jsonConfigReader) ReadConfiguration(environment string) iconfig.ApplicationConfig {
	var jsonConfig JsonConfig
	rootPath, err := getModuleRoot()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(fmt.Sprintf("%s/configs/application/%s.json", rootPath, environment))
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &jsonConfig)
	return jsonConfig
}

func New() iconfig.ConfigReader {
	return new(jsonConfigReader)
}
