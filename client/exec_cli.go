package client

import (
	"io/ioutil"
	"time"

	"github.com/33cn/chain33/client"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/queue"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/addrbook/executor"
	atypes "github.com/txchat/addrbook/types"
)

type ExecCli struct {
	ldb        db.KVDB
	sdb        db.DB
	height     int64
	blockTime  int64
	difficulty uint64
	q          queue.Queue
	cfg        *types.Chain33Config
	execAddr   string
}

func NewExecCli(cfgFile string) *ExecCli {
	sdb, _ := db.NewGoLevelDB("state", ".", 128)
	leveldb, err := db.NewGoLevelDB("local", ".", 128)
	if err != nil {
		panic(err)
	}
	ldb := db.NewKVDB(leveldb)

	data, _ := ioutil.ReadFile(cfgFile)
	cfg := types.NewChain33Config(string(data))
	executor.Init(atypes.ChatX, cfg, nil)
	execAddr := address.ExecAddress(atypes.ChatX)
	q := queue.New("channel")
	q.SetConfig(cfg)

	return &ExecCli{
		ldb:        ldb,
		sdb:        sdb,
		height:     0,
		blockTime:  time.Now().Unix(),
		difficulty: 1539918074,
		q:          q,
		cfg:        cfg,
		execAddr:   execAddr,
	}
}

func (c *ExecCli) Send(tx *types.Transaction) ([]*types.ReceiptLog, error) {
	var err error
	exec := executor.NewChat()
	api, _ := client.New(c.q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(c.sdb)
	exec.SetLocalDB(c.ldb)
	exec.SetEnv(c.height, c.blockTime, c.difficulty)
	if err := exec.CheckTx(tx, int(1)); err != nil {
		return nil, err
	}

	c.height++
	c.blockTime += 10
	c.difficulty++
	receipt, err := exec.Exec(tx, int(1))
	if err != nil {
		return nil, err
	}

	receiptDate := &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
	set, err := exec.ExecLocal(tx, receiptDate, int(1))
	if err != nil {
		return nil, err
	}

	util.SaveKVList(c.sdb, receipt.KV)
	for _, kv := range set.KV {
		c.ldb.Set(kv.Key, kv.Value)
	}

	//_, err = atypes.ParseReceiptLogs(receipt.Logs)
	return receipt.Logs, err
}

func (c *ExecCli) Query(fn string, msg proto.Message) ([]byte, error) {
	api, _ := client.New(c.q.Client(), nil)
	exec := executor.NewChat()
	exec.SetAPI(api)
	exec.SetStateDB(c.sdb)
	exec.SetLocalDB(c.ldb)
	exec.SetEnv(c.height, c.blockTime, c.difficulty)
	r, err := exec.Query(fn, types.Encode(msg))
	if err != nil {
		return nil, err
	}
	return types.Encode(r), nil
}

func (c *ExecCli) SendAsync(tx *types.Transaction) (txHash []byte, err error) {
	panic("not implement")
}

func (c *ExecCli) GetLastHeader() (int64, error) {
	return c.height, nil
}
