package cmd

import (
	"context"
	"fmt"

	"github.com/goProjects/simplpay/simplpay"
	"github.com/spf13/cobra"
)

// merchantCmd represents the merchant command
var merchantCmd = &cobra.Command{
	Use:   "merchant",
	Short: "commands related to merchant",
	Long:  `different commands related to merchant.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("merchant called")
	},
}

var merchantName string
var merchatDiscount float64

func init() {
	rootCmd.AddCommand(merchantCmd)

	createMerchantCmd.Flags().StringVarP(&merchantName, "name", "n", "", "merchant name")
	createMerchantCmd.Flags().Float64VarP(&merchatDiscount, "discountPercent", "d", 0, "merchant discount percent")

	merchantCmd.AddCommand(
		createMerchantCmd,
	)
}

var createMerchantCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a merchant with a given discount",
	Run:   createMerchant,
}

func createMerchant(command *cobra.Command, args []string) {
	m := &simplpay.Merchant{
		Name:            merchantName,
		DiscountPercent: merchatDiscount,
	}
	if err := simplpay.CreateMerchant(context.Background(), s, m); err != nil {
		fmt.Println(err)
	}
}
