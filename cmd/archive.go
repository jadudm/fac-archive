/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jadudm/fac-tool/internal/config"
	"github.com/jadudm/fac-tool/internal/fac"
	"github.com/jadudm/fac-tool/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func create_sqlite_db(db_name string) (*sql.DB, *sqlite.Queries, error) {
	db, queries, err := sqlite.CreateTables(db_name)
	return db, queries, err
}

func fac_get(url string) []byte {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("could not initialize new request", zap.Error(err))
	}

	req.Header = http.Header{
		"X-API-Key":      {viper.GetString("api.key")},
		"Accept-Profile": {"api_v1_1_0"},
	}

	res, err := client.Do(req)
	if err != nil {
		//Handle Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func archive_table(table string, max int, Q *sqlite.Queries) {

	for offset := 0; offset <= max; offset += fac.Limit {
		url := fmt.Sprintf("%s://%s/%s?limit=%d&offset=%d",
			viper.GetString("api.scheme"),
			viper.GetString("api.url"),
			table,
			fac.Limit,
			offset,
		)
		body := fac_get(url)

		// generals := []fac.General{}
		objects := make([]map[string]any, 0)
		jsonErr := json.Unmarshal(body, &objects)
		if jsonErr != nil {
			zap.L().Fatal("could not unmarshal response", zap.Error(jsonErr))
		}

		if len(objects) == 0 {
			break
		} else {
			ctx := context.Background()

			for _, g := range objects {
				b, err := json.Marshal(g)
				if err != nil {
					zap.L().Error("could not marshal to string", zap.Error(err))
				}
				// zap.L().Info("inserting raw", zap.String("json", string(b)))
				id, err := Q.RawInsert(ctx, sqlite.RawInsertParams{
					Source: table,
					Json:   string(b),
				})
				if err != nil {
					zap.L().Error("could not insert", zap.Error(err))
				} else {
					zap.L().Debug("inserted id", zap.Int64("id", id))
				}
			}
		}
	}
}

func archive(cmd *cobra.Command, args []string) {
	config.Init()

	current := time.Now()
	db_name := fmt.Sprintf("%s-fac.sqlite", current.Format("2006-01-02-15-04-05"))

	// fmt.Println("archive called")
	// fmt.Printf("api url: %s\n", viper.GetString("api.url"))

	_, queries, err := create_sqlite_db(db_name)
	if err != nil {
		zap.L().Fatal("could not create database. exiting.")
	}

	archive_table("general", 20, queries)
	archive_table("federal_awards", 20, queries)
	archive_table("notes_to_sefa", 20, queries)
	archive_table("findings", 20, queries)

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
