Install `protoc` command.
```
$ brew install protobuf
```

Install `protoc-gen-go`.
```
$ go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
```

Generate the Go code from the IDL.
```
$ protoc --go_out=. --go_opt=paths=source_relative --proto_path=. protos.proto
```