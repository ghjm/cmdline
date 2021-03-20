test:
	@go test -count=1 ./...

lint:
	@golint pkg/...

format:
	@find pkg/ -type f -name '*.go' -exec go fmt {} \;

fmt: format

pre-commit:
	@pre-commit run --all-files

