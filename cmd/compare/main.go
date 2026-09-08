package main

import (
	"log"
	"os"
	"runtime"
	"set-intersection/internal/comparator"
	"set-intersection/internal/config"
	"set-intersection/internal/datareader"
	"set-intersection/internal/dataset"
	"set-intersection/internal/datastore"
	"set-intersection/internal/report"
	"sync"
)

var version = "1.0"

func main() {
	config := config.LoadConfigurationCLI()

	if config.ShowVersion {
		log.Println(version)
		os.Exit(0)
	}

	if config.CpuNum != 0 {
		runtime.GOMAXPROCS(config.CpuNum)
	}

	if len(config.InputFilePaths) != 2 {
		log.Fatal("exactly two input files are required")
	}

	datasets := make([]dataset.Dataset, len(config.InputFilePaths))

	for i, path := range config.InputFilePaths {
		// Store and Reader are both interfaces. The constructor functions return pointers. The interfaces hold the pointers, so method calls affect the original objects.
		datasets[i] = dataset.Dataset{
			Store:  datastore.NewFrequencyStore(),
			Reader: datareader.NewCSVReader(path),
		}
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(datasets))

	for i := range datasets {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			// The Store interface holds a pointer. The read function will affect the original object.
			if err := datasets[i].Reader.Read(datasets[i].Store); err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Fatal(err)
		}
	}

	result, err := comparator.NewFrequencyComparator().Compare(datasets[0].Store, datasets[1].Store)
	if err != nil {
		log.Fatal(err)
	}

	comparisonReport := report.Report{
		Files: []report.FileStats{
			{
				Path:         config.InputFilePaths[0],
				TotalKeys:    datasets[0].Store.GetTotalCount(),
				DistinctKeys: datasets[0].Store.GetDistinctCount(),
			},
			{
				Path:         config.InputFilePaths[1],
				TotalKeys:    datasets[1].Store.GetTotalCount(),
				DistinctKeys: datasets[1].Store.GetDistinctCount(),
			},
		},
		Overlap: report.OverlapStats{
			Total:    result.TotalOverlap,
			Distinct: result.DistinctOverlap,
		},
	}

	if err := report.NewJSONReporter().Output(os.Stdout, comparisonReport); err != nil {
		log.Fatal(err)
	}
}
