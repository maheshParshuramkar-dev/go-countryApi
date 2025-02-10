test:
	go test -v -cover ./...

server:
	go run main.go

parallelBenchMark:
	go  test -benchmem -run=^$ -bench ^BenchmarkParallelGetCountryData$ github.com/go-countryApi/api -count=5

serialBenchMark:
	go  test -benchmem -run=^$ -bench ^BenchmarkSerialGetCountryData$ github.com/go-countryApi/api

.PHONY: test server parallelBenchMark serialBenchMark