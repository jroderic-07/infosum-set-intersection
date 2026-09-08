package datastore

import (
	"slices"
	"testing"
)

func TestFrequencyStore_AddAndGet(t *testing.T) {
	store := NewFrequencyStore()

	if err := store.Add("key-a"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if err := store.Add("key-a"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if err := store.Add("key-b"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	count, ok := store.Get("key-a")
	if !ok || count != 2 {
		t.Errorf("Get(key-a) = (%d, %v), want (2, true)", count, ok)
	}

	count, ok = store.Get("key-b")
	if !ok || count != 1 {
		t.Errorf("Get(key-b) = (%d, %v), want (1, true)", count, ok)
	}

	_, ok = store.Get("missing")
	if ok {
		t.Error("Get(missing) should not exist")
	}
}

func TestFrequencyStore_CountsAndKeys(t *testing.T) {
	store := NewFrequencyStore()

	for range 3 {
		if err := store.Add("x"); err != nil {
			t.Fatalf("Add returned error: %v", err)
		}
	}
	if err := store.Add("y"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 4 {
		t.Errorf("GetTotalCount() = %d, want 4", got)
	}

	if got := store.GetDistinctCount(); got != 2 {
		t.Errorf("GetDistinctCount() = %d, want 2", got)
	}

	keys := store.Keys()
	slices.Sort(keys)

	want := []string{"x", "y"}
	if !slices.Equal(keys, want) {
		t.Errorf("Keys() = %v, want %v", keys, want)
	}
}
