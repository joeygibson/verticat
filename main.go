package main

import (
	"fmt"
	"github.com/joeygibson/verticat/lib"
	"github.com/pborman/getopt/v2"
	"io"
	"os"
)

func main() {
	helpFlag := getopt.BoolLong("help", 'H', "show help")
	countFlag := getopt.BoolLong("count", 'c', "count rows")
	headRows := getopt.IntLong("head", 'h', 0,"take the first n rows")
	outFileName := getopt.StringLong("output", 'o', "", "write head/tail results to this file")
	versionFlag := getopt.BoolLong("version", 'v', "show version")

	getopt.Parse()
	args := getopt.Args()

	if *helpFlag {
		fmt.Println("Help!")
		getopt.PrintUsage(os.Stderr)
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println("Version!")
		os.Exit(0)
	}

	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "no file given")
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	inFile, err := os.Open(args[0])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error opening file: ", err)
		os.Exit(2)
	}

	result, err := lib.ProcessFile(inFile, *countFlag, *headRows)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error processing file: ", err)
		os.Exit(1)
	}

	if *countFlag {
		count := result.(int)
		fmt.Printf("%d %s\n", count, args[0])
	} else if *headRows > 0 {
		fragment := result.(lib.BinaryFileFragment)

		var output io.Writer

		if *outFileName != "" {
			_, err := os.Stat(*outFileName)
			if os.IsExist(err) {
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

		fragment.Write(output)
	}
}
