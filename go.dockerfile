FROM golang:1.20-alpine

RUN apk update && apk add --no-cache make protobuf-dev