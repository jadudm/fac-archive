-- name: RawInsert :one
INSERT INTO raw (source, json) VALUES (?, ?) RETURNING id;

-- name: ReportIdExists :one
SELECT EXISTS (SELECT 1 FROM general WHERE report_id = ?);

-- name: GetReportIdsBetween :many
SELECT report_id, fac_accepted_date FROM general 
WHERE 
  fac_accepted_date >= ?
  AND fac_accepted_date < ?;

-- name: GetRawIdFromReportId :one
SELECT raw_id
FROM general
WHERE 
  report_id = ?;

-- name: IsReportDownloaded :one
SELECT is_downloaded 
FROM pdfs
WHERE
  raw_id = ?;

-- name: SetReportDownloaded :exec
UPDATE pdfs
SET is_downloaded = 1
WHERE
    raw_id = ?;

-- name: UnsetReportDownloaded :exec
UPDATE pdfs
SET is_downloaded = 0
WHERE
    raw_id = ?;