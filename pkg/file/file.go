package file

import (
	"github.com/jszwec/csvutil"
	"os"
)

func LoadCsv[T any](path string) ([]T, error) {
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

func SaveCsv[T any](path string, records []T) error {
	bytes, err := csvutil.Marshal(records)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, bytes, 644); err != nil {
		return err
	}

	return nil
}
