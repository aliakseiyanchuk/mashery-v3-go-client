package transport

import (
	"fmt"
	"net/url"
)

// -----------------------------------
// Generic operations

// ResponseParserFunc Function that parses responses returned by JSON.
type ResponseParserFunc func(data []byte) (interface{}, int, error)

// FetchSpec Operation context
type FetchSpec struct {
	Pagination     PaginationType
	Resource       string
	Query          url.Values
	AppContext     string
	ResponseParser ResponseParserFunc
	Return404AsNil bool
}

type PaginationType int

const (
	PerPage PaginationType = iota
	PerItem
	NotRequired
)

type AsyncFetchResult struct {
	Data *WrappedResponse
	Err  error
}

// DestResource Resource that need to be called on the server. This method will return the resource and
// will append the query string, if specified
func (ctx *FetchSpec) DestResource() string {
	if ctx.Query == nil {
		return ctx.Resource
	} else {
		return fmt.Sprintf("%s?%s", ctx.Resource, ctx.Query.Encode())
	}
}

// WithQuery Append extra query parameters to the parent context
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
