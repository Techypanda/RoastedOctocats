package iconfig

import "techytechster.com/roastedoctocats/pkg/proto"

type ApplicationConfig interface {
	GetTestPem() string
	GetPAT() string
	GetGithubSecret() string
}

type ConfigReader interface {
	GetModelPrompt(promptType proto.ModelPromptType) string
	ReadConfiguration(environment string) ApplicationConfig
}
