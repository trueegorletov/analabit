PROTO_SRC=service/aggregator/proto/aggregator.proto
PRODUCER_PROTO_SRC=service/producer/proto/producer.proto
GO_OUT=./
MODULE=analabit
TABULA_VERSION=1.0.5
TABULA_JAR=tools/tabula.jar

.PHONY: check-deps
check-deps:
	@which pdftotext > /dev/null 2>&1 || (echo "Error: pdftotext not found. Please install poppler-utils package. On Ubuntu/Debian: sudo apt-get install poppler-utils. On Alpine: apk add poppler-utils" && exit 1)

.PHONY: proto
proto:
	protoc --go_out=$(GO_OUT) --micro_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_SRC)
	protoc --go_out=$(GO_OUT) --micro_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PRODUCER_PROTO_SRC)

.PHONY: tools
tools: $(TABULA_JAR)

$(TABULA_JAR):
	@mkdir -p tools
	@echo "Downloading Tabula $(TABULA_VERSION)..."
	@curl -L -o $(TABULA_JAR) https://github.com/tabulapdf/tabula-java/releases/download/v$(TABULA_VERSION)/tabula-$(TABULA_VERSION)-jar-with-dependencies.jar
	@echo "Tabula JAR downloaded to $(TABULA_JAR)"

.PHONY: dev
dev: check-deps tools
	bash scripts/dev.sh

.PHONY: test
test: check-deps
	go test ./...
