.PHONY: fmt test-run

fmt:
	gofmt -s -w -l ./
	goimports -w -local github.com/yuzuy/happi -l ./

test-run:
	@-touch stderr
	@go run main.go test.txt 2>> stderr
