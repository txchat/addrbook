package test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"github.com/stretchr/testify/assert"
	excli "github.com/txchat/addrbook/client"
	atypes "github.com/txchat/addrbook/types"
	"github.com/txchat/addrbook/util"
)

var client *excli.AddrbookClient
var priv crypto.PrivKey

func init() {
	priv = genKeyFromStr("29feed2cf0a9b09f0ec80e6def12f5010af6cd9dcf938215dd0c6fe44a242876") //1KnirEo67cuko2XxqCBo7VR7u3AeNLD9TS
	//cli := excli.NewExecCli("chain33.toml") // exec client
	cli := excli.NewGRPCCli(":8902") // grpc client
	client = excli.NewAddrbookCient(cli)
}

func TestAddFriend(t *testing.T) {
	addr := address.PubKeyToAddr(genKey().PubKey().Bytes())
	req := []*atypes.ReqUpdateFriend{
		{
			FriendAddress: addr,
			Type:          atypes.FriendAppend,
		},
	}
	logs, err := client.UpdateFriends(priv, req)
	assert.Nil(t, err)
	t.Log(logs)
}

func TestGetFriends(t *testing.T) {
	q := &atypes.ReqGetFriends{
		MainAddress: address.PubKeyToAddr(priv.PubKey().Bytes()),
		Count:       10,
		Time:        time.Now().UnixNano() / 1e6,
	}
	sum, err := util.GetSummary(q)
	assert.Nil(t, err)
	signature := priv.Sign(sum)
	q.Sign = &atypes.Sign{
		PublicKey: priv.PubKey().KeyString(),
		Signature: hex.EncodeToString(signature.Bytes()),
	}

	resp, err := client.QueryFriends(q)
	assert.Nil(t, err)
	t.Log(resp)
}

func genKey() crypto.PrivKey {
	c, _ := crypto.New(types.GetSignName("", types.SECP256K1))
	priv, _ := c.GenKey()
	return priv
}

func genKeyFromStr(key string) crypto.PrivKey {
	c, _ := crypto.New(types.GetSignName("", types.SECP256K1))
	bytes, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	priv, _ := c.PrivKeyFromBytes(bytes)
	return priv
}
