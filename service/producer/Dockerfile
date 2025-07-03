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
#   service/producer/proto/producer.proto

RUN go mod tidy

RUN go build -o producer-service ./service/producer/main.go

EXPOSE 8080

CMD [ "./producer-service" ]
