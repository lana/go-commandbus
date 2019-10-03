.DEFAULT_GOAL := test
# Run tests
test:
	@go test -v -coverprofile=./coverage.text -covermode=atomic $(shell go list ./...)
.PHONY: test
