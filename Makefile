build::	pulumi-go-capnp/schemas.capnp.go
	(cd pulumi-go-capnp && go build .)

pulumi-go-capnp/schemas.capnp.go:	local_gopath schemas/schemas.capnp
	PATH=$(PATH):$(PWD)/local_gopath/bin capnp compile -I$(PWD)/local_gopath/pkg/mod/zombiezen.com/go/capnproto2@v2.18.2+incompatible/std -ogo:./pulumi-go-capnp schemas/schemas.capnp --src-prefix=schemas

local_gopath:
	GOPATH=$(PWD)/local_gopath go install zombiezen.com/go/capnproto2/capnpc-go@latest


clean::
	chmod -R 0777 local_gopath
	rm -rf local_gopath
	rm pulumi-go-capnp/*.capnp.go
