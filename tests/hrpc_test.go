package tests

import (
	"encoding/json"
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

	socket, _ := hnet.ConnectSocket(hrpcAddr)
	client := hrpc.NewClient(socket)
	go client.Run()
	b.StartTimer()
	defer b.StopTimer()
	req := &helloReq{Msg: "Hello"}
	rsp := &helloRsp{}
	for i := 0; i < b.N; i++ {
		args, _ := json.Marshal(req)
		if resp, err := client.Call(1, args); err != nil {
			b.Fatal(err)
		} else {
			_ = json.Unmarshal(resp, rsp)
			if rsp.Msg != req.Msg {
				b.Errorf("fail hrpc %s %s", rsp.Msg, req.Msg)
			}
		}
	}
}

func Benchmark_hrpc_Go(b *testing.B) {
	startHRpcServer()

	socket, _ := hnet.ConnectSocket(hrpcAddr)
	client := hrpc.NewClient(socket)
	go client.Run()
	b.StartTimer()
	defer b.StopTimer()
	args := []byte{1}
	for i := 0; i < b.N; i++ {
		client.Go(1, args)
	}
}

func startHRpcServer() {
	hrpcOnce.Do(func() {
		go func() {
			hnet.ListenSocket(hrpcAddr, func(socket *hnet.Socket) {
				server := hrpc.NewServer(socket)
				server.SetHandlerCall(func(pid int32, seq uint64, args []byte) {
					req := &helloReq{}
					_ = json.Unmarshal(args, req)
					rsp := helloRsp{Msg: req.Msg}
					resp, _ := json.Marshal(rsp)
					_ = server.Reply(seq, resp)
				})
				go server.Run()
			})
		}()
		time.Sleep(100 * time.Millisecond)
	})
}
