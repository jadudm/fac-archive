CREATE TABLE IF NOT EXISTS pdfs (
  id INTEGER PRIMARY KEY,
	raw_id INTEGER,
	-- url TEXT,
	is_downloaded BOOLEAN DEFAULT 0
);

CREATE TABLE IF NOT EXISTS raw (
  id INTEGER PRIMARY KEY,
  source TEXT NOT NULL,
  json TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- No need to build this URL. I will need to track the boolean, though.
-- 'https://app.fac.gov/dissemination/report/pdf/' || json_extract(new.json, '$.report_id')
CREATE TRIGGER IF NOT EXISTS trig_raw_to_pdf
	AFTER INSERT
	ON raw
	WHEN new.source = 'general'
	BEGIN
		INSERT INTO pdfs
		(raw_id)
		VALUES
		(new.id)
		;
END;


----------------------------------------------------------
-- general
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS general (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	auditee_certify_name text NOT NULL,
	auditee_certify_title text NOT NULL,
	auditee_contact_name text NOT NULL,
	auditee_email text NOT NULL,
	auditee_name text NOT NULL,
	auditee_phone text NOT NULL,
	auditee_contact_title text NOT NULL,
	auditee_address_line_1 text NOT NULL,
	auditee_city text NOT NULL,
	auditee_state text NOT NULL,
	auditee_ein text NOT NULL,
	is_additional_ueis text NOT NULL,
	auditee_zip text NOT NULL,
	auditor_phone text NOT NULL,
	auditor_state text NOT NULL,
	auditor_city text NOT NULL,
	auditor_contact_title text NOT NULL,
	auditor_address_line_1 text NOT NULL,
	auditor_zip text NOT NULL,
	auditor_country text NOT NULL,
	auditor_contact_name text NOT NULL,
	auditor_email text NOT NULL,
	auditor_firm_name text NOT NULL,
	auditor_foreign_address text NOT NULL,
	auditor_ein text NOT NULL,
	cognizant_agency text NULL,
	oversight_agency text NULL,
	date_created date NOT NULL,
	ready_for_certification_date date NOT NULL,
	auditor_certified_date date NOT NULL,
	auditee_certified_date date NOT NULL,
	submitted_date date NOT NULL,
	fy_end_date date NOT NULL,
	fy_start_date date NOT NULL,
	audit_type text NOT NULL,
	gaap_results text NOT NULL,
	sp_framework_basis text NOT NULL,
	is_sp_framework_required text NOT NULL,
	sp_framework_opinions text NOT NULL,
	is_going_concern_included text NOT NULL,
	is_internal_control_deficiency_disclosed text NOT NULL,
	is_internal_control_material_weakness_disclosed text NOT NULL,
	is_material_noncompliance_disclosed text NOT NULL,
	is_aicpa_audit_guide_included text NOT NULL,
	dollar_threshold int8 NOT NULL,
	is_low_risk_auditee text NOT NULL,
	agencies_with_prior_findings text NOT NULL,
	entity_type text NOT NULL,
	number_months text NOT NULL,
	audit_period_covered text NOT NULL,
	total_amount_expended int8 NOT NULL,
	type_audit_code text NOT NULL,
	is_public bool NOT NULL,
	data_source text NOT NULL,
	fac_accepted_date date NOT NULL,
	auditor_certify_name text NOT NULL,
	auditor_certify_title text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_general 
  AFTER INSERT 
  ON raw
  WHEN new.source = 'general'
BEGIN
	INSERT INTO general 
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		auditee_certify_name,
		auditee_certify_title,
		auditee_contact_name,
		auditee_email,
		auditee_name,
		auditee_phone,
		auditee_contact_title,
		auditee_address_line_1,
		auditee_city,
		auditee_state,
		auditee_ein,
		is_additional_ueis,
		auditee_zip,
		auditor_phone,
		auditor_state,
		auditor_city,
		auditor_contact_title,
		auditor_address_line_1,
		auditor_zip,
		auditor_country,
		auditor_contact_name,
		auditor_email,
		auditor_firm_name,
		auditor_foreign_address,
		auditor_ein,
		cognizant_agency,
		oversight_agency,
		date_created,
		ready_for_certification_date,
		auditor_certified_date,
		auditee_certified_date,
		submitted_date,
		fy_end_date,
		fy_start_date,
		audit_type,
		gaap_results,
		sp_framework_basis,
		is_sp_framework_required,
		sp_framework_opinions,
		is_going_concern_included,
		is_internal_control_deficiency_disclosed,
		is_internal_control_material_weakness_disclosed,
		is_material_noncompliance_disclosed,
		is_aicpa_audit_guide_included,
		dollar_threshold,
		is_low_risk_auditee,
		agencies_with_prior_findings,
		entity_type,
		number_months,
		audit_period_covered,
		total_amount_expended,
		type_audit_code,
		is_public,
		data_source,
		fac_accepted_date,
		auditor_certify_name,
		auditor_certify_title
		) 
		VALUES
		(
			new.id, 
			json_extract(new.json, '$.report_id'),
			json_extract(new.json, '$.auditee_uei'),
			json_extract(new.json, '$.audit_year'),
			json_extract(new.json, '$.auditee_certify_name'),
			json_extract(new.json, '$.auditee_certify_title'),
			json_extract(new.json, '$.auditee_contact_name'),
			json_extract(new.json, '$.auditee_email'),
			json_extract(new.json, '$.auditee_name'),
			json_extract(new.json, '$.auditee_phone'),
			json_extract(new.json, '$.auditee_contact_title'),
			json_extract(new.json, '$.auditee_address_line_1'),
			json_extract(new.json, '$.auditee_city'),
			json_extract(new.json, '$.auditee_state'),
			json_extract(new.json, '$.auditee_ein'),
			json_extract(new.json, '$.is_additional_ueis'),
			json_extract(new.json, '$.auditee_zip'),
			json_extract(new.json, '$.auditor_phone'),
			json_extract(new.json, '$.auditor_state'),
			json_extract(new.json, '$.auditor_city'),
			json_extract(new.json, '$.auditor_contact_title'),
			json_extract(new.json, '$.auditor_address_line_1'),
			json_extract(new.json, '$.auditor_zip'),
			json_extract(new.json, '$.auditor_country'),
			json_extract(new.json, '$.auditor_contact_name'),
			json_extract(new.json, '$.auditor_email'),
			json_extract(new.json, '$.auditor_firm_name'),
			json_extract(new.json, '$.auditor_foreign_address'),
			json_extract(new.json, '$.auditor_ein'),
			json_extract(new.json, '$.cognizant_agency'),
			json_extract(new.json, '$.oversight_agency'),
			json_extract(new.json, '$.date_created'),
			json_extract(new.json, '$.ready_for_certification_date'),
			json_extract(new.json, '$.auditor_certified_date'),
			json_extract(new.json, '$.auditee_certified_date'),
			json_extract(new.json, '$.submitted_date'),
			json_extract(new.json, '$.fy_end_date'),
			json_extract(new.json, '$.fy_start_date'),
			json_extract(new.json, '$.audit_type'),
			json_extract(new.json, '$.gaap_results'),
			json_extract(new.json, '$.sp_framework_basis'),
			json_extract(new.json, '$.is_sp_framework_required'),
			json_extract(new.json, '$.sp_framework_opinions'),
			json_extract(new.json, '$.is_going_concern_included'),
			json_extract(new.json, '$.is_internal_control_deficiency_disclosed'),
			json_extract(new.json, '$.is_internal_control_material_weakness_disclosed'),
			json_extract(new.json, '$.is_material_noncompliance_disclosed'),
			json_extract(new.json, '$.is_aicpa_audit_guide_included'),
			json_extract(new.json, '$.dollar_threshold'),
			json_extract(new.json, '$.is_low_risk_auditee'),
			json_extract(new.json, '$.agencies_with_prior_findings'),
			json_extract(new.json, '$.entity_type'),
			json_extract(new.json, '$.number_months'),
			json_extract(new.json, '$.audit_period_covered'),
			json_extract(new.json, '$.total_amount_expended'),
			json_extract(new.json, '$.type_audit_code'),
			json_extract(new.json, '$.is_public'),
			json_extract(new.json, '$.data_source'),
			json_extract(new.json, '$.fac_accepted_date'),
			json_extract(new.json, '$.auditor_certify_name'),
			json_extract(new.json, '$.auditor_certify_title')
		);
END;

----------------------------------------------------------
-- federal_awards
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS federal_awards (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	additional_award_identification text NOT NULL,
	amount_expended int8 NOT NULL,
	award_reference text NOT NULL,
	cluster_name text NOT NULL,
	cluster_total int8 NOT NULL,
	federal_agency_prefix text NOT NULL,
	federal_award_extension text NOT NULL,
	federal_program_name text NOT NULL,
	federal_program_total int8 NOT NULL,
	findings_count int4 NOT NULL,
	is_direct text NOT NULL,
	is_loan text NOT NULL,
	is_major text NOT NULL,
	is_passthrough_award text NOT NULL,
	loan_balance text NOT NULL,
	audit_report_type text NOT NULL,
	other_cluster_name text NOT NULL,
	passthrough_amount int8 NULL,
	state_cluster_name text NOT NULL
);


CREATE TRIGGER IF NOT EXISTS trig_raw_to_federal_awards 
  AFTER INSERT 
  ON raw
  WHEN new.source = 'federal_awards'
BEGIN
	INSERT INTO federal_awards
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		additional_award_identification,
		amount_expended,
		award_reference,
		cluster_name,
		cluster_total,
		federal_agency_prefix,
		federal_award_extension,
		federal_program_name,
		federal_program_total,
		findings_count,
		is_direct,
		is_loan,
		is_major,
		is_passthrough_award,
		loan_balance,
		audit_report_type,
		other_cluster_name,
		passthrough_amount,
		state_cluster_name
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.additional_award_identification'),
		json_extract(new.json, '$.amount_expended'),
		json_extract(new.json, '$.award_reference'),
		json_extract(new.json, '$.cluster_name'),
		json_extract(new.json, '$.cluster_total'),
		json_extract(new.json, '$.federal_agency_prefix'),
		json_extract(new.json, '$.federal_award_extension'),
		json_extract(new.json, '$.federal_program_name'),
		json_extract(new.json, '$.federal_program_total'),
		json_extract(new.json, '$.findings_count'),
		json_extract(new.json, '$.is_direct'),
		json_extract(new.json, '$.is_loan'),
		json_extract(new.json, '$.is_major'),
		json_extract(new.json, '$.is_passthrough_award'),
		json_extract(new.json, '$.loan_balance'),
		json_extract(new.json, '$.audit_report_type'),
		json_extract(new.json, '$.other_cluster_name'),
		json_extract(new.json, '$.passthrough_amount'),
		json_extract(new.json, '$.state_cluster_name')
	);
	END;

