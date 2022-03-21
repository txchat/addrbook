package types

import (
	"fmt"

	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
)

const (
	TableChatPreFix                    = "LODB-chat"
	TableFriendName                    = "friend"
	TableFriendPrimary                 = "mainAddr_friAddr"
	TableFriendIndexMainAddr           = "mainAddr"
	TableFriendIndexMainAddrFriendAddr = "mainAddr_friAddr"

	TableBlackName                    = "black"
	TableBlackPrimary                 = "mainAddr_tagAddr"
	TableBlackIndexMainAddr           = "mainAddr"
	TableBlackIndexMainAddrTargetAddr = "mainAddr_tagAddr"

	TableUserName                   = "user"
	TableUserPrimary                = "mainAddr_fieldName"
	TableUserIndexMainAddr          = "mainAddr"
	TableUserIndexMainAddrFieldName = "mainAddr_fieldName"

	TableServerGroupsName                  = "server"
	TableServerGroupsPrimary               = "mainAddr_serverId"
	TableServerGroupsIndexMainAddr         = "mainAddr"
	TableServerGroupsIndexMainAddrServerId = "mainAddr_serverId"
)

var opt_friend_table = &table.Option{
	Prefix:  TableChatPreFix,
	Name:    TableFriendName,
	Primary: TableFriendPrimary,
	Index:   []string{TableFriendIndexMainAddr, TableFriendIndexMainAddrFriendAddr},
}

var opt_black_table = &table.Option{
	Prefix:  TableChatPreFix,
	Name:    TableBlackName,
	Primary: TableBlackPrimary,
	Index:   []string{TableBlackIndexMainAddr, TableBlackIndexMainAddrTargetAddr},
}

var opt_user_table = &table.Option{
	Prefix:  TableChatPreFix,
	Name:    TableUserName,
	Primary: TableUserPrimary,
	Index:   []string{TableUserIndexMainAddr, TableUserIndexMainAddrFieldName},
}

var opt_server_table = &table.Option{
	Prefix:  TableChatPreFix,
	Name:    TableServerGroupsName,
	Primary: TableServerGroupsPrimary,
	Index:   []string{TableServerGroupsIndexMainAddr, TableServerGroupsIndexMainAddrServerId},
}

func TableList(dbm db.KVDB, tableName, indexName string, data interface{}, count, direction int32) ([]*table.Row, error) {
	var query *table.Query
	var cur table.RowMeta
	var primaryName string
	tlog.Info("data Info", "tableName", tableName, "data", data)
	switch tableName {
	case TableFriendName:
		query = NewFriendTable(dbm).GetQuery(dbm)
		cur = &FriendRow{Friend: data.(*Friend)}
		primaryName = TableFriendIndexMainAddrFriendAddr
	case TableBlackName:
		query = NewBlackTable(dbm).GetQuery(dbm)
		cur = &BlackRow{Black: data.(*Black)}
		primaryName = TableBlackIndexMainAddrTargetAddr
	case TableUserName:
		query = NewUserTable(dbm).GetQuery(dbm)
		cur = &UserRow{Field: data.(*Field)}
		primaryName = TableUserIndexMainAddrFieldName
	case TableServerGroupsName:
		query = NewServerGroupsTable(dbm).GetQuery(dbm)
		cur = &ServerGroupsRow{ServerGroup: data.(*ServerGroup)}
		primaryName = TableServerGroupsIndexMainAddrServerId
	}

	var prefix []byte
	var err error
	if primaryName != "auto" {
		prefix, err = cur.Get(indexName)
		if err != nil && indexName != "" {
			tlog.Error("index get err", "indexName", indexName)
			return nil, err
		}
	}

	primary, err := cur.Get(primaryName)
	if err != nil {
		tlog.Error("primary get err", "primaryName", primaryName)
		return nil, err
	}
	if indexName == primaryName {
		primary = nil
	}
	tlog.Info("primary is", "primary name", primaryName, "primary", string(primary))
	tlog.Info("prefix is", "index name", indexName, "prefix", string(prefix))
	rows, err := query.ListIndex(indexName, prefix, primary, count, direction)
	if err != nil && err != types.ErrNotFound {
		tlog.Error("query List failed", "tableName", tableName, "indexName", indexName, "primaryName", primaryName, "prefix", string(prefix), "key", string(primary), "param", data, "err", err)
		return nil, err
	}
	tlog.Info("row data is", "rows", rows, "err", err)
	return rows, nil
}

//NewFriendTable 新建表
func NewFriendTable(kvdb db.KV) *table.Table {
	rowmeta := NewFriendRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_friend_table)
	if err != nil {
		tlog.Error("NewFriendTable err", "err", err)
		//panic(err)
	}
	return table
}

//FriendRow table meta 结构
type FriendRow struct {
	*Friend
}

//NewFriendRow 新建一个meta 结构
func NewFriendRow() *FriendRow {
	return &FriendRow{Friend: &Friend{}}
}

