.PHONY: install
install:
	go get ./...

.PHONY: cover
cover:
	go test ./... -race -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
