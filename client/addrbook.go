package client

import (
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	atypes "github.com/txchat/addrbook/types"
)

const (
	signType = types.SECP256K1
)

var (
	cfg = types.NewChain33Config(types.GetDefaultCfgstring())
	ety = atypes.NewType(nil)
	c   crypto.Crypto
)

func init() {
	cr, err := crypto.Load(types.GetSignName("", signType), -1)
	if err != nil {
		panic(err)
	}
	c = cr
}

type AddrbookClient struct {
	client   Cli
	txHeight int64
}

func NewAddrbookCient(cli Cli) *AddrbookClient {
	return &AddrbookClient{
		client:   cli,
		txHeight: types.LowAllowPackHeight,
	}
}

func (c *AddrbookClient) UpdateFriends(priv crypto.PrivKey, data []*atypes.ReqUpdateFriend) ([]*types.ReceiptLog, error) {
	req := &atypes.ChatAction{Value: &atypes.ChatAction_Update{Update: &atypes.UpdateFriends{Friends: data}}}
	return c.doSend(req, priv)
}

func (c *AddrbookClient) UpdateBlackList(priv crypto.PrivKey, data []*atypes.ReqUpdateBlackList) ([]*types.ReceiptLog, error) {
	req := &atypes.ChatAction{Value: &atypes.ChatAction_Black{Black: &atypes.UpdateBlackList{List: data}}}
	return c.doSend(req, priv)
}

func (c *AddrbookClient) UpdateUser(priv crypto.PrivKey, data []*atypes.ReqUpdateField) ([]*types.ReceiptLog, error) {
	req := &atypes.ChatAction{Value: &atypes.ChatAction_UpdateUser{UpdateUser: &atypes.UpdateFields{Fields: data}}}
	return c.doSend(req, priv)
}

func (c *AddrbookClient) UpdateServerGroup(priv crypto.PrivKey, data []*atypes.ReqUpdateServerGroup) ([]*types.ReceiptLog, error) {
	req := &atypes.ChatAction{Value: &atypes.ChatAction_UpdateServerGroup{UpdateServerGroup: &atypes.UpdateServerGroups{Groups: data}}}
	return c.doSend(req, priv)
}

func (c *AddrbookClient) QueryFriends(req *atypes.ReqGetFriends) (*atypes.ReplyGetFriends, error) {
	data, err := c.client.Query("GetFriends", req)
	if err != nil {
		return nil, err
	}

	var reply atypes.ReplyGetFriends
	err = types.Decode(data, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *AddrbookClient) QueryBlackList(req *atypes.ReqGetBlackList) (*atypes.ReplyGetBlackList, error) {
	data, err := c.client.Query("GetBlackList", req)
	if err != nil {
		return nil, err
	}

	var reply atypes.ReplyGetBlackList
	err = types.Decode(data, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *AddrbookClient) QueryUser(req *atypes.ReqGetUser) (*atypes.ReplyGetUser, error) {
	data, err := c.client.Query("QueryUser", req)
	if err != nil {
		return nil, err
	}

	var reply atypes.ReplyGetUser
	err = types.Decode(data, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *AddrbookClient) QueryServerGroup(req *atypes.ReqGetServerGroup) (*atypes.ReplyGetServerGroup, error) {
	data, err := c.client.Query("QueryServerGroup", req)
	if err != nil {
		return nil, err
	}

	var reply atypes.ReplyGetServerGroup
	err = types.Decode(data, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *AddrbookClient) doSend(wr *atypes.ChatAction, priv crypto.PrivKey) ([]*types.ReceiptLog, error) {
	var tx *types.Transaction
	var err error

	if wr.GetUpdate() != nil {
		tx, err = ety.Create("Update", wr.GetUpdate())
	} else if wr.GetBlack() != nil {
		tx, err = ety.Create("Black", wr.GetBlack())
	} else if wr.GetUpdateUser() != nil {
		tx, err = ety.Create("UpdateUser", wr.GetUpdateUser())
	} else if wr.GetUpdateServerGroup() != nil {
		tx, err = ety.Create("UpdateServerGroup", wr.GetUpdateServerGroup())
	}
	if err != nil {
		return nil, err
	}

	if c.txHeight > 0 {
		curHeight, err := c.client.GetLastHeader()
		if err != nil {
			return nil, err
		}
		// curHeight-90 <= x <= curHeight+30
		tx.Expire = types.TxHeightFlag + curHeight + c.txHeight
	}

	tx, err = types.FormatTx(cfg, cfg.GetTitle()+atypes.ChatX, tx)
	if err != nil {
		return nil, err
	}

	tx.Sign(signType, priv)

	return c.client.Send(tx)
}
