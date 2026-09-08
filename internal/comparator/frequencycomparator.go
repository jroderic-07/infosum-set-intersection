package comparator

import "set-intersection/internal/datastore"

type FrequencyComparator struct{}

func (c *FrequencyComparator) Compare(a datastore.DataStore, b datastore.DataStore) (Result, error) {
	var totalOverlap uint64
	var distinctOverlap uint64

	keysStore := a
	otherStore := b
	if b.GetDistinctCount() < a.GetDistinctCount() {
		keysStore = b
		otherStore = a
	}

	for _, key := range keysStore.Keys() {
		countKeys, ok := keysStore.Get(key)
		if !ok {
			continue
		}

		countOther, ok := otherStore.Get(key)
		if !ok {
			continue
		}

		distinctOverlap++
		totalOverlap += min(countKeys, countOther)
	}

	return Result{
		TotalOverlap:    totalOverlap,
		DistinctOverlap: distinctOverlap,
	}, nil
}

func NewFrequencyComparator() *FrequencyComparator {
	return &FrequencyComparator{}
}
