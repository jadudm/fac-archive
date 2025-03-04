/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jadudm/fac-archive/internal/archivedb"
	"github.com/jadudm/fac-archive/internal/fac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// FIXME:
// There is a bunch of copy-pasta between updateGeneral and
// updateTable. The query *might* be the abstraction point, but
// some of the logic is different, too. TBD.

func updateGeneral(cmd *cobra.Command, db *sql.DB, Q *archivedb.Queries) (int, int, int, []map[string]any) {

	table := "general"
	rows_retrieved := 0
	inserted := 0
	already_present := 0
	objects_retrieved := make([]map[string]any, 0)

	ctx := context.Background()
	last_updated, err := Q.GetMetadata(ctx, "last_updated")
	if err != nil {
		zap.L().Fatal("could not get last updated time", zap.Error(err))
	}

	last_updated_time, err := time.Parse("2006-01-02", last_updated)
	if err != nil {
		zap.L().Error("could not parse last updated time", zap.String("last_updated", last_updated))
	}
	// Always start a day before we last updated.
	last_updated_time = last_updated_time.Add(-25 * time.Hour)

	url := fmt.Sprintf("%s://%s/%s?fac_accepted_date=gte.%s",
		viper.GetString("api.scheme"),
		viper.GetString("api.url"),
		table,
		// cmd.Flag("fac-accepted-date").Value.String(),
		last_updated_time.Format("2006-01-02"),
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
		return rows_retrieved, already_present, inserted, objects_retrieved
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
				objects_retrieved = append(objects_retrieved, g)
			}
		}
		tx.Commit()
	}

	return rows_retrieved, already_present, inserted, objects_retrieved
}

func updateChunk(table string, chunk []string, db *sql.DB, Q *archivedb.Queries) (int, int, int) {
	rows_retrieved := 0
	inserted := 0
	already_present := 0
	comma_rids := strings.Join(chunk, ",")

	url := fmt.Sprintf("%s://%s/%s?report_id=in.(%s)",
		viper.GetString("api.scheme"),
		viper.GetString("api.url"),
		table,
		comma_rids,
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

		tx, err := db.Begin()
		if err != nil {
			zap.L().Error("could not initialize transaction", zap.Error(err))
		}

		defer tx.Rollback()
		qtx := Q.WithTx(tx)

		for _, g := range objects {
			inserted += 1
			zap.L().Debug("inserted", zap.String("report_id", g["report_id"].(string)))
			archivedb.RawJsonInsert(table, qtx, ctx, g)
		}
		tx.Commit()
	}

	return rows_retrieved, already_present, inserted
}

func updateTable(table string,
	db *sql.DB,
	Q *archivedb.Queries,
	generals []map[string]any) (int, int, int) {

	rows_retrieved := 0
	inserted := 0
	already_present := 0

	report_ids := make([]string, 0)
	for _, g := range generals {
		report_ids = append(report_ids, g["report_id"].(string))
	}

	// Now, chunk them 500 at a time.
	// Too many and Postgrest might fall over.
	for rids := range slices.Chunk(report_ids, 500) {
		rr, ap, i := updateChunk(table, rids, db, Q)
		rows_retrieved += rr
		already_present += ap
		inserted += i
	}

	return rows_retrieved, already_present, inserted
}

func update(cmd *cobra.Command, args []string) {
	db, Q, err := archivedb.GetSqliteDB(cmd.Flag("sqlite").Value.String())
	if err != nil {
		zap.L().Fatal("could not open database file", zap.String("sqlite filename", cmd.Flag("sqlite").Value.String()))
	}

	rows_retrieved, already_present, inserted, objects_retrieved := updateGeneral(cmd, db, Q)

	zap.L().Info("number of objects in general",
		zap.Int("rows_retrieved", rows_retrieved),
		zap.Int("already_present", already_present),
		zap.Int("inserted", inserted),
	)

	for _, table := range fac.Tables {
		if table != "general" {
			rows_retrieved, already_present, inserted = updateTable(table, db, Q, objects_retrieved)
			zap.L().Info("number of objects in "+table,
				zap.Int("rows_retrieved", rows_retrieved),
				zap.Int("already_present", already_present),
				zap.Int("inserted", inserted),
			)
		}
	}

	ctx := context.Background()
	Q.UpdateMetadata(ctx, archivedb.UpdateMetadataParams{
		Key:   "last_updated",
		Value: time.Now().Format("2006-01-02"),
	})
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Refreshes an existing archive",
	Long: `Given an existing archive, this command updates it with newer records.

Usage:

fac-archive update --sqlite <filename>

This will download all new records that were submitted since the archive was last updated.

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

	//updateCmd.Flags().String("fac-accepted-date", "", "FAC accepted date, inclusive")
	// updateCmd.MarkFlagRequired("fac-accepted-date")

}
