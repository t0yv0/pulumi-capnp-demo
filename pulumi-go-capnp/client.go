package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"zombiezen.com/go/capnproto2/rpc"
)

type callback struct {
	whenDone func(data []byte) error
}

func (cb callback) Done(p PromiseCallback_done) error {
	data, err := p.Params.Data()
	if err != nil {
		return err
	}
	return cb.whenDone(data)
}

var _ PromiseCallback_Server = callback{}

func then(p Promise, whenDone func(data []byte) error) Promise_then_Results_Promise {
	return p.Then(context.TODO(), func(p Promise_then_Params) error {
		return p.SetCallback(PromiseCallback_ServerToClient(callback{whenDone}))
	})
}

func resolve(p Promise, data []byte) Promise_resolve_Results_Promise {
	return p.Resolve(context.TODO(), func(p Promise_resolve_Params) error {
		return p.SetData(data)
	})

}

func runTestClient(ctx context.Context, c net.Conn) error {
	// Create a connection that we can use to get the PromiseBroker.
	conn := rpc.NewConn(rpc.StreamTransport(c))
	defer conn.Close()
	defer c.Close()

	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	pb := PromiseBroker{Client: conn.Bootstrap(ctx)}

	// Now we can call methods on pb, and they will be sent over c.
	p1 := pb.NewPromise(ctx, func(p PromiseBroker_newPromise_Params) error { return nil }).Promise()

	ch := make(chan int)

	then(p1, func(data []byte) error {
		fmt.Printf("CB1: Received: %s\n", string(data))
		return nil
	})

	then(p1, func(data []byte) error {
		fmt.Printf("CB2: Received: %s\n", string(data))

		ch <- 1
		return nil
	})

	resolve(p1, []byte("OK then.."))

	fmt.Printf("Waiting on the promise to be resolved\n")
	<-ch

	return nil
}

func clientMain() {
	// connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	err = runTestClient(context.Background(), conn)
	if err != nil {
		log.Fatal(err)
	}
}
