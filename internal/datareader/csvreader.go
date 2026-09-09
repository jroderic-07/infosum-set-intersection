package datareader

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"set-intersection/internal/datastore"
	"set-intersection/internal/validator"
	"strings"
)

type CSVReader struct {
	path      string
	validator validator.KeyValidator
}

func (r *CSVReader) Read(store datastore.DataStore) error {
	file, err := os.Open(r.path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if len(record) == 0 {
			continue
		}

		key := strings.TrimSpace(record[0])

		if r.validator != nil {
			if err := r.validator.Validate(key); err != nil {
				continue
			}
		}

		if err := store.Add(key); err != nil {
			return err
		}
	}

	return nil
}

func NewCSVReader(filePath string, keyValidator validator.KeyValidator) *CSVReader {
	return &CSVReader{
		path:      filePath,
		validator: keyValidator,
	}
}
