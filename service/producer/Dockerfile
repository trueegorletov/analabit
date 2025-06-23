FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN protoc --go_out=. --micro_out=. --go_opt=module=analabit --go-grpc_out=. --go-grpc_opt=module=analabit service/producer/proto/producer.proto

RUN go mod tidy

RUN go build -o producer-service ./service/producer/main.go

EXPOSE 8080

CMD [ "./producer-service" ]
