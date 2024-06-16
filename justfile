#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v

test:
    go test ./...

build:
    mkdir -p "dist/bin"
    go build -v -o dist/bin ./...

release:
    goreleaser release --clean

release-dry:
    goreleaser release --skip=publish --clean

snapshot:
    goreleaser release --snapshot --skip=publish --clean
