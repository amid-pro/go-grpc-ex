FROM golang:1.20-alpine

RUN apk update && apk add --no-cache make protobuf-dev


#WORKDIR /go/src/server

#COPY main.go .
#RUN go build -o server main.go
#RUN $(ls -la)
#CMD [". /server"]

#RUN apk update && apk add --no-cache make protobuf-dev

#RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#RUN go install github.com/golang/protobuf/protoc-gen-go@latest

#COPY main.go .

#RUN cat main.go

#CMD ["protoc"]

#RUN protoc