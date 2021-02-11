package lib

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadSignature(t *testing.T) {
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
}

func TestCount(t *testing.T) {
	file, err := os.Open("../test-data/sample")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}

	defer file.Close()

	res, err := countRows(file)
	if err != nil {
		t.Fatal("error processing file: ", err)
	}

	expected := 475
	actual := res

	if actual != expected {
		t.Fatalf("wrong count; expected %d, got %d", expected, actual)
	}
}

func TestPrintHeader(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantLen int
		want    []string
		wantErr bool
	}{
		{name: "all-types file",
			file:    "../test-data/all-types.bin",
			wantLen: 14,
			want:    []string{"8", "8", "10", "-1", "1", "8", "8", "8", "8", "8", "-1", "3", "24", "8"},
			wantErr: false},
		{name: "sample file",
			file:    "../test-data/sample",
			wantLen: 75,
			want: []string{"8", "8", "8", "4", "8", "8", "8", "8", "-1", "4", "-1", "8", "-1", "4",
				"-1", "8", "-1", "4", "4", "4", "8", "-1", "-1", "4", "4", "8", "8", "4", "4", "4", "4",
				"-1", "-1", "4", "4", "8", "8", "4", "4", "4", "4", "4", "4", "8", "8", "-1", "-1", "-1",
				"4", "-1", "4", "-1", "4", "4", "-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1",
				"-1", "-1", "-1", "-1", "2", "2", "2", "2", "4", "4", "1"},
			wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			var buf bytes.Buffer

			err = PrintHeader(file, &buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			outStr := strings.TrimSpace(buf.String())
			chunks := strings.Split(outStr, "\n")

			if len(chunks) != tt.wantLen {
				t.Errorf("wrong number of fields; got = %v, want %v", len(chunks), tt.wantLen)
			}

			for i, exp := range tt.want {
				if exp != chunks[i] {
					t.Errorf("PrintHeader() got = %v, want %v", chunks[i], exp)
				}
			}
		})
	}
}

func TestHead(t *testing.T) {
	type args struct {
		file                string
		rowsToTake          int
		shouldWriteMetaData bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "all rows",
			args:    args{file: "../test-data/sample", rowsToTake: 475, shouldWriteMetaData: false},
			want:    102728,
			wantErr: false},
		{name: "all rows, with metadata",
			args:    args{file: "../test-data/sample", rowsToTake: 475, shouldWriteMetaData: true},
			want:    103048,
			wantErr: false},
		{name: "5 rows",
			args:    args{file: "../test-data/sample", rowsToTake: 5, shouldWriteMetaData: false},
			want:    1102,
			wantErr: false},
		{name: "5 rows, with metadata",
			args:    args{file: "../test-data/sample", rowsToTake: 5, shouldWriteMetaData: true},
			want:    1422,
			wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			var buf bytes.Buffer

			err = Head(file, &buf, tt.args.rowsToTake, tt.args.shouldWriteMetaData, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Head() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actual := buf.Len()
			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("Head() got = %v, want %v", actual, tt.want)
			}
		})
	}
}

func TestTail(t *testing.T) {
	type args struct {
		file                string
		rowsToTake          int
		shouldWriteMetaData bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "all rows",
			args:    args{file: "../test-data/sample", rowsToTake: 475, shouldWriteMetaData: false},
			want:    102728,
			wantErr: false},
		{name: "all rows, with metadata",
			args:    args{file: "../test-data/sample", rowsToTake: 475, shouldWriteMetaData: true},
			want:    103048,
			wantErr: false},
		{name: "5 rows",
			args:    args{file: "../test-data/sample", rowsToTake: 5, shouldWriteMetaData: false},
			want:    1070,
			wantErr: false},
		{name: "5 rows, with metadata",
			args:    args{file: "../test-data/sample", rowsToTake: 5, shouldWriteMetaData: true},
			want:    1390,
			wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			var buf bytes.Buffer

			err = Tail(file, &buf, tt.args.rowsToTake, tt.args.shouldWriteMetaData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actual := buf.Len()

			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("Tail() got = %v, want %v", actual, tt.want)
			}
		})
	}
}

func TestCatWithReorder(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		newColumnOrder []uint
		shouldDiffer   bool
	}{
		{name: "no reordering, empty slice",
			fileName:       "../test-data/all-types.bin",
			newColumnOrder: nil,
			shouldDiffer:   false},
		{name: "no reordering",
			fileName:       "../test-data/all-types.bin",
			newColumnOrder: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			shouldDiffer:   false},
		{name: "reorder",
			fileName:       "../test-data/all-types.bin",
			newColumnOrder: []uint{1, 14, 2, 3, 4, 10, 5, 6, 7, 8, 9, 11, 12, 13},
			shouldDiffer:   true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.fileName)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			var buf bytes.Buffer

			err = Cat(file, &buf, true, tt.newColumnOrder)
			if err != nil {
				t.Fatal("error processing input file with reordering: ", err)
			}

			_, _ = file.Seek(0, io.SeekStart)

			var origOrderBuf bytes.Buffer

			err = Cat(file, &origOrderBuf, true, nil)
			if err != nil {
				t.Fatal("error processing input file with original ordering: ", err)
			}

			if reflect.DeepEqual(buf, origOrderBuf) && tt.shouldDiffer {
				t.Errorf("output should differ, but did not")
			}
		})
	}
}

func TestCat(t *testing.T) {
	type args struct {
		file                string
		shouldWriteMetaData bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "all rows",
			args:    args{file: "../test-data/sample", shouldWriteMetaData: false},
			want:    102728,
			wantErr: false},
		{name: "all rows, with metadata",
			args:    args{file: "../test-data/sample", shouldWriteMetaData: true},
			want:    103048,
			wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal("couldn't open file", err)
			}

			defer file.Close()

			var buf bytes.Buffer

			err = Cat(file, &buf, tt.args.shouldWriteMetaData, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actual := buf.Len()

			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("Tail() got = %v, want %v", actual, tt.want)
			}
		})
	}
}
