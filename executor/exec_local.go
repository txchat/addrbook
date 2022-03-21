package executor

import (
	"github.com/txchat/addrbook/executor/find"
	"sort"
	"strconv"
	"time"

	"github.com/33cn/chain33/types"
	chattypes "github.com/txchat/addrbook/types"
)

/*
 * 实现交易相关数据本地执行，数据不上链
 * 非关键数据，本地存储(localDB), 用于辅助查询，效率高
 */

func (c *Chat) ExecLocal_Update(payload *chattypes.UpdateFriends, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	fromAddr := tx.From()
	dbSet, err := c.execLocal(fromAddr, receiptData)
	if err != nil {
		return nil, err
	}
	return c.addAutoRollBack(tx, dbSet.KV), nil
}

func (c *Chat) ExecLocal_Black(payload *chattypes.UpdateBlackList, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	fromAddr := tx.From()
	dbSet, err := c.execLocal(fromAddr, receiptData)
	if err != nil {
		return nil, err
	}
	return c.addAutoRollBack(tx, dbSet.KV), nil
}

func (c *Chat) ExecLocal_UpdateUser(payload *chattypes.UpdateFields, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	fromAddr := tx.From()
	dbSet, err := c.execLocal(fromAddr, receiptData)
	if err != nil {
		return nil, err
	}
	return c.addAutoRollBack(tx, dbSet.KV), nil
}

func (c *Chat) ExecLocal_UpdateServerGroup(payload *chattypes.UpdateServerGroups, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	fromAddr := tx.From()
	dbSet, err := c.execLocal(fromAddr, receiptData)
	if err != nil {
		return nil, err
	}
	return c.addAutoRollBack(tx, dbSet.KV), nil
}

//设置自动回滚
func (c *Chat) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {

	dbSet := &types.LocalDBSet{}
	dbSet.KV = c.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}

func (c *Chat) execLocal(fromAddr string, receipt *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	for _, item := range receipt.Logs {
		switch item.Ty {
		case chattypes.TyLogUpdateFriends:
			return dbSet, c.updateFriends(dbSet, fromAddr, item)
		case chattypes.TyLogUpdateBlackList:
			return dbSet, c.updateBlackList(dbSet, fromAddr, item)
		case chattypes.TyLogUpdateUser:
			return dbSet, c.updateUser(dbSet, fromAddr, item)
		case chattypes.TyLogUpdateServerGroup:
			return dbSet, c.updateServerGroup(dbSet, fromAddr, item)
		}
	}
	return dbSet, nil
}

