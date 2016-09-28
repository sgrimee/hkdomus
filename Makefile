GO ?= go


default: mac

mac:
	GOOS=darwin GOARCH=amd64 $(GO) build

rpi:
	GOOS=linux GOARCH=arm $(GO) build -o hkdomus.rpi

