Install `protoc` command.
```
$ brew install protobuf
```

Install `protoc-gen-go`.
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Generate the Go code from the IDL.
```
$ protoc --go_out=. --go_opt=paths=source_relative protos.proto
```