----------------------------------------------------------
-- findings
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS findings (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	award_reference text NOT NULL,
	reference_number text NOT NULL,
	is_material_weakness text NOT NULL,
	is_modified_opinion text NOT NULL,
	is_other_findings text NOT NULL,
	is_other_matters text NOT NULL,
	is_questioned_costs text NOT NULL,
	is_repeat_finding text NOT NULL,
	is_significant_deficiency text NOT NULL,
	prior_finding_ref_numbers text NOT NULL,
	type_requirement text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_findings
  AFTER INSERT 
  ON raw
  WHEN new.source = 'findings'
BEGIN
	INSERT INTO findings
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		award_reference,
		reference_number,
		is_material_weakness,
		is_modified_opinion,
		is_other_findings,
		is_other_matters,
		is_questioned_costs,
		is_repeat_finding,
		is_significant_deficiency,
		prior_finding_ref_numbers,
		type_requirement
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.award_reference'),
		json_extract(new.json, '$.reference_number'),
		json_extract(new.json, '$.is_material_weakness'),
		json_extract(new.json, '$.is_modified_opinion'),
		json_extract(new.json, '$.is_other_findings'),
		json_extract(new.json, '$.is_other_matters'),
		json_extract(new.json, '$.is_questioned_costs'),
		json_extract(new.json, '$.is_repeat_finding'),
		json_extract(new.json, '$.is_significant_deficiency'),
		json_extract(new.json, '$.prior_finding_ref_numbers'),
		json_extract(new.json, '$.type_requirement')
	);
	END;

