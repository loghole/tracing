GOFLAGS = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GOTEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

.PHONY: test
test:
	go test -race -v -cover -coverprofile coverage.out $(GOTEST_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run -v
