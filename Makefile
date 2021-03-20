test:
	@go test -count=1 ./...

lint:
	@golint pkg/...

format:
	@find pkg/ -type f -name '*.go' -exec go fmt {} \;

fmt: format

pre-commit:
	@pre-commit run --all-files

example: cmd/example.go pkg/cmdline/cmdline.go
	@go build cmd/example.go

clean:
	@rm -f example

