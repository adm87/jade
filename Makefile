.PHONY: lint align fmt test test-integration check check-integration

test:
	go test ./...

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
