REGION ?= eu-west-2
BUCKET ?= pubclub-artifacts
SERVICE ?= confirmation
VERSION ?= 1.0.0

.PHONY: bucket
bucket:
	@if ! aws s3api head-bucket --bucket $(BUCKET) 2>/dev/null; then\
		aws s3api create-bucket --bucket $(BUCKET) --region $(REGION)\
		--create-bucket-configuration LocationConstraint=$(REGION) 1>/dev/null;\
	fi;

.PHONY: build
build:
	@cd services/$(SERVICE)/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go; \
	zip ../../builds/$(SERVICE).zip main

.PHONY: deploy
deploy: bucket build
	@aws s3 cp builds/$(SERVICE).zip s3://$(BUCKET)/$(SERVICE)/v$(VERSION)/main.zip
