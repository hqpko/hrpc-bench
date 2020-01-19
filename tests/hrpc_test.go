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

	socket, _ := hnet.ConnectSocket("tcp", hrpcAddr)
	client := hrpc.NewClient(socket)
	go client.Run()
	b.StartTimer()
	defer b.StopTimer()
	args := []byte{1}
	for i := 0; i < b.N; i++ {
		if _, err := client.Call(1, args); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_hrpc_Go(b *testing.B) {
	startHRpcServer()

	socket, _ := hnet.ConnectSocket("tcp", hrpcAddr)
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
			hnet.ListenSocket("tcp", hrpcAddr, func(socket *hnet.Socket) {
				server := hrpc.NewServer(socket)
				server.Register(1, func(seq uint64, args []byte) {
					server.Reply(seq, args)
				})
				go server.Run()
			})
		}()
		time.Sleep(100 * time.Millisecond)
	})
}
