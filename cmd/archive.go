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

func archiveTable(cmd *cobra.Command, table string, db *sql.DB, Q *archivedb.Queries) (int, error) {
	rows_retrieved := 0

	for offset := 0; offset <= fac.MaxRows; offset += fac.LimitPerQuery {
		// The query URL communicates the limit and offset, so that we
		// walk the entire dataset in a windowed manner. 0-20K, 20001-40K, etc.
		url := fmt.Sprintf("%s://%s/%s?limit=%d&offset=%d",
			viper.GetString("api.scheme"),
			viper.GetString("api.url"),
			table,
			fac.LimitPerQuery,
			offset,
		)

		zap.L().Info("fetching", zap.String("url", url))

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
	return -1, errors.New(fmt.Sprintf("could not archive %s", table))
}

func archive(cmd *cobra.Command, args []string) {

	db_name := cmd.Flag("sqlite").Value.String()
	// db_name := fmt.Sprintf("%s-fac.sqlite", current.Format("2006-01-02-15-04-05"))

	if _, err := os.Stat(db_name); err == nil {
		zap.L().Fatal("exiting; database already exists")
	} else if errors.Is(err, os.ErrNotExist) {

		db, queries, err := archivedb.CreateSqliteDB(db_name)
		if err != nil {
			zap.L().Fatal("could not create database. exiting.")
		}

		for _, table := range fac.Tables {
			start := time.Now()
			rows, err := archiveTable(cmd, table, db, queries)
			elapsed := time.Since(start)

			if err != nil {
				zap.L().Error("could not archive table", zap.String("table", table), zap.Error(err))
			}

			zap.L().Info("rows retrieved",
				zap.String("table", table),
				zap.Int("rows", rows),
				zap.Int64("duration", int64(elapsed.Seconds())))
		}

		ctx := context.Background()
		queries.AddMetadata(ctx, archivedb.AddMetadataParams{
			Key:   "last_updated",
			Value: time.Now().Format("2006-01-02"),
		})
	} else {
		zap.L().Fatal("Does the file exist? Does it not? I cannot tell. Exiting.")
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
