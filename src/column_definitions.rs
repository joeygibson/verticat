use std::error::Error;
use std::fs::File;

use crate::{read_u16, read_u32, read_u8};

#[derive(Debug)]
pub struct ColumnDefinitions {
    header_length: u32,
    version: u16,
    filler: u8,
    number_of_columns: u16,
    pub column_widths: Vec<u32>,
}

impl ColumnDefinitions {
    pub fn from_reader(reader: &mut File) -> Result<Self, Box<dyn Error>> {
        let header_length: u32 = read_u32(reader)?;
        let version = read_u16(reader)?;

        // drop the filler
        let filler = read_u8(reader)?;

        let number_of_columns = read_u16(reader)?;

        let mut column_widths: Vec<u32> = vec![];

        for _ in 0..number_of_columns {
            let value = read_u32(reader)?;
            column_widths.push(value);
        }

        Ok(ColumnDefinitions {
            header_length,
            version,
            filler,
            number_of_columns,
            column_widths,
        })
    }

    pub fn generate_output(&self) -> Result<Vec<u8>, Box<dyn Error>> {
        let mut record: Vec<u8> = vec![];

        record.append(&mut self.header_length.to_le_bytes().to_vec());
        record.append(&mut self.version.to_le_bytes().to_vec());
        record.append(&mut self.filler.to_le_bytes().to_vec());
        record.append(&mut self.number_of_columns.to_le_bytes().to_vec());

        let mut width_bytes: Vec<u8> = self
            .column_widths
            .iter()
            .flat_map(|bucket| bucket.to_le_bytes().to_vec())
            .collect();

        eprintln!("width_bytes: {}", &width_bytes.len());
        record.append(&mut width_bytes);

        Ok(record)
    }
}

#[cfg(test)]
mod tests {
    // use std::io::{Seek, SeekFrom};

    use std::io::{Seek, SeekFrom};

    use crate::column_definitions::ColumnDefinitions;

    #[test]
    fn test_read_from_good_file() {
        use std::fs::File;

        let mut file = File::open("test-data/all-types.bin").unwrap();

        file.seek(SeekFrom::Start(11)).unwrap();

        let expected_column_widths: [u32; 14] =
            [8, 8, 10, 4294967295, 1, 8, 8, 8, 8, 8, 4294967295, 3, 24, 8];

        let column_definitions = ColumnDefinitions::from_reader(&mut file).unwrap();

        for (index, expected_value) in expected_column_widths.iter().enumerate() {
            let value = column_definitions.column_widths[index];

            assert_eq!(*expected_value, value);
        }
    }
}
