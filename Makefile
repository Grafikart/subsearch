.PHONY: install
install:
	go get ./...

.PHONY: test
test:
	GO111MODULE=on go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html
