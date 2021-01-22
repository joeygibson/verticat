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
	//file, err := os.Open("../../verticareader/private-data/4k/flow_stats-0-1550160360")
	file, err := os.Open("../xxx.bin")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}

	defer file.Close()

	res, err := ProcessFile(file, true, 0)
	if err != nil {
		t.Fatal("error processing file: ", err)
	}

	expected := 475
	actual := res.(int)

	if actual != expected {
		t.Fatalf("wrong count; expected %d, got %d", expected, actual)
	}
}

func TestHead(t *testing.T) {
	file, err := os.Open("../../verticareader/private-data/4k/flow_stats-0-1550160360")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}

	defer file.Close()

	res, err := ProcessFile(file, false, 5)
	if err != nil {
		t.Fatal("error processing file: ", err)
	}

	expected := 1102
		//102728
	fragment := res.(BinaryFileFragment)
	actual := len(fragment.Data)

	if actual != expected {
		t.Fatalf("wrong byte count; expected %d, got %d", expected, actual)
	}

	file, err = os.Create("../xxx.bin")
	if err != nil {
		t.Fatal("error opening new file: ", err)
	}

	defer file.Close()

	fragment.Write(file)
}
