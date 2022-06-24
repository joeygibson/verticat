use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{Read, Write};
use std::path::Path;

use vertica_native_file::VerticaNativeFile;

mod column_conversion;
mod column_definitions;
mod column_type;
mod column_types;
mod file_signature;
mod vertica_native_file;

fn read_variable(reader: &mut impl Read, length: usize) -> Result<Vec<u8>, Box<dyn Error>> {
    let mut vec = vec![0u8; length];
    reader.read(vec.as_mut_slice())?;

    Ok(vec)
}

fn read_u32(reader: &mut impl Read) -> io::Result<u32> {
    let mut bytes: [u8; 4] = [0; 4];

    let res = reader.read(&mut bytes)?;

    if res != 4 {
        Ok(0)
    } else {
        Ok(u32::from_le_bytes(bytes))
    }
}

fn read_u16(reader: &mut impl Read) -> io::Result<u16> {
    let mut bytes: [u8; 2] = [0; 2];

    reader.read_exact(&mut bytes)?;

    Ok(u16::from_le_bytes(bytes))
}

fn read_u8(reader: &mut impl Read) -> io::Result<u8> {
    let mut bytes: [u8; 1] = [0; 1];

    reader.read_exact(&mut bytes)?;

    Ok(u8::from_le_bytes(bytes))
}

pub fn process_file(
    input: &str,
    count: bool,
    head: Option<i32>,
    tail: Option<i32>,
) -> Result<(), String> {
    if !Path::new(input).exists() {
        return Err(format!("input file {} does not exist", input));
    }

    let mut input_file = match File::open(&input) {
        Ok(i) => i,
        Err(e) => return Err(e.to_string()),
    };

    let native_file = match VerticaNativeFile::from_reader(&mut input_file) {
        Ok(i) => i,
        Err(e) => return Err(e.to_string()),
    };

    // `--count`
    if count {
        println!("{} {}", native_file.count(), input);
        return Ok(());
    }

    let mut stdout = io::stdout().lock();

    // `--head n`
    if head.is_some() {
        stdout.write_all(&native_file.generate_output().unwrap()).unwrap();

        native_file.take(head.unwrap() as usize).for_each(|row| {
            let data = row.generate_output().unwrap();
            stdout.write_all(&data).unwrap();
        });

        return Ok(());
    }

    // `--tail n`
    if tail.is_some() {
        let rows = native_file.count();

        let mut input_file = match File::open(&input) {
            Ok(i) => i,
            Err(e) => return Err(e.to_string()),
        };

        let native_file = match VerticaNativeFile::from_reader(&mut input_file) {
            Ok(i) => i,
            Err(e) => return Err(e.to_string()),
        };

        let rows_to_skip = rows - (tail.unwrap() as usize);

        native_file.skip(rows_to_skip).for_each(|row| {
            let data = row.generate_output().unwrap();
            stdout.write_all(&data).unwrap();
        });

        return Ok(());
    }

    return Ok(());
}
