#!/bin/sh
path=$(
  cd "$(dirname "$0")" || exit
  pwd
)

# machine x
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case ${OS} in
  "darwin")
  export PATH=${path}/osx_x86_64:$PATH
  ;;
  "linux")
  export PATH=${path}/linux_x86_64:$PATH
  ;;
esac

# proto生成命令，将pb.go文件生成到types/uwallet目录下
#protoc --go_out=plugins=grpc:../types/uwallet ./*.proto

chain33_path=$(go list -f '{{.Dir}}' "github.com/33cn/chain33")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${chain33_path}/types/proto/"
