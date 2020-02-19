.PHONY: test lint dep-update

lint:
	golangci-lint run --enable-all

test:
	go test

dep-update:
	go get -u ./...
	go test ./...
	go mod tidy
	go mod vendor
