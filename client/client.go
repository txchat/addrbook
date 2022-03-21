package client

import (
	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
)

type Cli interface {
	Send(tx *types.Transaction) ([]*types.ReceiptLog, error)    // sync
	SendAsync(tx *types.Transaction) (txHash []byte, err error) // async
	Query(fn string, msg proto.Message) ([]byte, error)
	GetLastHeader() (int64, error)
}
