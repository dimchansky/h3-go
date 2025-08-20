.PHONY: test bench lint gen ref

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	go vet ./...
	@echo "Note: golangci-lint not configured yet (external dependency)"

gen:
	@echo "Reserved for future table generation"

ref:
	@echo "Building H3 reference oracle for test validation..."
	$(MAKE) -C testref
	@echo "Oracle built successfully at testref/h3ref"
	@echo "Run 'make -C testref test' to verify installation"