PROTO_DIR := .
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all clean

all: generate

generate:
	protoc \
		--go_out=$(PROTO_DIR) \
		--go-grpc_out=$(PROTO_DIR)\
		$(PROTO_FILES)

clean:
	rm -f *.pb.go