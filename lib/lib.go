package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var FileSignature = []byte{0x4e, 0x41, 0x54, 0x49, 0x56, 0x45, 0x0a, 0xff, 0x0d, 0x0a, 0x00}

func ReadSignature(file *os.File) (bool, error) {
	signature := make([]byte, len(FileSignature))

	read, err := file.Read(signature)
	if err != nil {
		return false, err
	}

	if read != len(FileSignature) {
		return false, errors.New("wrong number of bytes read")
	}

	match := true

	for i, v := range signature {
		if v != FileSignature[i] {
			match = false
			break
		}
	}

	return match, nil
}

func readRow(file *os.File, definitions ColumnDefinitions) ([]byte, error) {
	var flatRow []byte
	row := make([][]byte, definitions.NumberOfColumns)

	var rowLen uint32
	err := binary.Read(file, binary.LittleEndian, &rowLen)
	if err != nil {
		return flatRow, err
	}

	bitfield, err := ReadBitfield(file, definitions.NumberOfColumns)
	if err != nil {
		return flatRow, err
	}

	nullValues := DecodeBitfield(bitfield)

	reorderedBitField := reorderBitfield(nullValues, definitions.NewColumnOrder)

	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, rowLen)
	flatRow = append(flatRow, lenBytes...)
	flatRow = append(flatRow, reorderedBitField...)

	for i, width := range definitions.Widths {
		if nullValues[i] {
			continue
		}

		var columnWidth uint32

		if width == math.MaxUint32 {
			err = binary.Read(file, binary.LittleEndian, &columnWidth)
			if err != nil {
				return flatRow, err
			}
		} else {
			columnWidth = width
		}

		var column = make([]byte, columnWidth)

		err = binary.Read(file, binary.LittleEndian, &column)
		if err != nil {
			return flatRow, err
		}

		var rowPos uint

		if definitions.NewColumnOrder != nil {
			rowPos = definitions.NewColumnOrder[i] - 1
		} else {
			rowPos = uint(i)
		}

		if width != columnWidth {
			widthBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(widthBytes, columnWidth)
			row[rowPos] = widthBytes
			row[rowPos] = append(row[rowPos], column...)
		} else {
			row[rowPos] = column
		}
	}

	for _, column := range row {
		flatRow = append(flatRow, column...)
	}

	return flatRow, nil
}

func reorderBitfield(nullValues []bool, newColumnOrder []uint) []byte {
	reorderedNullValues := make([]bool, len(nullValues))

	for i, val := range newColumnOrder {
		reorderedNullValues[i] = nullValues[val-1]
	}

	var bitfield []byte
	sliceLen := len(reorderedNullValues)

	chunkSize := 8

	for i := 0; i < sliceLen; i += chunkSize {
		var b byte

		end := i + chunkSize

		if end > sliceLen {
			end = sliceLen
		}

		for i1, value := range reorderedNullValues[i:end] {
			if value {
				b |= 1 << int8(math.Abs(float64(i1)-7))
			}
		}

		bitfield = append(bitfield, b)
	}

	return bitfield
}

func CountRows(inputFiles []*os.File) error {
	for _, file := range inputFiles {

		count, err := countRows(file)
		if err != nil {
			return errors.New(fmt.Sprintf("error counting rows for %s: %s",
				file.Name(), err.Error()))
		}

		fmt.Printf("%d %s\n", count, file.Name())
	}

	return nil
}

func countRows(file *os.File) (int, error) {
	_, err := ReadSignature(file)
	if err != nil {
		return -1, errors.New("error reading file signature: " + err.Error())
	}

	definitions, err := ReadColumnDefinitions(file, nil)
	if err != nil {
		return -1, errors.New("error reading column definitions: " + err.Error())
	}
	count := 0

	_, err = readRow(file, definitions)

	for err == nil {
		count += 1

		_, err = readRow(file, definitions)
	}

	if err == io.EOF {
		err = nil
	}

	return count, err
}

func resetFilePosition(file *os.File, pos int64) error {
	curPos, err := file.Seek(pos, io.SeekStart)
	if err != nil {
		return err
	}

	if curPos != pos {
		return errors.New("error resetting file position")
	}

	return nil
}

func Cat(file *os.File, writer io.Writer, shouldWriteMetaData bool, newColumnOrder []uint) error {
	err := Head(file, writer, math.MaxInt64, shouldWriteMetaData, newColumnOrder)
	if err != nil {
		return err
	}

	return nil
}

func PrintHeader(file *os.File, writer io.Writer) error {
	valid, err := ReadSignature(file)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New(fmt.Sprintf("invalid file signature: %s", file.Name()))
	}

	definitions, err := ReadColumnDefinitions(file, nil)
	if err != nil {
		return err
	}

	for _, width := range definitions.Widths {
		if width == math.MaxUint32 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, width)
		}
	}

	fmt.Fprintf(writer, "\n")

	return nil
}

func Head(file *os.File, writer io.Writer, rowsToTake int, shouldWriteMetaData bool,
	newColumnOrder []uint) error {
	valid, err := ReadSignature(file)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New(fmt.Sprintf("invalid file signature: %s", file.Name()))
	}

	definitions, err := ReadColumnDefinitions(file, newColumnOrder)
	if err != nil {
		return err
	}

	err = writeMetaData(shouldWriteMetaData, writer, definitions)
	if err != nil {
		return err
	}

	i := 0

	for i < rowsToTake {
		row, err := readRow(file, definitions)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		_, err = writer.Write(row)
		if err != nil {
			return err
		}

		i += 1
	}

	return nil
}

func Tail(file *os.File, writer io.Writer, rowsToTake int, shouldWriteMetaData bool,
	newColumnOrder []uint) error {
	totalRows, err := countRows(file)
	if err != nil {
		return err
	}

	err = resetFilePosition(file, 0)
	if err != nil {
		return err
	}

	valid, err := ReadSignature(file)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New(fmt.Sprintf("invalid file signature: %s", file.Name()))
	}

	definitions, err := ReadColumnDefinitions(file, newColumnOrder)
	if err != nil {
		return err
	}

	err = writeMetaData(shouldWriteMetaData, writer, definitions)
	if err != nil {
		return err
	}

	rowsToSkip := totalRows - rowsToTake

	i := 0

	for i < totalRows {
		row, err := readRow(file, definitions)
		if err != nil {
			return err
		}

		if i >= rowsToSkip {
			_, err = writer.Write(row)
			if err != nil {
				return err
			}
		}

		i += 1
	}

	return nil
}

func writeMetaData(shouldWrite bool, writer io.Writer, definitions ColumnDefinitions) error {
	if shouldWrite {
		_, err := writer.Write(FileSignature)
		if err != nil {
			return err
		}

		err = definitions.Write(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func DecodeBitfield(bitfield []byte) []bool {
	var nullValues []bool

	for _, b := range bitfield {
		for i := 7; i >= 0; i-- {
			isNull := b&(1<<i) != 0
			nullValues = append(nullValues, isNull)
		}
	}

	return nullValues
}

func ReadBitfield(file *os.File, numberOfColumns uint16) ([]byte, error) {
	bitfieldLength := numberOfColumns / 8
	if numberOfColumns%8 != 0 {
		bitfieldLength += 1
	}

	var bitfield = make([]byte, bitfieldLength)
	err := binary.Read(file, binary.BigEndian, &bitfield)
	if err != nil {
		return bitfield, err
	}

	return bitfield, nil
}
