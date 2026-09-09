# Set Intersection

CLI tool that compares two CSV files of UDPRN keys and prints overlap stats as JSON.

Written in Go 1.22. See [docs/architecture.md](docs/architecture.md) for design notes.

## Requirements

- Go 1.22+

## Quickstart

```bash
# build
make compare

# run
./bin/compare -in testdata/A_f.csv -in testdata/B_f.csv

# test
make test
```

Build output goes to `./bin/compare`. JSON goes to stdout; errors go to stderr.

## Flags

| Flag | Description |
|------|-------------|
| `-in` | Input CSV path. Pass twice, once per file. |
| `-version` | Print version and exit. |
| `-cpuNum` | Set `GOMAXPROCS` (default: runtime default). |

## Example

```bash
./bin/compare -in testdata/A_f.csv -in testdata/B_f.csv
```

```json
{
  "files": [
    {
      "path": "testdata/A_f.csv",
      "total_keys": 86535,
      "distinct_keys": 72798
    },
    {
      "path": "testdata/B_f.csv",
      "total_keys": 72846,
      "distinct_keys": 72814
    }
  ],
  "overlap": {
    "total": 58244,
    "distinct": 58221
  }
}
```

## Output fields

- `total_keys` — row count per file, including duplicates
- `distinct_keys` — unique keys per file
- `overlap.total` — shared key occurrences: sum of `min(countA, countB)` for each key in both files
- `overlap.distinct` — keys that appear in both files

## Data handling

Exactly two input files are required.

Each CSV is expected to contain a single column of UDPRN keys, one per row. Keys are stored as strings so leading zeros are preserved (e.g. `08034283`).

Only 8-digit numeric values are loaded. Rows that fail validation are skipped silently — this includes blank lines, empty values (`""`), the `udprn` header row, and malformed keys. The provided test files contain thousands of empty rows; these are excluded from all counts.

Only the first column is read. Extra columns, if present, are ignored. Overlap counts are symmetric — the order of `-in` flags affects report paths only, not the overlap figures.
