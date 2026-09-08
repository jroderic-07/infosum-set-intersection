package report

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJSONReporter_Output(t *testing.T) {
	input := Report{
		Files: []FileStats{
			{
				Path:         "a.csv",
				TotalKeys:    3,
				DistinctKeys: 2,
			},
			{
				Path:         "b.csv",
				TotalKeys:    2,
				DistinctKeys: 2,
			},
		},
		Overlap: OverlapStats{
			Total:    1,
			Distinct: 1,
		},
	}

	var buf bytes.Buffer
	reporter := NewJSONReporter()

	if err := reporter.Output(&buf, input); err != nil {
		t.Fatalf("Output returned error: %v", err)
	}

	var got Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if len(got.Files) != len(input.Files) {
		t.Fatalf("Files length = %d, want %d", len(got.Files), len(input.Files))
	}
	for i := range input.Files {
		if got.Files[i] != input.Files[i] {
			t.Errorf("Files[%d] = %+v, want %+v", i, got.Files[i], input.Files[i])
		}
	}
	if got.Overlap != input.Overlap {
		t.Errorf("Overlap = %+v, want %+v", got.Overlap, input.Overlap)
	}
}
