module github.com/hqpko/hrpc-bench

go 1.16

//replace github.com/hqpko/hrpc => ../hrpc

//replace google.golang.org/grpc => ../grpc-go

require (
	github.com/golang/protobuf v1.5.3
	github.com/hprose/hprose-golang v2.0.6+incompatible
	github.com/hqpko/hnet v0.6.0
	github.com/hqpko/hrpc v0.9.0
	github.com/lesismal/arpc v1.2.11
	github.com/smallnest/rpcx v1.8.3
	google.golang.org/grpc v1.53.0
)
