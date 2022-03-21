package types

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/txchat/addrbook/commands"
	"github.com/txchat/addrbook/executor"
	"github.com/txchat/addrbook/rpc"
	chattypes "github.com/txchat/addrbook/types"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     chattypes.ChatX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
