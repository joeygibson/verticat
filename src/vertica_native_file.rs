use std::error::Error;
use std::fs::File;
use std::io::Read;

use crate::column_definitions::ColumnDefinitions;
use crate::column_types::ColumnTypes;
use crate::file_signature::FileSignature;
use crate::{read_u32, read_u8, read_variable};

pub struct VerticaNativeFile<'a> {
    _signature: FileSignature,
    pub definitions: ColumnDefinitions,
    file: &'a File,
}

impl<'a> VerticaNativeFile<'a> {
    pub fn from_reader(reader: &'a mut File) -> Result<Self, Box<dyn Error>> {
        let signature = FileSignature::from_reader(reader)?;
        let definitions = ColumnDefinitions::from_reader(reader)?;

        Ok(VerticaNativeFile {
            _signature: signature,
            definitions,
            file: reader,
        })
    }
}

impl<'a> Iterator for VerticaNativeFile<'a> {
    type Item = Row;

    fn next(&mut self) -> Option<Self::Item> {
        let row_length = read_u32(&mut self.file).unwrap();

        if row_length <= 0 {
            return None;
        }

        match Row::from_reader(self.file, &self.definitions.column_widths) {
            Ok(row) => Some(row),
            Err(e) => {
                eprintln!("reading data: {}", e);
                None
            }
        }
    }
}

#[derive(Debug)]
pub struct Row {
    bitfield: Vec<u8>,
    data: Vec<u8>,
}

impl Row {
    fn from_reader(
        mut reader: impl Read,
        column_widths: &Vec<u32>,
    ) -> Result<Self, Box<dyn Error>> {
        let mut data: Vec<u8> = vec![];

        let (null_values, bitfield) = Row::read_bitfield(&mut reader, &column_widths)?;

        for (index, width) in column_widths.iter().enumerate() {
            if null_values[index] {
                continue;
            }

            let mut column: Vec<u8> = vec![];

            let column_width = if *width == std::u32::MAX {
                read_u32(&mut reader)?
            } else {
                *width
            };

            for _ in 0..column_width {
                let value = read_u8(&mut reader)?;

                column.push(value);
            }

            data.append(&mut column)
        }

        Ok(Row { bitfield, data })
    }

    fn read_bitfield(
        mut reader: &mut impl Read,
        column_widths: &Vec<u32>,
    ) -> Result<(Vec<bool>, Vec<u8>), Box<dyn Error>> {
        let mut null_values: Vec<bool> = vec![];

        let bitfield_length =
            (column_widths.len() / 8) + if column_widths.len() % 8 == 0 { 0 } else { 1 };
        let bitfield = read_variable(&mut reader, bitfield_length as usize)?;

        for byte in bitfield.iter() {
            for i in (0..8).rev() {
                null_values.push(byte & (1 << i) != 0);
            }
        }

        Ok((null_values, bitfield))
    }

    pub fn generate_output(&self) -> Result<Vec<u8>, Box<dyn Error>> {
        let mut record: Vec<u8> = vec![];

        record.append(&mut self.bitfield.clone());
        record.append(&mut self.data.clone());

        Ok(record)
    }
}

#[cfg(test)]
mod tests {
    use std::fs::File;

    use crate::vertica_native_file::VerticaNativeFile;

    #[test]
    fn test_read_from_good_file() {
        let mut file = File::open("data/all-types.bin").unwrap();

        let file = VerticaNativeFile::from_reader(&mut file).unwrap();

        for row in file {
            assert_eq!(14, row.data.len());
        }
    }

    #[test]
    fn test_read_from_good_file_with_nulls() {
        let mut file = File::open("data/all-types-with-nulls.bin").unwrap();

        let file = VerticaNativeFile::from_reader(&mut file).unwrap();

        for row in file {
            assert_eq!(14, row.data.len());
        }
    }
}
