package executor

import (
	"crypto/sha256"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	chattypes "github.com/txchat/addrbook/types"
	"github.com/txchat/addrbook/util"
	"sort"
	"strconv"
)

type ByIdGroup []*chattypes.ServerGroup

func (g ByIdGroup) Len() int      { return len(g) }
func (g ByIdGroup) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g ByIdGroup) Less(i, j int) bool {
	v1, err := strconv.Atoi(g[i].Id)
	if err != nil {
		panic(err)
	}
	v2, err := strconv.Atoi(g[j].Id)
	if err != nil {
		panic(err)
	}

	return v1 < v2
}
func (g ByIdGroup) Compare(i int, val interface{}) int {
	v1, err := strconv.Atoi(g[i].Id)
	if err != nil {
		panic(err)
	}
	v2, err := strconv.Atoi(val.(string))
	if err != nil {
		panic(err)
	}
	switch {
	case v1-v2 < 0:
		return -1
	case v1-v2 > 0:
		return 1
	default:
		return 0
	}
}

func indexServerGroups(gs []*chattypes.ServerGroup) map[string]*chattypes.ServerGroup {
	idxGroups := make(map[string]*chattypes.ServerGroup)
	for _, g := range gs {
		idxGroups[g.Id] = g
	}
	return idxGroups
}

func (r *Chat) getServerGroups(userAddr, index string) ([]*chattypes.ServerGroup, error) {
	rpData := &chattypes.ServerGroup{
		MainAddress: userAddr,
		Id:          index,
	}
	var indexName string
	indexName = chattypes.TableServerGroupsIndexMainAddr

	rpRows, err := chattypes.TableList(r.GetLocalDB(), chattypes.TableServerGroupsName, indexName, rpData, chattypes.MaxServerGroupsNumb, chattypes.ListDESC)
	if err != nil {
		elog.Error("tableList err", "err", err)
		return nil, err
	}
	groups := make([]*chattypes.ServerGroup, 0)
	for _, v := range rpRows {
		groups = append(groups, v.Data.(*chattypes.ServerGroup))
	}
	return groups, nil
}

func (r *Chat) getFriend(userAddr, friendAddr string) (*chattypes.Friend, error) {
	//query
	fd := &chattypes.Friend{
		MainAddress:   userAddr,
		FriendAddress: friendAddr,
	}
	cur := &chattypes.FriendRow{Friend: fd}
	primary, err := cur.Get(chattypes.TableFriendIndexMainAddrFriendAddr)
	if err != nil {
		return nil, err
	}
	fTable := chattypes.NewFriendTable(r.GetLocalDB())
	v, err := fTable.GetData(primary)
	if err != nil && err != types.ErrNotFound {
		elog.Error("query friend data failed", "err", err)
		return nil, err
	}
	if v == nil || v.Data == nil {
		return nil, nil
	}
	return v.Data.(*chattypes.Friend), nil
}

//Query_GetFriends 查询好友列表
func (r *Chat) Query_GetFriends(in *chattypes.ReqGetFriends) (types.Message, error) {
	rpData := &chattypes.Friend{
		MainAddress:   in.MainAddress,
		FriendAddress: in.Index,
	}
	var indexName string
	indexName = chattypes.TableFriendIndexMainAddr

	//TODO 查询配置
	queryTy := chattypes.Private
	//TODO debug模式
	if in.Sign == nil {
		queryTy = chattypes.Public
	}
	switch queryTy {
	case chattypes.Private:
		if in.Sign == nil {
			return nil, chattypes.ErrParams
		}

		sig, err := util.HexDecode(in.Sign.Signature)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		pubKey, err := util.HexDecode(in.Sign.PublicKey)

		sum, err := util.GetSummary(in)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		msg := sha256.Sum256(sum)
		if !util.Secp256k1Verify(msg[:], sig, pubKey) {
			return nil, chattypes.ErrSignErr
		}
		//检查时间是否过期
		if util.CheckTimeOut(in.Time) {
			return nil, chattypes.ErrQueryTimeOut
		}
		//检查地址是否正确
		if in.MainAddress != address.PubKeyToAddr(pubKey) {
			return nil, chattypes.ErrLackPermissions
		}
	case chattypes.Public:
	case chattypes.Protect:
		return nil, chattypes.ErrTypeNotExist
	default:
		return nil, chattypes.ErrTypeNotExist
	}

	rpRows, err := chattypes.TableList(r.GetLocalDB(), chattypes.TableFriendName, indexName, rpData, in.Count, chattypes.ListDESC)
	if err != nil {
		elog.Error("tableList err", "err", err)
		return nil, err
	}
	reply := &chattypes.ReplyGetFriends{}
	reply.Friends = make([]*chattypes.Friend, 0)
	for _, v := range rpRows {
		reply.Friends = append(reply.Friends, v.Data.(*chattypes.Friend))
	}
	return reply, nil
}

