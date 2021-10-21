package handlers

import (
	"github.com/cyberbono3/infura-coding/client-api/data"
	"github.com/cyberbono3/infura-coding/rpc"
	"github.com/hashicorp/go-hclog"
)

// Keys used for the object in the context
type KeyTransaction struct{}
type KeyBlock struct{}

// ClientHandler handler for apis
type ClientHandler struct {
	c *rpc.RPC
	l hclog.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with the given logger
func NewClientHandler(c *rpc.RPC, l hclog.Logger, v *data.Validation) *ClientHandler {
	return &ClientHandler{c, l, v}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
