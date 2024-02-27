test:
	go test ./...

test_coverage:
	go test ./... -coverprofile .testCoverage.txt

doc:
	godoc -http=:6060

release:
	$(shell goreleaser release --clean)

release-dry:
	$(shell goreleaser release --skip=publish)

snapshot:
	$(shell goreleaser release --snapshot --skip=publish --clean)
