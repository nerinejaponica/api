.PHONY: test cover cover-html check mocks install-tools

COVER_FILE := coverage.out
# mockgen LIVES IN GOPATH/bin, WHICH IS NOT ALWAYS ON PATH
MOCKGEN := $(shell go env GOPATH)/bin/mockgen

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

# INSTALL DEV TOOLS, PINNED TO A COMMIT HASH (v0.6.0)
install-tools:
	go install go.uber.org/mock/mockgen@2d1c58167e30f380cf78e44a43b100a14767e817

# GENERATE MOCKS FOR TESTS, COMMIT THE RESULT SO CI DOES NOT NEED mockgen
mocks:
	$(MOCKGEN) -source=internal/db/querier.go -destination=test/mock/querier_mock.go -package=mock