//CreateRow 新建数据行(注意index 数据一定也要保存到数据中,不能就保存eventid)
func (tx *FriendRow) CreateRow() *table.Row {
	return &table.Row{Data: &Friend{}}
}

//SetPayload 设置数据
func (tx *FriendRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*Friend); ok {
		tx.Friend = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *FriendRow) Get(key string) ([]byte, error) {
	switch key {
	/*	case TableFriendPrimary:
		return []byte(fmt.Sprintf("%020d", tx.Id)), nil*/
	case TableFriendIndexMainAddr:
		return []byte(tx.MainAddress), nil
	case TableFriendIndexMainAddrFriendAddr:
		if tx.FriendAddress == "" {
			return nil, nil
		}
		return []byte(fmt.Sprintf("%s_%s", tx.MainAddress, tx.FriendAddress)), nil
	}
	return nil, types.ErrNotFound
}

//黑名单
//NewFriendTable 新建表
func NewBlackTable(kvdb db.KV) *table.Table {
	rowmeta := NewBlackRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_black_table)
	if err != nil {
		tlog.Error("NewBlackTable err", "err", err)
		//panic(err)
	}
	return table
}

//FriendRow table meta 结构
type BlackRow struct {
	*Black
}

//NewFriendRow 新建一个meta 结构
func NewBlackRow() *BlackRow {
	return &BlackRow{Black: &Black{}}
}

//CreateRow 新建数据行(注意index 数据一定也要保存到数据中,不能就保存eventid)
func (tx *BlackRow) CreateRow() *table.Row {
	return &table.Row{Data: &Black{}}
}

//SetPayload 设置数据
func (tx *BlackRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*Black); ok {
		tx.Black = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *BlackRow) Get(key string) ([]byte, error) {
	switch key {
	/*	case TableFriendPrimary:
		return []byte(fmt.Sprintf("%020d", tx.Id)), nil*/
	case TableBlackIndexMainAddr:
		return []byte(tx.MainAddress), nil
	case TableBlackIndexMainAddrTargetAddr:
		if tx.TargetAddress == "" {
			return nil, nil
		}
		return []byte(fmt.Sprintf("%s_%s", tx.MainAddress, tx.TargetAddress)), nil
	}
	return nil, types.ErrNotFound
}

//NewUserTable 新建表
func NewUserTable(kvdb db.KV) *table.Table {
	rowmeta := NewUserRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_user_table)
	if err != nil {
		tlog.Error("NewUserTable err", "err", err)
		//panic(err)
	}
	return table
}

//FriendRow table meta 结构
type UserRow struct {
	*Field
}

//NewFriendRow 新建一个meta 结构
func NewUserRow() *UserRow {
	return &UserRow{Field: &Field{}}
}

//CreateRow 新建数据行(注意index 数据一定也要保存到数据中,不能就保存eventid)
func (tx *UserRow) CreateRow() *table.Row {
	return &table.Row{Data: &Field{}}
}

//SetPayload 设置数据
func (tx *UserRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*Field); ok {
		tx.Field = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *UserRow) Get(key string) ([]byte, error) {
	switch key {
	/*	case TableFriendPrimary:
		return []byte(fmt.Sprintf("%020d", tx.Id)), nil*/
	case TableUserIndexMainAddr:
		return []byte(tx.MainAddress), nil
	case TableUserIndexMainAddrFieldName:
		if tx.Name == "" {
			return nil, nil
		}
		return []byte(fmt.Sprintf("%s_%s", tx.MainAddress, tx.Name)), nil
	}
	return nil, types.ErrNotFound
}

//server groups
//NewServerGroupsTable 新建表
func NewServerGroupsTable(kvdb db.KV) *table.Table {
	rowmeta := NewServerGroupsRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_server_table)
	if err != nil {
		tlog.Error("NewServerGroupsTable err", "err", err)
		//panic(err)
	}
	return table
}

//FriendRow table meta 结构
type ServerGroupsRow struct {
	*ServerGroup
}

//NewFriendRow 新建一个meta 结构
func NewServerGroupsRow() *ServerGroupsRow {
	return &ServerGroupsRow{ServerGroup: &ServerGroup{}}
}

//CreateRow 新建数据行(注意index 数据一定也要保存到数据中,不能就保存eventid)
func (tx *ServerGroupsRow) CreateRow() *table.Row {
	return &table.Row{Data: &ServerGroup{}}
}

//SetPayload 设置数据
func (tx *ServerGroupsRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ServerGroup); ok {
		tx.ServerGroup = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *ServerGroupsRow) Get(key string) ([]byte, error) {
	switch key {
	/*	case TableFriendPrimary:
		return []byte(fmt.Sprintf("%020d", tx.Id)), nil*/
	case TableServerGroupsIndexMainAddr:
		return []byte(tx.MainAddress), nil
	case TableServerGroupsIndexMainAddrServerId:
		if tx.Id == "" {
			return nil, nil
		}
		return []byte(fmt.Sprintf("%s_%s", tx.MainAddress, tx.Id)), nil
	}
	return nil, types.ErrNotFound
}
