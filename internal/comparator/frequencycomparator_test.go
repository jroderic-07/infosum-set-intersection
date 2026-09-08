package comparator

import (
	"set-intersection/internal/datastore"
	"testing"
)

func storeWithKeys(keys ...string) *datastore.FrequencyStore {
	store := datastore.NewFrequencyStore()
	for _, key := range keys {
		_ = store.Add(key)
	}
	return store
}

func TestFrequencyComparator_Compare(t *testing.T) {
	tests := []struct {
		name         string
		aKeys        []string
		bKeys        []string
		wantTotal    uint64
		wantDistinct uint64
	}{
		{
			name:         "no overlap",
			aKeys:        []string{"a", "b"},
			bKeys:        []string{"c", "d"},
			wantTotal:    0,
			wantDistinct: 0,
		},
		{
			name:         "distinct overlap only",
			aKeys:        []string{"a", "b"},
			bKeys:        []string{"b", "c"},
			wantTotal:    1,
			wantDistinct: 1,
		},
		{
			name:         "duplicate keys",
			aKeys:        []string{"x", "x", "y"},
			bKeys:        []string{"x", "y", "y"},
			wantTotal:    2,
			wantDistinct: 2,
		},
		{
			name:         "identical sets",
			aKeys:        []string{"a", "a", "b"},
			bKeys:        []string{"a", "a", "b"},
			wantTotal:    3,
			wantDistinct: 2,
		},
		{
			name:         "both empty",
			aKeys:        nil,
			bKeys:        nil,
			wantTotal:    0,
			wantDistinct: 0,
		},
		{
			name:         "one store empty",
			aKeys:        []string{"a", "b"},
			bKeys:        nil,
			wantTotal:    0,
			wantDistinct: 0,
		},
		{
			name:         "extra keys in other store",
			aKeys:        []string{"a", "b"},
			bKeys:        []string{"b", "c"},
			wantTotal:    1,
			wantDistinct: 1,
		},
	}

	comparator := NewFrequencyComparator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := comparator.Compare(storeWithKeys(tt.aKeys...), storeWithKeys(tt.bKeys...))
			if err != nil {
				t.Fatalf("Compare returned error: %v", err)
			}

			if result.TotalOverlap != tt.wantTotal {
				t.Errorf("TotalOverlap = %d, want %d", result.TotalOverlap, tt.wantTotal)
			}
			if result.DistinctOverlap != tt.wantDistinct {
				t.Errorf("DistinctOverlap = %d, want %d", result.DistinctOverlap, tt.wantDistinct)
			}
		})
	}
}
