package tests

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/hprose/hprose-golang/rpc"
)

var (
	hproseAddr = "tcp://127.0.0.1:12007"
	hproseOnce = new(sync.Once)
)

func hproseHello(a int32) int32 {
	return a + 1
}

// Stub is ...
type HproseStub struct {
	Hello func(int32) (int32, error) `simple:"true" idempotent:"true" retry:"30"`
}

func Benchmark_hprose_Call(b *testing.B) {
	startHproseServer()

	client := rpc.NewClient(hproseAddr)
	defer client.Close()

	var stub *HproseStub
	client.UseService(&stub)

	b.StartTimer()
	defer b.StopTimer()
	for i := 0; i < b.N; i++ {
		if _, err := stub.Hello(1); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_hprose_Call_Concurrency(b *testing.B) {
	startHproseServer()

	client := rpc.NewClient(hproseAddr)
	defer client.Close()

	var stub *HproseStub
	client.UseService(&stub)

	b.StartTimer()
	defer b.StopTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := rand.Int31n(100)
		for pb.Next() {
			if resp, err := stub.Hello(req); err != nil {
				b.Fatal(err)
			} else if resp != req+1 {
				b.Fatal("resp!=req+1")
			}
		}
	})
}

func startHproseServer() {
	hproseOnce.Do(func() {
		server := rpc.NewTCPServer(hproseAddr)
		server.AddFunction("hello", hproseHello)
		server.Handle()
	})
}
