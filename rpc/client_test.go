package rpc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"math/big"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/cyberbono3/infura-coding/config"
)

type RPCTestSuite struct {
	suite.Suite
	client *RPC
}

func (s *RPCTestSuite) registerResponse(result string, callback func([]byte)) {
	httpmock.Reset()
	response := fmt.Sprintf(`{"jsonrpc":"2.0", "id":1, "result": %s}`, result)
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		callback(s.getBody(request))
		return httpmock.NewStringResponse(200, response), nil
	})
}

func (s *RPCTestSuite) getBody(request *http.Request) []byte {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	s.Require().Nil(err)

	return body
}

func (s *RPCTestSuite) TestEthGetTransactionByBlockNumberAndIndex() {
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getTransactionByBlockNumberAndIndex")
		s.paramsEqual(body, `["0x1f537da", "0xa"]`)
	})

	tx, err := s.client.EthGetTransactionByBlockNumberAndIndex(32847834, 10)
	s.Require().Nil(err)
	s.Require().NotNil(tx)
}

func (s *RPCTestSuite) methodEqual(body []byte, expected string) {
	value := gjson.GetBytes(body, "method").String()

	s.Require().Equal(expected, value)
}

func (s *RPCTestSuite) paramsEqual(body []byte, expected string) {
	value := gjson.GetBytes(body, "params").Raw
	if expected == "null" {
		s.Require().Equal(expected, value)
	} else {
		s.JSONEq(expected, value)
	}
}

func (s *RPCTestSuite) SetupSuite() {
	cfg := config.NewConfig(config.Mainnet, "10dea4f9680f440bb9eaf33ecc51c1ef")
	s.client = NewRPC(cfg.GetURL(), WithHttpClient(http.DefaultClient), WithLogger(nil), WithDebug(false))

	httpmock.Activate()
}

func (s *RPCTestSuite) TearDownSuite() {
	httpmock.Deactivate()
}

func (s *RPCTestSuite) TearDownTest() {
	httpmock.Reset()
}

func (s *RPCTestSuite) TestCall() {
	// Test http error
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		return nil, errors.New("Error")
	})

	_, err := s.client.Call("test")
	s.Require().NotNil(err)
	httpmock.Reset()

	// Test invalid response format
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, "{213"), nil
	})
	_, err = s.client.Call("test")
	s.Require().NotNil(err)
	httpmock.Reset()

	// Test eth error
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, `{"error": {"code": 21, "message": "eee"}}`), nil
	})
	_, err = s.client.Call("test")
	s.Require().NotNil(err)
	rpcError, ok := err.(RPCError)
	s.Require().True(ok)
	s.Require().Equal(21, rpcError.GetCode())
	s.Require().Equal("eee", rpcError.GetMessage())
}

func (s *RPCTestSuite) Test_call() {
	// Test http error
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("Error")
	})
	err := s.client.call("test", nil)
	s.Require().NotNil(err)

	// Test target is nil
	s.registerResponse(`{"foo": "bar"}`, func([]byte) {})
	err = s.client.call("test", nil)
	s.Require().Nil(err)

	// Test invalid target
	target := ""
	err = s.client.call("test", &target)
	s.Require().NotNil(err)
}


func (s *RPCTestSuite) registerResponseError(err error) {
	httpmock.Reset()
	httpmock.RegisterResponder("POST", s.client.URL(), func(request *http.Request) (*http.Response, error) {
		return nil, err
	})
}

