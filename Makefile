.PHONY: test cover cover-html check

COVER_FILE := coverage.out

# RUN ALL TESTS
test:
	go test ./... -count=1

# RUN TESTS AND PRINT ONE COVERAGE TABLE FOR ALL OF internal/
cover:
	@go test ./test/... -count=1 -coverpkg=./internal/... -coverprofile=$(COVER_FILE) > /dev/null
	@scripts/coverage-table.sh $(COVER_FILE)

# OPEN HTML COVERAGE REPORT IN BROWSER
cover-html: cover
	go tool cover -html=$(COVER_FILE)

# SAME CHECKS AS CI WILL RUN
check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)
	go vet ./...
	go test ./... -count=1
