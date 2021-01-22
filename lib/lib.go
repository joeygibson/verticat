package lib

import (
	"encoding/binary"
	"errors"
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

func ProcessFile(file *os.File) error {
	_, err := ReadSignature(file)
	if err != nil {
		return err
	}

	definitions, err := ReadColumnDefinitions(file)
	if err != nil {
		return err
	}

	var rowLen uint32
	err = binary.Read(file, binary.LittleEndian, &rowLen)
	if err != nil {
		return err
	}

	for rowLen > 0 {
		var data []byte

		bitfield, err := ReadBitfield(file, definitions.NumberOfColumns)
		if err != nil {
			return err
		}

		nullValues := DecodeBitfield(bitfield)

		for i, width := range definitions.Widths {
			if nullValues[i] {
				continue
			}

			var columnWidth uint32

			if width == math.MaxUint32 {
				err = binary.Read(file, binary.LittleEndian, &columnWidth)
				if err != nil {
					return err
				}
			} else {
				columnWidth = width
			}

			var column = make([]byte, columnWidth)

			err = binary.Read(file, binary.LittleEndian, &column)
			if err != nil {
				return err
			}

			data = append(data, column...)
		}

		err = binary.Read(file, binary.LittleEndian, &rowLen)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
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
