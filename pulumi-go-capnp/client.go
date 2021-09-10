package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"zombiezen.com/go/capnproto2/rpc"
)

func runTestClient(ctx context.Context, c net.Conn) error {
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.StreamTransport(c))
	defer conn.Close()
	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hf := HashFactory{Client: conn.Bootstrap(ctx)}

	// Now we can call methods on hf, and they will be sent over c.
	s := hf.NewSha1(ctx, func(p HashFactory_newSha1_Params) error {
		return nil
	}).Hash()

	// s refers to a remote Hash.  Method calls are delivered in order.
	s.Write(ctx, func(p Hash_write_Params) error {
		err := p.SetData([]byte("Hello, "))
		return err
	})
	s.Write(ctx, func(p Hash_write_Params) error {
		err := p.SetData([]byte("World!"))
		return err
	})
	// Get the sum, waiting for the result.
	result, err := s.Sum(ctx, func(p Hash_sum_Params) error {
		return nil
	}).Struct()
	if err != nil {
		return err
	}

	// Display the result.
	sha1Val, err := result.Hash()
	if err != nil {
		return err
	}
	fmt.Printf("sha1: %x\n", sha1Val)

	return c.Close()
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
