test:
	go test -v $(shell go list ./... | grep -v test_utils)

test_coverage:
	go test -v $(shell go list ./... | grep -v test_utils) -coverprofile .testCoverage.txt

doc:
	godoc -http=:6060
