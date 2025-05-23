// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package archivedb

import (
	"context"
	"time"
)

const addMetadata = `-- name: AddMetadata :exec
INSERT INTO metadata
(key, value)
VALUES (?, ?)
`

type AddMetadataParams struct {
	Key   string
	Value string
}

func (q *Queries) AddMetadata(ctx context.Context, arg AddMetadataParams) error {
	_, err := q.db.ExecContext(ctx, addMetadata, arg.Key, arg.Value)
	return err
}

const getMetadata = `-- name: GetMetadata :one
SELECT value
FROM metadata
WHERE key = ?
`

func (q *Queries) GetMetadata(ctx context.Context, key string) (string, error) {
	row := q.db.QueryRowContext(ctx, getMetadata, key)
	var value string
	err := row.Scan(&value)
	return value, err
}

const getRawIdFromReportId = `-- name: GetRawIdFromReportId :one
SELECT raw_id
FROM general
WHERE 
  report_id = ?
`

func (q *Queries) GetRawIdFromReportId(ctx context.Context, reportID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getRawIdFromReportId, reportID)
	var raw_id int64
	err := row.Scan(&raw_id)
	return raw_id, err
}

const getReportIdsBetween = `-- name: GetReportIdsBetween :many
SELECT report_id, fac_accepted_date FROM general 
WHERE 
  fac_accepted_date >= ?
  AND fac_accepted_date < ?
`

type GetReportIdsBetweenParams struct {
	FacAcceptedDate   time.Time
	FacAcceptedDate_2 time.Time
}

type GetReportIdsBetweenRow struct {
	ReportID        string
	FacAcceptedDate time.Time
}

func (q *Queries) GetReportIdsBetween(ctx context.Context, arg GetReportIdsBetweenParams) ([]GetReportIdsBetweenRow, error) {
	rows, err := q.db.QueryContext(ctx, getReportIdsBetween, arg.FacAcceptedDate, arg.FacAcceptedDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReportIdsBetweenRow
	for rows.Next() {
		var i GetReportIdsBetweenRow
		if err := rows.Scan(&i.ReportID, &i.FacAcceptedDate); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isReportDownloaded = `-- name: IsReportDownloaded :one
SELECT is_downloaded 
FROM pdfs
WHERE
  raw_id = ?
`

func (q *Queries) IsReportDownloaded(ctx context.Context, rawID int64) (bool, error) {
	row := q.db.QueryRowContext(ctx, isReportDownloaded, rawID)
	var is_downloaded bool
	err := row.Scan(&is_downloaded)
	return is_downloaded, err
}

const rawInsert = `-- name: RawInsert :one
INSERT INTO raw (source, json) VALUES (?, ?) RETURNING id
`

type RawInsertParams struct {
	Source string
	Json   string
}

func (q *Queries) RawInsert(ctx context.Context, arg RawInsertParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, rawInsert, arg.Source, arg.Json)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const reportIdExists = `-- name: ReportIdExists :one
SELECT EXISTS (SELECT 1 FROM general WHERE report_id = ?)
`

func (q *Queries) ReportIdExists(ctx context.Context, reportID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, reportIdExists, reportID)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const setReportDownloaded = `-- name: SetReportDownloaded :exec
UPDATE pdfs
SET is_downloaded = 1
WHERE
    raw_id = ?
`

func (q *Queries) SetReportDownloaded(ctx context.Context, rawID int64) error {
	_, err := q.db.ExecContext(ctx, setReportDownloaded, rawID)
	return err
}

const unsetReportDownloaded = `-- name: UnsetReportDownloaded :exec
UPDATE pdfs
SET is_downloaded = 0
WHERE
    raw_id = ?
`

func (q *Queries) UnsetReportDownloaded(ctx context.Context, rawID int64) error {
	_, err := q.db.ExecContext(ctx, unsetReportDownloaded, rawID)
	return err
}

const updateMetadata = `-- name: UpdateMetadata :exec
UPDATE metadata
SET value = ?
WHERE
  key = ?
`

type UpdateMetadataParams struct {
	Value string
	Key   string
}

func (q *Queries) UpdateMetadata(ctx context.Context, arg UpdateMetadataParams) error {
	_, err := q.db.ExecContext(ctx, updateMetadata, arg.Value, arg.Key)
	return err
}
