package inspection

import (
	"cmp"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"slices"
)

type Seed struct {
	Name string `csv:"name"`
	Url  string `csv:"url"`
}

type Result struct {
	No     int  `csv:"id"`
	Status int  `csv:"status"`
	Origin Seed `csv:"-"`
}

type AsyncResult struct {
	Value *Result
	Err   error
}

func inspect(seq int, record Seed) (*Result, error) {
	client := resty.New()
	slog.Info("request", "seq", seq)
	r, err := client.R().Get(record.Url)
	slog.Info("response", "seq", seq)
	if err != nil {
		return nil, err
	}

	return &Result{No: seq, Status: r.StatusCode(), Origin: record}, nil
}

func InspectRecords(records []Seed, concurrency int) (results []Result) {
	semChan := make(chan struct{}, concurrency)
	asyncResultsChan := make(chan AsyncResult, len(records))

	for i, record := range records {
		i := i
		record := record
		go func() {
			semChan <- struct{}{}
			r, err := inspect(i+1, record)
			asyncResultsChan <- AsyncResult{Value: r, Err: err}
			<-semChan
		}()
	}

	for _ = range records {
		result := <-asyncResultsChan
		if result.Err != nil {
			slog.Error(result.Err.Error())
		} else {
			results = append(results, *result.Value)
		}
	}

	slices.SortFunc(
		results,
		func(a, b Result) int { return cmp.Compare(a.No, b.No) },
	)

	return results
}
