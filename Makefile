.PHONY: test

test:
	go test ./opensubtitle -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
