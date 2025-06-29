/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/jadudm/fac-archive/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fac-archive",
	Short: "A tool for archiving and updating Federal Audit Clearinghouse data",
	Long: `

Example use:

Archive all data in the FAC:

fac-archive --sqlite fac.db 

Update everything from the past work week:

fac-update --start-date 2025-03-03 --end-date 2025-03-08 --sqlite fac.db

Download the PDFs for Monday, March 3rd:

fac-update reports --start-date 2025-03-03 --end-date 2025-03-04 --sqlite fac.db
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.Init()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.com.jadud.faccopy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("sqlite", "", "SQLite archive file")
	rootCmd.MarkPersistentFlagRequired("sqlite")
	viper.BindPFlag("sqlite", rootCmd.PersistentFlags().Lookup("sqlite"))

}
