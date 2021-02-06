package v3client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const tokenFile string = ".mashery-logon"

// Access token provider that supplies the access token, depending on the strategy.
// There are three strategies:
// - FixedTokenProvider yields a fixed token. This method is useful for short deployments where an access
// token is obtained by an outside process and would be stored e.g. in-memory.
// - FileSystemTokenProvider yields a token that that was previously saved in the file system, e.g. using the `mash-login`
// command
// - Both these methods have limited applicability time-span of 1 hour, since Mashery V3 token would expire after 1
// hour, and repeated logon would be necessary.
// - ClientCredentialsProvider can support operations of exceeding 1 hour by using Mashery V3 API to retrieve and refresh
// the access token.
//
// The calling code has to pick an appropriate provider depending on the context.
type V3AccessTokenProvider interface {
	// Yields an access token to be used in the next API call to Mashery
	AccessToken() (string, error)
}

type HttpTransport struct {
	mashEndpoint  string
	tokenProvider V3AccessTokenProvider
	avgNetLatency time.Duration
	sem           *semaphore.Weighted
	httpCl        *http.Client
}

func NewHttpClient(p V3AccessTokenProvider, qps int64, travelTimeComp time.Duration) Client {
	impl := HttpTransport{
		mashEndpoint:  "https://api.mashery.com/v3/rest",
		tokenProvider: p,
		sem:           semaphore.NewWeighted(qps),
		httpCl:        &http.Client{},
		avgNetLatency: travelTimeComp,
	}

	rv := PluggableClient{
		schema: ClientMethodSchema{
			GetApplicationContext:     GetApplication,
			GetApplicationPackageKeys: GetApplicationPackageKeys,
		},
		transport: &impl,
	}

	return &rv
}

func (c *HttpTransport) fetch(ctx context.Context, res string) (*http.Response, error) {
	get := fmt.Sprintf("%s%s", c.mashEndpoint, res)

	req, _ := http.NewRequest("GET", get, nil)
	return c.httpExec(ctx, req)
}

func (c *HttpTransport) delete(ctx context.Context, res string) (*http.Response, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%s", c.mashEndpoint, res), nil)
	return c.httpExec(ctx, req)
}

func (c *HttpTransport) post(ctx context.Context, res string, body interface{}) (*http.Response, error) {
	return c.send(ctx, "POST", res, body)
}

func (c *HttpTransport) put(ctx context.Context, res string, body interface{}) (*http.Response, error) {
	return c.send(ctx, "PUT", res, body)
}

