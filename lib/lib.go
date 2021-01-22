package lib

import (
	"encoding/binary"
	"errors"
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

func (binaryFile *BinaryFileFragment) Write(file io.Writer) error {
	_, err := file.Write(FileSignature)
	if err != nil {
		return err
	}

	err = binaryFile.Definitions.Write(file)
	if err != nil {
		return err
	}

	_, err = file.Write(binaryFile.Data)
	if err != nil {
		return err
	}

	return nil
}

func ProcessFile(file *os.File, countFlag bool, headRows int) (interface{}, error) {
	_, err := ReadSignature(file)
	if err != nil {
		return 0, err
	}

	definitions, err := ReadColumnDefinitions(file)
	if err != nil {
		return 0, err
	}

	result, err := iterateRows(file, definitions, countFlag, headRows)
	if err != nil {
		return result, err
	}

	return result, nil
}

func iterateRows(file *os.File, definitions ColumnDefinitions, countFlag bool, headRows int) (interface{}, error) {
	count := 0
	var data []byte

	var rowLen uint32
	err := binary.Read(file, binary.LittleEndian, &rowLen)
	if err != nil {
		return 0, err
	}

	iteration := 0

	for rowLen > 0 {
		if headRows > 0 && iteration >= headRows {
			break
		}

		var row []byte

		bitfield, err := ReadBitfield(file, definitions.NumberOfColumns)
		if err != nil {
			return 0, err
		}

		nullValues := DecodeBitfield(bitfield)

		if headRows > 0 {
			lenBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(lenBytes, rowLen)
			row = append(row, lenBytes...)
			row = append(row, bitfield...)
		}

		for i, width := range definitions.Widths {
			if nullValues[i] {
				continue
			}

			var columnWidth uint32

			if width == math.MaxUint32 {
				err = binary.Read(file, binary.LittleEndian, &columnWidth)
				if err != nil {
					return 0, err
				}
			} else {
				columnWidth = width
			}

			var column = make([]byte, columnWidth)

			err = binary.Read(file, binary.LittleEndian, &column)
			if err != nil {
				return 0, err
			}

			if width != columnWidth {
				widthBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(widthBytes, columnWidth)
				row = append(row, widthBytes...)
			}

			row = append(row, column...)
		}

		if countFlag {
			count += 1
		} else if headRows > 0 {
			data = append(data, row...)
		}

		err = binary.Read(file, binary.LittleEndian, &rowLen)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}

		iteration += 1
	}

	if countFlag {
		return count, nil
	} else if headRows > 0 {
		fragment := BinaryFileFragment{
			Definitions: definitions,
			Data:        data,
		}
		return fragment, nil
	} else {
		return -1, nil
	}
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