//Query_GetBlackList 查询黑名单列表
func (r *Chat) Query_GetBlackList(in *chattypes.ReqGetBlackList) (types.Message, error) {
	rpData := &chattypes.Black{
		MainAddress:   in.MainAddress,
		TargetAddress: in.Index,
	}
	var indexName string
	indexName = chattypes.TableBlackIndexMainAddr

	//TODO 查询配置
	queryTy := chattypes.Private
	//TODO debug模式
	if in.Sign == nil {
		queryTy = chattypes.Public
	}
	switch queryTy {
	case chattypes.Private:
		if in.Sign == nil {
			return nil, chattypes.ErrParams
		}

		sig, err := util.HexDecode(in.Sign.Signature)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		pubKey, err := util.HexDecode(in.Sign.PublicKey)

		sum, err := util.GetSummary(in)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		msg := sha256.Sum256(sum)
		if !util.Secp256k1Verify(msg[:], sig, pubKey) {
			return nil, chattypes.ErrSignErr
		}
		//检查时间是否过期
		if util.CheckTimeOut(in.Time) {
			return nil, chattypes.ErrQueryTimeOut
		}
		//检查地址是否正确
		if in.MainAddress != address.PubKeyToAddr(pubKey) {
			return nil, chattypes.ErrLackPermissions
		}
	case chattypes.Public:
	case chattypes.Protect:
		return nil, chattypes.ErrTypeNotExist
	default:
		return nil, chattypes.ErrTypeNotExist
	}

	rpRows, err := chattypes.TableList(r.GetLocalDB(), chattypes.TableBlackName, indexName, rpData, in.Count, chattypes.ListDESC)
	if err != nil {
		elog.Error("tableList err", "err", err)
		return nil, err
	}
	reply := &chattypes.ReplyGetBlackList{}
	reply.List = make([]*chattypes.Black, 0)
	for _, v := range rpRows {
		reply.List = append(reply.List, v.Data.(*chattypes.Black))
	}
	return reply, nil
}

//Query_GetUser 查询用户信息
func (r *Chat) Query_GetUser(in *chattypes.ReqGetUser) (types.Message, error) {
	//检查地址
	if err := address.CheckAddress(in.TargetAddress); err != nil {
		return nil, err
	}
	rpData := &chattypes.Field{
		MainAddress: in.TargetAddress,
		Name:        in.Index,
	}
	var indexName string
	indexName = chattypes.TableUserIndexMainAddr

	//TODO 查询配置
	queryTy := chattypes.Private
	//TODO debug模式
	if in.Sign == nil {
		queryTy = chattypes.Public
	}
	switch queryTy {
	case chattypes.Private:
		if in.Sign == nil {
			return nil, chattypes.ErrParams
		}

		sig, err := util.HexDecode(in.Sign.Signature)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		pubKey, err := util.HexDecode(in.Sign.PublicKey)

		sum, err := util.GetSummary(in)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		msg := sha256.Sum256(sum)
		if !util.Secp256k1Verify(msg[:], sig, pubKey) {
			return nil, chattypes.ErrSignErr
		}
		//检查时间是否过期
		if util.CheckTimeOut(in.Time) {
			return nil, chattypes.ErrQueryTimeOut
		}
		//检查地址是否正确
		if in.MainAddress != address.PubKeyToAddr(pubKey) {
			return nil, chattypes.ErrLackPermissions
		}
	case chattypes.Public:
	case chattypes.Protect:
		return nil, chattypes.ErrTypeNotExist
	default:
		return nil, chattypes.ErrTypeNotExist
	}

	isFriend := false
	isSelf := false
	fGroups := []string{}
	//检查是否是自己或者是好友
	if in.TargetAddress == in.MainAddress {
		isSelf = true
	} else {
		f, err := r.getFriend(in.TargetAddress, in.MainAddress)
		if err != nil {
			elog.Error("getFriend err", "err", err)
			return nil, err
		}
		if f != nil && f.FriendAddress == in.MainAddress && f.MainAddress == in.TargetAddress {
			isFriend = true
			if len(f.Groups) > 0 {
				fGroups = f.Groups
			}
		}
	}

	rpRows, err := chattypes.TableList(r.GetLocalDB(), chattypes.TableUserName, indexName, rpData, in.Count, chattypes.ListDESC)
	if err != nil {
		elog.Error("tableList err", "err", err)
		return nil, err
	}
	reply := &chattypes.ReplyGetUser{}
	//chat servers
	reply.ChatServers = make([]*chattypes.ReplyServerInfo, 0)
	//获取所有分组
	groups, err := r.getServerGroups(in.TargetAddress, "")
	sort.Sort(ByIdGroup(groups))
	idxGroups := indexServerGroups(groups)
	if isSelf {
		for _, g := range groups {
			reply.ChatServers = append(reply.ChatServers, &chattypes.ReplyServerInfo{
				Id:      g.Id,
				Name:    g.Name,
				Address: g.Value,
			})
		}
	} else if isFriend {
		for _, g := range fGroups {
			item := idxGroups[g]
			if item != nil {
				reply.ChatServers = append(reply.ChatServers, &chattypes.ReplyServerInfo{
					Id:      item.Id,
					Name:    item.Name,
					Address: item.Value,
				})
			}
		}
	}
	if len(reply.ChatServers) == 0 {
		//添加默认分组
		if len(groups) > 0 {
			reply.ChatServers = append(reply.ChatServers, &chattypes.ReplyServerInfo{
				Id:      groups[0].Id,
				Name:    groups[0].Name,
				Address: groups[0].Value,
			})
		}
	}

	//fields
	reply.Fields = make([]*chattypes.ReplyGetField, 0)
	for _, v := range rpRows {
		switch v.Data.(*chattypes.Field).Level {
		case chattypes.LvlPublic:
		case chattypes.LvlProtect:
			if !isSelf && !isFriend {
				continue
			}
		case chattypes.LvlPrivate:
			if !isSelf {
				continue
			}
		default:
			continue
		}
		val := v.Data.(*chattypes.Field)
		reply.Fields = append(reply.Fields, &chattypes.ReplyGetField{
			Name:  val.Name,
			Value: val.Value,
			Level: val.Level,
		})
	}

	//groups
	//查询好友分组列表
	if !isSelf {
		f, err := r.getFriend(in.MainAddress, in.TargetAddress)
		if err != nil {
			elog.Error("getFriend err", "err", err)
			return nil, err
		}
		if f != nil {
			reply.Groups = f.Groups
		}
	}
	return reply, nil
}

