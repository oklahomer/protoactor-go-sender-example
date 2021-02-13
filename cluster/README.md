# Setup
For the latest version of Protoactor-go, specify [gograinv2_out](https://github.com/AsynkronIT/protoactor-go/blob/dev/protobuf/protoc-gen-gograinv2/Makefile) instead of `gograin_out` to generate files.
```
protoc --gograinv2_out=. ./messages/protos.proto
protoc --gogoslick_out=. ./messages/protos.proto

docker build --rm -t protoactor-go-sample:latest consul
docker run -p 8500:8500  protoactor-go-sample
```
