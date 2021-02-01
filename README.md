# Verticat

![Build Status](https://github.com/joeygibson/verticat/workflows/build/badge.svg)
![Release Status](https://github.com/joeygibson/verticat/workflows/release/badge.svg)

A tool to count the rows, and take rows
from  [Vertica native binary files](https://www.vertica.com/docs/9.3.x/HTML/Content/Authoring/AdministratorsGuide/BinaryFilesAppendix/CreatingNativeBinaryFormatFiles.htm)
.

## Usage

```bash
Usage: verticat [-cfHv] [-h value] [-o value] [-t value] <file1> [file...]
 -c, --count       count rows
 -f, --force       force overwrite of output file
 -h, --head=value  take the first n rows
 -H, --help        show help
 -o, --output=value
                   write head/tail results to this file
 -t, --tail=value  take the last n rows
 -v, --version     show version
```

Running with `--count` will print out the number of data rows in the file(s). The header does not count as a row.

Running with `--head` or `--tail` will copy the first `n` rows, or the last `n` rows, respectively, to the output. When
either is used with multiple files, `n` rows will be taken from each file, and written to the output.

If `--output` is not specified, `stdout` is used. Be careful with this, since binary data echoed to the console may be
undesirable.

If combining multiple files with `--head`, `--tail`, or no options, they must all share the same column layout.

If no options are given, `verticat` acts like the standard `cat` program, reading from `stdin` if no files are given. Since
Vertica native files start with a metadata header, if you want to cat multiple files together, specify them as arguments
to `verticat` itself.
