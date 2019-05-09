TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

sanitycheck:
	$(MAKE) fmtcheck

build-darwin:
	GOOS=darwin GOARCH=amd64 go install

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install

build: build-linux

test: sanitycheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=300s -parallel=4

testacc: sanitycheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -parallel=8

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build build-darwin sanitycheck test testacc fmt fmtcheck errcheck test-compile
