package types

import (
	"reflect"

	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

/*
 * 交易相关类型定义
 * 交易action通常有对应的log结构，用于交易回执日志记录
 * 每一种action和log需要用id数值和name名称加以区分
 */

var (
	//ChatX 执行器名称定义
	ChatX = "chat"
	//定义actionMap
	actionMap = map[string]int32{
		NameUpdateFriendsAction:     TyUpdateFriendsAction,
		NameUpdateBlackListAction:   TyUpdateBlackListAction,
		NameUpdateUserAction:        TyUpdateUserAction,
		NameUpdateServerGroupAction: TyUpdateServerGroupAction,
	}
	//定义log的id和具体log类型及名称，填入具体自定义log类型
	logMap = map[int64]*types.LogInfo{
		//LogID:	{Ty: reflect.TypeOf(LogStruct), Name: LogName},
		TyLogUpdateFriends:     {Ty: reflect.TypeOf(UpdateFriends{}), Name: "TyUpdateFriendsLog"},
		TyLogUpdateBlackList:   {Ty: reflect.TypeOf(UpdateBlackList{}), Name: "TyUpdateBlackListLog"},
		TyLogUpdateUser:        {Ty: reflect.TypeOf(UpdateFields{}), Name: "TyUpdateUserLog"},
		TyLogUpdateServerGroup: {Ty: reflect.TypeOf(UpdateServerGroups{}), Name: "TyUpdateServerGroupLog"},
	}
	tlog = log.New("module", "chat.types")
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(ChatX))
	//注册合约启用高度
	types.RegFork(ChatX, InitFork)
	types.RegExec(ChatX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.DisableCheckFork(true)
	cfg.RegisterDappFork(ChatX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(ChatX, NewType(cfg))
}

type chatType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *chatType {
	c := &chatType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload 获取合约action结构
func (c *chatType) GetPayload() types.Message {
	return &ChatAction{}
}

// GeTypeMap 获取合约action的id和name信息
func (c *chatType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap 获取合约log相关信息
func (c *chatType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
