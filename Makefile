.DEFAULT_GOAL := test

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./coverage.text -o ./cover.html
	@rm ./coverage.text
.PHONY: cover

# Run tests
test:
	@go test -v -coverprofile=./coverage.text -covermode=atomic $(shell go list ./...)
.PHONY: test
