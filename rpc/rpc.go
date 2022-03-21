package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"

	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"

	chattypes "github.com/txchat/addrbook/types"
)

/*
 * 实现json rpc和grpc service接口
 * json rpc用Jrpc结构作为接收实例
 * grpc使用channelClient结构作为接收实例
 */

func (c *channelClient) Update(ctx context.Context, v *chattypes.UpdateFriends) (*types.UnsignTx, error) {
	payload := &chattypes.ChatAction{
		Ty:    chattypes.TyUpdateFriendsAction,
		Value: &chattypes.ChatAction_Update{Update: v},
	}
	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(chattypes.ChatX), types.Encode(payload))
	if err != nil {
		return nil, err
	}
	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) Black(ctx context.Context, v *chattypes.UpdateBlackList) (*types.UnsignTx, error) {
	payload := &chattypes.ChatAction{
		Ty:    chattypes.TyUpdateBlackListAction,
		Value: &chattypes.ChatAction_Black{Black: v},
	}
	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(chattypes.ChatX), types.Encode(payload))
	if err != nil {
		return nil, err
	}
	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) UpdateUser(ctx context.Context, v *chattypes.UpdateFields) (*types.UnsignTx, error) {
	payload := &chattypes.ChatAction{
		Ty:    chattypes.TyUpdateUserAction,
		Value: &chattypes.ChatAction_UpdateUser{UpdateUser: v},
	}
	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(chattypes.ChatX), types.Encode(payload))
	if err != nil {
		return nil, err
	}
	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) UpdateServerGroup(ctx context.Context, v *chattypes.UpdateServerGroups) (*types.UnsignTx, error) {
	payload := &chattypes.ChatAction{
		Ty:    chattypes.TyUpdateServerGroupAction,
		Value: &chattypes.ChatAction_UpdateServerGroup{UpdateServerGroup: v},
	}
	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(chattypes.ChatX), types.Encode(payload))
	if err != nil {
		return nil, err
	}
	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}

// GetFriends 获得好友记录
func (c *channelClient) GetFriends(ctx context.Context, in *chattypes.ReqGetFriends) (*chattypes.ReplyGetFriends, error) {
	cfg := c.GetConfig()
	v, err := c.Query(cfg.ExecName(chattypes.ChatX), chattypes.FuncNameGetFriends, in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*chattypes.ReplyGetFriends); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetFriends 获得好友记录
func (c *channelClient) GetBlackList(ctx context.Context, in *chattypes.ReqGetBlackList) (*chattypes.ReplyGetBlackList, error) {
	cfg := c.GetConfig()
	v, err := c.Query(cfg.ExecName(chattypes.ChatX), chattypes.FuncNameGetBlackList, in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*chattypes.ReplyGetBlackList); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetFriends 获得好友记录
func (c *channelClient) GetUser(ctx context.Context, in *chattypes.ReqGetUser) (*chattypes.ReplyGetUser, error) {
	cfg := c.GetConfig()
	v, err := c.Query(cfg.ExecName(chattypes.ChatX), chattypes.FuncNameGetUser, in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*chattypes.ReplyGetUser); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetFriends 获得好友记录
func (c *channelClient) GetServerGroup(ctx context.Context, in *chattypes.ReqGetServerGroup) (*chattypes.ReplyGetServerGroups, error) {
	cfg := c.GetConfig()
	v, err := c.Query(cfg.ExecName(chattypes.ChatX), chattypes.FuncNameGetServerGroup, in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*chattypes.ReplyGetServerGroups); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

func (j *Jrpc) GetFriends(in *chattypes.ReqGetFriends, result *interface{}) error {
	log.Info("GetFriends jrpc call")
	v, err := j.cli.GetFriends(context.Background(), in)
	if err != nil {
		return err
	}
	*result = v
	return nil
}

func (j *Jrpc) GetBlackList(in *chattypes.ReqGetBlackList, result *interface{}) error {
	log.Info("GetBlackList jrpc call")
	v, err := j.cli.GetBlackList(context.Background(), in)
	if err != nil {
		return err
	}
	*result = v
	return nil
}

func (j *Jrpc) GetUser(in *chattypes.ReqGetUser, result *interface{}) error {
	log.Info("GetUser jrpc call")
	v, err := j.cli.GetUser(context.Background(), in)
	if err != nil {
		return err
	}
	*result = v
	return nil
}

func (j *Jrpc) GetServerGroup(in *chattypes.ReqGetServerGroup, result *interface{}) error {
	log.Info("GetServerGroup jrpc call")
	v, err := j.cli.GetServerGroup(context.Background(), in)
	if err != nil {
		return err
	}
	*result = v
	return nil
}

func (j *Jrpc) CreateRawUpdateFriendTx(in json.RawMessage, result *interface{}) error {
	log.Info("CreateRawUpdateFriendTx jrpc call")
	cfg := j.cli.GetConfig()
	data, err := types.CallCreateTxJSON(cfg, cfg.ExecName(chattypes.ChatX), chattypes.NameUpdateFriendsAction, in)
	if err != nil {
		return err
	}
	//创建交易通常返回十六进制格式原数据
	*result = hex.EncodeToString(data)
	return nil
}

func (j *Jrpc) CreateRawUpdateBlackTx(in json.RawMessage, result *interface{}) error {
	log.Info("CreateRawUpdateBlackTx jrpc call")
	cfg := j.cli.GetConfig()
	data, err := types.CallCreateTxJSON(cfg, cfg.ExecName(chattypes.ChatX), chattypes.NameUpdateBlackListAction, in)
	if err != nil {
		return err
	}
	//创建交易通常返回十六进制格式原数据
	*result = hex.EncodeToString(data)
	return nil
}

func (j *Jrpc) CreateRawUpdateUserTx(in json.RawMessage, result *interface{}) error {
	log.Info("CreateRawUpdateUserTx jrpc call")
	cfg := j.cli.GetConfig()
	data, err := types.CallCreateTxJSON(cfg, cfg.ExecName(chattypes.ChatX), chattypes.NameUpdateUserAction, in)
	if err != nil {
		return err
	}
	//创建交易通常返回十六进制格式原数据
	*result = hex.EncodeToString(data)
	return nil
}

func (j *Jrpc) CreateRawUpdateServerGroupTx(in json.RawMessage, result *interface{}) error {
	log.Info("CreateRawUpdateServerGroupTx jrpc call")
	cfg := j.cli.GetConfig()
	data, err := types.CallCreateTxJSON(cfg, cfg.ExecName(chattypes.ChatX), chattypes.NameUpdateServerGroupAction, in)
	if err != nil {
		return err
	}
	//创建交易通常返回十六进制格式原数据
	*result = hex.EncodeToString(data)
	return nil
}
