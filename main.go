package main

import (
	"fmt"
	"github.com/joeygibson/verticat/lib"
	"github.com/pborman/getopt/v2"
	"io"
	"os"
)

var (
	gitTag    string // Git tag
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	helpFlag := getopt.BoolLong("help", 'H', "show help")
	countFlag := getopt.BoolLong("count", 'c', "count rows")
	headRows := getopt.IntLong("head", 'h', 0, "take the first n rows")
	tailRows := getopt.IntLong("tail", 't', 0, "take the last n rows")
	outFileName := getopt.StringLong("output", 'o', "", "write head/tail results to this file")
	forceFlag := getopt.BoolLong("force", 'f', "force overwrite of output file")
	versionFlag := getopt.BoolLong("version", 'v', "show version")
	printHeaderFlag := getopt.BoolLong("print-header", 'p', "print out header and exit")

	getopt.SetParameters("[file...]")
	getopt.Parse()

	args := getopt.Args()

	if *helpFlag {
		printUsage()
		os.Exit(0)
	}

	if *versionFlag {
		var msg string

		if gitTag == "" && sha1ver == "" && buildTime == "" {
			msg = "development version"
		} else if gitTag != "" {
			msg = fmt.Sprintf("%s - %s - %s", gitTag, sha1ver, buildTime)
		} else {
			msg = fmt.Sprintf("%s - %s", sha1ver, buildTime)
		}

		fmt.Println(msg)
		os.Exit(0)
	}

	if *headRows > 0 && *tailRows > 0 {
		fmt.Fprintf(os.Stderr, "Error: --head and --tail are mutually exclusive\n\n")
		printUsage()
		os.Exit(1)
	}

	var inputFiles []*os.File

	if len(args) == 0 {
		if *tailRows > 0 {
			fmt.Fprintf(os.Stderr, "Error: --tail doesn't work when reading from stdin.\n\n")
			printUsage()
			os.Exit(1)
		}

		inputFiles = append(inputFiles, os.Stdin)
	} else {
		for _, fileName := range args {
			inFile, err := os.Open(fileName)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "error opening file: ", err)
				os.Exit(2)
			}

			defer inFile.Close()

			inputFiles = append(inputFiles, inFile)
		}
	}

	if *countFlag {
		err := lib.CountRows(inputFiles)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		os.Exit(0)
	}

	var output io.Writer

	if *outFileName != "" {
		_, err := os.Stat(*outFileName)
		if !os.IsNotExist(err) && !*forceFlag {
			fmt.Fprintln(os.Stderr, "output file exists; overwrite with --force")
			os.Exit(2)
		}

		file, err := os.Create(*outFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening new file: ", err)
			os.Exit(2)
		}

		defer file.Close()
		output = file
	} else {
		output = os.Stdout
	}

	firstTime := true

	for _, inputFile := range inputFiles {
		if *headRows != 0 {
			err := lib.Head(inputFile, output, *headRows, firstTime)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error processing file: ", err)
				os.Exit(1)
			}
		} else if *tailRows != 0 {
			err := lib.Tail(inputFile, output, *tailRows, firstTime)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error processing file: ", err)
				os.Exit(1)
			}
		} else if *printHeaderFlag {
			err := lib.PrintHeader(inputFile, output)

			if err != nil {
				fmt.Fprintln(os.Stderr, "error processing file: ", err)
				os.Exit(1)
			}
		} else {
			err := lib.Cat(inputFile, output, firstTime)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error processing file: ", err)
				os.Exit(1)
			}
		}

		firstTime = false
	}
}

func printUsage() {
	getopt.PrintUsage(os.Stderr)
	fmt.Fprintln(os.Stderr, "\nIf --head or --tail are used with multiple files, that many rows")
	fmt.Fprintln(os.Stderr, "will be taken from each file, and written to the output.")
	fmt.Fprintln(os.Stderr, "\nIf operating on multiple files, they must all share the same column layout.")
	fmt.Fprintln(os.Stderr, "\nWith no options passed, verticat acts like the standard cat program,")
	fmt.Fprintln(os.Stderr, "reading from stdin if no files are given. Since Vertica native files")
	fmt.Fprintln(os.Stderr, "start with a metadata header, if you want to cat multiple files together,")
	fmt.Fprintln(os.Stderr, "specify them as arguments to verticat itself.")
	fmt.Fprintln(os.Stderr, "\nThe -print-header outputs just the column widths for each file. This might")
	fmt.Fprintln(os.Stderr, "be useful in determining what an unknown file contains.")
}
