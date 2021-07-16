### GRPC

```bash
# ref: https://www.grpc.io/docs/languages/go/quickstart/#get-the-example-code
# ref: https://github.com/philips/grpc-gateway-example

apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

go get -u  google.golang.org/protobuf/cmd/protoc-gen-go
go get -u  google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greeting.proto

protoc -I. \
    --proto_path=/home/augustu/Software/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.9.5/third_party/googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=logtostderr=true:. \
    --swagger_out=logtostderr=true:. \
    greeting.proto


```



