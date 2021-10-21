package rpc

import "fmt"

// EthError - ethereum error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err RPCError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

func (err RPCError) GetCode() int {
	return err.Code
}

func (err RPCError) GetMessage() string {
	return err.Message
}
