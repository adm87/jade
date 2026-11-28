.PHONY: lint align fmt test bench test-integration check check-integration

test:
	go test ./...

## Benchmarks only, no tests
bench:
	go test -run '^$$' -bench . -benchmem ./...

## Opens a window; needs a display
test-integration:
	go test -tags integration -run '^TestIntegration' ./...

## Dagger (runs on a virtual display)
check:
	dagger check

check-integration:
	dagger check --check=integration

## Check only (CI-safe)
lint:
	betteralign ./...
	golangci-lint run ./...

## Apply fixes
align:
	betteralign -apply ./...

fmt:
	golangci-lint fmt ./...
