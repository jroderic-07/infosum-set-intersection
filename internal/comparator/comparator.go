package comparator

import "set-intersection/internal/datastore"

// Structure for the comparison results.
type Result struct {
	TotalOverlap    uint64
	DistinctOverlap uint64
}

// Interface defines how data stores are compared.
// Implemented by FrequencyComparator.
type Comparator interface {
	Compare(a, b datastore.DataStore) (Result, error)
}
