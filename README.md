# fac-tool


## about

`fac-tool` provides utilities for working with Federal Audit Clearinghouse data (FAC). With it, you can:

1. make a complete copy of the *data* in the FAC
2. update the data incrementally
3. download of all of the public audit reports from the FAC

There is roughly 6GB of data in the FAC. There are roughly 2.5TB of audit reports (Feb 2025).

`fac-copy` is written in Go, and should work on Mac, Windows, and Linux.

## design

`fac-copy` has three subcommands:

* `fac-copy archive` will create a new archive of all of the FAC data 
* `fac-copy update` updates an existing archive
* `fac-copy reports` will download PDFs of audit reports




## resources

* https://appliedgo.com/blog/go-project-layout