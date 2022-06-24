# Verticat

![Build Status](https://github.com/joeygibson/verticat/workflows/build/badge.svg)
![Release Status](https://github.com/joeygibson/verticat/workflows/release/badge.svg)

A tool to count the rows, and take rows
from  [Vertica native binary files](https://www.vertica.com/docs/9.3.x/HTML/Content/Authoring/AdministratorsGuide/BinaryFilesAppendix/CreatingNativeBinaryFormatFiles.htm)
.

## Usage

```bash
count/head/tail Vertica native binary files

USAGE:
    verticat [OPTIONS] [filename]...

ARGS:
    <filename>...    The file(s) to process

OPTIONS:
    -c, --count       count rows
    -h, --head <n>    take the first n rows
        --help        Print help information
    -t, --tail <n>    take the last n rows
    -V, --version     Print version information
```

Running with `--count` will print out the number of data rows in the file(s). The header does not count as a row.

Running with `--head` or `--tail` will copy the first `n` rows, or the last `n` rows, respectively, to the output. When
either is used with multiple files, `n` rows will be taken from each file, and written to the output.

If combining multiple files with `--head`, `--tail`, or no options, they must all share the same column layout.

If no options are given, `verticat` acts like the standard `cat` program, reading from `stdin` if no files are given. Since
Vertica native files start with a metadata header, if you want to cat multiple files together, specify them as arguments
to `verticat` itself.

