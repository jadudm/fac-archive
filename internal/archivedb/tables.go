package archivedb

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/jadudm/fac-tool/internal/fac"
	_ "modernc.org/sqlite"

	"go.uber.org/zap"
)

//go:embed schema.sql
var ddl string

func CreateTables(db_name string) (*sql.DB, *Queries, error) {
	ctx := context.Background()

	zap.L().Info("creating database", zap.String("filename", db_name))

	db, err := sql.Open(fac.SqliteDriver, db_name)
	if err != nil {
		zap.L().Error("could not create database file", zap.Error(err))
		return nil, nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		zap.L().Error("could not create tables", zap.Error(err))
		return nil, nil, err
	}

	return db, New(db), nil
}
