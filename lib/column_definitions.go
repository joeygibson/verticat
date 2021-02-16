package lib

import (
	"encoding/binary"
	"io"
	"os"
)

type ColumnDefinitions struct {
	HeaderLength    uint32
	Version         uint16
	Filler          byte
	NumberOfColumns uint16
	Widths          []uint32
	NewColumnOrder  []uint
}

func (c ColumnDefinitions) Write(file io.Writer) (err error) {
	err = binary.Write(file, binary.LittleEndian, c.HeaderLength)
	if err != nil {
		return
	}

	err = binary.Write(file, binary.LittleEndian, c.Version)
	if err != nil {
		return
	}

	err = binary.Write(file, binary.LittleEndian, c.Filler)
	if err != nil {
		return
	}

	err = binary.Write(file, binary.LittleEndian, c.NumberOfColumns)
	if err != nil {
		return
	}

	if c.NewColumnOrder == nil {
		err = binary.Write(file, binary.LittleEndian, c.Widths)
		if err != nil {
			return
		}
	} else {
		// Re-order the column widths based on the specified order
		orderedWidths := make([]uint32, len(c.Widths))

		for i, val := range c.NewColumnOrder {
			orderedWidths[i] = c.Widths[val - 1]
		}

		err = binary.Write(file, binary.LittleEndian, orderedWidths)
		if err != nil {
			return
		}
	}

	return nil
}

func ReadColumnDefinitions(file *os.File, newColumnOrder []uint) (ColumnDefinitions, error) {
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
		NewColumnOrder: newColumnOrder,
	}

	return definitions, nil
}
