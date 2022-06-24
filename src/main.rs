use clap::{Arg, ArgMatches, Command};

use verticat::process_file;

const VERSION: &'static str = env!("CARGO_PKG_VERSION");

const COUNT: &'static str = "count";
const HEAD: &'static str = "head";
const INPUT: &'static str = "input";
const TAIL: &'static str = "tail";

fn main() {
    let app = Command::new("verticat")
        .version(VERSION)
        .about("count/head/tail Vertica native binary files")
        .arg_required_else_help(true)
        .arg(
            Arg::new(INPUT)
                .takes_value(true)
                .multiple_values(true)
                .value_name("filename")
                .help("The file(s) to process"),
        )
        .arg(
            Arg::new(COUNT)
                .takes_value(false)
                .short('c')
                .long(COUNT)
                .help("count rows"),
        )
        .arg(
            Arg::new(HEAD)
                .takes_value(true)
                .short('h')
                .long(HEAD)
                .value_name("n")
                .help("take the first n rows"),
        )
        .arg(
            Arg::new(TAIL)
                .takes_value(true)
                .short('t')
                .long(TAIL)
                .value_name("n")
                .help("take the last n rows"),
        );

    let args = app.get_matches();

    let input: Vec<&str> = args.values_of(INPUT).unwrap().collect();
    let count = args.is_present(COUNT);
    let head = get_param(&args, HEAD);
    let tail = get_param(&args, TAIL);

    if (head.is_some() && tail.is_some()) || (count && (head.is_some() || tail.is_some())) {
        eprintln!("count, head, and tail are mutually exclusive");
        return;
    }

    for filename in input {
        match process_file(filename, count, head, tail) {
            Ok(_) => {}
            Err(e) => eprintln!("Error: {}", e),
        }
    }
}

fn get_param(args: &ArgMatches, name: &str) -> Option<i32> {
    return if args.is_present(name) {
        match args.value_of(name) {
            Some(rows) => Some(rows.parse::<i32>().expect("number out of range")),
            None => Some(5),
        }
    } else {
        None
    };
}
