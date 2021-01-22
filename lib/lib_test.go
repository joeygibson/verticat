package lib

import (
	"os"
	"testing"
)

func TestReadSignature(t *testing.T) {
	file, err := os.Open("../private-data/sample")
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
}

func TestCount(t *testing.T) {
	file, err := os.Open("../../verticareader/private-data/4k/flow_stats-0-1550160360")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}

	defer file.Close()

	res, err := ProcessFile(file, true)
	if err != nil {
		t.Fatal("error processing file: ", err)
	}

	expected := 475

	if res.(int) != expected {
		t.Fatalf("wrong count; expected %d, got %d", expected, res.(int))
	}
}
