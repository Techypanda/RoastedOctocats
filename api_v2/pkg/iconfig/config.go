package iconfig

type ApplicationConfig interface {
	GetTestPem() string
	GetPAT() string
	GetGithubSecret() string
}

type ConfigReader interface {
	GetModelPrompt() string
	ReadConfiguration(environment string) ApplicationConfig
}
