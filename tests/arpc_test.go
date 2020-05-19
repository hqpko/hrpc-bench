package tests

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/lesismal/arpc"
)

var (
	arpcAddr = "127.0.0.1:12008"
	arpcOnec = new(sync.Once)
)

func Benchmark_arpc_Call(b *testing.B) {
	startArpcServer()

	client, err := arpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	client.Run()
	defer client.Stop()

	req := &helloReq{Msg: "Hello"}
	rsp := &helloRsp{}
	for i := 0; i < b.N; i++ {
		if err = client.Call("Hello", req, rsp, time.Second*5); err != nil {
			b.Fatal(err)
		}
	}
}
func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", arpcAddr, time.Second*3)
}

type helloReq struct {
	Msg string
}

type helloRsp struct {
	Msg string
}

func onHello(ctx *arpc.Context) {
	req := &helloReq{}
	rsp := &helloRsp{}

	ctx.Bind(req)
	// log.Printf("OnHello: \"%v\"", req.Msg)

	rsp.Msg = req.Msg
	ctx.Write(rsp)
}
func startArpcServer() {
	arpcOnec.Do(func() {
		go func() {
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
