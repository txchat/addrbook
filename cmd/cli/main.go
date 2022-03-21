package main

import (
	_ "github.com/33cn/chain33/system"
	"github.com/33cn/chain33/util/cli"
	"github.com/33cn/plugin/cli/buildflags"
	_ "github.com/33cn/plugin/plugin"
	_ "github.com/txchat/addrbook"
)

func main() {
	if buildflags.RPCAddr == "" {
		buildflags.RPCAddr = "http://localhost:8801"
	}
	cli.Run(buildflags.RPCAddr, buildflags.ParaName, "")
}
