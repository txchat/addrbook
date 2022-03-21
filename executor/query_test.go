package executor

import (
	"encoding/hex"
	"sort"
	"testing"

	"github.com/33cn/chain33/types"
	"github.com/haltingstate/secp256k1-go"
	chattypes "github.com/txchat/addrbook/types"
	"github.com/txchat/addrbook/util"
)

func Test_Query_GetFriends(t *testing.T) {
	c := NewChat().(*Chat)
	in := chattypes.ReqGetFriends{
		MainAddress: "12VvgbuFKwuPMwVD6DjmyTBmWeL7HNvczF",
		Count:       20,
		Index:       "",
		Time:        1581499625,
	}
	//模拟客户端签名
	sum, err := util.GetSummary(&in)
	if err != nil {
		t.Error(err)
		return
	}
	priKeyStr := "2be91095f403d219060d257d910b5eada78a17b5a525897c387be4de1993dfb7"
	privateKey, err := hex.DecodeString(priKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	signature := secp256k1.Sign(sum, privateKey)
	t.Log("signature", hex.EncodeToString(signature))

	/*in.Sign = &chattypes.Sign{
		PublicKey: "02be5910fb49cb6e36f939080afd583be3676db96e7ed5177df6a401e1419f7375",
		Signature: hex.EncodeToString(signature),
	}*/

	in.Sign = &chattypes.Sign{
		PublicKey: "",
		Signature: "",
	}

	m, err := c.Query_GetFriends(&in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success", m)
}

func Test_PluginSignVerify(t *testing.T) {
	msg := []byte("123")
	//RedPacketX := "chat"
	pubKeyStr := "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa"
	//pubKeyStr := "02be5910fb49cb6e36f939080afd583be3676db96e7ed5177df6a401e1419f7375"
	priKeyStr := "16d7ccefb9d1abc06e2e6d6567105a1f24f2e069569360ada51837d45ec21b71"

	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := hex.DecodeString(priKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	signature := secp256k1.Sign(msg, privateKey)
	t.Log("signature", hex.EncodeToString(signature))

	if 1 == secp256k1.VerifySignature(msg, signature, publicKey) {
		t.Log("success")
	}

	data := msg
	sign := &types.Signature{
		Ty:        1,
		Pubkey:    ([]byte)(publicKey),
		Signature: ([]byte)(signature),
	}
	//rlt := types.CheckSign(data, types.ExecName(RedPacketX), sign)
	rlt := types.CheckSign(data, GetName(), sign, 0)
	t.Log("结果", rlt)
}

func Test_PluginVerify(t *testing.T) {
	msg := []byte("1613977615909*yes")
	//RedPacketX := "chat"
	pubKeyStr := "0310c3474b518dd9ed05a36372705a35d29b854d06fd39f7d026d518b6865e1f7b"

	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	signature, err := hex.DecodeString("00967bc443f8904378316319957f4791104f4eba5e3ea2165ebccc9fe2574313600ecd1431641f236ab66bfae4594f2327ba4aa54bba1690fcd9b5406fe6a95401")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature", hex.EncodeToString(signature))

	if 1 == secp256k1.VerifySignature(msg, signature, publicKey) {
		t.Log("success")
	}

	data := msg
	sign := &types.Signature{
		Ty:        1,
		Pubkey:    ([]byte)(publicKey),
		Signature: ([]byte)(signature),
	}
	//rlt := types.CheckSign(data, types.ExecName(RedPacketX), sign)
	rlt := types.CheckSign(data, GetName(), sign, 0)
	t.Log("结果", rlt)
}

func Test_SortGroupById(t *testing.T) {
	groups := []*chattypes.ServerGroup{
		{
			Id:          "3",
			Name:        "test1",
			Value:       "",
			MainAddress: "",
		}, {
			Id:          "2",
			Name:        "test2",
			Value:       "",
			MainAddress: "",
		}, {
			Id:          "10",
			Name:        "test3",
			Value:       "",
			MainAddress: "",
		},
	}

	sort.Sort(ByIdGroup(groups))
	t.Log("after sorted:")
	for i, g := range groups {
		t.Logf("num:%v:%v\n", i, g.Id)
	}
}
