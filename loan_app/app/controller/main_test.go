package controller_test

import (
	"os"
	"testing"

	"github.com/goProjects/loan_app/app/test_helper"
)

func TestMain(m *testing.M) {
	test_helper.InitializeLogger()
	test_helper.SetupDatabase()

	exitCode := m.Run()

	test_helper.TeardownDatabase()

	os.Exit(exitCode)
}
