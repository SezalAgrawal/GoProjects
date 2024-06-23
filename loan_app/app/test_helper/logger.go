package test_helper

import (
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
)

func InitializeLogger() {
	logger.Init(logger.FATAL, utils.TestingEnv)
}
