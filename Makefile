
pulumi-go-capnp/books.capnp.go:	local_gopath schemas/books.capnp
	PATH=$(PATH):$(PWD)/local_gopath/bin capnp compile -I$(PWD)/local_gopath/pkg/mod/zombiezen.com/go/capnproto2@v2.18.2+incompatible/std -ogo:./pulumi-go-capnp schemas/books.capnp --src-prefix=schemas

local_gopath:
	GOPATH=$(PWD)/local_gopath go install zombiezen.com/go/capnproto2/capnpc-go@latest


clean::
	chmod -R 0777 local_gopath
	rm -rf local_gopath
	rm pulumi-go-capnp/*.capnp.go
