package lib

import (
	"os"
	"testing"
)

func TestReadColumnDefinitions(t *testing.T) {
	file, err := os.Open("../test-data/sample")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}

	defer file.Close()

	match, err := ReadSignature(file)
	if err != nil {
		t.Fatal("error matching file: ", err)
	}

	if !match {
		t.Fatal("signature didn't match")
	}

	definitions, err := ReadColumnDefinitions(file, nil)
	if err != nil {
		t.Fatal("error reading column definitions: ", err)
	}

	if definitions.HeaderLength != 305 {
		t.Errorf("wrong header length: %d", definitions.HeaderLength)
	}

	if definitions.Version != 1 {
		t.Errorf("wrong version: %d", definitions.Version)
	}

	if definitions.NumberOfColumns != 75 {
		t.Errorf("wrong column count: %d", definitions.NumberOfColumns)
	}

	if len(definitions.Widths) != int(definitions.NumberOfColumns) {
		t.Errorf("wrong column count: %d != %d", len(definitions.Widths),
			definitions.NumberOfColumns)
	}
}
