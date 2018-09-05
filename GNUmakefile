TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
DEP := $(shell command -v dep 2> /dev/null)

default: build

sanitycheck:
	$(MAKE) depensure
	$(MAKE) fmtcheck


build-darwin: sanitycheck
	GOOS=darwin GOARCH=amd64 go install

build: sanitycheck
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=300s -parallel=4

testacc: sanitycheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -parallel=8

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

depensure:
ifndef DEP
  $(error "No dep in $(PATH), install: https://github.com/golang/dep#setup")
endif
	@sh -c "dep ensure"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@dep status

test-compile: depensure
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build build-darwin sanitycheck depensure test testacc vet fmt fmtcheck errcheck vendor-status test-compile
