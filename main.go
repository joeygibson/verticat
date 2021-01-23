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

	getopt.SetParameters("<file>")
	getopt.Parse()

	args := getopt.Args()

	if *helpFlag {
		getopt.PrintUsage(os.Stderr)
		os.Exit(0)
	}

	if *versionFlag {
		var msg string

		if gitTag != "" {
			msg = fmt.Sprintf("%s - %s - %s", gitTag, sha1ver, buildTime)
		} else {
			msg = fmt.Sprintf("%s - %s", sha1ver, buildTime)
		}

		fmt.Println(msg)
		os.Exit(0)
	}

	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: no file given\n\n")
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	if *headRows > 0 && *tailRows > 0 {
		fmt.Fprintf(os.Stderr, "Error: --head and --tail are mutually exclusive\n\n")
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	inFile, err := os.Open(args[0])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error opening file: ", err)
		os.Exit(2)
	}

	result, err := lib.ProcessFile(inFile, *countFlag, *headRows, *tailRows)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error processing file: ", err)
		os.Exit(1)
	}

	if *countFlag {
		count := result.(int)
		fmt.Printf("%d %s\n", count, args[0])
	} else if *headRows > 0 || *tailRows > 0 {
		fragment := result.(lib.BinaryFileFragment)

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

		err = fragment.Write(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error writing new file: ", err)
			os.Exit(2)
		}
	}
}
