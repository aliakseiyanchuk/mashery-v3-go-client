package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpTransport struct {
	MashEndpoint  string
	Authorizer    Authorizer
	AvgNetLatency time.Duration
	Sem           *semaphore.Weighted
	HttpClient    *http.Client
}

func (c *HttpTransport) Fetch(ctx context.Context, res string) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", c.MashEndpoint, res)

	if req, err := http.NewRequest("GET", uri, nil); err != nil {
		return nil, err
	} else {
		return c.httpExec(ctx, req)
	}
}

func (c *HttpTransport) Delete(ctx context.Context, res string) (*http.Response, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%s", c.MashEndpoint, res), nil)
	return c.httpExec(ctx, req)
}

func (c *HttpTransport) Post(ctx context.Context, res string, body interface{}) (*http.Response, error) {
	return c.Send(ctx, "POST", res, body)
}

func (c *HttpTransport) Put(ctx context.Context, res string, body interface{}) (*http.Response, error) {
	return c.Send(ctx, "PUT", res, body)
}

func (c *HttpTransport) Send(ctx context.Context, meth string, res string, body interface{}) (*http.Response, error) {
	if dat, err := json.Marshal(body); err == nil {
		req, _ := http.NewRequest(meth, fmt.Sprintf("%s%s", c.MashEndpoint, res), bytes.NewReader(dat))

		// With the client, only JSON is sent up and down.
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		rv, rvErr := c.httpExec(ctx, req)
		_ = req.Body.Close()

		return rv, rvErr
	} else {
		return nil, err
	}
}

func (c *HttpTransport) httpExec(ctx context.Context, req *http.Request) (*http.Response, error) {
	// TODO: add check for the cancelled context

	var lastErr error

	for i := 0; i < 10; i++ {
		err := c.Sem.Acquire(ctx, 1)

		if err != nil {
			return nil, err
		} else {
			go c.releaseSemaphoreLater()
		}

		if c.Authorizer != nil {
			if tkn, err := c.Authorizer.Authorization(); err != nil {
				return nil, err
			} else if len(tkn) > 0 {
				for k, v := range tkn {
					req.Header.Add(k, v)
				}
			}
		}

		resp, lastErr := c.HttpClient.Do(req)

		// If, for whatever reason, the request still gets over QPS, re-try with progressive
		// back-offs could be tried.
		if lastErr == nil && resp.StatusCode == 403 {
			if str := resp.Header.Get("X-Mashery-Error-Code"); str == "ERR_403_DEVELOPER_OVER_QPS" {
				d := time.Duration(1+i) * time.Second
				time.Sleep(d)
				continue
			}
		}

		// Where the response is successful or cannot be re-tried, the both
		// are returned to the caller
		return resp, lastErr
	}

	return nil, lastErr
}

func (c *HttpTransport) releaseSemaphoreLater() {
	time.Sleep(time.Second + c.AvgNetLatency)
	c.Sem.Release(1)
}

// ReadResponseBody Reads the response body of the response
func ReadResponseBody(r *http.Response) ([]byte, error) {
	if r.Body != nil {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		return b, err
	} else {
		return []byte{}, nil
	}
}
