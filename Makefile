.PHONY: fmt

fmt:
	gofmt -s -w -l ./
	goimports -w -local github.com/yuzuy/happi -l ./
