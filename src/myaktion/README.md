## Development

### Ensure dependencies in IDE

The following steps are only necessary for resolving dependencies in the
IDE and for local testing. The docker build will work nevertheless.

1. Copy protobuf file from `banktransfer` service:

       # cp ../banktransfer/grpc/banktransfer/banktransfer.proto ./client/banktransfer/

2. Run `go generate` in this folder:

       # go generate ./...


