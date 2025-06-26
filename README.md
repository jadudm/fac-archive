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

## grabbing the tool

This repository builds releases for multiple platforms. You need to download the file that is appropriate for you.

* For a Mac with Apple Silicon (most recent Macs), you want `fac-archive-mac-arm64`
* For a Mac with an Intel processor, you want `fac-archive-mac-amd64`
* For Linux with Intel/AMD, you want `fac-archive-linux-amd64`
* For Linux with an ARM processor (some cloud environments), use `fac-archive-linux-arm64`
* For Windows with Intel (most users), `fac-archive-windows-amd64`
* For Windows with ARM (Surface devices and some cloud environments), use `fac-archive-windows-arm64`

See [BUILDING.md](BUILDING.md) for how to build and run the software yourself. 

We will refer to the tool as `fac-archive` for brevity.

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

Other config keys that can be set:

```
api:
  accept_profile: api_v1_2_0
  limit_per_query: 5000
```

These let the tool grab from different API endpoints, and limit how many records come back per call.

```
copy_json: false
```

This flag determines whether or not the triggers will run to copy JSON objects from their raw form into the structured/relational tables. `true` will create relational tables from the JSON objects; `false` will not. The `false` flag is useful if you want to only work with the raw JSON values, and/or if you're concerned the JSON values might violate the constranits imposed on the relational tables (e.g. `NOT NULL`).

### obtaining a key

Follow the instructions [here](https://www.fac.gov/api/) for obtaining an API key. It will be mailed to you, and you can then paste it into your `config.yaml` file. 

* Never share your API key.
* Never post your API key to a public website.

## generating an archive

Running

```
fac-archive archive --sqlite fac.db
```

will create an SQLite file called `fac.db` and download all of the data from the FAC into that database:

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

## updating an archive

Once you have downloaded a complete archive (which takes around 30-40m), you can incrementally update the archive. 

```
fac-archive archive --sqlite fac.db
```

The database knows when it was last updated, and will download all new records since its last update. The output will look something like this:

```
{"level":"info","msg":"number of objects in general","rows_retrieved":422,"already_present":243,"inserted":179}
{"level":"info","msg":"number of objects in federal_awards","rows_retrieved":3196,"already_present":0,"inserted":3196}
{"level":"info","msg":"number of objects in findings","rows_retrieved":136,"already_present":0,"inserted":136}
{"level":"info","msg":"number of objects in findings_text","rows_retrieved":85,"already_present":0,"inserted":85}
{"level":"info","msg":"number of objects in notes_to_sefa","rows_retrieved":440,"already_present":0,"inserted":440}
{"level":"info","msg":"number of objects in corrective_action_plans","rows_retrieved":68,"already_present":0,"inserted":68}
{"level":"info","msg":"number of objects in passthrough","rows_retrieved":2390,"already_present":0,"inserted":2390}
{"level":"info","msg":"number of objects in secondary_auditors","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in additional_ueis","rows_retrieved":9,"already_present":0,"inserted":9}
{"level":"info","msg":"number of objects in additional_eins","rows_retrieved":13,"already_present":0,"inserted":13}
```

If there are no new records:

```
{"level":"info","msg":"number of objects in general","rows_retrieved":417,"already_present":417,"inserted":0}
{"level":"info","msg":"number of objects in federal_awards","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in findings","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in findings_text","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in notes_to_sefa","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in corrective_action_plans","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in passthrough","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in secondary_auditors","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in additional_ueis","rows_retrieved":0,"already_present":0,"inserted":0}
{"level":"info","msg":"number of objects in additional_eins","rows_retrieved":0,"already_present":0,"inserted":0}
```

## downloading report PDFs

Downloading PDFs takes time, and therefore requires you to add more parameters.

```
fac-archive reports --sqlite fac.db --start-date 2025-03-04 --end-date 2025-03-05 --report-destination pdfs/
```

will download all PDFs that were submitted March 4th (inclusive) through March 5th (exclusive). This means that 

* All reports submitted on the start date (March 4th) will be downloaded
* All reports submitted on the end date (March 5th) will NOT be downloaded

The `--report-destination` flag choses the directory into which the reports will be downloaded. Subdirectories will be created for the date of each report; for example, if you point to the directory `pdfs/`, the archiver will create a set of folders that looks like:

```
pdfs/
 |- 2025-03-01/
 |---- 2023-06-GSAFAC-0000353213.pdf 
 |---- 2024-06-GSAFAC-0000352510.pdf 
 |---- 2024-06-GSAFAC-0000355942.pdf
 |- 2025-03-02/
 |---- 2024-06-GSAFAC-0000356942.pdf
 | ...
```

and so on.

In order to download *all reports*, enter a date before 2016 for the start date, and after the current date for the end date.

> [!CAUTION]
> There are (as of March 2025) roughly 2.5TB of PDFs in the FAC. It will take *several days* over an institutional or fiber-based internet connection to download everything. The download process keeps track of which reports were successfully downloaded, and can be stopped and restarted without having to re-download all of the previously downloaded reports.

Reports never change, so there is no reasons to re-download them repeatedly.

# license

This project is released into the Public Domain in the United States of America.

Everywhere else, it is made available under the [CC0 1.0 Universal license](https://creativecommons.org/publicdomain/zero/1.0/legalcode
).


