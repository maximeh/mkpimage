gopath = "$(CURDIR)/third_party"
cover = $(COVER)

.PNONY: all test deps fmt clean

all: clean fmt deps test
	@echo "==> Compiling source code."
	@env GOPATH=$(gopath) go build -v -o ./bin/mkpimage ./mkpimage

race: clean fmt deps test
	@echo "==> Compiling source code with race detection enabled."
	@env GOPATH=$(gopath) go build -race -o ./bin/mkpimage ./mkpimage

test:
	@echo "==> Running tests."
	@env GOPATH=$(gopath) go test $(cover) ./mkpimage

deps:
	@echo "==> Downloading dependencies."
	@env GOPATH=$(gopath) go get -d -v ./mkpimage/...
	@echo "==> Removing SCM files from third_party."
	@find ./third_party -type d -name .git | xargs rm -rf
	@find ./third_party -type d -name .bzr | xargs rm -rf
	@find ./third_party -type d -name .hg | xargs rm -rf

fmt:
	@echo "==> Formatting source code."
	@gofmt -w ./mkpimage

clean:
	@echo "==> Cleaning up previous builds."
	@rm -rf "$(GOPATH)/pkg" ./third_party/pkg ./bin
