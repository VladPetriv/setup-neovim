.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: build
build:
	go build -o setup-nvim ./cmd/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: install
install:
	bash ./scripts/install.sh

.PHONY: uninstall
uninstall:
	bash ./scripts/uninstall.sh

