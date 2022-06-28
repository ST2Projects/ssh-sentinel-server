test:
	go test ./...

test_coverage:
	go test ./... -coverprofile .testCoverage.txt

doc:
	godoc -http=:6060
