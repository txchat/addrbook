package types

import "time"

const (
	FriendAppend = 1
	FriendDelete = 2

	BlackAppend = 1
	BlackDelete = 2

	UserAppend = 1
	UserDelete = 2

	GroupAppend = 1
	GroupDelete = 2
	GroupEdit   = 3
)

const MaxServerGroupsNumb = 20
const DefaultServerGroupName = "default"
const DefaultServerGroupId = "1"
const MaxGroupsLength = 20

const (
	Private = 1
	Protect = 2
	Public  = 3

	blacklist = 1
	Whitelist = 2

	QueryTimeOut = time.Second * 120
)

const (
	LvlPublic  = "public"
	LvlProtect = "protect"
	LvlPrivate = "private"
)

const (
	//ListDESC 表示记录降序排列
	ListDESC = int32(0)

	//ListASC 表示记录升序排列
	ListASC = int32(1)

	FuncNameGetFriends     = "GetFriends"
	FuncNameGetBlackList   = "GetBlackList"
	FuncNameGetUser        = "GetUser"
	FuncNameGetServerGroup = "GetServerGroup"
)

// action类型id和name，这些常量可以自定义修改
const (
	TyUpdateFriendsAction = iota + 100
	TyUpdateBlackListAction
	TyUpdateUserAction
	TyUpdateServerGroupAction

	NameUpdateFriendsAction     = "Update"
	NameUpdateBlackListAction   = "Black"
	NameUpdateUserAction        = "UpdateUser"
	NameUpdateServerGroupAction = "UpdateServerGroup"
)

// log类型id值
const (
	TyUnknownLog = iota + 100
	TyLogUpdateFriends
	TyLogUpdateBlackList
	TyLogUpdateUser
	TyLogUpdateServerGroup
)
