package util

import (
	"encoding/hex"
	"testing"

	"github.com/haltingstate/secp256k1-go"
	chattypes "github.com/txchat/addrbook/types"
)

func Test_Sign(t *testing.T) {
	in := chattypes.ReqGetFriends{
		MainAddress: "12mXWbj3qzB7oiSq8iCN9c3MLbuJLi3mrc",
		Count:       20,
		Index:       "=",
		Time:        1581587923000,
	}
	//模拟客户端签名
	sum, err := GetSummary(&in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(sum))
	priKeyStr := "2be91095f403d219060d257d910b5eada78a17b5a525897c387be4de1993dfb7"
	privateKey, err := hex.DecodeString(priKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	signature := secp256k1.Sign(sum, privateKey)
	t.Log("signature", hex.EncodeToString(signature))
}

func Test_Verify2(t *testing.T) {
	in := chattypes.ReqGetFriends{
		MainAddress: "18uQM7bvNHrLD7ZnWEWfuCpWKKiF3ERqgm",
		Count:       40,
		Index:       "",
		Time:        1599104451440,
	}
	//获取摘要
	sum, err := GetSummary(&in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(sum))

	pubKeyStr := "02d7a5f5f9434785d2a8fe4fa220db70b23a17b1fc2f7a97913ee03d8586f826b2"
	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	sig, err := hex.DecodeString("2a396299f2ef7bfb800579d221f9d30732b2304f460ba72a841c44841c502b5038856c574e163f4f1aaf64af0115a51aca63768718813804f4a6cf1352b465ff01")
	if err != nil {
		t.Error(err)
		return
	}
	ret := secp256k1.VerifySignature(sum, sig, publicKey)
	t.Log("verify", ret)
}

func Test_VerifyChain33(t *testing.T) {
	in := chattypes.ReqGetFriends{
		MainAddress: "18uQM7bvNHrLD7ZnWEWfuCpWKKiF3ERqgm",
		Count:       40,
		Index:       "",
		Time:        1599104451440,
	}
	//获取摘要
	sum, err := GetSummary(&in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(sum))

	pubKeyStr := "02d7a5f5f9434785d2a8fe4fa220db70b23a17b1fc2f7a97913ee03d8586f826b2"
	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	sig, err := hex.DecodeString("2a396299f2ef7bfb800579d221f9d30732b2304f460ba72a841c44841c502b5038856c574e163f4f1aaf64af0115a51aca63768718813804f4a6cf1352b465ff01")
	if err != nil {
		t.Error(err)
		return
	}
	ret := Secp256k1VerifyChain33(sum, sig, publicKey)
	t.Log("verify", ret)
}

func Test_Verify(t *testing.T) {
	sum := []byte("count%3D40%26index%3D%26mainAddress%3D1JoFzozbxvst22c2K7MBYwQGjCaMZbC5Qm%26time%3D1581594811061")

	pubKeyStr := "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa"
	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	sig, err := hex.DecodeString("4d333dca5f2703cb97610b9ec8749e0edfda41f53773c29c372ae69bd91b9b74410b5d4fda79d26962bf057da0d40c587c5f45e5cdb69a4a59f9fbb8d2f9d05800")
	if err != nil {
		t.Error(err)
		return
	}
	ret := secp256k1.VerifySignature(sum, sig, publicKey)
	t.Log("verify", ret)
}

func Test_VerifyChain33Two(t *testing.T) {
	sum := []byte("1624503307409*aaaaaaaaaaaaaaaaaa")

	pubKeyStr := "030681fa07b6b5172ba4f868364d414308c092feccddda59178132a83025ca5c58"
	publicKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		t.Error(err)
		return
	}

	sig, err := hex.DecodeString("93181cf41ebc1b65ff3fe464a2e088711fb376c0ac27802311b183b2995fe3d57bb9908cf4dd3de7ce3fffbdeee9ce0e0526daa4a4b245432380cb6c90fdafbe01")
	if err != nil {
		t.Error(err)
		return
	}
	ret := Secp256k1VerifyChain33(sum, sig, publicKey)
	t.Log("verify", ret)
}
