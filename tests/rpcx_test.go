package tests

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
)

var (
	rpcxAddr = "127.0.0.1:12004"
	rpcxOnec = new(sync.Once)
)

type RPCXReq struct {
}

func (r *RPCXReq) Mul(ctx context.Context, args *Req, reply *Resp) error {
	reply.B = args.A + 1
	return nil
}

func Benchmark_rpcx_Call(b *testing.B) {
	startRpcxServer()

	c := client.NewClient(client.DefaultOption)
	if err := c.Connect("tcp", rpcxAddr); err != nil {
		b.Fatal(err)
	}
	req := &Req{A: 1}
	reply := &Resp{}
	for i := 0; i < b.N; i++ {
		if err := c.Call(context.Background(), "RPCXReq", "Mul", req, reply); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_rpcx_Call_Concurrency(b *testing.B) {
	startRpcxServer()

	c := client.NewClient(client.DefaultOption)
	if err := c.Connect("tcp", rpcxAddr); err != nil {
		b.Fatal(err)
	}
	b.RunParallel(func(pb *testing.PB) {
		req := &HelloRequest{A: rand.Int31n(100)}
		reply := &Resp{}
		for pb.Next() {
			if err := c.Call(context.Background(), "RPCXReq", "Mul", req, reply); err != nil {
				b.Fatal(err)
			} else if reply.B != req.A+1 {
				b.Fatal("resp.B != req.A+1")
			}
			req.A++
		}
	})
}

func Benchmark_rpcx_Go(b *testing.B) {
	startRpcxServer()

	c := client.NewClient(client.DefaultOption)
	if err := c.Connect("tcp", rpcxAddr); err != nil {
		b.Fatal(err)
	}
	req := &Req{A: 1}
	reply := &Resp{}
	for i := 0; i < b.N; i++ {
		_ = c.Go(context.Background(), "RPCXReq", "Mul", req, reply, make(chan *client.Call, 1))
	}
}

func startRpcxServer() {
	rpcxOnec.Do(func() {
		go func() {
			s := server.NewServer()
			_ = s.Register(new(RPCXReq), "")
			_ = s.Serve("tcp", rpcxAddr)
		}()

		time.Sleep(100 * time.Millisecond)
	})
}
