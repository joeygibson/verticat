package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	helpFlag := getopt.BoolLong("help", 'h', "show help")
	countFlag := getopt.BoolLong("count", 'c', "count rows")
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

	fmt.Println(*countFlag)

	_, err := os.Open(args[0])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error opening file: ", err)
		os.Exit(2)
	}
}
