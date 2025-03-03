-- name: RawInsert :one
INSERT INTO raw (source, json) VALUES (?, ?) RETURNING id;

-- name: ReportIdExists :one
SELECT EXISTS (SELECT 1 FROM general WHERE report_id = ?);

-- name: GetReportIdsBetween :many
SELECT report_id, fac_accepted_date FROM general 
WHERE 
  fac_accepted_date >= ?
  AND fac_accepted_date < ?;