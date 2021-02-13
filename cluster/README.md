# Setup Consul
```
protoc --gograin_out=. ./messages/protos.proto
protoc --gogoslick_out=. ./messages/protos.proto

docker build --rm -t protoactor-go-sample:latest consul
docker run -p 8500:8500  protoactor-go-sample
```
