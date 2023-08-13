package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jszwec/csvutil"
	"log"
	"os"
)

type Record struct {
	Name string `csv:"name"`
	Url  string `csv:"url"`
}

func loadCsv[T any](path string) ([]T, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var records []T
	if err := csvutil.Unmarshal(bytes, &records); err != nil {
		return nil, err
	}

	return records, nil
}

func saveCsv[T any](path string, records []T) error {
	bytes, err := csvutil.Marshal(records)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, bytes, 644); err != nil {
		return err
	}

	return nil
}

type InspectionResult struct {
	No     int    `csv:"id"`
	Status int    `csv:"status"`
	Origin Record `csv:"-"`
}

func inspect(records []Record) (results []InspectionResult) {
	client := resty.New()

	for i, record := range records {
		r, err := client.R().Get(record.Url)
		if err != nil {
			log.Print(err)
			continue
		}

		results = append(results, InspectionResult{No: i + 1, Status: r.StatusCode(), Origin: record})
	}

	return results
}

type Args struct {
	source string
	dst    string
}

func parseArgs() Args {
	source := flag.String("source", "./input.csv", "入力CSVのパス")
	dst := flag.String("dst", "./output.csv", "出力CSVのパス")
	flag.Parse()
	return Args{source: *source, dst: *dst}
}

func main() {
	args := parseArgs()

	records, err := loadCsv[Record](args.source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", records)

	results := inspect(records)
	fmt.Printf("%+v\n", results)

	if err := saveCsv(args.dst, results); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success!")
}
