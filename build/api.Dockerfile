FROM golang:1.26.0 AS builder

RUN apt-get update && apt-get install -y unzip curl

RUN apt-get install -y \
    protobuf-compiler \
    libprotobuf-dev \
    unzip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /app/chronflow
COPY . .

RUN protoc \
    --proto_path=proto \ 
    --go_out=internal/pb \
    --go_opt=paths=source_relative \
    --go-grpc_out=internal/pb \
    --go-grpc_opt=paths=source_relative \
    proto/*.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o chronflow-api ./cmd/api/main.go

FROM alpine:3.21.6
WORKDIR /app/chronflow
COPY --from=builder /app/chronflow/chronflow-api .

CMD ["./chronflow-api"]