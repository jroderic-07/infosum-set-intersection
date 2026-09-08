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
      "total_keys": 95001,
      "distinct_keys": 72800
    },
    {
      "path": "testdata/B_f.csv",
      "total_keys": 80001,
      "distinct_keys": 72816
    }
  ],
  "overlap": {
    "total": 65399,
    "distinct": 58223
  }
}
```

## Output fields

- `total_keys` — row count per file, including duplicates
- `distinct_keys` — unique keys per file
- `overlap.total` — shared key occurrences: sum of `min(countA, countB)` for each key in both files
- `overlap.distinct` — keys that appear in both files

## Assumptions

- Exactly two input files are required.
- Each CSV has one column, one key per row.
- The header row (`udprn` in the sample files) is counted as a key.
- Blank lines are skipped.
