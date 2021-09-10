using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("main");
$Go.import("main");

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
