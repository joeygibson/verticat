package lib

import (
	"os"
	"reflect"
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
	file, err := os.Open("../private-data/sample")
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
	type args struct {
		file      string
		countFlag bool
		headRows  int
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "all rows",
			args:    args{file: "../private-data/sample", countFlag: false, headRows: 475},
			want:    102728,
			wantErr: false},
		{name: "5 rows",
			args:    args{file: "../private-data/sample", countFlag: false, headRows: 5},
			want:    1102,
			wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			got, err := ProcessFile(file, tt.args.countFlag, tt.args.headRows)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fragment := got.(BinaryFileFragment)
			actual := len(fragment.Data)

			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("ProcessFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
