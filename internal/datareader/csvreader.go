package datareader

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"set-intersection/internal/datastore"
)

type CSVReader struct {
	path string
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

		key := record[0]

		if err := store.Add(key); err != nil {
			return err
		}
	}

	return nil
}

func NewCSVReader(filePath string) *CSVReader {
	return &CSVReader{
		path: filePath,
	}
}
