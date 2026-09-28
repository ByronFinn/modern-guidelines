# Maintenance helpers for the modern-guidelines hub.
#
# `make test`              runs the full test suite.
# `make generate-features` rebuilds docs/features/<lang>.md from the
#                          guideline datasets via internal/featuresgen
#                          (same as `go generate ./internal/guidelines`).
# `make ingest`            re-ingests the JetBrains Go dataset from the
#                          upstream commit locked in internal/ingest,
#                          overwriting internal/guidelines/data/go.json.
# `make check`             test + gofmt + vet, the pre-push gate.

.PHONY: test generate-features ingest check

test:
	@go test ./...

generate-features:
	@go run ./internal/featuresgen internal/guidelines/data docs/features

ingest:
	@go run ./scripts/ingest-jetbrains-go

check: test
	@unformatted=$$(gofmt -l main.go internal scripts 2>/dev/null); \
	if [ -n "$$unformatted" ]; then \
		echo "files need gofmt -w:"; echo "$$unformatted"; exit 1; \
	fi
	@go vet ./...
