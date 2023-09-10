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
	"time"
)

type V3Transport struct {
	HttpTransport
}

// AsyncFetch Perform a fetch asynchronously, returning the response in the provided channel.
func (c *V3Transport) AsyncFetch(ctx context.Context, opContext FetchSpec, comm chan AsyncFetchResult[int]) {
	//rv, err := c.Fetch(ctx, opContext.DestResource())

	// Send the communication back.
	//comm <- AsyncFetchResult{
	//	Data: rv,
	//	Err:  err,
	//}
}

func GetObject[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) (T, bool, error) {
	rv, resp, err := performGenericObjectCRUDWithResponse(ctx, opCtx, "get", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	})
	return rv, resp.StatusCode == 200, err
}

func CreateObject[T any](ctx context.Context, opCtx ObjectUpsertSpec[T], c *HttpTransport) (T, error) {
	return performGenericObjectCRUD(ctx, opCtx.ObjectFetchSpec, "post", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Post(ctx, opCtx.DestResource(), opCtx.Upsert)
	})
}

func UpdateObject[T any](ctx context.Context, opCtx ObjectUpsertSpec[T], c *HttpTransport) (T, error) {
	return performGenericObjectCRUD(ctx, opCtx.ObjectFetchSpec, "put", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Put(ctx, opCtx.DestResource(), opCtx.Upsert)
	})
}

// ExchangeObject exchange an input object for the output one.
func ExchangeObject[TIn, TOut any](ctx context.Context, opCtx ObjectExchangeSpec[TIn, TOut], verb string, c *HttpTransport) (TOut, error) {
	return performGenericObjectCRUD(ctx, opCtx.ObjectFetchSpec, verb, func(ctx context.Context) (*WrappedResponse, error) {
		return c.Send(ctx, verb, opCtx.DestResource(), opCtx.Body)
	})
}

func DeleteObject[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) error {
	_, err := performGenericObjectCRUD(ctx, opCtx, "delete", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Delete(ctx, opCtx.DestResource())
	})

	return err
}

// Count the number of objects that match the specified criteria
func Count[T any](ctx context.Context, opCtx ObjectListFetchSpec[T], c *HttpTransport) (int64, error) {
	limitQuery := url.Values{
		"limit": []string{"1"},
	}

	builder := opCtx.ToBuilder()
	builder.WithMergedQuery(limitQuery)

	if _, wr, err := performGenericObjectCRUDWithResponse(ctx, builder.Build().AsObjectFetchSpec(), "get", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	}); err != nil {
		return -1, err
	} else {
		return extractTotalCount(wr), nil
	}
}

func Exists[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) (bool, error) {
	if _, wr, err := performGenericObjectCRUDWithResponse(ctx, opCtx, "get", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	}); err != nil {
		return false, err
	} else {
		return wr.StatusCode == 200, err
	}
}

// AsyncFetch Perform a fetch asynchronously, returning the response in the provided channel.
func AsyncFetch[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport, comm chan AsyncFetchResult[T]) {
	rv, err := performGenericObjectCRUD(ctx, opCtx, "get", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	})

	// Send the communication back.
	comm <- AsyncFetchResult[T]{
		Data: rv,
		Err:  err,
	}
}

func FetchAll[T any](ctx context.Context, opCtx ObjectListFetchSpec[T], c *HttpTransport) ([]T, error) {
	rv, _, err := FetchAllWithExists(ctx, opCtx, c)
	return rv, err
}

