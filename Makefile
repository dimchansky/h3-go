.PHONY: test bench lint gen

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	go vet ./...
	@echo "Note: golangci-lint not configured yet (external dependency)"

gen:
	@echo "Reserved for future table generation"