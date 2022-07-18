package main

import (
	"flag"
	"runtime/debug"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin"
	_ "github.com/txchat/addrbook"

	"github.com/33cn/chain33/util/cli"
)

var percent = flag.Int("p", 0, "SetGCPercent")

func main() {
	flag.Parse()
	if *percent < 0 || *percent > 100 {
		*percent = 0
	}
	if *percent > 0 {
		debug.SetGCPercent(*percent)
	}
	cli.RunChain33("", "")
}
