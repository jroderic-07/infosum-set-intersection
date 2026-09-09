# Architecture and Design Decisions

## Overview
This repository follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) pattern. 

* **Entrypoint:** `cmd/compare/main.go` orchestrates the core workflow.
* **Pipeline:** Ingests CSV files into flat frequency tables, compares them, and outputs a JSON report.

## Component Interfaces
Core behavior is abstracted behind Go interfaces so individual components can be swapped without changing domain logic:

| Component | Interface | Implementation |
|-----------|-----------|----------------|
| Reader | `DataReader` | `CSVReader` |
| Validator | `KeyValidator` | `RegexValidator` |
| Storage | `DataStore` | `FrequencyStore` |
| Comparison | `Comparator` | `FrequencyComparator` |
| Reporting | `Reporter` | `JSONReporter` |

## Concurrency
Files are read concurrently using Goroutines. Each ingestion stream writes to its own isolated `FrequencyStore` instance to prevent shared state and lock contention.

## Scaling

The provided datasets are roughly **95,000** and **80,000** rows each. After validation, approximately **86,500** and **72,800** valid UDPRN keys are loaded per file. At this size, the full datasets fit comfortably in memory and the comparison completes in under a second.

Given these file sizes, a simple approach is sufficient:

* **`CSVReader`** reads each file row by row using Go's standard library
* Both files are loaded **concurrently** via goroutines, since each ingestion pipeline is independent
* Valid keys are stored in an in-memory **`FrequencyStore`** (a map of key → count), which supports duplicate detection and both overlap metrics

This is a deliberate tradeoff: optimise for clarity and correctness at the current scale, not for hypothetical multi-terabyte inputs.

The codebase is written behind interfaces (`DataReader`, `DataStore`, `Comparator`, `Reporter`, `KeyValidator`) so that if file size requirements changed significantly, alternative implementations could be introduced without rewriting the CLI workflow or comparison logic.

If datasets grew to hundreds of gigabytes or beyond, likely next steps would be:

* **Streaming ingestion** — process each CSV in a single pass without holding the entire file in memory; write directly into a backing store as rows arrive
* **External storage** — replace `FrequencyStore` with a disk-backed `DataStore` implementation (e.g. BoltDB, RocksDB, or an embedded database) so key counts are not limited by RAM
* **Sorted merge join** — sort keys on disk, then walk both sorted streams in parallel to compute overlap without loading both full datasets into memory at once
* **Approximate methods** — use structures such as HyperLogLog for distinct counts when exact overlap is not required and memory must stay bounded

Other extensions — alternative input formats, additional report outputs, or a new HTTP entry point (`cmd/api/main.go`) — follow the same pattern: new implementation, existing interface.

## Testing and CI/CD
* **Local Testing:** Run unit and end-to-end tests via the Makefile (`make test` and `make e2e`).
* **Planned CI/CD Pipeline:** Automate test execution on pull requests, cross-compile binaries, build Docker containers for API deployments, and publish release artifacts.