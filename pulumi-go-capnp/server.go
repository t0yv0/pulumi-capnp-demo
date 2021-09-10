package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"zombiezen.com/go/capnproto2/rpc"
)

// promiseBroker is a local implementation of PromiseBroker.
type promiseBroker struct{}

func (pb promiseBroker) NewPromise(call PromiseBroker_newPromise) error {
	fmt.Printf("NewPromise\n")
	// Create a new locally implemented Promise capability.
	p := Promise_ServerToClient(promiseImpl{state: &promiseState{}})
	return call.Results.SetPromise(p)
}

// promiseImpl a local implementation of Promise.
// TODO mutex
type promiseImpl struct {
	state *promiseState
}

type promiseState struct {
	gotBytes  bool
	callbacks []PromiseCallback
	bytes     []byte
}

func (ps promiseImpl) Resolve(call Promise_resolve) error {
	data, err := call.Params.Data()
	if err != nil {
		return err
	}

	st := ps.state

	fmt.Printf("promise.Resolve(data: %s)\n", string(data))

	if st.gotBytes {
		return fmt.Errorf("Cannot Resolve a promise twice")
	}

	// dataCopy := make([]byte, len(data))
	// copy(dataCopy, data)

	st.bytes = data
	st.gotBytes = true

	fmt.Printf("Found %d callbacks\n", len(st.callbacks))

	for _, cb := range st.callbacks {
		fmt.Printf("calling cb.Done()\n")
		cb.Done(context.TODO(),
			func(params PromiseCallback_done_Params) error {
				params.SetData(st.bytes)
				return nil
			})
	}

	st.callbacks = nil

	return nil
}

func (ps promiseImpl) Then(call Promise_then) error {
	fmt.Printf("promise.Then\n")
	cb := call.Params.Callback()
	st := ps.state
	if st.gotBytes {
		cb.Done(context.TODO(),
			func(params PromiseCallback_done_Params) error {
				params.SetData(st.bytes)
				return nil
			})

	} else {
		st.callbacks = append(st.callbacks, cb)
	}
	return nil
}

func runServer(c net.Conn, pb *PromiseBroker) error {
	// Listen for calls, using the PromiseBroker as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(pb.Client))
	// Wait for connection to abort.
	err := conn.Wait()
	return err
}

func serverMain() {
	fmt.Println("Starting server on port :8000...\n")

	fmt.Println("Allocating PromiseBroker")
	pb := PromiseBroker_ServerToClient(promiseBroker{})

	// listen on port 8000
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error connecting: %v\n", err)
			continue
		}

		fmt.Printf("Serving to %s\n", conn.RemoteAddr().String())

		err = runServer(conn, &pb)

		if err == io.EOF {
			fmt.Printf("Done serving to %s\n", conn.RemoteAddr().String())
			continue
		}

		if err != nil {
			fmt.Printf("Error serving to %s: %v\n", conn.RemoteAddr().String(), err)
			continue
		}

		fmt.Printf("Done serving to %s\n", conn.RemoteAddr().String())
	}
}
