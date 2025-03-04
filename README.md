# fac-archive


## about

`fac-archive` provides utilities for working with Federal Audit Clearinghouse data (FAC). With it, you can:

1. make a complete copy of the *data* in the FAC
2. update the data incrementally
3. download of all of the public audit reports from the FAC

The archive tool can be used to:

* Back up roughly 10GB of data to an SQLite file
* Download approximately 2.5TB of audit reports (PDFs)

`fac-archive` is written in Go, and should work on Mac, Windows, and Linux.

## design

`fac-archive` has three subcommands:

* `fac-archive archive` will create a new archive of all of the FAC data 
* `fac-archive update` updates an existing archive
* `fac-archive reports` will download PDFs of audit reports

## configuration

You will need a `config.yaml` file. It must be placed in one of two places:

1. `$HOME/.factool/config.yaml`
2. A file called `config.yaml` in the same directory as the `fac-archive` executable

The file should have the form:

```
api:
  scheme: https
  url: api.fac.gov
  key: <your-api-key>
  
# This can be DEBUG, INFO, WARN, or ERROR.
# The default is INFO.
debug_level: DEBUG
```

### obtaining a key


## generating an archive

Running

```
fac-archive archive
```

will create a timestamped SQLite file in the same directory that the tool is run. During archiving, a small amount of information is logged along the way:

```
{"level":"info","msg":"creating database","filename":"2025-02-28-12-19-57-fac.sqlite"}
{"level":"info","msg":"rows retrieved","table":"general","rows":346123,"duration":124}
{"level":"info","msg":"rows retrieved","table":"federal_awards","rows":5933059,"duration":2640}
{"level":"info","msg":"rows retrieved","table":"findings","rows":534306,"duration":133}
{"level":"info","msg":"rows retrieved","table":"findings_text","rows":120015,"duration":48}
{"level":"info","msg":"rows retrieved","table":"notes_to_sefa","rows":528469,"duration":297}
{"level":"info","msg":"rows retrieved","table":"corrective_action_plans","rows":116443,"duration":35}
{"level":"info","msg":"rows retrieved","table":"passthrough","rows":4042221,"duration":1388}
{"level":"info","msg":"rows retrieved","table":"secondary_auditors","rows":1804,"duration":0}
{"level":"info","msg":"rows retrieved","table":"additional_ueis","rows":14984,"duration":3}
{"level":"info","msg":"rows retrieved","table":"additional_eins","rows":59343,"duration":11}
```

## resources

* https://appliedgo.com/blog/go-project-layout