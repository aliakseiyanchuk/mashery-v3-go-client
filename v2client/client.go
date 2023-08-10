package v2client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
	"sync"
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

	HttpStatusCode int `json:"-"`
}

type V2Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type Client interface {
	Invoke(ctx context.Context, method string, obj interface{}) (V2Result, error)
	InvokeDirect(ctx context.Context, req V2Request) (V2Result, error)
	GetRawResponse(ctx context.Context, req V2Request) (*transport.WrappedResponse, error)

	Close(ctx context.Context)
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

	return ci.InvokeDirect(ctx, req)
}

func (ci *ClientImpl) InvokeDirect(ctx context.Context, req V2Request) (V2Result, error) {
	if resp, err := ci.GetRawResponse(ctx, req); err != nil {
		return V2Result{}, err
	} else if body, err := resp.Body(); err != nil {
		return V2Result{}, err
	} else {
		var rv V2Result
		rv.HttpStatusCode = resp.StatusCode

		err := json.Unmarshal(body, &rv)
		return rv, err
	}
}

func (ci *ClientImpl) Close(ctx context.Context) {
	ci.transport.HttpClient.CloseIdleConnections()
}

func (ci *ClientImpl) GetRawResponse(ctx context.Context, req V2Request) (*transport.WrappedResponse, error) {
	// Implement rate-controls
	time.Sleep(ci.transport.DelayBeforeCall())

	m, _ := ci.transport.Authorizer.QueryStringAuthorization(ctx)
	qs := url.Values{}
	for k, v := range m {
		qs[k] = []string{v}
	}

	if resp, err := ci.transport.Post(ctx, "?"+qs.Encode(), req); err != nil {
		return nil, &errwrap.WrappedError{
			Context: "sending V2 post request",
			Cause:   err,
		}
	} else {
		return resp, err
	}
}

type Params struct {
	transport.HTTPClientParams
	AreaNID        int
	Authorizer     transport.Authorizer
	QPS            int64
	TravelTimeComp time.Duration

	MasheryEndpoint string
}

func (h *Params) FillDefaults() error {
	if h.Authorizer == nil {
		return errors.New("v2 client requires a non-nil Authorizer")
	}
	if len(h.MasheryEndpoint) == 0 {
		if h.AreaNID > 0 {
			h.MasheryEndpoint = fmt.Sprintf("https://api.mashery.com/v2/json-rpc/%d", h.AreaNID)
		} else {
			return errors.New("for an empty MasheryEndpoint, input must supply AreaNID")
		}
	}
	if h.TravelTimeComp == 0 {
		h.TravelTimeComp = time.Millisecond * 147
	}
	if h.QPS <= 0 {
		h.QPS = 2
	}

	if h.Timeout == 0 {
		h.Timeout = time.Second * 60
	}

	return nil
}

// NewHTTPClient Create a new V2 HTTP client to invoke Mashery V2 API
func NewHTTPClient(params Params) Client {
	if err := params.FillDefaults(); err != nil {
		panic(err)
	}

	return &ClientImpl{
		transport: &transport.HttpTransport{
			MashEndpoint: params.MasheryEndpoint,
			Authorizer:   params.Authorizer,

			AvgNetLatency: params.TravelTimeComp,

			HttpClient: params.CreateClient(),
			Mutex:      &sync.Mutex{},
			MaxQPS:     params.QPS,
		}}
}
