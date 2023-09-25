#!/bin/sh

export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin

rm go.mod

go mod init server

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/golang/protobuf/protoc-gen-go@latest

go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest
go mod tidy
go mod vendor

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server.proto

go build -mod=vendor



./server