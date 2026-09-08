package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// These are end-to-end CLI tests.
// These tests build the compare binary, run it as a subprocess and checks the output/exit code.

var (
	projectRoot string
	compareBin  string
)

// Creates a temporary directory for the binary. Used by all tests.
func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getwd failed: %v\n", err)
		os.Exit(1)
	}

	projectRoot = filepath.Join(wd, "..", "..")
	tempDir, err := os.MkdirTemp("", "compare-test-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "mkdir temp failed: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	compareBin = filepath.Join(tempDir, "compare")
	build := exec.Command("go", "build", "-o", compareBin, "./cmd/compare")
	build.Dir = projectRoot
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, output)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// Runs binary and returns stdout.
func runCompare(args ...string) ([]byte, error) {
	cmd := exec.Command(compareBin, args...)
	cmd.Dir = projectRoot
	return cmd.Output()
}

// Runs binary and returns stdout and stderr.
func runCompareCombined(args ...string) ([]byte, error) {
	cmd := exec.Command(compareBin, args...)
	cmd.Dir = projectRoot
	return cmd.CombinedOutput()
}

// Expected output structure.
type cliReport struct {
	Files []struct {
		Path         string `json:"path"`
		TotalKeys    uint64 `json:"total_keys"`
		DistinctKeys uint64 `json:"distinct_keys"`
	} `json:"files"`
	Overlap struct {
		Total    uint64 `json:"total"`
		Distinct uint64 `json:"distinct"`
	} `json:"overlap"`
}

// Runs binary for two large files. Doesn't check exact results. Checks pipeline works as expected.
func TestCompareCLI(t *testing.T) {
	output, err := runCompare("-in", "testdata/A_f.csv", "-in", "testdata/B_f.csv")
	if err != nil {
		t.Fatalf("compare failed: %v", err)
	}

	var report cliReport
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, output)
	}

	if len(report.Files) != 2 {
		t.Fatalf("got %d files in report, want 2", len(report.Files))
	}

	if report.Files[0].TotalKeys == 0 || report.Files[1].TotalKeys == 0 {
		t.Fatalf("expected non-zero file counts, got %+v", report.Files)
	}

	if report.Overlap.Total == 0 || report.Overlap.Distinct == 0 {
		t.Fatalf("expected non-zero overlap, got %+v", report.Overlap)
	}
}

// Runs binary for two small files. Checks exact results. Expected values are hard coded.
func TestCompareCLISmallFixtureExactValues(t *testing.T) {
	output, err := runCompare("-in", "testdata/small_a.csv", "-in", "testdata/small_b.csv")
	if err != nil {
		t.Fatalf("compare failed: %v", err)
	}

	var report cliReport
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, output)
	}

	want := cliReport{}
	want.Files = []struct {
		Path         string `json:"path"`
		TotalKeys    uint64 `json:"total_keys"`
		DistinctKeys uint64 `json:"distinct_keys"`
	}{
		{Path: "testdata/small_a.csv", TotalKeys: 3, DistinctKeys: 2},
		{Path: "testdata/small_b.csv", TotalKeys: 3, DistinctKeys: 2},
	}
	want.Overlap.Total = 2
	want.Overlap.Distinct = 2

	if len(report.Files) != len(want.Files) {
		t.Fatalf("got %d files, want %d", len(report.Files), len(want.Files))
	}
	for i := range want.Files {
		if report.Files[i] != want.Files[i] {
			t.Errorf("Files[%d] = %+v, want %+v", i, report.Files[i], want.Files[i])
		}
	}
	if report.Overlap != want.Overlap {
		t.Errorf("Overlap = %+v, want %+v", report.Overlap, want.Overlap)
	}
}

// Runs binary for two small files. Checks exact results. Expected values are stored in a file.
func TestCompareCLIGoldenFile(t *testing.T) {
	output, err := runCompare("-in", "testdata/small_a.csv", "-in", "testdata/small_b.csv")
	if err != nil {
		t.Fatalf("compare failed: %v", err)
	}

	expected, err := os.ReadFile(filepath.Join(projectRoot, "testdata", "expected_small.json"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	var got, want cliReport
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, output)
	}
	if err := json.Unmarshal(expected, &want); err != nil {
		t.Fatalf("invalid expected JSON: %v", err)
	}

	if got.Files[0] != want.Files[0] || got.Files[1] != want.Files[1] {
		t.Errorf("report = %+v, want %+v", got, want)
	}
	if got.Overlap != want.Overlap {
		t.Errorf("Overlap = %+v, want %+v", got.Overlap, want.Overlap)
	}
}

// Runs binary with only one input file.
func TestCompareCLIRequiresTwoFiles(t *testing.T) {
	cmd := exec.Command(compareBin, "-in", "testdata/A_f.csv")
	cmd.Dir = projectRoot
	err := cmd.Run()
	if err == nil {
		t.Fatal("compare should fail when only one input file is provided")
	}

	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expected exit error, got %T: %v", err, err)
	}
}

// Runs binary with non-existent file.
func TestCompareCLIMissingFile(t *testing.T) {
	cmd := exec.Command(compareBin, "-in", "missing.csv", "-in", "testdata/A_f.csv")
	cmd.Dir = projectRoot
	err := cmd.Run()
	if err == nil {
		t.Fatal("compare should fail when an input file is missing")
	}

	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expected exit error, got %T: %v", err, err)
	}
}

// Runs binary with version flag and checks something was printed.
func TestCompareCLIVersion(t *testing.T) {
	output, err := runCompareCombined("-version")
	if err != nil {
		t.Fatalf("compare -version failed: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("expected version output")
	}
}
