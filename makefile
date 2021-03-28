test:
	go test -v -race ./...

release:
	goreleaser release
