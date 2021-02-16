# Verticat

![Build Status](https://github.com/joeygibson/verticat/workflows/build/badge.svg)
![Release Status](https://github.com/joeygibson/verticat/workflows/release/badge.svg)

A tool to count the rows, and take rows
from  [Vertica native binary files](https://www.vertica.com/docs/9.3.x/HTML/Content/Authoring/AdministratorsGuide/BinaryFilesAppendix/CreatingNativeBinaryFormatFiles.htm)
.

## Usage

```bash
Usage: verticat [-cfHpv] [-h value] [-o value] [-r value] [-t value] [file...]
 -c, --count       count rows
 -f, --force       force overwrite of output file
 -h, --head=value  take the first n rows
 -H, --help        show help
 -o, --output=value
                   write head/tail results to this file
 -p, --print-header
                   print out header and exit
 -r, --reorder=value
                   reorder columns based on this file
 -t, --tail=value  take the last n rows
 -v, --version     show version
```

## Options

`--count` or `-c` will print out the number of data rows in the file(s). 
The header does not count as a row.

`--force`, or `-f` will overwrite an output file, if it exists.

`--head`, or `-h` copies the first `n` rows of the given file(s). If multiple files are
given, they must all have an identical layout, as the header of the first will be 
copied with all of the output.

`--help`, or `-H` prints usage information.

`--output`, or `-o` gives a filename to send the output to. If this is not specified,
output goes to `stdout`.

`--print-header`, or `-p` will print the list of column widths, in order, and exit.

`--reorder`, or `-r`, specifies a file with the the desired column reordering. The 
indices can be separated by any sort of whitespace (spaces, tabs, newslines, etc.).

`--tail`, or `-t`, copies the last `n` rows of the given file(s). As with `--head`, 
if multiple files are given, they must all have an identical layout, as the header 
of the first will be copied with all of the output.

`--version`, or `-v` display the program version.

## Notes

If no filenames are given, `verticat` acts like the standard `cat` program, reading 
from `stdin`. **This only works with a single file**, since Vertica native files
have a header. If you need to combine multiple files, give them as arguments to the
program itself.

## Examples

To read the first 5 lines of `foo.bin`, and write the output to `bar.bin`,

```bash
verticat --head 5 -o bar.bin foo.bin
```

To read the last 5 lines of `foo.bin`, and write the output to `bar.bin`,

```bash
verticat --tail 5 -o bar.bin foo.bin
```

To read the first 5 lines of `foo.bin`, `bar.bin`, and `baz.bin`, and write the
output to `quux.bin`

```bash
verticat --head 5 -o quux.bin foo.bin bar.bin baz.bin
```

To combine all of `foo.bin` and `bar.bin` and write the output to `baz.bin`

```bash
verticat -o baz.bin foo.bin bar.bin
```

To reorder `foo.bin` using the ordering specified in `reo.txt`, and write to `stdout`.
In this example, the first five columns will be written in reverse order, and the final 
four in their original order. 

```bash
$> cat reo.txt
5, 4, 3, 2, 1, 6, 7, 8, 9

$> verticat -r reo.txt foo.bin
```
