package main

import (
	"github.com/goProjects/simplpay/cmd"
)

func main() {
	cmd.Execute()
	defer cmd.Cleanup()
}
