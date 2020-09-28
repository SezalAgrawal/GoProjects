package cmd

import (
	"fmt"
	"os"

	"github.com/goProjects/simplpay/db"
	"github.com/goProjects/simplpay/simplpay"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	database db.Provider
	s        *simplpay.Service
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simplpay",
	Short: "simplpay commands",
	Long:  `long description`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mustInitDB() {
	url := "mongodb://localhost:27017"
	mainDB := "simplpay"
	database = db.NewProvider(url, mainDB)
	if ok, err := database.Ok(); !ok {
		panic(fmt.Sprintf("can not establish connection to database: %s reason: %s", mainDB, err))
	}
	database.MustCreateIndexes()
}

// Cleanup closes db connection
func Cleanup() {
	if err := database.Close(); err != nil {
		log.Errorf("error while closing database connection: %s", err)
	}
}

func init() {
	mustInitDB()
	s = &simplpay.Service{
		DB: database,
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simplpay.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".simplpay" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".simplpay")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
