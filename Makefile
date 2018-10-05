# Name of the server executable
SERVER = server
# Name of the collector executable
COLLECTOR = collector

# Output binary directory
BIN_DIR = bin
# Protobuf schemas directory
PROTO_DIR = schemas

GOOS ?= $(shell go tool dist env | grep GOOS | sed 's/"//g' | sed 's/.*=//g')
GOARCH ?= $(shell go tool dist env | grep GOARCH | sed 's/"//g' | sed 's/.*=//g')

all: $(SERVER) $(COLLECTOR)

schemas: $(PROTO_DIR)/%.pb.go

$(SERVER): schemas
	go build -o $(BIN_DIR)/zephyrus-server-$(GOOS)-$(GOARCH) cmd/$(SERVER)/main.go

$(COLLECTOR): schemas
	go build -o $(BIN_DIR)/zephyrus-collector-$(GOOS)-$(GOARCH) cmd/$(COLLECTOR)/main.go

$(PROTO_DIR)/%.pb.go: $(wildcard $(PROTO_DIR)/*.proto)
	protoc -I $(PROTO_DIR) $(PROTO_DIR)/*.proto --go_out=plugins=grpc:$(PROTO_DIR)

lint:
	.ci/lint.sh

clean:
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f $(BIN_DIR)/*

.PHONY: lint clean
