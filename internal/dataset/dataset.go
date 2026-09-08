package dataset

import (
	"set-intersection/internal/datareader"
	"set-intersection/internal/datastore"
)

type Dataset struct {
	Store  datastore.DataStore
	Reader datareader.DataReader
}
