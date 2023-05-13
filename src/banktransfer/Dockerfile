FROM golang:1.20-buster

RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

WORKDIR /go/src/app
COPY ./banktransfer .

RUN go mod download
RUN go generate ./...
RUN go install

COPY ./wait-for-it.sh .
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["banktransfer"]

EXPOSE 9111
