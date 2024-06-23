package utils

type Environment string

const (
	DevEnv        Environment = "development"
	TestingEnv    Environment = "test"
	StagingEnv    Environment = "staging"
	UnicornEnv    Environment = "unicorn"
	SandboxEnv    Environment = "sandbox"
	ProductionEnv Environment = "production"
	QAEnv         Environment = "qa"
)