func (s *RPCTestSuite) TestGetBlock() {
	s.registerResponseError(errors.New("Error"))
	block, err := s.client.getBlock("eth_getBlockByHash", true)
	s.Require().NotNil(err)

	// Test with transactions
	result := ` {
        "difficulty": "0x81299d4dbde29",
        "extraData": "0x706f6f6c2e65746866616e732e6f726720284d4e323729",
        "gasLimit": "0x667900",
        "gasUsed": "0x639fa0",
        "hash": "0x2bdda43f649c564642101fc990f569dd855e60f88bf83e931f509a92c62700f9",
        "logsBloom": "0x111",
        "miner": "0x1e9939daaad6924ad004c2560e90804164900341",
        "mixHash": "0xa6b69fa82eaea8674236170a2d8ea41d80c176315a579138b718f3bcaa4c39ab",
        "nonce": "0xefd7ef000d0b78b8",
        "number": "0x4055d5",
        "parentHash": "0x913f938dcb4ff83b2b6b42a0cf6517d438a3ce95174e9342c780fd20c84dfd03",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x2fc6",
        "stateRoot": "0xab9287d3b8864338892d1d572198933979e39bfcfbde569ea52be15a9691b4c1",
        "timestamp": "0x59a556bd",
        "totalDifficulty": "0x2b5f79e86aaf701c81",
        "transactions": [
			{
				"blockHash": "0x2bdda43f649c564642101fc990f569dd855e60f88bf83e931f509a92c62700f9",
				"blockNumber": "0x4055d5",
				"from": "0xa95350d70b18fa29f6b5eb8d627ceeeee499340d",
				"gas": "0x5208",
				"gasPrice": "0x6edf2a079e",
				"hash": "0xf519ca0e9ceeb0405dfeb95544179f557e3221213f07e33709af7ced60ab61b9",
				"input": "0x",
				"nonce": "0x289b",
				"to": "0xb595f3390fcec074237c8264b908fc73d4aedc93",
				"transactionIndex": "0x0",
				"value": "0xdbd2fc137a30000"
			},
			{
				"blockHash": "0x2bdda43f649c564642101fc990f569dd855e60f88bf83e931f509a92c62700f9",
				"blockNumber": "0x4055d5",
				"from": "0x0f1b76410215ed963ea2c3d3eaddd4a56350b422",
				"gas": "0x3d090",
				"gasPrice": "0x1176592e00",
				"hash": "0xa72743a3608e2ae7b3d1cc1f0e3ceed9a1c78d803eba5f28d5d6908adfaa211c",
				"input": "0x278b8c0e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004551ce090d138000000000000000000000000000006bea43baa3f7a6f765f14f10a1a1b08334ef4500000000000000000000000000000000000000000000003627e8f712373c000000000000000000000000000000000000000000000000000000000000004059b200000000000000000000000000000000000000000000000000000000418e8e7d000000000000000000000000000000000000000000000000000000000000001b64b1fee882b69969c9395a095e45e4b0abb3b19806ba040a6765194f966ae64e24a3d44a837de95e014b6b0f7eea075e30cca0414a18c0a27a7f349271689f3d",
				"nonce": "0x1c2",
				"to": "0x8d12a197cb00d4747a1fe03395095ce2a5cc6819",
				"transactionIndex": "0x1",
				"value": "0x0"
			}
		],
        "transactionsRoot": "0x97849642410701c38f904912238eb78d3aa854e72c5ae39394c7217f4f9474bc",
        "uncles": ["0xf14cdb8a75de31dcf3da7a3a52c1fffcbaa3d56de9f50f86767fa411c10f4397"]
	}`
	hash := "0x2bdda43f649c564642101fc990f569dd855e60f88bf83e931f509a92c62700f9"
	s.registerResponse(result, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByHash")
	})

	block, err = s.client.getBlock("eth_getBlockByHash", true)
	s.Require().Nil(err)
	s.Require().NotNil(block)
	s.Require().Equal(hash, block.Hash)
	s.Require().Equal(4216277, block.Number)
	s.Require().Equal("0x913f938dcb4ff83b2b6b42a0cf6517d438a3ce95174e9342c780fd20c84dfd03", block.ParentHash)
	s.Require().Equal("0xefd7ef000d0b78b8", block.Nonce)
	s.Require().Equal("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", block.Sha3Uncles)
	s.Require().Equal("0x111", block.LogsBloom)
	s.Require().Equal("0x97849642410701c38f904912238eb78d3aa854e72c5ae39394c7217f4f9474bc", block.TransactionsRoot)
	s.Require().Equal("0xab9287d3b8864338892d1d572198933979e39bfcfbde569ea52be15a9691b4c1", block.StateRoot)
	s.Require().Equal("0x1e9939daaad6924ad004c2560e90804164900341", block.Miner)
	s.Require().Equal(newBigInt("2272251724160553"), block.Difficulty)
	s.Require().Equal(newBigInt("800089780620203400321"), block.TotalDifficulty)
	s.Require().Equal("0x706f6f6c2e65746866616e732e6f726720284d4e323729", block.ExtraData)
	s.Require().Equal(12230, block.Size)
	s.Require().Equal(6715648, block.GasLimit)
	s.Require().Equal(6528928, block.GasUsed)
	s.Require().Equal(1504007869, block.Timestamp)
	s.Require().Equal([]string{"0xf14cdb8a75de31dcf3da7a3a52c1fffcbaa3d56de9f50f86767fa411c10f4397"}, block.Uncles)
	s.Require().Equal(2, len(block.Transactions))

	s.Require().Equal(Transaction{
		Hash:             "0xf519ca0e9ceeb0405dfeb95544179f557e3221213f07e33709af7ced60ab61b9",
		Nonce:            10395,
		BlockHash:        block.Hash,
		BlockNumber:      &block.Number,
		TransactionIndex: ptrInt(0),
		From:             "0xa95350d70b18fa29f6b5eb8d627ceeeee499340d",
		To:               "0xb595f3390fcec074237c8264b908fc73d4aedc93",
		Value:            newBigInt("990000000000000000"),
		Gas:              21000,
		GasPrice:         newBigInt("476190476190"),
		Input:            "0x",
	}, block.Transactions[0])

	s.Require().Equal(Transaction{
		Hash:             "0xa72743a3608e2ae7b3d1cc1f0e3ceed9a1c78d803eba5f28d5d6908adfaa211c",
		Nonce:            450,
		BlockHash:        block.Hash,
		BlockNumber:      &block.Number,
		TransactionIndex: ptrInt(1),
		From:             "0x0f1b76410215ed963ea2c3d3eaddd4a56350b422",
		To:               "0x8d12a197cb00d4747a1fe03395095ce2a5cc6819",
		Value:            newBigInt("0"),
		Gas:              250000,
		GasPrice:         newBigInt("75000000000"),
		Input:            "0x278b8c0e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004551ce090d138000000000000000000000000000006bea43baa3f7a6f765f14f10a1a1b08334ef4500000000000000000000000000000000000000000000003627e8f712373c000000000000000000000000000000000000000000000000000000000000004059b200000000000000000000000000000000000000000000000000000000418e8e7d000000000000000000000000000000000000000000000000000000000000001b64b1fee882b69969c9395a095e45e4b0abb3b19806ba040a6765194f966ae64e24a3d44a837de95e014b6b0f7eea075e30cca0414a18c0a27a7f349271689f3d",
	}, block.Transactions[1])

	httpmock.Reset()
	// Test without transactions
	result = `{
		"difficulty": "0x7feab8ef4d978",
		"extraData": "0xd58301050b8650617269747986312e31352e31826c69",
		"gasLimit": "0x665f6b",
		"gasUsed": "0x1d71b",
		"hash": "0x23be1464d0e805fe3cec49039a9cf7fae7c09d2efacbed2abb10ef7ddae960ab",
		"logsBloom": "0x222",
		"miner": "0x6a7a43be33ba930fe58f34e07d0ad6ba7adb9b1f",
		"mixHash": "0xa8f339af405f7f3a7b7c163f8889f44343abfbbeda13c41e06923de349ea6483",
		"nonce": "0x19a48ee424b5088f",
		"number": "0x4105f3",
		"parentHash": "0xbc3e37984a619008d75e7f73865247fb420ae5ed2c921599d099ab5f20519396",
		"receiptsRoot": "0xa1384524d42ff86fdf4e44eeea853aba4e772a52240037cbfddc22782bad017e",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size": "0x490",
		"stateRoot": "0xbe7e86ee05a5d49ba64b3d9f3f0129bab90308032e42307a1a2ef5c8971c5f5c",
		"timestamp": "0x59b62713",
		"totalDifficulty": "0x30e3d47fb9d7a43f7c",
		"transactions": [
			"0x160e19780a24f3d78492c7ac7228e0220d4b96878fec19daf182e1d8c4b3d94e"
		],
		"transactionsRoot": "0x1bcd58c2420d63c5e8ed3182afd33c01737be38a4a8c10a81dfb70b692e8f286",
		"uncles": []
	}`

	hash = "0x23be1464d0e805fe3cec49039a9cf7fae7c09d2efacbed2abb10ef7ddae960ab"
	s.registerResponse(result, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByHash")
	})

	block, err = s.client.getBlock("eth_getBlockByHash", false)
	s.Require().Nil(err)
	s.Require().NotNil(block)
	s.Require().Equal(hash, block.Hash)
	s.Require().Equal(4261363, block.Number)
	s.Require().Equal("0xbc3e37984a619008d75e7f73865247fb420ae5ed2c921599d099ab5f20519396", block.ParentHash)
	s.Require().Equal("0x19a48ee424b5088f", block.Nonce)
	s.Require().Equal("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", block.Sha3Uncles)
	s.Require().Equal("0x222", block.LogsBloom)
	s.Require().Equal("0x1bcd58c2420d63c5e8ed3182afd33c01737be38a4a8c10a81dfb70b692e8f286", block.TransactionsRoot)
	s.Require().Equal("0xbe7e86ee05a5d49ba64b3d9f3f0129bab90308032e42307a1a2ef5c8971c5f5c", block.StateRoot)
	s.Require().Equal("0x6a7a43be33ba930fe58f34e07d0ad6ba7adb9b1f", block.Miner)
	s.Require().Equal(newBigInt("2250337628248440"), block.Difficulty)
	s.Require().Equal(newBigInt("901860602515894321020"), block.TotalDifficulty)
	s.Require().Equal("0xd58301050b8650617269747986312e31352e31826c69", block.ExtraData)
	s.Require().Equal(1168, block.Size)
	s.Require().Equal(6709099, block.GasLimit)
	s.Require().Equal(120603, block.GasUsed)
	s.Require().Equal(1505109779, block.Timestamp)
	s.Require().Equal([]string{}, block.Uncles)
	s.Require().Equal(1, len(block.Transactions))
	s.Require().Equal(Transaction{
		Hash:             "0x160e19780a24f3d78492c7ac7228e0220d4b96878fec19daf182e1d8c4b3d94e",
		Nonce:            0,
		BlockHash:        "",
		BlockNumber:      nil,
		TransactionIndex: nil,
		From:             "",
		To:               "",
		Value:            big.Int{},
		Gas:              0,
		GasPrice:         big.Int{},
		Input:            "",
	}, block.Transactions[0])

	s.registerResponse("null", func(body []byte) {})

	block, err = s.client.getBlock("eth_getBlockByHash", false)
	s.Require().Nil(block)
	s.Require().Nil(err)
}

func TestRPCTestSuite(t *testing.T) {
	suite.Run(t, new(RPCTestSuite))
}
