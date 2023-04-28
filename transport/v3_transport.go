package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type V3Transport struct {
	HttpTransport
}

// AsyncFetch Perform a fetch asynchronously, returning the response in the provided channel.
func (c *V3Transport) AsyncFetch(ctx context.Context, opContext FetchSpec, comm chan AsyncFetchResult) {
	rv, err := c.Fetch(ctx, opContext.DestResource())

	// Send the communication back.
	comm <- AsyncFetchResult{
		Data: rv,
		Err:  err,
	}
}

func (c *V3Transport) GetObject(ctx context.Context, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.Fetch(ctx, opCtx.DestResource()); err == nil {
		if dat, err := resp.Body(); err != nil {
			return nil, &errwrap.WrappedError{
				Context: fmt.Sprintf("get %s->read server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			if resp.StatusCode == 200 {
				// Ignore page when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &errwrap.WrappedError{
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

func (c *V3Transport) DeleteObject(ctx context.Context, opCtx FetchSpec) error {
	if resp, err := c.Delete(ctx, opCtx.Resource); err == nil {
		if resp.StatusCode == 200 {
			return nil
		} else {
			return errors.New(fmt.Sprintf("delete %s->response code %d", opCtx.AppContext, resp.StatusCode))
		}
	} else {
		return &errwrap.WrappedError{
			Context: fmt.Sprintf("delete %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// CreateObject Create a new service.
func (c *V3Transport) CreateObject(ctx context.Context, objIn interface{}, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.Post(ctx, opCtx.DestResource(), objIn); err == nil {
		if dat, err := resp.Body(); err != nil {
			return nil, &errwrap.WrappedError{
				Context: fmt.Sprintf("create %s->read server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			if resp.StatusCode == 200 {
				// Ignore page size when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &errwrap.WrappedError{
						Context: fmt.Sprintf("create %s->unmarshal response (%s)", opCtx.AppContext, dat),
						Cause:   jsonErr,
					}
				} else {
					return rv, nil
				}
			} else {
				return nil, v3ErrorFromResponse(fmt.Sprintf("create %s", opCtx.AppContext), resp.StatusCode, resp.Header, dat)
			}
		}
	} else {
		return nil, &errwrap.WrappedError{
			Context: fmt.Sprintf("create %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// UpdateObject Update existing object
func (c *V3Transport) UpdateObject(ctx context.Context, objIn interface{}, opCtx FetchSpec) (interface{}, error) {
	if resp, err := c.Put(ctx, opCtx.DestResource(), objIn); err == nil {
		if dat, err := resp.Body(); err != nil {
			return nil, &errwrap.WrappedError{
				Context: fmt.Sprintf("update %s->read server response", opCtx.AppContext),
				Cause:   err,
			}

		} else {
			if resp.StatusCode == 200 {
				// Ignoring page size when retrieving an object
				if rv, _, jsonErr := opCtx.ResponseParser(dat); jsonErr != nil {
					return nil, &errwrap.WrappedError{
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
		return nil, &errwrap.WrappedError{
			Context: fmt.Sprintf("update %s->connect", opCtx.AppContext),
			Cause:   err,
		}
	}
}

// Count the number of objects that match the specified criteria
func (c *V3Transport) Count(ctx context.Context, opCtx FetchSpec) (int64, error) {
	countSpec := opCtx.WithQuery(url.Values{
		"limit": {"1"},
	})

	if cnt, err := c.Fetch(ctx, countSpec.DestResource()); err != nil {
		return -1, &errwrap.WrappedError{
			Context: fmt.Sprintf("count %s->fetch count", countSpec.AppContext),
			Cause:   err,
		}
	} else {
		return extractTotalCount(cnt), nil
	}

}

// FetchAll Fetch all Mashery objects, including the handling for the pagination
func (c *V3Transport) FetchAll(ctx context.Context, opCtx FetchSpec) ([]interface{}, error) {

	firstPage, err := c.Fetch(ctx, opCtx.DestResource())
	if err != nil {
		return nil, &errwrap.WrappedError{
			Context: fmt.Sprintf("fetch all %s->fetch first page", opCtx.AppContext),
			Cause:   err,
		}
	}

	if firstPage.StatusCode == 200 {
		if dat, err := firstPage.Body(); err != nil {
			return nil, &errwrap.WrappedError{
				Context: fmt.Sprintf("fetch all %s->read first page server response", opCtx.AppContext),
				Cause:   err,
			}
		} else {
			fp, pageSize, err := opCtx.ResponseParser(dat)
			if err != nil {
				return nil, &errwrap.WrappedError{
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

						go c.AsyncFetch(ctx, opCtx.WithQuery(qs), commChan)
					}

					for p := 1; p <= allFetches; p++ {
						asyncRead := <-commChan
						if asyncRead.Err != nil {
							collErr = asyncRead.Err
							// TODO: if error occurred, we might need to terminate the rest
							// of the fetching operations.
						} else {
							if pageDat, pageReadErr := asyncRead.Data.Body(); pageReadErr != nil {
								collErr = &errwrap.WrappedError{
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
	} else if firstPage.StatusCode == 404 && opCtx.Return404AsNil {
		return nil, nil
	}

	return nil, &errwrap.WrappedError{
		Context: fmt.Sprintf("fetchAll %s->fetch first page->response", opCtx.AppContext),
		Cause:   errors.New(fmt.Sprintf("received status code %d", firstPage.StatusCode)),
	}
}

func (c *V3Transport) V3FilteringParams(params map[string]string, fields []string) url.Values {
	qs := url.Values{}
	if len(params) > 0 {
		qs["filter"] = []string{V3FilterExpression(params)}
	}

	if len(fields) > 0 {
		qs["fields"] = []string{strings.Join(fields, ",")}
	}

	return qs
}

func v3ErrorFromResponse(context string, code int, headers http.Header, data []byte) error {
	uCtx := fmt.Sprintf("%s->api call", context)

	// Did we receive a generic error?
	var rv masherytypes.V3GenericErrorResponse
	if err := json.Unmarshal(data, &rv); err == nil && rv.HasData() {
		return &errwrap.WrappedError{
			Context: uCtx,
			Cause:   &rv,
		}
	}

	// Did we receive at least one error?
	var propRv masherytypes.V3PropertyErrorMessages
	if err := json.Unmarshal(data, &propRv); err == nil && len(propRv.Errors) > 0 {
		return &errwrap.WrappedError{
			Context: uCtx,
			Cause:   &propRv,
		}
	}

	// The error is not really know; so the output would be printed in the output
	return &errwrap.WrappedError{
		Context: uCtx,
		Cause: &masherytypes.V3UndeterminedError{
			Code:   code,
			Header: headers,
			Body:   data,
		},
	}
}

// Extract Mashery-supplied total count of elements from this response
func extractTotalCount(resp *WrappedResponse) int64 {
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

// V3FilterExpression Converts query parameters to V3-required filter string.
func V3FilterExpression(params map[string]string) string {
	filterTokens := make([]string, len(params))
	idx := 0
	for k, v := range params {
		filterTokens[idx] = fmt.Sprintf("%s:%s", k, v)
		idx++
	}
	return strings.Join(filterTokens, ",")
}