func (c *HttpTransport) send(ctx context.Context, meth string, res string, body interface{}) (*http.Response, error) {
	if dat, err := json.Marshal(body); err == nil {
		req, _ := http.NewRequest(meth, fmt.Sprintf("%s%s", c.mashEndpoint, res), bytes.NewReader(dat))

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

func readResponseBody(r *http.Response) ([]byte, error) {
	if r.Body != nil {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		return b, err
	} else {
		return []byte{}, nil
	}
}

func (c *HttpTransport) v3FilteringParams(params map[string]string, fields []string) url.Values {
	qs := url.Values{}
	if len(params) > 0 {
		qs["filter"] = []string{toV3FilterExpression(params)}
	}

	if len(fields) > 0 {
		qs["fields"] = []string{strings.Join(fields, ",")}
	}

	return qs
}

// TODO: Need to define the method for collectAll

func (c *HttpTransport) httpExec(ctx context.Context, req *http.Request) (*http.Response, error) {
	// TODO: add check for the cancelled context

	var lastErr error

	for i := 0; i < 10; i++ {
		err := c.sem.Acquire(ctx, 1)

		if err != nil {
			return nil, err
		} else {
			go c.releaseSemaphoreLater()
		}

		tkn, err := c.tokenProvider.AccessToken()
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

		resp, lastErr := c.httpCl.Do(req)

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
		// are returned back to the caller
		return resp, lastErr
	}

	return nil, lastErr
}

func (c *HttpTransport) releaseSemaphoreLater() {
	time.Sleep(time.Second + c.avgNetLatency)
	c.sem.Release(1)
}

type AsyncFetchResult struct {
	Data *http.Response
	Err  error
}

type WrappedError struct {
	Context string
	Cause   error
}

func (w *WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", w.Context, w.Cause)
}

func (m *WrappedError) Unwrap() error {
	return m.Cause
}

func v3ErrorFromResponse(context string, code int, headers http.Header, data []byte) error {
	uCtx := fmt.Sprintf("%s->api call", context)

	// Did we receive a generic error?
	var rv V3GenericErrorResponse
	if err := json.Unmarshal(data, &rv); err == nil && rv.hasData() {
		return &WrappedError{
			Context: uCtx,
			Cause:   &rv,
		}
	}

	// Did we receive at least one error?
	var propRv V3PropertyErrorMessages
	if err := json.Unmarshal(data, &propRv); err == nil && len(propRv.Errors) > 0 {
		return &WrappedError{
			Context: uCtx,
			Cause:   &propRv,
		}
	}

	// The error is not really know; so the output would be printed in the output
	return &WrappedError{
		Context: uCtx,
		Cause: &V3UndeterminedError{
			Code:   code,
			Header: headers,
			Body:   data,
		},
	}
}

// -----------------------------------
// Generic operations

// Function that parses responses returned by JSON.
type ResponseParserFunc func(data []byte) (interface{}, int, error)

// Operation context
type FetchSpec struct {
	Pagination     PaginationType
	Resource       string
	Query          url.Values
	AppContext     string
	ResponseParser ResponseParserFunc
}

// Resource that need to be called on the server. This method will return the resource and
// will append the query string, if specified
func (ctx *FetchSpec) DestResource() string {
	if ctx.Query == nil {
		return ctx.Resource
	} else {
		return fmt.Sprintf("%s?%s", ctx.Resource, ctx.Query.Encode())
	}
}

// Append extra query parameters to the parent context
func (ctx *FetchSpec) WithQuery(qs url.Values) FetchSpec {
	return FetchSpec{
		Pagination:     ctx.Pagination,
		Resource:       ctx.Resource,
		Query:          merge(ctx.Query, qs),
		AppContext:     ctx.AppContext,
		ResponseParser: ctx.ResponseParser,
	}
}

func merge(qs ...url.Values) url.Values {
	rv := url.Values{}

	for _, q := range qs {
		for k, v := range q {
			rv[k] = v
		}
	}

	return rv
}

type PaginationType int

const (
	PerPage PaginationType = iota
	PerItem
	NotRequired
)

// Perform a fetch asynchronously, returning the response in the provided channel.
func (c *HttpTransport) asyncFetch(ctx context.Context, opContext FetchSpec, comm chan AsyncFetchResult) {
	rv, err := c.fetch(ctx, opContext.DestResource())

	// Send the communication back.
	comm <- AsyncFetchResult{
		Data: rv,
		Err:  err,
	}
}

func (c *HttpTransport) getObject(ctx context.Context, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.fetch(ctx, opCtx.DestResource()); err == nil {
		if dat, err := readResponseBody(resp); err != nil {
			return nil, &WrappedError{
				Context: fmt.Sprintf("get %s->read server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			if resp.StatusCode == 200 {
				// Ignore page when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &WrappedError{
						Context: fmt.Sprintf("get %s->unmarshal response", opCtx.AppContext),
						Cause:   jsonErr,
					}
				} else {
					return rv, nil
				}
			} else if resp.StatusCode == 404 {
				return nil, nil
			} else {
				return nil, v3ErrorFromResponse(fmt.Sprintf("get %s", opCtx.AppContext), resp.StatusCode, resp.Header, dat)
			}
		}
	} else {
		return nil, err
	}
}

func (c *HttpTransport) deleteObject(ctx context.Context, opCtx FetchSpec) error {
	if resp, err := c.delete(ctx, opCtx.Resource); err == nil {
		if resp.StatusCode == 200 {
			return nil
		} else {
			return errors.New(fmt.Sprintf("delete %s->response code %d", opCtx.AppContext, resp.StatusCode))
		}
	} else {
		return &WrappedError{
			Context: fmt.Sprintf("delete %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// Create a new service.
func (c *HttpTransport) createObject(ctx context.Context, objIn interface{}, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.post(ctx, opCtx.DestResource(), objIn); err == nil {
		if dat, err := readResponseBody(resp); err != nil {
			return nil, &WrappedError{
				Context: fmt.Sprintf("create %s->read server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			if resp.StatusCode == 200 {
				// Ignore page size when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &WrappedError{
						Context: fmt.Sprintf("create %s->unmarshal response", opCtx.AppContext),
						Cause:   err,
					}
				} else {
					return rv, nil
				}
			} else {
				return nil, v3ErrorFromResponse(fmt.Sprintf("create %s", opCtx.AppContext), resp.StatusCode, resp.Header, dat)
			}
		}
	} else {
		return nil, &WrappedError{
			Context: fmt.Sprintf("create %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// Update existing object
func (c *HttpTransport) updateObject(ctx context.Context, objIn interface{}, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.put(ctx, opCtx.DestResource(), objIn); err == nil {
		if dat, err := readResponseBody(resp); err != nil {
			return nil, &WrappedError{
				Context: fmt.Sprintf("update %s->read server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			if resp.StatusCode == 200 {
				// Ignoring page size when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &WrappedError{
						Context: fmt.Sprintf("update %s->unmarshal response", opCtx.AppContext),
						Cause:   err,
					}
				} else {
					return &rv, nil
				}
			} else {
				return nil, v3ErrorFromResponse(opCtx.AppContext, resp.StatusCode, resp.Header, dat)
			}
		}
	} else {
		return nil, &WrappedError{
			Context: fmt.Sprintf("update %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// Count the number of objects that match the specified criteria
func (c *HttpTransport) count(ctx context.Context, opCtx FetchSpec) (int64, error) {
	countSpec := opCtx.WithQuery(url.Values{
		"limit": {"1"},
	})

	if cnt, err := c.fetch(ctx, countSpec.DestResource()); err != nil {
		return -1, &WrappedError{
			Context: fmt.Sprintf("count %s->fetch count", countSpec.AppContext),
			Cause:   err,
		}
	} else {
		return extractTotalCount(cnt), nil
	}

}

// Extract Masherty-supplied total count of elements from this response
func extractTotalCount(resp *http.Response) int64 {
	totalCountHdr := resp.Header.Get("X-Total-Count")

	if len(totalCountHdr) > 0 {
		if totalCountI, err := strconv.ParseInt(totalCountHdr, 10, 0); err != nil {
			return -1
		} else {
			return totalCountI
		}
	}

	return 0
}

// Fetch all Mashery objects, including the handling for the pagination
func (c *HttpTransport) fetchAll(ctx context.Context, opCtx FetchSpec) ([]interface{}, error) {

	firstPage, err := c.fetch(ctx, opCtx.DestResource())
	if err != nil {
		return nil, &WrappedError{
			Context: fmt.Sprintf("fetch all %s->fetch first page", opCtx.AppContext),
			Cause:   err,
		}
	}

	if firstPage.StatusCode == 200 {
		if dat, err := readResponseBody(firstPage); err != nil {
			return nil, &WrappedError{
				Context: fmt.Sprintf("fetch all %s->read first page server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			fp, pageSize, err := opCtx.ResponseParser(dat)
			if err != nil {
				return nil, &WrappedError{
					Context: fmt.Sprintf("fetch all %s->unmarshal first page", opCtx.AppContext),
					Cause:   err,
				}
			}

			// Store the first page to be returned
			rv := []interface{}{fp}
			var collErr error

			// Check if reading further pages is necessary
			totalCountHdr := firstPage.Header.Get("X-Total-Count")
			if len(totalCountHdr) > 0 {
				totalCountI, _ := strconv.ParseInt(totalCountHdr, 10, 0)

				totalCount := int(totalCountI)
				if totalCount > pageSize {
					allFetches := totalCount / pageSize

					commChan := make(chan AsyncFetchResult)
					defer close(commChan)

					for p := 1; p <= allFetches; p++ {
						offset := p
						if opCtx.Pagination == PerItem {
							offset *= pageSize
						}

						qs := url.Values{
							"offset": {strconv.Itoa(offset)},
						}

						go c.asyncFetch(ctx, opCtx.WithQuery(qs), commChan)
					}

					for p := 1; p <= allFetches; p++ {
						asyncRead := <-commChan
						if asyncRead.Err != nil {
							collErr = asyncRead.Err
							// TODO: if error occurred, we might need to terminate the rest
							// of the fetching operations.
						} else {
							if pageDat, pageReadErr := readResponseBody(asyncRead.Data); pageReadErr != nil {
								collErr = &WrappedError{
									Context: fmt.Sprintf("fetch all %s->read async response", opCtx.AppContext),
									Cause:   pageReadErr,
								}
							} else {
								fp, _, jsonErr := opCtx.ResponseParser(pageDat)

								if jsonErr != nil {
									collErr = jsonErr
								} else {
									rv = append(rv, fp)
								}
							}
						}
					}
				}
			}

			return rv, collErr
		}
	} else {
		return nil, &WrappedError{
			Context: fmt.Sprintf("fetchAll %s->fetch first page->response", opCtx.AppContext),
			Cause:   errors.New(fmt.Sprintf("received status code %d", firstPage.StatusCode)),
		}
	}
}
