package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	chattypes "github.com/txchat/addrbook/types"
)

var (
	//日志
	elog = log.New("module", "Chat.executor")
)

var driverName = chattypes.ChatX

// Init 重命名执行器名称
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	drivers.Register(cfg, GetName(), NewChat, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Chat{}))
}

type Chat struct {
	drivers.DriverBase
}

func NewChat() drivers.Driver {
	t := &Chat{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName get driver name
func GetName() string {
	return NewChat().GetName()
}

func (c *Chat) GetDriverName() string {
	return driverName
}
