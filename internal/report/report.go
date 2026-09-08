package report

import "io"

// Format for the report output structure.
type Report struct {
	Files   []FileStats  `json:"files"`
	Overlap OverlapStats `json:"overlap"`
}

type FileStats struct {
	Path         string `json:"path"`
	TotalKeys    uint64 `json:"total_keys"`
	DistinctKeys uint64 `json:"distinct_keys"`
}

type OverlapStats struct {
	Total    uint64 `json:"total"`
	Distinct uint64 `json:"distinct"`
}

// Interface defines how a report is written to an io writer.
// Implemented by JSONReporter.
type Reporter interface {
	Output(w io.Writer, report Report) error
}
