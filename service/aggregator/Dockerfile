FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN protoc --go_out=. --micro_out=. --go_opt=module=analabit --go-grpc_out=. --go-grpc_opt=module=analabit service/aggregator/proto/aggregator.proto

RUN go mod tidy

RUN go build -o aggregator-service ./service/aggregator/main.go

EXPOSE 8081

CMD [ "./aggregator-service" ]
