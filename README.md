# fac-tool


## about

`fac-tool` provides utilities for working with Federal Audit Clearinghouse data (FAC). With it, you can:

1. make a complete copy of the *data* in the FAC
2. update the data incrementally
3. download of all of the public audit reports from the FAC

There is roughly 6GB of data in the FAC. There are roughly 2.5TB of audit reports (Feb 2025).

`fac-tool` is written in Go, and should work on Mac, Windows, and Linux.

## design

`fac-tool` has three subcommands:

* `fac-tool archive` will create a new archive of all of the FAC data 
* `fac-tool update` updates an existing archive
* `fac-tool reports` will download PDFs of audit reports

## configuration

You will need a `config.yaml` file. It must be placed in one of two places:

1. `$HOME/.factool/config.yaml`
2. In the same directory as the `fac-tool` executable

The file should have the form:

```
api:
  scheme: https
  url: api.fac.gov
  key: <your-api-key>
```

### obtaining a key



## resources

* https://appliedgo.com/blog/go-project-layout