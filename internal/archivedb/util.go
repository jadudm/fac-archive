package archivedb

import (
	"context"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

func CreateSqliteDB(db_name string) (*sql.DB, *Queries, error) {
	db, queries, err := CreateTables(db_name)
	return db, queries, err
}

func GetSqliteDB(db_name string) (*sql.DB, *Queries, error) {
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		zap.L().Error("could not open database file", zap.Error(err))
		return nil, nil, err
	}

	return db, New(db), nil
}

func RawJsonInsert(table string, qTx *Queries, ctx context.Context, g map[string]any) {
	b, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("could not marshal to string", zap.Error(err))
	}
	id, err := qTx.RawInsert(ctx, RawInsertParams{
		Source: table,
		Json:   string(b),
	})
	if err != nil {
		zap.L().Error("could not insert", zap.String("table", table), zap.Error(err))
	} else {
		zap.L().Debug("inserted id", zap.Int64("id", id))
	}
}
