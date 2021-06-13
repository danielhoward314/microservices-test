## GRPC and protoc

1. set up go module in root directory: `go mod init github.com/<uname>/<module>`
2. in `.proto` file, include the following after the `syntax = "proto3";` line: `option go_package = "github.com/<uname>/<module>";`
3. `go get google.golang.org/grpc/cmd/protoc-gen-go-grpc`
4. assuming a dir `github.com/<uname>/<module>/protos` that contains the target `.proto` files to compile, run the following command from `github.com/<uname>/<module>`:

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/<filename>.proto
```

This will output two files:

1. A file with name `<filename>.pb.go` that handles only the binary serializing and deserializing.
2. A file with name `<filename>_go.pb.go` that exposes the client and server interfaces/methods for grpc.

Note the grpc server struct that will get registered has:

```
type Currency struct {
	protos.UnimplementedCurrencyServer
    // ...
}
```

curl-like equivalent for grpc: `brew install grpcurl`
list services and methods: `grpcurl --plaintext localhost:9092 list`
get signature of `<Service>.<Method>`: `grpcurl --plaintext localhost:9092 describe Currency.GetRate`
get signature of `message`: `grpcurl --plaintext localhost:9092 describe .RateRequest`
call a method (note case sensitivity of json key names): 
```
grpcurl --plaintext -d '{"Base": "GBP", "Destination": "USD"}' localhost:9092 Currency.GetRate
```