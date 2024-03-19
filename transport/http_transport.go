package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type MiddlewareFuncPipeline struct {
	f     MiddlewareFunc
	order int
}
type HttpTransport struct {
	MashEndpoint  string
	Authorizer    Authorizer
	AvgNetLatency time.Duration
	HttpExecutor  HttpExecutor

	PlannedSecond  int64
	AllocatedCalls int64
	MaxQPS         int64

	Mutex *sync.Mutex

	ExchangeListener ExchangeListener
	Pipeline         MiddlewareFunc
}

func (c *HttpTransport) DelayBeforeCall() time.Duration {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// A call to Mashery will be received at this time
	nextServTime := time.Now().Add(c.AvgNetLatency)
	nextServSecond := nextServTime.Unix()

	if nextServSecond > c.PlannedSecond {
		c.PlannedSecond = nextServSecond
		c.AllocatedCalls = 1
		return time.Duration(0)
	} else if nextServSecond == c.PlannedSecond && c.AllocatedCalls < c.MaxQPS {
		c.AllocatedCalls++
		return time.Duration(0)
	} else {
		wait := c.PlannedSecond - nextServSecond
		if c.AllocatedCalls < c.MaxQPS {
			c.AllocatedCalls++
		} else {
			wait++
			c.PlannedSecond++
			c.AllocatedCalls = 1
		}
		return time.Second * time.Duration(wait)
	}
}

func (c *HttpTransport) Fetch(ctx context.Context, res string) (*WrappedResponse, error) {
	uri := fmt.Sprintf("%s%s", c.MashEndpoint, res)

	if req, err := http.NewRequest("GET", uri, nil); err != nil {
		return nil, err
	} else {
		return c.httpExec(ctx, &WrappedRequest{Request: req})
	}
}

func (c *HttpTransport) Delete(ctx context.Context, res string) (*WrappedResponse, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%s", c.MashEndpoint, res), nil)
	return c.httpExec(ctx, &WrappedRequest{Request: req})
}

func (c *HttpTransport) Post(ctx context.Context, res string, body interface{}) (*WrappedResponse, error) {
	return c.Send(ctx, "POST", res, body)
}

func (c *HttpTransport) Put(ctx context.Context, res string, body interface{}) (*WrappedResponse, error) {
	return c.Send(ctx, "PUT", res, body)
}

func (c *HttpTransport) Send(ctx context.Context, meth string, res string, body interface{}) (*WrappedResponse, error) {
	if dat, err := json.Marshal(body); err == nil {
		req, reqCreteErr := http.NewRequest(meth, fmt.Sprintf("%s%s", c.MashEndpoint, res), bytes.NewReader(dat))
		if reqCreteErr != nil {
			return nil, reqCreteErr
		}

		defer req.Body.Close()
		// With the client, only JSON is sent up and down.
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		wr := &WrappedRequest{
			Request: req,
			Body:    body,
		}

		rv, rvErr := c.httpExec(ctx, wr)

		return rv, rvErr
	} else {
		return nil, err
	}
}

func (c *HttpTransport) httpExec(ctx context.Context, wrq *WrappedRequest) (*WrappedResponse, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if c.Authorizer != nil {
		if tkn, err := c.Authorizer.HeaderAuthorization(ctx); err != nil {
			return nil, err
		} else if len(tkn) > 0 {
			for k, v := range tkn {
				wrq.Request.Header.Add(k, v)
			}
		}
	}

	var wrs *WrappedResponse
	resp, lastErr := c.HttpExecutor.Do(wrq.Request)
	if lastErr == nil {
		wrs = &WrappedResponse{
			Request:    wrq,
			Response:   resp,
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
		}
	}

	if c.ExchangeListener != nil {
		c.ExchangeListener(ctx, wrq, wrs, lastErr)
	}

	// Where the response is successful or cannot be re-tried, the both
	// are returned to the caller
	return wrs, lastErr
}

// ReadResponseBody Reads the response body of the response
func ReadResponseBody(r *http.Response) ([]byte, error) {
	if r.Body != nil {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		return b, err
	} else {
		return []byte{}, nil
	}
}
