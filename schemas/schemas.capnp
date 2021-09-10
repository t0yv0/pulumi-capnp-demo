using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("main");
$Go.import("foo/books");

interface HashFactory {
	newSha1 @0 () -> (hash :Hash);
}

interface Hash {
	write @0 (data :Data) -> ();
	sum @1 () -> (hash :Data);
}
