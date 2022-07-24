Install `protoc` command.
```
$ brew install protobuf
```

Install `protoc-gen-gograinv2`.
```
$ go install github.com/asynkron/protoactor-go/protobuf/protoc-gen-gograinv2
```

Generate the Go code from the IDL.
```
# For remoting
$ protoc --go_out=. --go_opt=paths=source_relative protos.proto

# For Cluster grain
$ protoc --gograinv2_out=paths=source_relative:. protos.proto
```
