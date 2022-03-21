package executor

import (
	"github.com/33cn/chain33/types"
	chattypes "github.com/txchat/addrbook/types"
)

/*
 * 实现交易的链上执行接口
 * 关键数据上链（statedb）并生成交易回执（log）
 */

func (c *Chat) Exec_Update(payload *chattypes.UpdateFriends, tx *types.Transaction, index int) (*types.Receipt, error) {
	//implement code
	action := NewAction(c, tx, index)
	return action.UpdateFriend(payload)
}

func (c *Chat) Exec_Black(payload *chattypes.UpdateBlackList, tx *types.Transaction, index int) (*types.Receipt, error) {
	//implement code
	action := NewAction(c, tx, index)
	return action.UpdateBlackList(payload)
}

func (c *Chat) Exec_UpdateUser(payload *chattypes.UpdateFields, tx *types.Transaction, index int) (*types.Receipt, error) {
	//implement code
	action := NewAction(c, tx, index)
	return action.UpdateUser(payload)
}

func (c *Chat) Exec_UpdateServerGroup(payload *chattypes.UpdateServerGroups, tx *types.Transaction, index int) (*types.Receipt, error) {
	//implement code
	action := NewAction(c, tx, index)
	return action.UpdateServerGroup(payload)
}
