package dataset

import (
	"set-intersection/internal/datareader"
	"set-intersection/internal/datastore"
	"testing"
)

func TestDataset(t *testing.T) {
	store := datastore.NewFrequencyStore()
	reader := datareader.NewCSVReader("example.csv", nil)

	d := Dataset{
		Store:  store,
		Reader: reader,
	}

	if d.Store == nil {
		t.Fatal("Store should not be nil")
	}
	if d.Reader == nil {
		t.Fatal("Reader should not be nil")
	}
}
