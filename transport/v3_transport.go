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
)

type MiddlewareFunc func(ctx context.Context, transport *HttpTransport) (*WrappedResponse, error)

// ChainedMiddlewareFunc is a middleware function
type ChainedMiddlewareFunc func(ctx context.Context, transport *HttpTransport, middlewareFunc MiddlewareFunc) (*WrappedResponse, error)

func GetObject[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) (T, bool, error) {
	rv, resp, err := performGenericObjectCRUDWithResponse[T](ctx, c, opCtx, opCtx.FetchFunc())
	return rv, resp.StatusCode == 200, err
}

func CreateObject[T any](ctx context.Context, opCtx ObjectUpsertSpec[T], c *HttpTransport) (T, error) {
	return performGenericObjectCRUD(ctx, c, opCtx.ObjectFetchSpec, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Post(ctx, opCtx.DestResource(), opCtx.Upsert)
	})
}

func UpdateObject[T any](ctx context.Context, opCtx ObjectUpsertSpec[T], c *HttpTransport) (T, error) {
	return performGenericObjectCRUD(ctx, c, opCtx.ObjectFetchSpec, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Put(ctx, opCtx.DestResource(), opCtx.Upsert)
	})
}

// ExchangeObject exchange an input object for the output one.
func ExchangeObject[TIn, TOut any](ctx context.Context, opCtx ObjectExchangeSpec[TIn, TOut], verb string, c *HttpTransport) (TOut, error) {
	return performGenericObjectCRUD(ctx, c, opCtx.ObjectFetchSpec, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Send(ctx, verb, opCtx.DestResource(), opCtx.Body)
	})
}

func DeleteObject[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) error {
	_, err := performGenericObjectCRUD(ctx, c, opCtx, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
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

	if _, wr, err := performGenericObjectCRUDWithResponse(ctx, c, builder.Build().AsObjectFetchSpec(), func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	}); err != nil {
		return -1, err
	} else {
		return extractTotalCount(wr), nil
	}
}

func Exists[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport) (bool, error) {
	if _, wr, err := performGenericObjectCRUDWithResponse(ctx, c, opCtx, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Fetch(ctx, opCtx.DestResource())
	}); err != nil {
		return false, err
	} else {
		return wr.StatusCode == 200, err
	}
}

// AsyncFetch Perform a fetch asynchronously, returning the response in the provided channel.
func AsyncFetch[T any](ctx context.Context, opCtx ObjectFetchSpec[T], c *HttpTransport, comm chan AsyncFetchResult[T]) {
	rv, err := performGenericObjectCRUD(ctx, c, opCtx, func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
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

// FetchAllWithExists Fetch all Mashery objects, including the handling for the pagination
func FetchAllWithExists[T any](ctx context.Context, opCtx ObjectListFetchSpec[T], c *HttpTransport) ([]T, bool, error) {

	firstPageData, firstPageResponse, firstPageFetchErr := performGenericObjectCRUDWithResponse[[]T](ctx, c, opCtx.AsObjectFetchSpec(), func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
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
	//fmt.Println(totalCountHdr)

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
				//fmt.Println(offsetParam)

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

func performGenericObjectCRUD[T any](ctx context.Context, c *HttpTransport, opCtx ObjectFetchSpec[T], f MiddlewareFunc) (T, error) {
	rv, _, err := performGenericObjectCRUDWithResponse(ctx, c, opCtx, f)
	return rv, err
}

// BuildPipeline builds a call execution pipeline from the supplied chained middleware functions.
func BuildPipeline(mf MiddlewareFunc, funcs []ChainedMiddlewareFunc) MiddlewareFunc {
	rv := mf
	for _, cmf := range funcs {
		rv = wrapMiddleware(cmf, rv)
	}

	return rv
}

func wrapMiddleware(p1 ChainedMiddlewareFunc, mf MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, transport *HttpTransport) (*WrappedResponse, error) {
		return p1(ctx, transport, mf)
	}
}

func executeCallPipeline(ctx context.Context, c *HttpTransport, execFunc MiddlewareFunc) (*WrappedResponse, error) {
	cCtx := context.WithValue(ctx, LeafExecutor, execFunc)

	return c.Pipeline(cCtx, c)
}

func performGenericObjectCRUDWithResponse[T any](entryCtx context.Context, c *HttpTransport, opCtx ObjectFetchSpec[T], f MiddlewareFunc) (T, *WrappedResponse, error) {
	ctx := entryCtx
	if !opCtx.Return404AsNil {
		ctx = context.WithValue(ctx, SendErrorOn404, true)
	}

	if wr, err := executeCallPipeline(ctx, c, f); err != nil {
		return opCtx.ValueFactory(), wr, err
	} else {
		rv := opCtx.ValueFactory()

		if wr.StatusCode == 200 && !opCtx.IgnoreResponse {
			if jsonErr := json.Unmarshal(wr.MustBody(), &rv); jsonErr != nil {
				fmt.Println(string(wr.MustBody()))

				return rv, wr, &errwrap.WrappedError{
					Context: fmt.Sprintf("%s %s->unmarshal response", wr.Request.Request.Method, opCtx.AppContext),
					Cause:   jsonErr,
				}
			}
		}

		return rv, wr, err
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

func v3BasicError(wr *WrappedResponse) error {
	// Did we receive a generic error?
	var rv masherytypes.V3GenericErrorResponse
	if err := json.Unmarshal(wr.MustBody(), &rv); err == nil && rv.HasData() {
		return &rv
	}

	// Did we receive at least one error?
	var propRv masherytypes.V3PropertyErrorMessages
	if err := json.Unmarshal(wr.MustBody(), &propRv); err == nil && len(propRv.Errors) > 0 {
		return &propRv
	}

	// The error is not really know; so the output would be printed in the output

	return &masherytypes.V3UndeterminedError{
		Code:   wr.StatusCode,
		Header: wr.Header,
		Body:   wr.MustBody(),
	}
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
