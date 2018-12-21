GO = $(shell echo ${GOPATH})
CFLAGS = -v -work -x
admin_sdk: deps
	@echo "Building..."
	@echo "GOPATH=$(GO)"
	@echo "CFLAGS=$(CFLAGS)"
	@go build -o admin_sdk $(CFLAGS) src/main.go
	@echo "Done!"

deps:
	@echo "Gethering dependencies..."
	@echo "GOPATH=$(GO)"
	@echo Admin SDK
	@go get firebase.google.com/go
	@echo "Done!"
	@echo "MQTT client"
	@go get github.com/goiiot/libmqtt
	@echo "Done!" 

clean:
	@echo "Cleaning..."
	@echo "GOPATH=$(GO)"
	@go clean src/main.go 
	@rm -f admin_sdk
	@echo "Done!"

clean_all:
	@echo "Sweepping up..."
	@echo "GOPATH=$(GO)"
	@rm -rf pkg
	$(shell for item in $$(ls $$GOPATH/src); do [ -f "$$GOPATH/src/$$item" ] || rm -rf "$$GOPATH/src/$$item"; done)
	@echo "Done!"

