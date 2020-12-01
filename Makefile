TEST?=$$(go list ./... | grep -v vendor )
SWEEP?=all
SWEEP_DIR?=./rollbar
HOSTNAME=github.com
NAMESPACE=rollbar
NAME=rollbar
BINARY=terraform-provider-${NAME}
VERSION=0.2.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

install012: build
	mkdir -p ~/.terraform.d/plugins/${OS_ARCH}/
	mv ${BINARY} ~/.terraform.d/plugins/${OS_ARCH}/terraform-provider-${NAME}_v${VERSION}

test: 
	go test -covermode=atomic -coverprofile=coverage.out $(TEST) || exit 1
	@#echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test -covermode=atomic -coverprofile=coverage.out $(TEST) -v $(TESTARGS) -timeout 120m   

slscan:
	./.slscan.sh

SHELL=bash
sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	@read -p "Are you sure? " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; \
	then \
		go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m; \
	fi


#-------------------------------------------------------------------------------
#
# Terraform - easily run Terraform commands with a the latest provider code.
#
#-------------------------------------------------------------------------------
apply: install _terraform_cleanup _terraform_init _terraform_apply _terraform_log
plan: install _terraform_cleanup _terraform_init _terraform_plan _terraform_log
destroy: install _terraform_cleanup _terraform_init _terraform_destroy _terraform_log

_terraform_cleanup:
	# Cleanup last run
	rm -vrf example/.terraform /tmp/terraform-provider-rollbar.log
_terraform_init:
	# Initialize terraform
	(cd example && terraform init)
_terraform_log:
	# Print the debug log
	cat /tmp/terraform-provider-rollbar.log
_terraform_apply:
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform apply) || true
_terraform_apply_nodebug:
	(cd example && terraform apply) || true
_terraform_apply_auto:
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform apply --auto-approve) || true
_terraform_plan:
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform plan)
_terraform_destroy:
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform destroy)

docker12:
	docker build . --build-arg version=0.12.5 -t terraform-0.12-provider-rollbar 
	docker run terraform-0.12-provider-rollbar plan -var rollbar_token=$$ROLLBAR_API_KEY

docker13:
	docker build . --build-arg version=0.13.5 -t terraform-0.13-provider-rollbar 
	docker run terraform-0.13-provider-rollbar plan -var rollbar_token=$$ROLLBAR_API_KEY
