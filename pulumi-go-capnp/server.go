package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"log"
	"net"

	//"golang.org/x/net/context"
	"zombiezen.com/go/capnproto2/rpc"
)

// hashFactory is a local implementation of HashFactory.
type hashFactory struct{}

func (hf hashFactory) NewSha1(call HashFactory_newSha1) error {
	// Create a new locally implemented Hash capability.
	hs := Hash_ServerToClient(hashServer{sha1.New()})
	// Notice that methods can return other interfaces.
	return call.Results.SetHash(hs)
}

// hashServer is a local implementation of Hash.
type hashServer struct {
	h hash.Hash
}

func (hs hashServer) Write(call Hash_write) error {
	data, err := call.Params.Data()
	if err != nil {
		return err
	}
	_, err = hs.h.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (hs hashServer) Sum(call Hash_sum) error {
	s := hs.h.Sum(nil)
	return call.Results.SetHash(s)
}

func runServer(c net.Conn) error {
	// Create a new locally implemented HashFactory.
	main := HashFactory_ServerToClient(hashFactory{})
	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(main.Client))
	// Wait for connection to abort.
	err := conn.Wait()
	return err
}

func serverMain() {
	fmt.Println("Starting server on port :8000...\n")

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

		err = runServer(conn)

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
