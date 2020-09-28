package cmd

import (
	"context"
	"fmt"

	"github.com/goProjects/simplpay/simplpay"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "commands related to user",
	Long:  `different commands related to user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user called")
	},
}

var userName string
var userEmail string
var userCreditLimit float64

func init() {
	rootCmd.AddCommand(userCmd)

	createUserCmd.Flags().StringVarP(&userName, "name", "n", "", "user name")
	createUserCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	createUserCmd.Flags().Float64VarP(&userCreditLimit, "creditLimit", "c", 0, "user credit limit")

	userCmd.AddCommand(
		createUserCmd,
	)
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a user with a given credit limit",
	Run:   createUser,
}

func createUser(command *cobra.Command, args []string) {
	u := &simplpay.User{
		Name:        userName,
		Email:       userEmail,
		CreditLimit: userCreditLimit,
	}
	if err := simplpay.CreateUser(context.Background(), s, u); err != nil {
		fmt.Println(err)
	}
}
