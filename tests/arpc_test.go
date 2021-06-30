package tests

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/lesismal/arpc"
	log2 "github.com/lesismal/arpc/log"
)

var (
	arpcAddr = "127.0.0.1:12008"
	arpcOnce = new(sync.Once)
)

func Benchmark_arpc_Call(b *testing.B) {
	startArpcServer()

	client, err := arpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	defer client.Stop()

	req := &helloReq{A: 1}
	rsp := &helloRsp{}
	for i := 0; i < b.N; i++ {
		if err = client.Call("Hello", req, rsp, time.Second*5); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_arpc_Call_Concurrency(b *testing.B) {
	startArpcServer()

	client, err := arpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	defer client.Stop()

	b.RunParallel(func(pb *testing.PB) {
		req := &HelloRequest{A: rand.Int31n(100)}
		rsp := &helloRsp{}
		for pb.Next() {
			if err = client.Call("Hello", req, rsp, time.Second*5); err != nil {
				b.Fatal(err)
			} else if rsp.B != req.A+1 {
				b.Fatal("resp.B != req.A+1")
			}
			req.A++
		}
	})
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", arpcAddr, time.Second*3)
}

type helloReq struct {
	A int32
}

type helloRsp struct {
	B int32
}

func onHello(ctx *arpc.Context) {
	req := &helloReq{}
	rsp := &helloRsp{}

	ctx.Bind(req)
	// log.Printf("OnHello: \"%v\"", req.Msg)

	rsp.B = req.A + 1
	ctx.Write(rsp)
}
func startArpcServer() {
	arpcOnce.Do(func() {
		go func() {
			log2.SetLevel(log2.LevelNone)
			ln, err := net.Listen("tcp", arpcAddr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			svr := arpc.NewServer()
			svr.Handler.Handle("Hello", onHello)
			svr.Serve(ln)
		}()

		time.Sleep(100 * time.Millisecond)
	})
}
