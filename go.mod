module protoactor-go-sender-example

go 1.12

// See https://github.com/etcd-io/etcd/issues/12124
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/asynkron/protoactor-go v0.0.0-20220616142548-afd2d973a1d1
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/mitchellh/go-testing-interface v1.14.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.4
)
