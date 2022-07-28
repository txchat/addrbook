# 去中心化聊天合约

## 开始

### 从源码编译运行

环境要求:
1. Golang 1.17 or later, 参考[golang官方安装文档](https://go.dev/doc/install)

```shell
# 编译本机系统和指令集的可执行文件
$ make build

# 编译目标机器的可执行文件,例如
$ make build_linux_amd64
```
编译成功后目标执行文件在工程目录的`build`文件夹下。

运行服务命令：
```shell
$ ./build/addrbook -f ./build/addrbook.toml
```

## 接口文档
```
统一 jsonrpc 2.0 格式，所有请求均为 `POST` 方法"
```

| 接口                                                       | 功能      |
|----------------------------------------------------------|---------|
| [GetFriends](#getfriends好友列表查询)                          | 好友列表    |
| [CreateRawUpdateFriendTx](#updatefriends更新好友列表)          | 更新好友    |
| [GetBlackList](#getblacklist黑名单查询)                       | 黑名单列表   |
| [CreateRawUpdateBlackTx](#updateblacklist更新黑名单列表)        | 更新黑名单列表 |
| [GetUser](#getuser查询用户信息)                                | 查询用户信息  |
| [GetServerGroup](#getservergroup查询分组信息)                  | 查询分组信息  |
| [CreateRawUpdateUserTx](#updateuser更新用户信息)               | 更新用户信息  |
| [CreateRawUpdateServerGroupTx](#updateservergroup更新分组信息) | 更新分组信息  |

---
### GetFriends(好友列表查询)
[返回概览](#去中心化聊天合约)

请求参数

| **参数**      | **名字** | **类型** | **约束** | **说明**                     |
|-------------|--------|--------|--------|----------------------------|
| method      | 方法     | string | 必填     | chain33框架上的方法, Query 是查询操作 |
| execer      | 执行器    | string | 必填     |                            |
| funcName    | 方法     | string | 必填     |                            |
| mainAddress | 自己的地址  | string | 必填     |                            |
| count       | 查询数量   | int    | 必填     |                            |
| index       | 索引开始地址 | string | 必填     | 空为从头查询                     |
| time        | 时间     | int64  | 必填     | 当前时间戳,单位 毫秒                |
| sign        | 签名信息   | object | 必填     | 详情如下                       |
| publicKey   | 公钥     | string | 必填     | hex字符串表示                   |
| signature   | 签名     | string | 必填     | hex字符串表示                   |

```json
{
	"jsonrpc": "2.0",
	"id": 1,
	"method": "Chain33.Query",
	"params": [{
		"execer": "chat",
		"funcName": "GetFriends",
		"payload": {
			"mainAddress": "12VvgbuFKwuPMwVD6DjmyTBmWeL7HNvczF",
			"count":20,
			"index":"",
			"time": 1581587512000,
			"sign": {
				"publicKey" : "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa",
				"signature" : "93cb8bca06e73f8f2702eceb5cdee419edc47eb45a5c23d20f4bb3bcd28b215c70f3099d035b566c3046d036da2ef3b2f7cc9fbe053b00cae859f5509efb4d2e00"
			}
		}
	}]
}
```

返回参数

| **参数**        | **名字** | **类型**   | **说明** |
|---------------|--------|----------|--------|
| mainAddress   | 地址     | string   | 自己的地址  |
| friendAddress | 地址     | string   | 好友地址   |
| createTime    | 创建时间   | int      | 单位：毫秒  |
| groups        | 好友分组   | string[] | 分组id   |

```
{
	"id": 1,
	"result": {
		"friends": [{
			"mainAddress": "1KiWf5jtovYc2a6NJV5sr5cpE9wENYgnpW",
            "friendAddress": "19hREx1Z6ueTpphyatZC3z1sw4e1Jyxhsk",
            "createTime": "1580817800",
            "groups": ["1","2"]
		}]
	}
}
```

### UpdateFriends(更新好友列表)
[返回概览](#去中心化聊天合约)

请求参数　

| **参数**        | **名字** | **类型**   | 约束   | **说明**   |
|---------------|--------|----------|------|----------|
| friends       | 好友地址   | object[] | true | 详情如下     |
| friendAddress | 地址     | string   | true | 好友地址     |
| type          | 更新类型   | int      | true | 1:添加2:删除 |
| groups        | 好友分组   | string[] | true | 分组id     |

```json
{
	"jsonrpc": "2.0",
	"method": "chat.CreateRawUpdateFriendTx",
	"params": [{
		"friends": [{
			"friendAddress": "1BH9k5vPypTMMHZWtx83zspF4ScBDePXGR",
			"type": 1,
            "groups": ["1"]
		}, {
			"friendAddress": "1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs",
			"type": 1,
            "groups": ["1"]
		}]
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **说明** |
|--------|--------|--------|--------|
| id     | 地址     | object |        |
| result | 结果hash | string |        |
| error  | 错误     | object |        |

```json
{
    "id": 1,
    "result": "0a0463686174123050640a2c0a0210010a260a2231436245565439526e4d356f5a68574d6a34667855724a5839345674526f747a7673100120a08d0630b1b78ddaeda3dbbf653a22313241485774504e68785a4b5152515443504b61334a596e5a374b59356e4776326b",
    "error": null
}
```

### GetBlackList(黑名单查询)
[返回概览](#去中心化聊天合约)

请求参数　　

| **参数**      | **名字** | **类型** | **约束** | **说明**                     |
|-------------|--------|--------|--------|----------------------------|
| method      | 方法     | string | 必填     | chain33框架上的方法, Query 是查询操作 |
| execer      | 执行器    | string | 必填     |                            |
| funcName    | 方法     | string | 必填     |                            |
| mainAddress | 自己的地址  | string | 必填     |                            |
| count       | 查询数量   | int    | 必填     |                            |
| index       | 索引开始地址 | string | 必填     | 空为从头查询                     |
| time        | 时间     | int64  | 必填     | 当前时间戳,单位 毫秒                |
| sign        | 签名信息   | object | 必填     | 详情如下                       |
| publicKey   | 公钥     | string | 必填     | hex字符串表示                   |
| signature   | 签名     | string | 必填     | hex字符串表示                   |

```json
{
	"jsonrpc": "2.0",
	"id": 1,
	"method": "Chain33.Query",
	"params": [{
		"execer": "chat",
		"funcName": "GetBlackList",
		"payload": {
			"mainAddress": "12VvgbuFKwuPMwVD6DjmyTBmWeL7HNvczF",
			"count":20,
			"index":"",
			"time": 1581587512000,
			"sign": {
				"publicKey" : "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa",
				"signature" : "93cb8bca06e73f8f2702eceb5cdee419edc47eb45a5c23d20f4bb3bcd28b215c70f3099d035b566c3046d036da2ef3b2f7cc9fbe053b00cae859f5509efb4d2e00"
			}
		}
	}]
}
```

返回参数

| **参数**        | **名字** | **类型** | **说明**  |
|---------------|--------|--------|---------|
| mainAddress   | 地址     | string | 自己的地址   |
| targetAddress | 地址     | string | 黑名单用户地址 |
| createTime    | 创建时间   | int    | 单位：毫秒   |

```json
{
	"id": 1,
	"result": {
		"list": [{
			"mainAddress": "1KiWf5jtovYc2a6NJV5sr5cpE9wENYgnpW",
            "targetAddress": "19hREx1Z6ueTpphyatZC3z1sw4e1Jyxhsk",
            "createTime": "1580817800"
		}]
	}
}
```

### UpdateBlackList(更新黑名单列表)
[返回概览](#去中心化聊天合约)

请求参数　

| **参数**        | **名字** | **类型**   | **约束**   | **说明** |
|---------------|--------|----------|----------|--------|
| list          | 黑民单列表  | object[] | 详情如下     |        |
| targetAddress | 地址     | string   | 好友地址     |        |
| type          | 更新类型   | int      | 1:添加2:删除 |        |

```json
{
	"jsonrpc": "2.0",
	"method": "chat.CreateRawUpdateBlackTx",
	"params": [{
		"list": [{
			"targetAddress": "1BH9k5vPypTMMHZWtx83zspF4ScBDePXGR",
			"type": 1
		}, {
			"targetAddress": "1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs",
			"type": 1
		}]
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **说明** |
|--------|--------|--------|--------|
| id     | 地址     | object |        |
| result | 结果hash | string |        |
| error  | 错误     | object |        |

```json
{
    "id": 1,
    "result": "0a0463686174123050640a2c0a0210010a260a2231436245565439526e4d356f5a68574d6a34667855724a5839345674526f747a7673100120a08d0630b1b78ddaeda3dbbf653a22313241485774504e68785a4b5152515443504b61334a596e5a374b59356e4776326b",
    "error": null
}
```

### GetUser(查询用户信息)

[返回概览](#去中心化聊天合约)

请求参数　　

|   **参数**    |   **名字**   | **类型** | **约束** |               **说明**                |
| :-----------: | :----------: | :------: | :------: | :-----------------------------------: |
|    method     |     方法     |  string  |   必填   | chain33框架上的方法, Query 是查询操作 |
|    execer     |    执行器    |  string  |   必填   |                                       |
|   funcName    |     方法     |  string  |   必填   |                                       |
|  mainAddress  |  自己的地址  |  string  |   必填   |                                       |
| targetAddress |   目标地址   |  string  |   必填   |                                       |
|     count     |   查询数量   |   int    |   必填   |                                       |
|     index     | 索引开始地址 |  string  |   必填   |             空为从头查询              |
|     time      |     时间     |  int64   |   必填   |         当前时间戳,单位 毫秒          |
|     sign      |   签名信息   |  object  |   必填   |               详情如下                |
|   publicKey   |     公钥     |  string  |   必填   |             hex字符串表示             |
|   signature   |     签名     |  string  |   必填   |             hex字符串表示             |

```json
{
	"jsonrpc": "2.0",
	"id": 1,
	"method": "Chain33.Query",
	"params": [{
		"execer": "chat",
		"funcName": "GetUser",
		"payload": {
			"mainAddress": "12VvgbuFKwuPMwVD6DjmyTBmWeL7HNvczF",
			"targetAddress": "1BH9k5vPypTMMHZWtx83zspF4ScBDePXGR",
			"count":20,
			"index":"",
			"time": 1581587512000,
			"sign": {
				"publicKey" : "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa",
				"signature" : "93cb8bca06e73f8f2702eceb5cdee419edc47eb45a5c23d20f4bb3bcd28b215c70f3099d035b566c3046d036da2ef3b2f7cc9fbe053b00cae859f5509efb4d2e00"
			}
		}
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **约束** |            **说明**            |
| :------: | :------: | :------: | :------: | :----------------------------: |
|  groups  | 分组列表 | string[] |   true   |     如果查询的是好友，返回好友分组列表    |
|  chatServers  | 服务列表 | object[] |   true   |            详情如下            |
|  id   | 编号 |  string  |   true   |  |
|   name   |   名称   |  string  |   true   |                                |
|  address   |    地址    |  string  |   true   |                                |
|  fields  | 字段列表 | object[] |   true   |            详情如下            |
|   name   |   名称   |  string  |   true   |                                |
|  value   |    值    |  string  |   true   |                                |
|  level   | 安全级别 |  string  |   true   | [private]、[protect]、[public] |

```json
{
	"id": 1,
	"result": {
        "groups":[],
		"chatServers": [{
            "id":"1",
			"name": "",
			"address": "127.0.0.1:8088"
		}],
		"fields": [{
			"name": "nickname",
			"value": "张三",
			"level": "public"
		}, {
			"name": "phone",
			"value": "157********",
			"level": "private"
		}]
	}
}
```

安全级别描述：

- private：只有自己能查询
- protect：好友能查询
- public：所有人能查询

### GetServerGroup(查询分组信息)

[返回概览](#去中心化聊天合约)

请求参数　　

|  **参数**   |   **名字**   | **类型** | **约束** |               **说明**                |
| :---------: | :----------: | :------: | :------: | :-----------------------------------: |
|   method    |     方法     |  string  |   必填   | chain33框架上的方法, Query 是查询操作 |
|   execer    |    执行器    |  string  |   必填   |                                       |
|  funcName   |     方法     |  string  |   必填   |                                       |
| mainAddress |  自己的地址  |  string  |   必填   |                                       |
|    count    |   查询数量   |   int    |   必填   |                                       |
|    index    | 索引开始地址 |  string  |   必填   |             空为从头查询              |
|    time     |     时间     |  int64   |   必填   |         当前时间戳,单位 毫秒          |
|    sign     |   签名信息   |  object  |   必填   |               详情如下                |
|  publicKey  |     公钥     |  string  |   必填   |             hex字符串表示             |
|  signature  |     签名     |  string  |   必填   |             hex字符串表示             |

```json
{
	"jsonrpc": "2.0",
	"id": 1,
	"method": "Chain33.Query",
	"params": [{
		"execer": "chat",
		"funcName": "GetServerGroup",
		"payload": {
			"mainAddress": "12VvgbuFKwuPMwVD6DjmyTBmWeL7HNvczF",
			"count":20,
			"index":"",
			"time": 1581587512000,
			"sign": {
				"publicKey" : "026372967a9885e1248c4d06407582d3808dde75b931f2186e68246698aa92dfaa",
				"signature" : "93cb8bca06e73f8f2702eceb5cdee419edc47eb45a5c23d20f4bb3bcd28b215c70f3099d035b566c3046d036da2ef3b2f7cc9fbe053b00cae859f5509efb4d2e00"
			}
		}
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **约束** | **说明** |
| :------: | :------: | :------: | :------: | :------: |
|  groups  | 字段列表 | object[] |   true   | 详情如下 |
|    id    |   索引   |  string  |   true   |          |
|   name   |   名称   |  string  |   true   |          |
|  value   |    值    |  string  |   true   |          |

```json
{
	"id": 1,
	"result": {
		"groups": [{
            "id": "",
            "name": "g1",
            "value": "192...."
        }, {
            "id": "",
            "name": "g2",
            "value": "183."
        }]
	}
}
```

### UpdateUser(更新用户信息)

[返回概览](#去中心化聊天合约)

请求参数　

| **参数** | **名字** | **类型** | **约束** |            **说明**            |
| :------: | :------: | :------: | :------: | :----------------------------: |
|  fields  | 字段列表 | object[] |   true   |            详情如下            |
|   name   |   名称   |  string  |   true   |                                |
|  value   |    值    |  string  |   true   |                                |
|   type   | 操作类型 |   int    |   true   |     1->新增/修改；2->删除      |
|  level   | 安全级别 |  string  |   true   | [private]、[protect]、[public] |

```json
{
	"jsonrpc": "2.0",
	"method": "chat.CreateRawUpdateUserTx",
	"params": [{
		"fields": [{
			"name": "nickname",
			"value": "张三",
            "type": 1,
			"level": "public"
		}, {
			"name": "phone",
			"value": "157********",
            "type": 1,
			"level": "private"
		}]
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **说明** |
| :------: | :------: | :------: | :------: |
|    id    |   地址   |  object  |          |
|  result  | 结果hash |  string  |          |
|  error   |   错误   |  object  |          |

```json
{
    "id": 1,
    "result": "0a0463686174123050640a2c0a0210010a260a2231436245565439526e4d356f5a68574d6a34667855724a5839345674526f747a7673100120a08d0630b1b78ddaeda3dbbf653a22313241485774504e68785a4b5152515443504b61334a596e5a374b59356e4776326b",
    "error": null
}
```

字段枚举：

```
昵称：nickname
头像：avatar
手机号：phone
邮箱：email
公钥：pubKey

链信息：模式匹配`chain\.\w+$`
```

### UpdateServerGroup(更新分组信息)

[返回概览](#去中心化聊天合约)

请求参数　

| **参数** | **名字** | **类型** | **约束** |         **说明**          |
| :------: | :------: | :------: | :------: | :-----------------------: |
|  groups  |   分组   | object[] |   true   |         详情如下          |
|    id    |   序号   |  string  |   true   |                           |
|   type   | 操作类型 |   int    |   true   | 1->新增；2->删除；3->修改 |
|   name   |   名称   |  string  |   true   |                           |
|  value   |    值    |  string  |   true   |                           |

```json
{
	"jsonrpc": "2.0",
	"method": "chat.CreateRawUpdateServerGroupTx",
	"params": [{
        "groups": [{
            "id": "1",
            "type": 1,
            "name": "g1",
            "value": "192...."
        }, {
            "id": "2",
            "type": 1,
            "name": "g2",
            "value": "183."
        }]
	}]
}
```

返回参数

| **参数** | **名字** | **类型** | **说明** |
| :------: | :------: | :------: | :------: |
|    id    |   地址   |  object  |          |
|  result  | 结果hash |  string  |          |
|  error   |   错误   |  object  |          |

```json
{
    "id": 1,
    "result": "0a0463686174123050640a2c0a0210010a260a2231436245565439526e4d356f5a68574d6a34667855724a5839345674526f747a7673100120a08d0630b1b78ddaeda3dbbf653a22313241485774504e68785a4b5152515443504b61334a596e5a374b59356e4776326b",
    "error": null
}
```

注意：默认分组名称约定为："default"

### 查询接口的签名和验证

1. 将参数（参数不包含Sign字段）按照《键名》字典序升序排序
2. 将参数按照 k1=v1&k2=v2.... 的形式组合成字符串（内容为空的也要带上）
3. 将第二步的结果进行sha256
4. 使用secp256k1加密第三步的结果
5. 将第四步hex编码成十六进制字符串填入结果中