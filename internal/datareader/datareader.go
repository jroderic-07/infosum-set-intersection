package datareader

import "set-intersection/internal/datastore"

// Interface defines how data is read into a data store.
// Implemented by CSVReader. Reads data from a CSV file into FrequencyStore.
type DataReader interface {
	Read(store datastore.DataStore) error
}
