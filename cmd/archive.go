/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jadudm/fac-tool/internal/config"
	"github.com/jadudm/fac-tool/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func create_sqlite_db() {
	current := time.Now()
	db_name := fmt.Sprintf("%s-fac.sqlite", current.Format("2006-01-02-15-04-05"))
	_, queries := sqlite.CreateTables(db_name)
	ctx := context.Background()
	gs, err := queries.GetGenerals(ctx)
	if err != nil {
		zap.L().Error("could not fetch generals", zap.Error(err))
	}
	log.Println(gs)
}

func archive(cmd *cobra.Command, args []string) {
	config.Init()
	fmt.Println("archive called")
	fmt.Printf("api url: %s\n", viper.GetString("api.url"))
	create_sqlite_db()
}

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: archive,
}

func init() {
	rootCmd.AddCommand(archiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// archiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// archiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
