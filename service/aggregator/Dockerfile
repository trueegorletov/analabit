FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# RUN go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest  # disabled, code already generated

COPY . .

# RUN protoc --proto_path=. \
#   --go_out=paths=source_relative:. \
#   --go-grpc_out=paths=source_relative:. \
#   --micro_out=. \
#   service/aggregator/proto/aggregator.proto

RUN go mod tidy

RUN go build -o aggregator-service ./service/aggregator/main.go

EXPOSE 8081

CMD [ "./aggregator-service" ]
