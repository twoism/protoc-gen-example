package clients

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/twoism/protoc-gen-example/clients/ruby"
)

func GenerateAll(gen *generator.Generator) {
	for _, file := range gen.Request.GetProtoFile() {
		for _, srv := range file.GetService() {
			// Generate http python client and server code
			client := ruby.New(srv, file)
			gen.Response.File = append(gen.Response.File, client.File())
		}
	}
}
