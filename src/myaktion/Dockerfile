FROM golang:1.20-buster

# non-go modules dependencies
RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# copy code and protobuf
WORKDIR /go/src/app
COPY ./myaktion .
COPY ./banktransfer/grpc/banktransfer/banktransfer.proto ./client/banktransfer/
COPY ./wait-for-it.sh .

RUN go mod download
RUN go generate ./...
RUN go install

RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["myaktion"]
EXPOSE 8000