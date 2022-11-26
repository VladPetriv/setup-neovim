.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run --enable-all ./...