// FetchAll Fetch all Mashery objects, including the handling for the pagination
func FetchAllWithExists[T any](ctx context.Context, opCtx ObjectListFetchSpec[T], c *HttpTransport) ([]T, bool, error) {

	firstPageData, firstPageResponse, firstPageFetchErr := performGenericObjectCRUDWithResponse[[]T](ctx, opCtx.AsObjectFetchSpec(), "get", func(ctx context.Context) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	})

	if firstPageFetchErr != nil {
		return nil, false, firstPageFetchErr
	} else if firstPageData == nil {
		return nil, firstPageResponse.StatusCode != 404, nil
	}

	rv := firstPageData
	var collErr error

	pageSize := len(rv)
	// Don't try doing anything else if the response contains no data.
	if pageSize == 0 {
		return rv, firstPageResponse.StatusCode != 404, nil
	}

	totalCountHdr := firstPageResponse.Header.Get("X-Total-Count")
	if len(totalCountHdr) > 0 {
		totalCount, _ := strconv.ParseInt(totalCountHdr, 10, 0)

		if totalCount > int64(pageSize) {
			allFetches := int(totalCount / int64(pageSize))

			commChan := make(chan AsyncFetchResult[[]T])
			defer close(commChan)

			for p := 1; p <= allFetches; p++ {
				offset := p
				if opCtx.Pagination == PerItem {
					offset *= pageSize
				}

				offsetParam := url.Values{
					"offset": {strconv.Itoa(offset)},
				}

				pageFetchSpec := opCtx.ToBuilder()
				pageFetchSpec.WithMergedQuery(offsetParam)

				go AsyncFetch(ctx, pageFetchSpec.Build().AsObjectFetchSpec(), c, commChan)
			}

			for p := 1; p <= allFetches; p++ {
				asyncRead := <-commChan
				if asyncRead.Err != nil {
					collErr = asyncRead.Err
					// TODO: if error occurred, we might need to terminate the rest
					// of the fetching operations.
				} else {
					if asyncRead.Data == nil {
						collErr = &errwrap.WrappedError{
							Context: fmt.Sprintf("fetch all %s->read async response", opCtx.AppContext),
							Cause:   errors.New("nil response received"),
						}
						// TODO Terminate making any remaining calls to the server
					} else {
						rv = append(rv, asyncRead.Data...)
					}
				}
			}
		}
	}

	return rv, true, collErr
}

type CallFunc func(ctx context.Context) (*WrappedResponse, error)

func performGenericObjectCRUD[T any](ctx context.Context, opCtx ObjectFetchSpec[T], verb string, f CallFunc) (T, error) {
	rv, _, err := performGenericObjectCRUDWithResponse(ctx, opCtx, verb, f)
	return rv, err
}

func performGenericObjectCRUDWithResponse[T any](ctx context.Context, opCtx ObjectFetchSpec[T], verb string, f CallFunc) (T, *WrappedResponse, error) {
	for i := 0; i < 10; i++ {
		if resp, err := f(ctx); err == nil {
			if dat, err := resp.Body(); err != nil {
				return opCtx.ValueFactory(), resp, &errwrap.WrappedError{
					Context: fmt.Sprintf("%s %s->read server response", verb, opCtx.AppContext),
					Cause:   err,
				}
			} else {
				if resp.StatusCode == 200 {
					// In specific situations (e.g. with the DELETE operation) the body of the response
					// will be ignored, as the server is not supplying any useful information in the body.
					// The calling code would specify a primitive (e.g. int) value that will be ignored by
					// the calling code.
					rv := opCtx.ValueFactory()

					if !opCtx.IgnoreResponse {
						if jsonErr := json.Unmarshal(dat, &rv); jsonErr != nil {
							return rv, resp, &errwrap.WrappedError{
								Context: fmt.Sprintf("%s %s->unmarshal response", verb, opCtx.AppContext),
								Cause:   jsonErr,
							}
						}
					}

					return rv, resp, nil
				} else if resp.StatusCode == 403 {
					if str := resp.Header.Get("X-Mashery-Error-Code"); str == "ERR_403_DEVELOPER_OVER_QPS" {
						d := time.Duration(1+i) * time.Second
						time.Sleep(d)
						continue
					} else {
						return opCtx.ValueFactory(), resp, &errwrap.WrappedError{
							Context: fmt.Sprintf("%s %s->not authorized (http code 403)", verb, opCtx.AppContext),
							Cause:   errors.New("call denied by server"),
						}
					}
				} else if resp.StatusCode == 404 {
					if opCtx.Return404AsNil {
						return opCtx.ValueFactory(), resp, nil
					} else {
						return opCtx.ValueFactory(), resp, &errwrap.WrappedError{
							Context: fmt.Sprintf("%s %s->no such resource", verb, opCtx.AppContext),
							Cause:   errors.New("error code 404 is not an expected response to this request"),
						}
					}
				} else {
					return opCtx.ValueFactory(), resp, v3ErrorFromResponse(fmt.Sprintf("%s %s", verb, opCtx.AppContext), resp.StatusCode, resp.Header, dat)
				}
			}
		} else {
			return opCtx.ValueFactory(), nil, err
		}
	}

	return opCtx.ValueFactory(), nil, &errwrap.WrappedError{
		Context: fmt.Sprintf("%s %s->unsatisfiable operation", verb, opCtx.AppContext),
		Cause:   errors.New("operation unsuccessful after all available retries"),
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
				if opCtx.ResponseParser == nil {
					return nil, nil
				}

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
func (c *HttpTransport) Count(ctx context.Context, opCtx CommonFetchSpec) (int64, error) {
	builder := opCtx.ToBuilder()
	countSpec := builder.WithQuery(url.Values{
		"limit": {"1"},
	}).Build()

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

	//
	return nil, nil
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
