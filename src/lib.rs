use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{stdout, BufWriter, Read, Write};
use std::path::Path;
use std::str::FromStr;

use column_types::ColumnTypes;
use vertica_native_file::VerticaNativeFile;
use crate::vertica_native_file::Row;

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
    force: bool,
) -> Result<(), String> {
    println!(
        "input: {:?}, count: {}, head: {:?}, tail: {:?}, force: {}",
        input, count, head, tail, force
    );

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

    if count {
        println!("{} {}", native_file.count(), input);
        return Ok(());
    }

    let mut stdout = io::stdout().lock();

    if head.is_some() {
        native_file.take(head.unwrap() as usize)
            .for_each(|row| {
                let data = row.generate_output().unwrap();
                stdout.write_all(&data);
            });

        return Ok(());
    }

    if tail.is_some() {
        // native_file.drop(head.unwrap() as usize)
        //     .for_each(|row| {
        //         let data = row.generate_output().unwrap();
        //         stdout.write_all(&data);
        //     });

        return Ok(());
    }

    return Ok(());
}
