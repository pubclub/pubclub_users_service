SERVICE ?= confirmation

.PHONY: build
build:
	@cd services/$(SERVICE)/; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go; \
	zip ../../builds/$(SERVICE).zip main
