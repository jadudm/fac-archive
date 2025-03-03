/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jadudm/fac-tool/internal/archivedb"
	"github.com/jadudm/fac-tool/internal/fac"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func getReports(cmd *cobra.Command) (int, int64) {
	t1, err := time.Parse("2006-01-02", cmd.Flag("start-date").Value.String())
	if err != nil {
		zap.L().Error("could not parse date", zap.Error(err))
	}
	t2, err := time.Parse("2006-01-02", cmd.Flag("end-date").Value.String())
	if err != nil {
		zap.L().Error("could not parse date", zap.Error(err))
	}

	_, Q, err := archivedb.GetSqliteDB(cmd.Flag("sqlite").Value.String())
	if err != nil {
		zap.L().Fatal("could not get database", zap.Error(err))
	}

	ctx := context.Background()
	// FIXME: Because of how dates are parsed, subtract an hour so that we get the
	// right inclusive/exclusive behavior from the query. Also, this could
	// be a timezone issue, and a possible incorrect way to handle things...
	rids, err := Q.GetReportIdsBetween(ctx, archivedb.GetReportIdsBetweenParams{
		FacAcceptedDate:   t1.Add(-1 * time.Hour),
		FacAcceptedDate_2: t2.Add(-1 * time.Hour),
	})
	if err != nil {
		zap.L().Error("could not fetch report ids", zap.Error(err))
	}

	var bytes_downloaded int64 = 0
	reports_downloaded := 0

	for _, rid := range rids {
		url := fac.PdfBase + rid.ReportID
		// Make subdirs, or we could end up with too many files in one place.
		// Literally... `ls` and `dir` break on some platforms with too many files.
		path := filepath.Join(cmd.Flag("report-destination").Value.String(),
			rid.FacAcceptedDate.Format("2006-01-02"))

		err := os.MkdirAll(path, 0755)
		if err != nil {
			zap.L().Error("could not make destination directory", zap.Error(err))
		}

		out, err := os.Create(
			filepath.Join(path, fmt.Sprintf("%s.pdf", rid.ReportID)))
		if err != nil {
			zap.L().Error("could not create report file", zap.Error(err))
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			zap.L().Error("could not GET report", zap.Error(err))
		}
		defer resp.Body.Close()

		n, err := io.Copy(out, resp.Body)
		if err != nil {
			zap.L().Error("could not copy bytes to disk", zap.Error(err))
		}

		zap.L().Debug("report downloaded",
			zap.String("report_id", rid.ReportID),
			zap.String("fac_accepted_date", rid.FacAcceptedDate.Format("2006-01-02")))

		out.Close()
		resp.Body.Close()
		bytes_downloaded += n
		reports_downloaded += 1
	}

	return reports_downloaded, bytes_downloaded
}

func reports(cmd *cobra.Command, args []string) {
	reports, bytes := getReports(cmd)
	zap.L().Info("downloaded",
		zap.Int("reports", reports),
		zap.Int64("bytes", bytes))
}

// reportsCmd represents the reports command
var reportsCmd = &cobra.Command{
	Use:   "reports",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: reports,
}

func init() {
	rootCmd.AddCommand(reportsCmd)
	// at the root
	// archiveCmd.Flags().String("sqlite", "", "SQLite archive file")
	// archiveCmd.MarkFlagRequired("sqlite")

	reportsCmd.Flags().String("start-date", "", "FAC acceptance date to begin PDF downloads (inclusive)")
	reportsCmd.Flags().String("end-date", "", "FAC acceptance date to end PDF downloads (exclusive)")
	reportsCmd.Flags().String("report-destination", "", "Directory for PDF download")

	reportsCmd.MarkFlagRequired("start-date")
	reportsCmd.MarkFlagRequired("end-date")
	reportsCmd.MarkFlagRequired("report-destination")

}
