export PATH := $(PATH):.

.PHONY: protos-example
protos-example:
	protoc -I api-common-protos -I . --go_out=generated --go-resource_out=generated ./examples/external/external.proto
	protoc -I api-common-protos -I . --go_out=generated --go-resource_out=generated ./examples/multipattern/multipattern.proto
	protoc -I api-common-protos -I . --go_out=generated --go-resource_out=generated ./examples/namefield/namefield.proto
	protoc -I api-common-protos -I . --go_out=generated --go-resource_out=generated ./examples/simple/simple.proto

# This is done differently to test that it's possible to generate code
# referencing a resource, even if that resource is not explicitly imported. All
# that matters is that it's part of the set of files to generate.
	protoc -I api-common-protos -I . --go_out=generated --go-resource_out=generated ./examples/reference/reference.proto ./examples/simple/simple.proto ./examples/external/external.proto

.PHONY: clean
clean:
	rm -rf generated/github.com
	rm -rf generated/google.golang.org
	rm -rf generated/example.com/external
	rm -rf generated/example.com/multipattern
	rm -rf generated/example.com/namefield
	rm -rf generated/example.com/reference
	rm -rf generated/example.com/simple
