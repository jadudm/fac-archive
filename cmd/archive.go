/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jadudm/fac-archive/internal/archivedb"
	"github.com/jadudm/fac-archive/internal/fac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func archiveTable(table string, db *sql.DB, Q *archivedb.Queries) (int, error) {
	rows_retrieved := 0
	var limit_per_query int
	if viper.GetInt("api.limit_per_query") != 0 {
		limit_per_query = viper.GetInt("api.limit_per_query")
	} else {
		limit_per_query = fac.LimitPerQuery
	}

	for offset := 0; offset <= fac.MaxRows; offset += limit_per_query {
		// The query URL communicates the limit and offset, so that we
		// walk the entire dataset in a windowed manner. 0-20K, 20001-40K, etc.
		url := fmt.Sprintf("%s://%s/%s?limit=%d&offset=%d",
			viper.GetString("api.scheme"),
			viper.GetString("api.url"),
			table,
			limit_per_query,
			offset,
		)

		zap.L().Info("fetching",
			zap.String("accept_profile", viper.GetString("api.accept_profile")),
			zap.String("url", url))

		// Fetch the JSON body
		body, err := fac.FacGet(url)
		if err != nil {
			zap.L().Fatal("could not fetch body. exiting.")
		}

		// To make this generic, the tables are pulled as untyped hashmaps.
		objects := make([]map[string]any, 0)
		jsonErr := json.Unmarshal(body, &objects)
		if jsonErr != nil {
			zap.L().Fatal("could not unmarshal response", zap.Error(jsonErr))
		}

		// When we hit zero objects, we've pulled everything there is to pull
		if len(objects) == 0 {
			// If we get nothing back, we're done.
			return rows_retrieved, nil
		} else {
			// Count the rows retrieved.
			rows_retrieved += len(objects)

			ctx := context.Background()

			// Using transactions lets us insert 20K rows at a time.
			// Doing this means 350K rows come down in ~2m. Otherwise, it takes
			// many, many, many hours when inserts are row-by-row.
			tx, err := db.Begin()
			if err != nil {
				zap.L().Error("could not initialize transaction", zap.Error(err))
			}

			defer tx.Rollback()
			qtx := Q.WithTx(tx)
			for _, g := range objects {
				archivedb.RawJsonInsert(table, qtx, ctx, g)
			}
			// Commit after building up 20K inserts.
			tx.Commit()
		}
	}

	// Should not get here.
	return -1, fmt.Errorf("could not archive %s", table)
}

func archive(cmd *cobra.Command, args []string) {

	// var db_name string
	// REQUIRED FLAGS
	// cmd.Flags().StringVar(&db_name, "sqlite", "", "SQLite database name")
	// cmd.MarkFlagRequired("sqlite")

	// OPTIONAL FLAGS
	// var api_version string
	// cmd.Flags().StringVar(&api_version, "api_version", "api_v1_1_0", "FAC API version")

	if _, err := os.Stat(viper.GetString("sqlite")); err == nil {
		zap.L().Fatal("exiting; database already exists")
	} else if errors.Is(err, os.ErrNotExist) {

		db, queries, err := archivedb.CreateSqliteDB(viper.GetString("sqlite"))
		if err != nil {
			zap.L().Fatal("could not create database. exiting.")
		}

		// Set the internal flag so we know if we should run the triggers
		// that copy the JSON to structured tables
		ctx := context.Background()
		queries.AddMetadata(ctx, archivedb.AddMetadataParams{
			Key:   "copy_json",
			Value: viper.GetString("copy_json"),
		})

		for _, table := range fac.Tables {
			start := time.Now()
			rows, err := archiveTable(table, db, queries)
			elapsed := time.Since(start)

			if err != nil {
				zap.L().Error("could not archive table", zap.String("table", table), zap.Error(err))
			}

			zap.L().Info("rows retrieved",
				zap.String("table", table),
				zap.Int("rows", rows),
				zap.Int64("duration", int64(elapsed.Seconds())))
		}

		ctx = context.Background()
		queries.AddMetadata(ctx, archivedb.AddMetadataParams{
			Key:   "last_updated",
			Value: time.Now().Format("2006-01-02"),
		})

	} else {
		zap.L().Fatal("Does the database exist? Does it not? I cannot tell. Exiting.")
	}
}

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archives all data from the FAC",
	Long: `When launched, this command creates an archive of all data in the FAC. 
This must be run before any other commands (update, pdfs) can be run.

Usage:

fac-archive archive --sqlite <filename>

The --sqlite flag is required, naming the database that will be created and written to.
`,
	Run: archive,
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
