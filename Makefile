GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build: deps
	$(GOBUILD) -o server $(GOPATH)/src/github.com/cyyber/FindMaxNumberRPC/cmd/server/main.go
	$(GOBUILD) -o client $(GOPATH)/src/github.com/cyyber/FindMaxNumberRPC/cmd/client/main.go
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm $(GOPATH)/src/github.com/cyyber/FindMaxNumberRPC/server
	rm $(GOPATH)/src/github.com/cyyber/FindMaxNumberRPC/client
deps:
	@if [ -z "$(GOPATH)" ]; then \
	  echo "GOPATH Not Set" && \
	  exit 1; \
	fi
	$(GOGET) ./...
