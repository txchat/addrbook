#!/bin/sh
# proto生成命令，将pb.go文件生成到types/uwallet目录下
#protoc --go_out=plugins=grpc:../types/uwallet ./*.proto

chain33_path=$(go list -f '{{.Dir}}' "github.com/33cn/chain33")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${chain33_path}/types/proto/"
