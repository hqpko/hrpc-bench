module github.com/hqpko/hrpc-bench

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/hprose/hprose-golang v2.0.6+incompatible
	github.com/hqpko/hnet v0.5.3
	github.com/hqpko/hrpc v0.8.0
	github.com/lesismal/arpc v1.1.0
	github.com/smallnest/rpcx v1.6.4
	google.golang.org/grpc v1.38.0
	google.golang.org/grpc/examples v0.0.0-20210628165121-83f9def5feb3 // indirect
)

//replace github.com/hqpko/hrpc => ../hrpc

//replace google.golang.org/grpc => ../grpc-go