----------------------------------------------------------
-- notes_to_sefa
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS notes_to_sefa (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	accounting_policies text NOT NULL,
	is_minimis_rate_used text NOT NULL,
	rate_explained text NOT NULL,
	content text NOT NULL,
	title text NOT NULL,
	contains_chart_or_table text NOT NULL
);


CREATE TRIGGER IF NOT EXISTS trig_raw_to_notes_to_sefa
  AFTER INSERT 
  ON raw
  WHEN new.source = 'notes_to_sefa'
BEGIN
	INSERT INTO notes_to_sefa
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		accounting_policies,
		is_minimis_rate_used,
		rate_explained,
		content,
		title,
		contains_chart_or_table
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.accounting_policies'),
		json_extract(new.json, '$.is_minimis_rate_used'),
		json_extract(new.json, '$.rate_explained'),
		json_extract(new.json, '$.content'),
		json_extract(new.json, '$.title'),
		json_extract(new.json, '$.contains_chart_or_table')
	);
	END;

----------------------------------------------------------
-- findings_text
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS findings_text (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	finding_ref_number text NOT NULL,
	contains_chart_or_table text NOT NULL,
	finding_text text NOT NULL
);


CREATE TRIGGER IF NOT EXISTS trig_raw_to_findings_text
  AFTER INSERT 
  ON raw
  WHEN new.source = 'findings_text'
