package datareader

import (
	"os"
	"path/filepath"
	"set-intersection/internal/datastore"
	"set-intersection/internal/validator"
	"testing"
)

func TestCSVReader_Read(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.csv")

	content := "key-a\nkey-a\nkey-b\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(path, nil)

	if err := reader.Read(store); err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 3 {
		t.Errorf("GetTotalCount() = %d, want 3", got)
	}

	if got := store.GetDistinctCount(); got != 2 {
		t.Errorf("GetDistinctCount() = %d, want 2", got)
	}

	count, ok := store.Get("key-a")
	if !ok || count != 2 {
		t.Errorf("Get(key-a) = (%d, %v), want (2, true)", count, ok)
	}
}

func TestCSVReader_ReadMissingFile(t *testing.T) {
	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(filepath.Join(t.TempDir(), "missing.csv"), nil)

	if err := reader.Read(store); err == nil {
		t.Fatal("Read should return an error for a missing file")
	}
}

func TestCSVReader_ReadEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.csv")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(path, nil)

	if err := reader.Read(store); err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 0 {
		t.Errorf("GetTotalCount() = %d, want 0", got)
	}
}

func TestCSVReader_ReadSingleRow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "single.csv")
	if err := os.WriteFile(path, []byte("only-key\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(path, nil)

	if err := reader.Read(store); err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 1 {
		t.Errorf("GetTotalCount() = %d, want 1", got)
	}

	count, ok := store.Get("only-key")
	if !ok || count != 1 {
		t.Errorf("Get(only-key) = (%d, %v), want (1, true)", count, ok)
	}
}

func TestCSVReader_ReadIgnoresBlankLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "blank-line.csv")
	if err := os.WriteFile(path, []byte("key-a\n\nkey-b\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(path, nil)

	if err := reader.Read(store); err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 2 {
		t.Errorf("GetTotalCount() = %d, want 2", got)
	}
}

func TestCSVReader_SkipsInvalidKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid-keys.csv")
	content := "08034283\n\nudprn\n71842328\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	udprnValidator, err := validator.NewUDPRNValidator()
	if err != nil {
		t.Fatalf("NewUDPRNValidator returned error: %v", err)
	}

	store := datastore.NewFrequencyStore()
	reader := NewCSVReader(path, udprnValidator)

	if err := reader.Read(store); err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if got := store.GetTotalCount(); got != 2 {
		t.Errorf("GetTotalCount() = %d, want 2", got)
	}
}
