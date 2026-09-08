package datastore

// Interface defines how keys and their counts are stored and queried.
// Implemented by FrequencyStore to store UDPRN values and their counts. Used by reader, comparator, and reporter.
type DataStore interface {
	Add(key string) error
	Get(key string) (uint64, bool)
	GetDistinctCount() uint64
	GetTotalCount() uint64
	Keys() []string
}
