package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hqpko/hnet"
	"github.com/hqpko/hrpc"
)

var (
	hrpcAddr = "127.0.0.1:12003"
	hrpcOnce = new(sync.Once)
)

func Benchmark_hrpc_Call(b *testing.B) {
	startHRpcServer()

	socket, _ := hnet.ConnectSocket(hrpcAddr)
	client := hrpc.NewClient(socket)
	go client.Run()
	b.StartTimer()
	defer b.StopTimer()
	req := &Req{A: 1}
	for i := 0; i < b.N; i++ {
		args, _ := proto.Marshal(req)
		resp := &Resp{}
		if respData, err := client.Call(1, args); err != nil {
			b.Fatal(err)
		} else {
			_ = proto.Unmarshal(respData, resp)
		}
	}
}

func Benchmark_hrpc_Go(b *testing.B) {
	startHRpcServer()

	socket, _ := hnet.ConnectSocket(hrpcAddr)
	client := hrpc.NewClient(socket)
	req := &Req{A: 1}
	for i := 0; i < b.N; i++ {
		args, _ := proto.Marshal(req)
		if err := client.OneWay(2, args); err != nil {
			b.Fatal(err)
		}
	}
}

func startHRpcServer() {
	hrpcOnce.Do(func() {
		go func() {
			hnet.ListenSocket(hrpcAddr, func(socket *hnet.Socket) {
				server := hrpc.NewServer(socket)
				server.SetHandlerCall(func(pid int32, seq uint64, args []byte) {
					_ = server.Reply(seq, args)
				})
				server.SetHandlerOneWay(func(pid int32, args []byte) {})
				go server.Run()
			})
		}()
		time.Sleep(100 * time.Millisecond)
	})
}
