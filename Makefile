SOURCES := $(shell find . -name '*.go')
BINARY := kubectl-get_yaml

build: kubectl-get_yaml

test: $(SOURCES)
	go test -v -short -race -timeout 30s ./...

clean:
	@rm -rf $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags="-s -w" ./cmd/$(BINARY).go
