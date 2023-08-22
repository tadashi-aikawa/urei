package inspection

import (
	"cmp"
	"github.com/cheggaaa/pb/v3"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"slices"
	"strconv"
)

type Seed struct {
	Name string `csv:"name"`
	Url  string `csv:"url"`
}

type Result struct {
	No     int    `csv:"id"`
	Name   string `csv:"name"`
	Url    string `csv:"url"`
	Status string `csv:"status"`
	Origin Seed   `csv:"-"`
}

func inspect(seq int, record Seed) Result {
	client := resty.New()
	slog.Debug("request", "seq", seq)
	r, err := client.R().Get(record.Url)
	slog.Debug("response", "seq", seq)
	if err != nil {
		slog.Error(err.Error())
		return Result{No: seq, Name: record.Name, Url: record.Url, Status: "error", Origin: record}
	}

	return Result{No: seq, Name: record.Name, Url: record.Url, Status: strconv.Itoa(r.StatusCode()), Origin: record}
}

func InspectRecords(records []Seed, concurrency int) (results []Result) {
	semChan := make(chan struct{}, concurrency)
	resultChan := make(chan Result, len(records))
	progress := pb.StartNew(len(records))

	for i, record := range records {
		i := i
		record := record
		go func() {
			semChan <- struct{}{}

			progress.Increment()
			resultChan <- inspect(i+1, record)

			<-semChan
		}()
	}

	for _ = range records {
		result := <-resultChan
		results = append(results, result)
	}

	slices.SortFunc(
		results,
		func(a, b Result) int { return cmp.Compare(a.No, b.No) },
	)

	return results
}
