package main

import (
	"flag"
	"github.com/tadashi-aikawa/urei/app/inspection"
	"github.com/tadashi-aikawa/urei/pkg/file"
	"log"
)

type Args struct {
	source      string
	dst         string
	concurrency int
}

func parseArgs() Args {
	source := flag.String("source", "./input.csv", "入力CSVのパス")
	dst := flag.String("dst", "./output.csv", "出力CSVのパス")
	concurrency := flag.Int("concurrency", 1, "同時実行数")
	flag.Parse()
	return Args{source: *source, dst: *dst, concurrency: *concurrency}
}

func main() {
	args := parseArgs()

	records, err := file.LoadCsv[inspection.Seed](args.source)
	if err != nil {
		log.Fatal(err)
	}

	results := inspection.InspectRecords(records, args.concurrency)
	if err := file.SaveCsv(args.dst, results); err != nil {
		log.Fatal(err)
	}
}
