# pulumi-capnp-demo

This demo looks at using higher-order RPC programming model provided
by [Cap’n Proto](https://capnproto.org/) for implementing distributed
promises.

We have been exploring how Pulumi can benefit from distributed promise
broker capabilities. That is, we want promises to be constructed,
cross-referenced and resolved across as set of subprocesses in
different languages.

Implementing this over gRPC would involve a `<promise-id: str>` being
passsed through the RPC layer, since gRPC is first-order. In contrast,
here we have more of an illusion of a single-process programming
model:


```
interface PromiseBroker {
	newPromise @0 () -> (promise :Promise);
}

interface PromiseCallback {
        done @0 (data: Data) -> ();
}

interface Promise {
        then @0 (callback :PromiseCallback) -> ();
        resolve @1 (data: Data) -> ();
}
```

It is an interesting programming model and the current implementation
is just a quick demo and is not ready to be relied on. Some questions
remain around lifetimes of the Promise and PromiseCallback proxy
objects and interoperability of N clients.


## Building

```
$ nix-shell --pure
$ make
```

## Running

In one shell:

```
$ ./pulumi-go-capnp/pulumi-capnp-demo server
Starting server on port :8000...

Allocating PromiseBroker
...
```

In another shell:

```
$ ./pulumi-capnp-demo client
Waiting on the promise to be resolved
CB2: Received: OK then..
CB1: Received: OK then..
```
