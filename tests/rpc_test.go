package tests

import (
	"math/rand"
	"net/rpc"
	"sync"
	"testing"
	"time"

	"github.com/hqpko/hnet"
)

var (
	rpcAddr = "127.0.0.1:12005"
	rpcOnce = new(sync.Once)
)

type RPCReq struct {
}

func (r *RPCReq) Add(req *Req, resp *Resp) error {
	resp.B = req.A + 1
	return nil
}

func Benchmark_go_rpc_Call(b *testing.B) {
	startRpcServer()

	s, _ := hnet.ConnectSocket(rpcAddr)
	client := rpc.NewClient(s)
	b.StartTimer()
	defer b.StopTimer()
	req := &Req{A: 1}
	reply := &Resp{}
	for i := 0; i < b.N; i++ {
		if err := client.Call("RPCReq.Add", req, reply); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_go_rpc_Call_Concurrency(b *testing.B) {
	startRpcServer()

	s, _ := hnet.ConnectSocket(rpcAddr)
	client := rpc.NewClient(s)
	b.StartTimer()
	defer b.StopTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := &HelloRequest{A: rand.Int31n(100)}
		reply := &Resp{}
		for pb.Next() {
			if err := client.Call("RPCReq.Add", req, reply); err != nil {
				b.Fatal(err)
			} else if reply.B != req.A+1 {
				b.Fatal("resp.B != req.A+1")
			}
		}
		req.A++
	})
}

func Benchmark_go_rpc_Go(b *testing.B) {
	startRpcServer()

	s, _ := hnet.ConnectSocket(rpcAddr)
	client := rpc.NewClient(s)
	b.StartTimer()
	defer b.StopTimer()
	req := &Req{A: 1}
	reply := &Resp{}
	for i := 0; i < b.N; i++ {
		_ = client.Go("RPCReq.Add", req, reply, make(chan *rpc.Call, 1))
	}
}

func startRpcServer() {
	rpcOnce.Do(func() {
		go func() {
			_ = hnet.ListenSocket(rpcAddr, func(socket *hnet.Socket) {
				server := rpc.NewServer()
				_ = server.Register(new(RPCReq))
				go server.ServeConn(socket)
			})
		}()

		time.Sleep(100 * time.Millisecond)
	})
}
