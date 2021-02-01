package lib

import (
	"bytes"
	"os"
	"reflect"
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

//func TestCat(t *testing.T) {
//	file, err := os.Open("../test-data/sample")
//	if err != nil {
//		t.Fatal("couldn't open file", err)
//	}
//
//	fileInfo, _ := file.Stat()
//	fileInfo.Size()
//
//	defer file.Close()
//
//	var buf bytes.Buffer
//
//	err = Cat(file, &buf, true)
//	if err != nil {
//		t.Fatal("error reading rows", err)
//	}
//
//	fullLen := buf.Len()
//
//	resetFilePosition(file, 0)
//	buf.Reset()
//
//	err = Cat(file, &buf, false)
//	if err != nil {
//		t.Fatal("error reading rows", err)
//	}
//
//	noMetaLen := buf.Len()
//
//	if fullLen == noMetaLen {
//		t.Fatal("lengths should differ")
//	}
//}

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

			err = Head(file, &buf, tt.args.rowsToTake, tt.args.shouldWriteMetaData)
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

			err = Cat(file, &buf, tt.args.shouldWriteMetaData)
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
