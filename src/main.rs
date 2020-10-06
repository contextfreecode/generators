use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};

type Result<T> = std::result::Result<T, Box<dyn Error>>;

fn read_rows(reader: BufReader<File>)
        -> impl Iterator<Item = Result<Vec<String>>> {
    use genawaiter::{rc::gen, yield_};
    gen!({
        for line_result in reader.lines() {
            match line_result {
                Ok(line) => {
                    let row = line
                        .split('\t')
                        .map(|s| s.to_string())
                        .collect::<Vec<String>>();
                    yield_!(Ok(row));
                }
                Err(err) => {
                    yield_!(Err(err.into()));
                }
            }
        }
    })
    .into_iter()
}

fn main() -> Result<()> {
    let file = File::open("data/cities500.txt")?;
    let reader = BufReader::new(file);
    let mut population_south = 0i64;
    let mut population_total = 0i64;
    let mut row_count = 0;
    for row_result in read_rows(reader) {
        let row = row_result?;
        let latitude = row[4].parse::<f64>()?;
        let population = row[14].parse::<i64>()?;
        population_total += population;
        if latitude < 0.0 {
            population_south += population;
        }
        row_count += 1;
    }
    println!("South population: {}", population_south);
    println!("Total population: {}", population_total);
    println!("# of rows: {}", row_count);
    Ok(())
}
