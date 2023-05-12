FROM golang:1.20-buster

WORKDIR /go/src/app
COPY ./banktransfer .

RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

RUN go mod download
RUN go generate ./...
RUN go install

CMD ["banktransfer"]

EXPOSE 9111
