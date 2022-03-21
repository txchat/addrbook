package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	atypes "github.com/txchat/addrbook/types"
	"google.golang.org/grpc"
)

var (
	conns sync.Map
)

type GRPCCli struct {
	client types.Chain33Client
}

func NewGRPCCli(grpcAddr string) *GRPCCli {
	client, err := GetClient(grpcAddr)
	if err != nil {
		return nil
	}
	return &GRPCCli{client: client}
}

func (c *GRPCCli) Send(tx *types.Transaction) ([]*types.ReceiptLog, error) {
	logs, err := c.sendAndWaitReceipt(tx)
	if err != nil {
		return nil, ParseError(err)
	}
	for _, l := range logs {
		if l.Ty == types.TyLogErr {
			return nil, errors.New(string(l.Log))
		}
	}
	return logs, nil
}

func (c *GRPCCli) Query(fn string, msg proto.Message) ([]byte, error) {
	ss := strings.Split(fn, ".")
	var in types.ChainExecutor
	if len(ss) == 2 {
		in.Driver = ss[0]
		in.FuncName = ss[1]
	} else {
		in.Driver = cfg.GetTitle() + atypes.ChatX
		in.FuncName = fn
	}
	in.Param = types.Encode(msg)

	r, err := c.client.QueryChain(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	if !r.IsOk {
		return nil, errors.New(string(r.Msg))
	}
	return r.Msg, nil
}

// 发送交易并等待执行结果
// 如果交易非法，返回错误信息
// 如果交易执行成功，返回 交易哈希、回报
func (c *GRPCCli) sendAndWaitReceipt(tx *types.Transaction) (logs []*types.ReceiptLog, err error) {
	hash, err := c.SendAsync(tx)
	if err != nil {
		return nil, err
	}

	timeOut := time.NewTimer(time.Second * 5)
	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			d, _ := c.client.QueryTransaction(context.Background(), &types.ReqHash{Hash: hash})
			if d != nil {
				//_, err := etypes.ParseReceiptLogs(d.Receipt.Logs)
				return d.Receipt.Logs, err
			}
		case <-timeOut.C:
			return nil, fmt.Errorf("timeout query receipt, tx: %x", hash)
		}
	}
}

func (c *GRPCCli) SendAsync(tx *types.Transaction) (txHash []byte, err error) {
	r, err := c.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, err
	}
	if !r.IsOk {
		return nil, errors.New(string(r.Msg))
	}
	return r.Msg, nil
}

func (c *GRPCCli) GetLastHeader() (int64, error) {
	header, err := c.client.GetLastHeader(context.Background(), &types.ReqNil{})
	if err != nil {
		return 0, err
	}
	return header.Height, nil
}

func ParseError(err error) error {
	// rpc error: code = Unknown desc = ErrNotBank
	str := err.Error()
	sep := "desc = "
	i := strings.Index(str, sep)
	if i != -1 {
		return errors.New(str[i+len(sep):])
	}
	return err
}

func GetClient(target string) (types.Chain33Client, error) {
	val, ok := conns.Load(target)
	if !ok {
		conn, err := newGrpcConn(target)
		if err != nil {
			return nil, err
		}
		client := types.NewChain33Client(conn)
		conns.Store(target, client)
		return client, nil
	} else {
		return val.(types.Chain33Client), nil
	}
}

func newGrpcConn(target string) (*grpc.ClientConn, error) {
	return grpc.Dial(target, grpc.WithInsecure())
}
