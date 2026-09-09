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

## Extensibility Options
The interface architecture makes it straightforward to add:
* Alternative storage backends (e.g., external databases)
* Different comparison logic or output formats (e.g., CSV, plain text)
* New entry points, such as an HTTP API (`cmd/api/main.go`)

## Testing and CI/CD
* **Local Testing:** Run unit and end-to-end tests via the Makefile (`make test` and `make e2e`).
* **Planned CI/CD Pipeline:** Automate test execution on pull requests, cross-compile binaries, build Docker containers for API deployments, and publish release artifacts.