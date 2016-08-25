#!/bin/bash

set -e

cd $(dirname $0)/..
rm -rf generated/clients/*

go build
#protoc --proto_path=/var/lib/protobuf/src --plugin=protoc-gen-client=./protoc-gen-client --client_out=plugins=pyclient:./generated/clients -I. ./protos/test.proto ./protos/test_http.proto ./protos/pb/lyft/lyft_http_options.proto
protoc  --plugin=protoc-gen-client=./protoc-gen-example \
        --client_out=plugins=protoc-gen-client=:./generated/clients \
        -I. ./protos/hello.proto
protoc --ruby_out=./generated/clients/example/hello -I. ./protos/hello.proto
