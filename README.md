# fac-copy


## about

`fac-copy` serves to:

1. Make a complete copy of the *data* in the Federal Audit Clearinghouse (FAC)
2. Update the data incrementally
3. Download of all of the public audit reports from the FAC

There is roughly 6GB of data in the FAC. There are roughly 2.5TB of audit reports (Feb 2025).

`fac-copy` is written in Go, and should work on Mac, Windows, and Linux.

## design

`fac-copy` has three subcommands:

* `fac-copy archive` will create a new archive of all of the FAC data 
* `fac-copy update` updates an existing archive
* `fac-copy pdfs` will download PDFs of audit reports




## resources

* https://appliedgo.com/blog/go-project-layout