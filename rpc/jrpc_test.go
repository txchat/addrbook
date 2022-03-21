package rpc

import (
	"strings"
	"testing"

	commonlog "github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/rpc/jsonclient"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
	"github.com/stretchr/testify/assert"
	cty "github.com/txchat/addrbook/types"
)

func init() {
	commonlog.SetLogLevel("error")
}

func Test_JRPCChannel(t *testing.T) {
	// 启动RPCmocker
	mocker := testnode.New("--notset--", nil)
	defer func() {
		mocker.Close()
	}()
	mocker.Listen()

	jrpcClient := mocker.GetJSONC()
	assert.NotNil(t, jrpcClient)

	testCases := []struct {
		fn func(*testing.T, *jsonclient.JSONClient) error
	}{
		{fn: testAssetCmd},
	}
	for index, testCase := range testCases {
		err := testCase.fn(t, jrpcClient)
		if err == nil {
			continue
		}
		assert.NotEqualf(t, err, types.ErrActionNotSupport, "test index %d", index)
		if strings.Contains(err.Error(), "rpc: can't find") {
			assert.FailNowf(t, err.Error(), "test index %d", index)
		}
	}
}

func testAssetCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := cty.ReqGetFriends{}
	return jrpc.Call("chat.GetFriends", params, nil)
}
