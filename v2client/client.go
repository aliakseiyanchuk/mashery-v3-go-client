package v2client

import (
	"context"
	"encoding/json"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

type JSONRPCQueryResult struct {
	TotalItems   int           `json:"total_items"`
	TotalPages   int           `json:"total_pages"`
	ItemsPerPage int           `json:"items_per_page"`
	CurrentPage  int           `json:"current_page"`
	Items        []interface{} `json:"items"`
}

type JSONRPCError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

type V2Result struct {
	Version string        `json:"jsonrpc"`
	Id      *int          `json:"id"`
	Result  interface{}   `json:"result"`
	Error   *JSONRPCError `json:"error"`
}

type V2Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type Client interface {
	Invoke(ctx context.Context, method string, obj interface{}) (V2Result, error)
}

type ClientImpl struct {
	transport *transport.HttpTransport
}

func (ci *ClientImpl) Invoke(ctx context.Context, method string, obj interface{}) (V2Result, error) {
	req := V2Request{
		Version: "2.0",
		Method:  method,
		Params:  obj,
		Id:      1,
	}

	if resp, err := ci.transport.Post(ctx, "", req); err != nil {
		return V2Result{}, &errwrap.WrappedError{
			Context: "sending V2 post request",
			Cause:   err,
		}
	} else {
		if body, err := transport.ReadResponseBody(resp); err != nil {
			return V2Result{}, &errwrap.WrappedError{
				Context: "reading V2 response",
				Cause:   err,
			}
		} else {
			var rv V2Result
			err := json.Unmarshal(body, &rv)
			return rv, err
		}
	}
}
