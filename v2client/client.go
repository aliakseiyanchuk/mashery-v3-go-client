package v2client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"golang.org/x/sync/semaphore"
	"net/http"
	"net/url"
	"time"
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
	V2Authorizer transport.Authorizer
	transport    *transport.HttpTransport
}

func (ci *ClientImpl) Invoke(ctx context.Context, method string, obj interface{}) (V2Result, error) {
	req := V2Request{
		Version: "2.0",
		Method:  method,
		Params:  obj,
		Id:      1,
	}

	m, _ := ci.V2Authorizer.Authorization()
	qs := url.Values{}
	for k, v := range m {
		qs[k] = []string{v}
	}

	if resp, err := ci.transport.Post(ctx, "?"+qs.Encode(), req); err != nil {
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

// NewHTTPClient Create a new V2 HTTP client to invoke Mashery V2 API
func NewHTTPClient(areaNID int, p transport.Authorizer, qps int64, travelTimeComp time.Duration) Client {
	if p == nil {
		panic("v2 HTTP client requires an authorizer")
	}

	return &ClientImpl{
		V2Authorizer: p,
		transport: &transport.HttpTransport{
			MashEndpoint:  fmt.Sprintf("https://api.mashery.com/v2/json-rpc/%d", areaNID),
			Authorizer:    nil,
			Sem:           semaphore.NewWeighted(qps),
			AvgNetLatency: travelTimeComp,
			HttpClient: &http.Client{
				Timeout: time.Second * 60,
			},
		}}
}
