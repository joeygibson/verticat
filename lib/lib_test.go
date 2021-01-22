package lib

import (
	"os"
	"testing"
)

func TestReadSignature(t *testing.T) {
	file, err := os.Open("../private-data/flow_stats-0-1550160360")
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