func (c *Chat) updateFriends(dbSet *types.LocalDBSet, fromAddr string, item *types.ReceiptLog) error {
	var receipt chattypes.UpdateFriends
	err := types.Decode(item.Log, &receipt)
	if err != nil {
		return err
	}
	fTable := chattypes.NewFriendTable(c.GetLocalDB())
	for _, f := range receipt.Friends {
		switch f.Type {
		case chattypes.FriendDelete:
			//delete
			fd := &chattypes.Friend{
				MainAddress:   fromAddr,
				FriendAddress: f.FriendAddress,
			}

			cur := &chattypes.FriendRow{Friend: fd}
			primary, err := cur.Get(chattypes.TableFriendIndexMainAddrFriendAddr)
			if err != nil {
				return err
			}
			err = fTable.Del(primary)
			if err != nil {
				return err
			}
		case chattypes.FriendAppend:
			//query
			fd := &chattypes.Friend{
				MainAddress:   fromAddr,
				FriendAddress: f.FriendAddress,
			}
			cur := &chattypes.FriendRow{Friend: fd}
			primary, err := cur.Get(chattypes.TableFriendIndexMainAddrFriendAddr)
			if err != nil {
				return err
			}

			if r, err := fTable.GetData(primary); r != nil && err == nil {
				//更新
				elog.Info("get friends", "raw", r)
				err = fTable.Update(primary, &chattypes.Friend{
					MainAddress:   fromAddr,
					FriendAddress: f.FriendAddress,
					//CreateTime: NowMillionSecond(),
					CreateTime: c.GetBlockTime() * 1000,
					Groups:     f.Groups,
				})
				if err != nil {
					return err
				}
			} else {
				//新增
				err = fTable.Add(&chattypes.Friend{
					MainAddress:   fromAddr,
					FriendAddress: f.FriendAddress,
					//CreateTime: NowMillionSecond(),
					CreateTime: c.GetBlockTime() * 1000,
					Groups:     f.Groups,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	kv, err := fTable.Save()
	if err != nil {
		return err
	}
	dbSet.KV = append(dbSet.KV, kv...)
	return nil
}

func (c *Chat) updateBlackList(dbSet *types.LocalDBSet, fromAddr string, item *types.ReceiptLog) error {
	var receipt chattypes.UpdateBlackList
	err := types.Decode(item.Log, &receipt)
	if err != nil {
		return err
	}
	bTable := chattypes.NewBlackTable(c.GetLocalDB())
	for _, b := range receipt.List {
		switch b.Type {
		case chattypes.BlackDelete:
			//delete
			fd := &chattypes.Black{
				MainAddress:   fromAddr,
				TargetAddress: b.TargetAddress,
			}

			cur := &chattypes.BlackRow{Black: fd}
			primary, err := cur.Get(chattypes.TableBlackIndexMainAddrTargetAddr)
			if err != nil {
				return err
			}
			err = bTable.Del(primary)
			if err != nil {
				return err
			}
		case chattypes.BlackAppend:
			//query
			fd := &chattypes.Black{
				MainAddress:   fromAddr,
				TargetAddress: b.TargetAddress,
			}
			cur := &chattypes.BlackRow{Black: fd}
			primary, err := cur.Get(chattypes.TableBlackIndexMainAddrTargetAddr)
			if err != nil {
				return err
			}

			if r, err := bTable.GetData(primary); r != nil && err == nil {
				//更新
				elog.Info("get black list", "raw", r)
				err = bTable.Update(primary, &chattypes.Black{
					MainAddress:   fromAddr,
					TargetAddress: b.TargetAddress,
					//CreateTime: NowMillionSecond(),
					CreateTime: c.GetBlockTime() * 1000,
				})
				if err != nil {
					return err
				}
			} else {
				//新增
				err = bTable.Add(&chattypes.Black{
					MainAddress:   fromAddr,
					TargetAddress: b.TargetAddress,
					//CreateTime: NowMillionSecond(),
					CreateTime: c.GetBlockTime() * 1000,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	kv, err := bTable.Save()
	if err != nil {
		return err
	}
	dbSet.KV = append(dbSet.KV, kv...)
	return nil
}

func (c *Chat) updateUser(dbSet *types.LocalDBSet, fromAddr string, item *types.ReceiptLog) error {
	var receipt chattypes.UpdateFields
	err := types.Decode(item.Log, &receipt)
	if err != nil {
		return err
	}
	fTable := chattypes.NewUserTable(c.GetLocalDB())
	for _, f := range receipt.Fields {
		switch f.Type {
		case chattypes.UserDelete:
			//delete
			fd := &chattypes.Field{
				MainAddress: fromAddr,
				Name:        f.Name,
			}

			cur := &chattypes.UserRow{Field: fd}
			primary, err := cur.Get(chattypes.TableUserIndexMainAddrFieldName)
			if err != nil {
				return err
			}
			err = fTable.Del(primary)
			if err != nil {
				return err
			}
		case chattypes.UserAppend:
			//check
			if !CheckFieldName(f.Name) {
				return chattypes.ErrFieldType
			}
			if !CheckLevel(f.Level) {
				return chattypes.ErrLevelTypeUnDefine
			}
			//query
			fd := &chattypes.Field{
				MainAddress: fromAddr,
				Name:        f.Name,
			}
			cur := &chattypes.UserRow{Field: fd}
			primary, err := cur.Get(chattypes.TableUserIndexMainAddrFieldName)
			if err != nil {
				return err
			}

			if r, err := fTable.GetData(primary); r != nil && err == nil {
				//更新
				elog.Info("get user fields", "raw", r)
				err = fTable.Update(primary, &chattypes.Field{
					MainAddress: fromAddr,
					Name:        f.Name,
					Value:       f.Value,
					Level:       f.Level,
				})
				if err != nil {
					return err
				}
			} else {
				//新增
				err = fTable.Add(&chattypes.Field{
					MainAddress: fromAddr,
					Name:        f.Name,
					Value:       f.Value,
					Level:       f.Level,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	kv, err := fTable.Save()
	if err != nil {
		return err
	}
	dbSet.KV = append(dbSet.KV, kv...)
	return nil
}

func (c *Chat) updateServerGroup(dbSet *types.LocalDBSet, fromAddr string, item *types.ReceiptLog) error {
	var receipt chattypes.UpdateServerGroups
	err := types.Decode(item.Log, &receipt)
	if err != nil {
		return err
	}
	fTable := chattypes.NewServerGroupsTable(c.GetLocalDB())
	//获取所有分组信息
	groups, err := c.getServerGroups(fromAddr, "")
	sort.Sort(ByIdGroup(groups))
	for _, f := range receipt.Groups {
		switch f.Type {
		case chattypes.GroupDelete:
			//delete
			fd := &chattypes.ServerGroup{
				MainAddress: fromAddr,
				Id:          f.Id,
			}

			cur := &chattypes.ServerGroupsRow{ServerGroup: fd}
			primary, err := cur.Get(chattypes.TableServerGroupsIndexMainAddrServerId)
			if err != nil {
				return err
			}
			err = fTable.Del(primary)
			if err != nil {
				return err
			}
			idx := find.Find(ByIdGroup(groups), f.Id)
			if idx > 0 && idx < len(groups) {
				groups = append(groups[:idx], groups[idx+1:]...)
			}
		case chattypes.GroupAppend:
			if len(groups) >= chattypes.MaxGroupsLength {
				return chattypes.ErrGroupsArrayOut
			}
			gid, err := GenerateGroupId(groups)
			if err != nil {
				return err
			}
			//未创建分组
			g := &chattypes.ServerGroup{
				MainAddress: fromAddr,
				Id:          gid,
				Name:        f.Name,
				Value:       f.Value,
			}
			err = fTable.Add(g)
			if err != nil {
				return err
			}
			//add
			groups = append(groups, g)
		case chattypes.GroupEdit:
			fd := &chattypes.ServerGroup{
				MainAddress: fromAddr,
				Id:          f.Id,
			}
			cur := &chattypes.ServerGroupsRow{ServerGroup: fd}
			primary, err := cur.Get(chattypes.TableServerGroupsIndexMainAddrServerId)
			if err != nil {
				return err
			}

			if r, err := fTable.GetData(primary); r != nil && err == nil {
				//更新
				elog.Info("get server groups", "raw", r)
				err = fTable.Update(primary, &chattypes.ServerGroup{
					MainAddress: fromAddr,
					Id:          f.Id,
					Name:        f.Name,
					Value:       f.Value,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	kv, err := fTable.Save()
	if err != nil {
		return err
	}
	dbSet.KV = append(dbSet.KV, kv...)
	return nil
}

func NowMillionSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GenerateGroupId(groups []*chattypes.ServerGroup) (string, error) {
	if len(groups) < 1 {
		return "1", nil
	}
	idx, err := strconv.ParseInt(groups[len(groups)-1].Id, 10, 64)
	if err != nil {
		return "", chattypes.ErrGroupIdFailed
	}
	return strconv.FormatInt(int64(idx+1), 10), nil
}
