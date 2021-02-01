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

type BinaryFileFragment struct {
	Definitions ColumnDefinitions
	Data        []byte
}

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

func (binaryFile *BinaryFileFragment) Write(file io.Writer, writeHeaders bool) error {
	if writeHeaders {
		_, err := file.Write(FileSignature)
		if err != nil {
			return err
		}

		err = binaryFile.Definitions.Write(file)
		if err != nil {
			return err
		}
	}
	_, err := file.Write(binaryFile.Data)
	if err != nil {
		return err
	}

	return nil
}

//func ProcessFile(file *os.File, outFile *os.File, countFlag bool, headRows int, tailRows int) (interface{}, error) {
//	_, err := ReadSignature(file)
//	if err != nil {
//		return 0, err
//	}
//
//	definitions, err := ReadColumnDefinitions(file)
//	if err != nil {
//		return 0, err
//	}
//
//	result, err := iterateRows(file, outFile, definitions, countFlag, headRows, tailRows)
//	if err != nil {
//		return result, err
//	}
//
//	return result, nil
//}

func readRow(file *os.File, definitions ColumnDefinitions) ([]byte, error) {
	var row []byte
	var rowLen uint32
	err := binary.Read(file, binary.LittleEndian, &rowLen)
	if err != nil {
		return row, err
	}

	bitfield, err := ReadBitfield(file, definitions.NumberOfColumns)
	if err != nil {
		return row, err
	}

	nullValues := DecodeBitfield(bitfield)

	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, rowLen)
	row = append(row, lenBytes...)
	row = append(row, bitfield...)

	for i, width := range definitions.Widths {
		if nullValues[i] {
			continue
		}

		var columnWidth uint32

		if width == math.MaxUint32 {
			err = binary.Read(file, binary.LittleEndian, &columnWidth)
			if err != nil {
				return row, err
			}
		} else {
			columnWidth = width
		}

		var column = make([]byte, columnWidth)

		err = binary.Read(file, binary.LittleEndian, &column)
		if err != nil {
			return row, err
		}

		if width != columnWidth {
			widthBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(widthBytes, columnWidth)
			row = append(row, widthBytes...)
		}

		row = append(row, column...)
	}

	return row, nil
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

	definitions, err := ReadColumnDefinitions(file)
	if err != nil {
		return -1, errors.New("error reading column definitions: " + err.Error())
	}
	count := 0

	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return count, err
	}

	_, err = readRow(file, definitions)

	for err == nil {
		count += 1

		_, err = readRow(file, definitions)
	}

	if err == io.EOF {
		err = nil
	}

	err = resetFilePosition(file, pos)

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

func Cat(file *os.File, writer io.Writer, shouldWriteMetaData bool) error {
	err := Head(file, writer, math.MaxInt64, shouldWriteMetaData)
	if err != nil {
		return err
	}

	return nil
}

func Head(file *os.File, writer io.Writer, rowsToTake int, shouldWriteMetaData bool) error {
	valid, err := ReadSignature(file)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New(fmt.Sprintf("invalid file signature: %s", file.Name()))
	}

	definitions, err := ReadColumnDefinitions(file)
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

func Tail(file *os.File, writer io.Writer, rowsToTake int, shouldWriteMetaData bool) error {
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

	definitions, err := ReadColumnDefinitions(file)
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

//func iterateRows(file *os.File, outFile *os.File, definitions ColumnDefinitions, countFlag bool,
//	headRows int, tailRows int) error {
//	count := 0
//
//	if tailRows > 0 {
//		pos, err := file.Seek(0, io.SeekCurrent)
//		if err != nil {
//			return err
//		}
//
//		err := iterateRows(file, outFile, definitions, true, 0, 0)
//
//		count = result.(int)
//
//		curPos, err := file.Seek(pos, io.SeekStart)
//		if err != nil {
//			return err
//		}
//
//		if curPos != pos {
//			return 0, errors.New("error resetting file position")
//		}
//	}
//
//	var rowLen uint32
//	err := binary.Read(file, binary.LittleEndian, &rowLen)
//	if err != nil {
//		return 0, err
//	}
//
//	iteration := 0
//
//	for rowLen > 0 {
//		if headRows > 0 && iteration >= headRows {
//			break
//		}
//
//		var row []byte
//
//		bitfield, err := ReadBitfield(file, definitions.NumberOfColumns)
//		if err != nil {
//			return 0, err
//		}
//
//		nullValues := DecodeBitfield(bitfield)
//
//		if headRows > 0 || (tailRows > 0 && iteration >= (count-tailRows)) {
//			lenBytes := make([]byte, 4)
//			binary.LittleEndian.PutUint32(lenBytes, rowLen)
//			row = append(row, lenBytes...)
//			row = append(row, bitfield...)
//		}
//
//		for i, width := range definitions.Widths {
//			if nullValues[i] {
//				continue
//			}
//
//			var columnWidth uint32
//
//			if width == math.MaxUint32 {
//				err = binary.Read(file, binary.LittleEndian, &columnWidth)
//				if err != nil {
//					return 0, err
//				}
//			} else {
//				columnWidth = width
//			}
//
//			var column = make([]byte, columnWidth)
//
//			err = binary.Read(file, binary.LittleEndian, &column)
//			if err != nil {
//				return 0, err
//			}
//
//			if width != columnWidth {
//				widthBytes := make([]byte, 4)
//				binary.LittleEndian.PutUint32(widthBytes, columnWidth)
//				row = append(row, widthBytes...)
//			}
//
//			row = append(row, column...)
//		}
//
//		if countFlag {
//			count += 1
//		} else if headRows > 0 || (tailRows > 0 && iteration >= (count-tailRows)) {
//			outFile.Write(row)
//		}
//
//		err = binary.Read(file, binary.LittleEndian, &rowLen)
//		if err != nil {
//			if err == io.EOF {
//				break
//			} else {
//				return 0, err
//			}
//		}
//
//		iteration += 1
//	}
//
//	if countFlag {
//		fmt.Printf("%d %s\n", count, file.Name())
//	}
//}

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
