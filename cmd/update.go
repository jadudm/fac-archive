/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jadudm/fac-tool/internal/archivedb"
	"github.com/jadudm/fac-tool/internal/fac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func doUpdate(cmd *cobra.Command, db *sql.DB, Q *archivedb.Queries) (int, int, int) {
	// fmt.Println("update called")
	// fmt.Println(cmd.Flag("sqlite").Value)
	table := "general"
	rows_retrieved := 0
	inserted := 0
	already_present := 0

	url := fmt.Sprintf("%s://%s/%s?fac_accepted_date=gte.%s",
		viper.GetString("api.scheme"),
		viper.GetString("api.url"),
		table,
		cmd.Flag("fac-accepted-date").Value.String(),
	)

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
		return rows_retrieved, already_present, inserted
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

		// ctx := context.Background()
		for _, g := range objects {
			// Is the Report ID already present?
			ctx2 := context.Background()
			y_or_n, err := Q.ReportIdExists(ctx2, g["report_id"].(string))
			if err != nil {
				zap.L().Error("could not look up report id", zap.Error(err))
			}
			// If the report ID is present
			if y_or_n == 1 {
				already_present += 1
				zap.L().Debug("present", zap.String("report_id", g["report_id"].(string)))
			} else {
				inserted += 1
				zap.L().Debug("inserted", zap.String("report_id", g["report_id"].(string)))
				archivedb.RawJsonInsert(table, qtx, ctx, g)
				// Commit after building up 20K inserts.
			}
		}
		tx.Commit()

	}

	return rows_retrieved, already_present, inserted
}

func update(cmd *cobra.Command, args []string) {
	db, Q, err := archivedb.GetSqliteDB(cmd.Flag("sqlite").Value.String())
	if err != nil {
		zap.L().Fatal("could not open database file", zap.String("sqlite filename", cmd.Flag("sqlite").Value.String()))
	}

	rows, present, inserted := doUpdate(cmd, db, Q)

	zap.L().Info("number of objects",
		zap.Int("rows_retrieved", rows),
		zap.Int("already_present", present),
		zap.Int("inserted", inserted),
	)

}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Refreshes an existing archive",
	Long: `Given an existing archive, this command updates it with newer records.

Usage:

fac-tool update --sqlite <filename> --fac-accepted-date 2025-03-02

This will download all new records that were submitted since the archive was created.

The filename for the database that will be updated is required. It must exist.

The date must be in the format YYYY-MM-DD, and all audits on-and-after that date
will be retrieved, and new audits added to the database.

This does not download PDFs.
`,
	Run: update,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	updateCmd.Flags().String("fac-accepted-date", "", "FAC accepted date, inclusive")
	updateCmd.MarkFlagRequired("fac-accepted-date")

}
