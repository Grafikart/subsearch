.PHONY: install
install:
	go get ./...

.PHONY: test
test:
	go test ./... -race -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
