package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/hqpko/hnet"

	"github.com/hqpko/hrpc"
)

var (
	hrpcAddr = "127.0.0.1:12003"
	hrpcOnce = new(sync.Once)
)

func Benchmark_hrpc_Call(b *testing.B) {
	startHRpcServer()

	client := hrpc.NewStream()
	go client.Run(getSocket(hrpcAddr))
	time.Sleep(100 * time.Millisecond)
	defer client.Close()
	b.StartTimer()
	defer b.StopTimer()
	reply := &Resp{}
	req := &Req{A: 1}
	for i := 0; i < b.N; i++ {
		if err := client.Call(1, req, reply); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_hrpc_Go(b *testing.B) {
	startHRpcServer()

	client := hrpc.NewStream()
	go client.Run(getSocket(hrpcAddr))
	time.Sleep(100 * time.Millisecond)
	defer client.Close()
	b.StartTimer()
	defer b.StopTimer()
	req := &Req{A: 1}
	reply := &Resp{}
	for i := 0; i < b.N; i++ {
		_ = client.Go(1, req, reply, false)
	}
}

func startHRpcServer() {
	hrpcOnce.Do(func() {
		go func() {
			_ = hnet.ListenSocket("tcp", hrpcAddr, func(socket *hnet.Socket) {
				s := hrpc.NewStream()
				s.Register(1, func(args *Req, reply *Resp) error {
					reply.B = args.A + 1
					return nil
				})
				go func() {
					_ = s.Run(socket)
				}()
			})
		}()
		time.Sleep(100 * time.Millisecond)
	})
}

func getSocket(addr string) *hnet.Socket {
	socket, _ := hnet.ConnectSocket("tcp", addr)
	return socket
}
