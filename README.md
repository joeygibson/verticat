# Verticat

![Build Status](https://github.com/joeygibson/verticat/workflows/build/badge.svg)
![Release Status](https://github.com/joeygibson/verticat/workflows/release/badge.svg)

A tool to count the rows, and take rows from  [Vertica native binary files](https://www.vertica.com/docs/9.3.x/HTML/Content/Authoring/AdministratorsGuide/BinaryFilesAppendix/CreatingNativeBinaryFormatFiles.htm).

## Usage

```bash
Usage: verticat [-cfHv] [-h value] [-o value] [-t value] <file>
 -c, --count       count rows
 -f, --force       force overwrite of output file
 -h, --head=value  take the first n rows
 -H, --help        show help
 -o, --output=value
                   write head/tail results to this file
 -t, --tail=value  take the last n rows
 -v, --version     show version
```

Running with just `-c` will print out the number of data rows in the file. The header does not
count as a row.

Running with `--head` or `--tail` will copy the first `n` rows, or the last `n` rows, 
respectively, to the file specified with `-o`. If no `-o` is given, `stdout` is used.