BEGIN
	INSERT INTO findings_text
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		finding_ref_number,
		contains_chart_or_table,
		finding_text
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.finding_ref_number'),
		json_extract(new.json, '$.contains_chart_or_table'),
		json_extract(new.json, '$.finding_text')
	);
	END;

----------------------------------------------------------
-- additional_ueis
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS additional_ueis (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	additional_uei text NOT NULL
);


CREATE TRIGGER IF NOT EXISTS trig_raw_to_additional_ueis
  AFTER INSERT 
  ON raw
  WHEN new.source = 'additional_ueis'
BEGIN
	INSERT INTO additional_ueis
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		additional_uei
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.additional_uei')
	);
	END;


----------------------------------------------------------
-- corrective_action_plans
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS corrective_action_plans (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	finding_ref_number text NOT NULL,
	contains_chart_or_table text NOT NULL,
	planned_action text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_corrective_action_plans
  AFTER INSERT 
  ON raw
  WHEN new.source = 'corrective_action_plans'
BEGIN
	INSERT INTO corrective_action_plans
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		finding_ref_number,
		contains_chart_or_table,
		planned_action
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.finding_ref_number'),
		json_extract(new.json, '$.contains_chart_or_table'),
		json_extract(new.json, '$.planned_action')
	);
	END;

----------------------------------------------------------
-- passthrough
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS passthrough (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	award_reference text NOT NULL,
	passthrough_id text NOT NULL,
	passthrough_name text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_passthrough
  AFTER INSERT 
  ON raw
  WHEN new.source = 'passthrough'
BEGIN
	INSERT INTO passthrough
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		award_reference,
		passthrough_id,
		passthrough_name
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.award_reference'),
		json_extract(new.json, '$.passthrough_id'),
		json_extract(new.json, '$.passthrough_name')
	);
	END;

----------------------------------------------------------
-- secondary_auditors
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS secondary_auditors (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	auditor_ein text NOT NULL,
	auditor_name text NOT NULL,
	contact_name text NOT NULL,
	contact_title text NOT NULL,
	contact_email text NOT NULL,
	contact_phone text NOT NULL,
	address_street text NOT NULL,
	address_city text NOT NULL,
	address_state text NOT NULL,
	address_zipcode text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_secondary_auditors
  AFTER INSERT 
  ON raw
  WHEN new.source = 'secondary_auditors'
BEGIN
	INSERT INTO secondary_auditors
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		auditor_ein,
		auditor_name,
		contact_name,
		contact_title,
		contact_email,
		contact_phone,
		address_street,
		address_city,
		address_state,
		address_zipcode
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.auditor_ein'),
		json_extract(new.json, '$.auditor_name'),
		json_extract(new.json, '$.contact_name'),
		json_extract(new.json, '$.contact_title'),
		json_extract(new.json, '$.contact_email'),
		json_extract(new.json, '$.contact_phone'),
		json_extract(new.json, '$.address_street'),
		json_extract(new.json, '$.address_city'),
		json_extract(new.json, '$.address_state'),
		json_extract(new.json, '$.address_zipcode')
	);
	END;

----------------------------------------------------------
-- additional_eins
----------------------------------------------------------
CREATE TABLE IF NOT EXISTS additional_eins (
	id INTEGER PRIMARY KEY,
	raw_id INTEGER NOT NULL,
	report_id text NOT NULL,
	auditee_uei text NOT NULL,
	audit_year integer NOT NULL,
	additional_ein text NOT NULL
);

CREATE TRIGGER IF NOT EXISTS trig_raw_to_additional_eins
  AFTER INSERT 
  ON raw
  WHEN new.source = 'additional_eins'
BEGIN
	INSERT INTO additional_eins
	(
		raw_id,
		report_id,
		auditee_uei,
		audit_year,
		additional_ein
	) 
	VALUES
	(
		new.id,
		json_extract(new.json, '$.report_id'),
		json_extract(new.json, '$.auditee_uei'),
		json_extract(new.json, '$.audit_year'),
		json_extract(new.json, '$.additional_ein')
	);
	END;