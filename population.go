package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func newTabReader(reader io.Reader) (tabReader *csv.Reader) {
	tabReader = csv.NewReader(reader)
	tabReader.Comma = '\t'
	tabReader.LazyQuotes = true
	return
}

// ReadRows reads tab-delimited rows into memory.
func ReadRows(reader io.Reader) ([][]string, error) {
	return newTabReader(reader).ReadAll()
}

// RowErr is a row and/or an error. Really just or, but whatever.
type RowErr struct {
	row []string
	err error
}

func readRowsRoutine(reader io.Reader, rows chan RowErr) {
	defer close(rows)
	tabReader := newTabReader(reader)
	for {
		row, err := tabReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			rows <- RowErr{nil, err}
		} else {
			rows <- RowErr{row, nil}
		}
	}
}

// ReadRowsChan reads tab-delimited through a channel.
func ReadRowsChan(reader io.Reader) (rows chan RowErr) {
	rows = make(chan RowErr)
	go readRowsRoutine(reader, rows)
	return
}

func main() {
	// Open and manage file.
	reader, err := os.Open("data/cities500.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	// Init counts.
	populationSouth := int64(0)
	populationTotal := int64(0)
	rowCount := 0
	// Loop on rows, needing different forms depending on mode.
	// rows, err := ReadRows(reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, row := range rows {
	rows := ReadRowsChan(reader)
	for rowErr := range rows {
		if rowErr.err != nil {
			log.Fatal(rowErr.err)
		}
		row := rowErr.row
		// Parse fields, showing Go's love for trees before the forest.
		latitude, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		population, err := strconv.ParseInt(row[14], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		// Actual logic.
		populationTotal += population
		if latitude < 0 {
			populationSouth += population
		}
		rowCount++
	}
	// Print results.
	fmt.Printf("South population: %d\n", populationSouth)
	fmt.Printf("Total population: %d\n", populationTotal)
	fmt.Printf("# of rows: %d\n", rowCount)
}
