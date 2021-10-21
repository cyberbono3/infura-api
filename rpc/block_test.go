package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cyberbono3/infura-coding/config"
)

// Good
func TestGetBlockCall(t *testing.T) {

	cfg := config.NewConfig(config.Mainnet, "10dea4f9680f440bb9eaf33ecc51c1ef")
	rpc := NewRPC(cfg.GetURL(), WithHttpClient(http.DefaultClient), WithLogger(nil), WithDebug(false))

	tt := []struct {
		name      string
		block     int
		show      bool
		blockHash string
		err       bool
	}{

		{"valid input 1", 114996, true, "0x6cb46b8ff3cfa8f58b21e04bfe99cc75d71b04c6d15953ed9c73026acef4bf68", false},
		{"valid input 2", 114955, false, "0xe51453b66daf532b24bc0938be40e4c36d8a6f0b89853891bb079b6c6dd278a7", false},
		{"valid input 3", 11496784, false, "0xc56b20ce7a7246664b26dd24444177db4b27111098fe8eab66ceb19e8d2cc292", false},
		{"valid input 4", 11496784, true, "0xc56b20ce7a7246664b26dd24444177db4b27111098fe8eab66ceb19e8d2cc292", false},

		{"negative block number", -11246512, true, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", true},
		{"missing block number", 11246512, false, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", true},
		{"missing showTransFlag", 11246512, false, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var params []interface{}

			switch tc.name {
			case "missing block number":
				params = []interface{}{tc.show}
			case "missing showTransFlag":
				params = []interface{}{IntToHex(tc.block)}
			default:
				params = []interface{}{IntToHex(tc.block), tc.show}
			}

			req := RPCRequest{
				JSONRPC: "2.0",
				Method:  "eth_getBlockByNumber",
				Params:  params,
				ID:      1,
			}

			body, err := json.Marshal(req)
			assert.Nil(t, err)
			assert.NotNil(t, body)

			response, err := rpc.client.Post(cfg.GetURL(), "application/json", bytes.NewBuffer(body))
			assert.NotNil(t, response)
			assert.Nil(t, err)

			data, err := ioutil.ReadAll(response.Body)
			assert.NotNil(t, data)
			assert.Nil(t, err)

			rpcResponse := new(RPCResponse)
			if err := json.Unmarshal(data, rpcResponse); err != nil {
				t.Errorf("Unable to unmarshal data to RPCResponse, error: %v", err)
			}

			blockjson := rpcResponse.Result

			var responseBlock proxyBlock
			if tc.show {
				responseBlock = new(proxyBlockWithTransactions)
			} else {
				responseBlock = new(proxyBlockWithoutTransactions)
			}

			if tc.err == true {
				assert.Nil(t, blockjson)
				if err := json.Unmarshal(blockjson, responseBlock); err == nil {
					t.Error("Unmarshalling error must not be nil")
				}
				block := responseBlock.toBlock()
				assert.NotEqual(t, block.Hash, tc.blockHash)
				return
			}

			assert.NotNil(t, blockjson)

			if err = json.Unmarshal(blockjson, responseBlock); err != nil {
				t.Errorf("Unable to unmarshal json to Block, error: %v", err)
			}

			block := responseBlock.toBlock()
			t.Log("hash", block.Hash)

			assert.Equal(t, block.Hash, tc.blockHash)

		})
	}
}

//Good
func TestGetBlockAPI(t *testing.T) {
	endpoint := "http://18.222.141.251:8080/block"

	tt := []struct {
		name      string
		block     int
		show      bool
		blockHash string
		err       bool
	}{

		{"valid input 1", 114996, false, "0x6cb46b8ff3cfa8f58b21e04bfe99cc75d71b04c6d15953ed9c73026acef4bf68", false},
		{"valid input 2", 114955, false, "0xe51453b66daf532b24bc0938be40e4c36d8a6f0b89853891bb079b6c6dd278a7", false},
		{"valid input 3", 11496784, false, "0xc56b20ce7a7246664b26dd24444177db4b27111098fe8eab66ceb19e8d2cc292", false},
		{"valid input 4", 11496784, false, "0xc56b20ce7a7246664b26dd24444177db4b27111098fe8eab66ceb19e8d2cc292", false},
		{"negative block number", -11246512, true, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", true},
		{"missing block number", 11246512, false, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", true},
		{"missing showTransFlag", 11246512, false, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", false},
		{"invalid arguments", 11246512, false, "0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0", false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var m map[string]interface{}

			switch tc.name {
			case "missing block number":
				m = map[string]interface{}{"show": tc.show}
			case "missing showTransFlag":
				m = map[string]interface{}{"block": tc.block}
			case "invalid arguments":
				m = map[string]interface{}{"block": "test", "show": "test"}
			default:
				m = map[string]interface{}{"block": tc.block, "show": tc.show}
			}

			body, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Unable to marshal payload json to JSON, err: %v", err)
			}
			req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("Unable to return request object, err: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Unable to execute post request, err: %v", err)
			}
			defer resp.Body.Close()

			if tc.name == "invalid arguments" {
				if http.StatusBadRequest != resp.StatusCode {
					t.Errorf("Response status code is not valid, expected: %v, got : %v ", http.StatusBadRequest, resp.StatusCode)
				}
				return
			}

			if tc.err {
				if http.StatusUnprocessableEntity != resp.StatusCode {
					t.Errorf("Response status code is not valid, expected: %v, got : %v ", http.StatusUnprocessableEntity, resp.StatusCode)
				}
				return
			}

			if http.StatusOK != resp.StatusCode {
				t.Errorf("Response status code is not OK, expected: %v, got: %v ", http.StatusOK, resp.StatusCode)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Unable to read response body, err: %v", err)
			}

			var responseBlock proxyBlock
			if tc.show {
				responseBlock = new(proxyBlockWithTransactions)
			} else {
				responseBlock = new(proxyBlockWithoutTransactions)
			}

			_ = json.Unmarshal(data, responseBlock)

			block := responseBlock.toBlock()

			assert.Equal(t, block.Hash, tc.blockHash)

		})
	}
}
