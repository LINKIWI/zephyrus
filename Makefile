# Name of the server executable
SERVER = server
# Name of the collector executable
COLLECTOR = collector

# Output binary directory
BIN_DIR = bin
# Protobuf schemas directory
PROTO_DIR = schemas

all: $(SERVER) $(COLLECTOR)

schemas: $(PROTO_DIR)/%.pb.go

$(SERVER): schemas
	go build -o $(BIN_DIR)/server cmd/$(SERVER)/main.go

$(COLLECTOR): schemas
	go build -o $(BIN_DIR)/collector cmd/$(COLLECTOR)/main.go

$(PROTO_DIR)/%.pb.go: $(wildcard $(PROTO_DIR)/*.proto)
	protoc -I $(PROTO_DIR) $(PROTO_DIR)/*.proto --go_out=plugins=grpc:$(PROTO_DIR)
