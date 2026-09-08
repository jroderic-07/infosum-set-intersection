test:
	go test -v ./...
	go vet ./...

compare:
	go build -o ./bin/compare cmd/compare/main.go

e2e: compare
	./bin/compare -in testdata/A_f.csv -in testdata/B_f.csv

.PHONY: test compare e2e