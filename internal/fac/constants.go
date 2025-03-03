package fac

var PdfBase = "https://app.fac.gov/dissemination/report/pdf/"
var LimitPerQuery = 20000
var MaxRows = 10000000

var Tables = [...]string{
	"general",
	"federal_awards",
	"findings",
	"findings_text",
	"notes_to_sefa",
	"corrective_action_plans",
	"passthrough",
	"secondary_auditors",
	"additional_ueis",
	"additional_eins",
}
