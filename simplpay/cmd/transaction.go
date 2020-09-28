package cmd

import (
	"context"
	"fmt"

	"github.com/goProjects/simplpay/simplpay"
	"github.com/spf13/cobra"
)

// transactionCmd represents the transaction command
var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "commands related to transaction",
	Long:  `commands related to transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("transaction called")
	},
}

var transactionAmount float64

func init() {
	rootCmd.AddCommand(transactionCmd)

	createTransactionCmd.Flags().StringVarP(&userName, "userName", "u", "", "user name")
	createTransactionCmd.Flags().StringVarP(&merchantName, "merchantName", "m", "", "merchant name")
	createTransactionCmd.Flags().Float64VarP(&transactionAmount, "amount", "a", 0, "amount given by user to merchant")

	transactionCmd.AddCommand(
		createTransactionCmd,
	)
}

var createTransactionCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a tx",
	Run:   createTransaction,
}

func createTransaction(command *cobra.Command, args []string) {
	u := &simplpay.Transaction{
		UserName:     userName,
		MerchantName: merchantName,
		Amount:       transactionAmount,
	}
	if err := simplpay.CreateTransaction(context.Background(), s, u); err != nil {
		fmt.Println(err)
	}
}
