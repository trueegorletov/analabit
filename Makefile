PROTO_SRC=service/aggregator/proto/aggregator.proto
PRODUCER_PROTO_SRC=service/producer/proto/producer.proto
GO_OUT=./
MODULE=analabit

.PHONY: proto
proto:
	protoc --go_out=$(GO_OUT) --micro_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_SRC)
	protoc --go_out=$(GO_OUT) --micro_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PRODUCER_PROTO_SRC)