//Query_GetServerGroup 查询服务分组列表
func (r *Chat) Query_GetServerGroup(in *chattypes.ReqGetServerGroup) (types.Message, error) {
	rpData := &chattypes.ServerGroup{
		MainAddress: in.MainAddress,
		Id:          in.Index,
	}
	var indexName string
	indexName = chattypes.TableServerGroupsIndexMainAddr

	//TODO 查询配置
	queryTy := chattypes.Private
	//TODO debug模式
	if in.Sign == nil {
		queryTy = chattypes.Public
	}
	switch queryTy {
	case chattypes.Private:
		if in.Sign == nil {
			return nil, chattypes.ErrParams
		}

		sig, err := util.HexDecode(in.Sign.Signature)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		pubKey, err := util.HexDecode(in.Sign.PublicKey)

		sum, err := util.GetSummary(in)
		if err != nil {
			return nil, chattypes.ErrParams
		}
		msg := sha256.Sum256(sum)
		if !util.Secp256k1Verify(msg[:], sig, pubKey) {
			return nil, chattypes.ErrSignErr
		}
		//检查时间是否过期
		if util.CheckTimeOut(in.Time) {
			return nil, chattypes.ErrQueryTimeOut
		}
		//检查地址是否正确
		if in.MainAddress != address.PubKeyToAddr(pubKey) {
			return nil, chattypes.ErrLackPermissions
		}
	case chattypes.Public:
	case chattypes.Protect:
		return nil, chattypes.ErrTypeNotExist
	default:
		return nil, chattypes.ErrTypeNotExist
	}

	rpRows, err := chattypes.TableList(r.GetLocalDB(), chattypes.TableServerGroupsName, indexName, rpData, in.Count, chattypes.ListDESC)
	if err != nil {
		elog.Error("tableList err", "err", err)
		return nil, err
	}
	groups := make([]*chattypes.ServerGroup, 0)
	for _, v := range rpRows {
		if val := v.Data.(*chattypes.ServerGroup); val != nil {
			groups = append(groups, val)
		} else {
			elog.Error("query server groups val empty", "val", val)
		}
	}
	sort.Sort(ByIdGroup(groups))
	reply := &chattypes.ReplyGetServerGroups{}
	reply.Groups = make([]*chattypes.ReplyGetServerGroup, len(groups))
	for i, group := range groups {
		reply.Groups[i] = &chattypes.ReplyGetServerGroup{
			Id:    group.Id,
			Name:  group.Name,
			Value: group.Value,
		}
	}
	return reply, nil
}
