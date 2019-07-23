.PHONY: cover

cover:
	go test ./opensubtitle -race -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
