# директория с proto (относительно Makefile)
REPO_ROOT := $(shell git rev-parse --show-toplevel)
PROTO_DIR := $(REPO_ROOT)/protos


OUT_DIR := .

.PHONY: all clean

all: generate

generate:

	protoc \
		-I$(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go-grpc_out=$(OUT_DIR)\
		$(PROTO_DIR)/remna.proto\
		$(PROTO_DIR)/order.proto

clean:
	rm -f $(OUT_DIR)*.pb.go