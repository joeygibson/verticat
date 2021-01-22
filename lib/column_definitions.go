package lib

import (
	"encoding/binary"
	"os"
)

type ColumnDefinitions struct {
	HeaderLength    uint32
	Version         uint16
	Filler          byte
	NumberOfColumns uint16
	Widths          []uint32
}

func ReadColumnDefinitions(file *os.File) (ColumnDefinitions, error) {
	var headerLength uint32
	var version uint16
	var filler byte
	var numberOfColumns uint16

	err := binary.Read(file, binary.LittleEndian, &headerLength)
	if err != nil {
		return ColumnDefinitions{}, err
	}

	err = binary.Read(file, binary.LittleEndian, &version)
	if err != nil {
		return ColumnDefinitions{}, err
	}

	err = binary.Read(file, binary.LittleEndian, &filler)
	if err != nil {
		return ColumnDefinitions{}, err
	}

	err = binary.Read(file, binary.LittleEndian, &numberOfColumns)
	if err != nil {
		return ColumnDefinitions{}, err
	}

	widths := make([]uint32, numberOfColumns)

	err = binary.Read(file, binary.LittleEndian, &widths)
	if err != nil {
		return ColumnDefinitions{}, err
	}

	definitions := ColumnDefinitions{
		HeaderLength:    headerLength,
		Version:         version,
		Filler:          filler,
		NumberOfColumns: numberOfColumns,
		Widths:          widths,
	}

	return definitions, nil
}
