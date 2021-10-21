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

func TestGetTransactionCall(t *testing.T) {

	cfg := config.NewConfig(config.Mainnet, "10dea4f9680f440bb9eaf33ecc51c1ef")
	rpc := NewRPC(cfg.GetURL(), WithHttpClient(http.DefaultClient), WithLogger(nil), WithDebug(false))

	tt := []struct {
		name         string
		block        int
		index        int
		expectedHash string
		err          bool
	}{
		{"valid input 1", 11537751, 4, "0x29aee213e5dafbb3c58975b2610d9d3c89e290f15cd7511c0ee90edf93211da5", false},
		{"valid input 2", 11536651, 10, "0xa7c68d639b8528129355e7f1d13a8d9607c676d76d2bdd985cc21de16ed2f9b7", false},
		{"valid input 3", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", false},
		{"missing transaction index", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"negative transaction index", 11246512, -10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"negative block number", -11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"missing block number", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var params []interface{}

			switch tc.name {
			case "missing transaction index":
				params = []interface{}{IntToHex(tc.block)}
			case "missing block number":
				params = []interface{}{IntToHex(tc.index)}
			default:
				params = []interface{}{IntToHex(tc.block), IntToHex(tc.index)}
			}

			req := RPCRequest{
				JSONRPC: "2.0",
				Method:  "eth_getTransactionByBlockNumberAndIndex",
				Params:  params,
				ID:      1,
			}

			body, err := json.Marshal(req)
			assert.Nil(t, err)
			assert.NotNil(t, body)

			response, err := rpc.client.Post(rpc.URL(), "application/json", bytes.NewBuffer(body))
			assert.NotNil(t, response)
			assert.Nil(t, err)

			data, err := ioutil.ReadAll(response.Body)
			assert.NotNil(t, data)
			assert.Nil(t, err)

			resp := new(RPCResponse)
			if err := json.Unmarshal(data, resp); err != nil {
				t.Errorf("Unable to unmarshal data to RPCResponse, error: %v", err)
			}

			txjson := resp.Result
			tx := new(Transaction)

			if tc.err == true {
				assert.Nil(t, txjson)
				if err := json.Unmarshal(txjson, tx); err == nil {
					t.Error("Unmarshalling error must be not nil")
				}
				assert.NotEqual(t, tx.Hash, tc.expectedHash)
				return
			}

			assert.NotNil(t, txjson)
			if err := json.Unmarshal(txjson, tx); err != nil {
				t.Errorf("Unable to unmarshal json to Transaction, error: %v", err)
			}

			assert.Equal(t, tx.Hash, tc.expectedHash)

		})
	}
}

func TestGetTransactionAPI(t *testing.T) {
	endpoint := "http://18.222.141.251:8080/transaction"

	tt := []struct {
		name         string
		block        int
		index        int
		expectedHash string
		err          bool
	}{

		{"valid input 1", 11537751, 4, "0x29aee213e5dafbb3c58975b2610d9d3c89e290f15cd7511c0ee90edf93211da5", false},
		{"valid input 2", 11536651, 10, "0xa7c68d639b8528129355e7f1d13a8d9607c676d76d2bdd985cc21de16ed2f9b7", false},
		{"valid input 3", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", false},
		{"missing transaction index", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"negative transaction index", 11246512, -10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"negative block number", -11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"missing block number", 11246512, 10, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
		{"invalid arguments", 0, 0, "0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932", true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var m map[string]interface{}

			switch tc.name {
			case "missing transaction index":
				m = map[string]interface{}{"block": tc.block}
			case "missing block number":
				m = map[string]interface{}{"index": tc.index}
			case "invalid arguments":
				m = map[string]interface{}{"block": "test", "index": "test"}
			default:
				m = map[string]interface{}{"block": tc.block, "index": tc.index}
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

			tx := new(Transaction)

			if err = json.Unmarshal(data, tx); err != nil {
				t.Errorf("Unable unmarshal json to Transaction, err: %v ", err)
			}

			assert.Equal(t, tx.Hash, tc.expectedHash)

		})
	}
}
