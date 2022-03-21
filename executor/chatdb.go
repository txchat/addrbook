package executor

import (
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	chattypes "github.com/txchat/addrbook/types"
)

//Action 具体动作执行
type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	localDB      dbm.KVDB
	index        int
	mainHeight   int64
}

//NewAction 生成Action对象
func NewAction(Chat *Chat, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromAddr := tx.From()

	return &Action{
		coinsAccount: Chat.GetCoinsAccount(),
		db:           Chat.GetStateDB(),
		txhash:       hash,
		fromaddr:     fromAddr,
		blocktime:    Chat.GetBlockTime(),
		height:       Chat.GetHeight(),
		execaddr:     dapp.ExecAddress(string(tx.Execer)),
		localDB:      Chat.GetLocalDB(),
		index:        index,
		mainHeight:   Chat.GetMainHeight(),
	}
}

//UpdateFriend 更新好友信息
func (action *Action) UpdateFriend(send *chattypes.UpdateFriends) (*types.Receipt, error) {
	var logs []*types.ReceiptLog

	//生成日志
	receiptLog := action.getUpdateReceiptLog(send)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}, nil
}

//UpdateBlackList 更新黑名单信息
func (action *Action) UpdateBlackList(send *chattypes.UpdateBlackList) (*types.Receipt, error) {
	var logs []*types.ReceiptLog

	//生成日志
	receiptLog := action.getUpdateBlackListReceiptLog(send)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}, nil
}

//getUpdateReceiptLog 更新好友日志
func (action *Action) getUpdateReceiptLog(send *chattypes.UpdateFriends) *types.ReceiptLog {
	log := &types.ReceiptLog{
		Ty: chattypes.TyLogUpdateFriends,
	}

	r := &chattypes.UpdateFriends{
		Friends: send.Friends,
	}

	//dapp.HeightIndexStr(action.height, int64(action.index))
	log.Log = types.Encode(r)
	return log
}

//getUpdateReceiptLog 更新黑名单列表日志
func (action *Action) getUpdateBlackListReceiptLog(send *chattypes.UpdateBlackList) *types.ReceiptLog {
	log := &types.ReceiptLog{
		Ty: chattypes.TyLogUpdateBlackList,
	}

	r := &chattypes.UpdateBlackList{
		List: send.List,
	}

	//dapp.HeightIndexStr(action.height, int64(action.index))
	log.Log = types.Encode(r)
	return log
}

//UpdateFriend 更新好友信息
func (action *Action) UpdateUser(send *chattypes.UpdateFields) (*types.Receipt, error) {
	var logs []*types.ReceiptLog

	//生成日志
	receiptLog := action.getUpdateUserReceiptLog(send)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}, nil
}

//UpdateBlackList 更新黑名单信息
func (action *Action) UpdateServerGroup(send *chattypes.UpdateServerGroups) (*types.Receipt, error) {
	var logs []*types.ReceiptLog

	//生成日志
	receiptLog := action.getUpdateServerGroupReceiptLog(send)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}, nil
}

//getUpdateUserReceiptLog 更新用户信息日志
func (action *Action) getUpdateUserReceiptLog(send *chattypes.UpdateFields) *types.ReceiptLog {
	log := &types.ReceiptLog{
		Ty: chattypes.TyLogUpdateUser,
	}

	r := &chattypes.UpdateFields{
		Fields: send.Fields,
	}

	//dapp.HeightIndexStr(action.height, int64(action.index))
	log.Log = types.Encode(r)
	return log
}

//getUpdateServerGroupReceiptLog 更新服务分组列表日志
func (action *Action) getUpdateServerGroupReceiptLog(send *chattypes.UpdateServerGroups) *types.ReceiptLog {
	log := &types.ReceiptLog{
		Ty: chattypes.TyLogUpdateServerGroup,
	}

	r := &chattypes.UpdateServerGroups{
		Groups: send.Groups,
	}

	//dapp.HeightIndexStr(action.height, int64(action.index))
	log.Log = types.Encode(r)
	return log
}
