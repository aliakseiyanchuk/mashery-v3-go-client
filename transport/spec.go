package transport

import (
	"context"
	"fmt"
	"net/url"
)

// -----------------------------------
// Generic operations

// ResponseParserFunc Function that parses responses returned by JSON.
type ResponseParserFunc func(data []byte) (interface{}, int, error)

type Supplier[T any] func() T
type ContextSupplier[T any] func(ctx context.Context) T

type CommonFetchSpec struct {
	Pagination PaginationType
	Resource   string
	Query      url.Values
	AppContext string
	// TODO: This needs to be renamed into Return404AsError
	Return404AsNil bool
	IgnoreResponse bool

	FetchMiddleware []ChainedMiddlewareFunc
}

func (cfs *CommonFetchSpec) ToBuilder() CommonFetchSpecBuilder {
	return CommonFetchSpecBuilder{
		Pagination:     cfs.Pagination,
		Resource:       cfs.Resource,
		Query:          cloneQueryString(cfs.Query),
		AppContext:     cfs.AppContext,
		Return404AsNil: cfs.Return404AsNil,
		IgnoreResponse: cfs.IgnoreResponse,
	}
}

func cloneQueryString(in url.Values) url.Values {
	rv := url.Values{}
	for k, v := range in {
		rv[k] = v
	}

	return rv
}

// FetchSpec Operation context
type FetchSpec struct {
	// Left temporarily: this type has to disappear
	Pagination     PaginationType
	Resource       string
	Query          url.Values
	AppContext     string
	ResponseParser ResponseParserFunc
	Return404AsNil bool
}

type ObjectFetchSpec[T any] struct {
	CommonFetchSpec
	ValueFactory Supplier[T]
}

type ObjectUpsertSpec[T any] struct {
	ObjectFetchSpec[T]
	Upsert T
}

type ObjectExchangeSpec[TIn, TOut any] struct {
	ObjectFetchSpec[TOut]
	Body TIn
}

// TODO: Need to add exchange

type ObjectListFetchSpec[T any] struct {
	CommonFetchSpec
	ValueFactory Supplier[[]T]
}

func (olfs ObjectListFetchSpec[T]) AsObjectFetchSpec() ObjectFetchSpec[[]T] {
	return ObjectFetchSpec[[]T]{
		CommonFetchSpec: olfs.CommonFetchSpec,
		ValueFactory:    olfs.ValueFactory,
	}
}

func (olfs ObjectListFetchSpec[T]) ToBuilder() ObjectListFetchSpecBuilder[T] {
	return ObjectListFetchSpecBuilder[T]{
		CommonFetchSpecBuilder: olfs.CommonFetchSpec.ToBuilder(),
		ValueFactory:           olfs.ValueFactory,
	}
}

type PaginationType int

const (
	PerPage PaginationType = iota
	PerItem
	NotRequired
)

type AsyncFetchResult[T any] struct {
	Data T
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

// DestResource Resource that need to be called on the server. This method will return the resource and
// will append the query string, if specified
func (ctx *CommonFetchSpec) DestResource() string {
	if ctx.Query == nil {
		return ctx.Resource
	} else {
		//fmt.Println(ctx.Query.Encode())
		return fmt.Sprintf("%s?%s", ctx.Resource, ctx.Query.Encode())
	}
}

func (cfs *CommonFetchSpec) FetchFunc() func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
	return func(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
		return c.Fetch(ctx, cfs.DestResource())
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

type CommonFetchSpecBuilder struct {
	Pagination     PaginationType
	Resource       string
	Query          url.Values
	AppContext     string
	Return404AsNil bool
	IgnoreResponse bool
}

func (b *CommonFetchSpecBuilder) WithPagination(p PaginationType) *CommonFetchSpecBuilder {
	b.Pagination = p
	return b
}

func (b *CommonFetchSpecBuilder) WithResource(p string, params ...interface{}) *CommonFetchSpecBuilder {
	if len(params) == 0 {
		b.Resource = p
	} else {
		b.Resource = fmt.Sprintf(p, params...)
	}
	return b
}

func (b *CommonFetchSpecBuilder) WithQuery(p url.Values) *CommonFetchSpecBuilder {
	b.Query = p
	return b
}

func (b *CommonFetchSpecBuilder) WithMergedQuery(p url.Values) *CommonFetchSpecBuilder {
	if p == nil {
		return b
	}

	if b.Query == nil {
		return b.WithQuery(p)
	} else {
		for k, v := range p {
			b.Query[k] = v
		}
	}

	return b
}

func (b *CommonFetchSpecBuilder) WithAppContext(p string) *CommonFetchSpecBuilder {
	b.AppContext = p
	return b
}

func (b *CommonFetchSpecBuilder) WithReturn404AsNil(p bool) *CommonFetchSpecBuilder {
	b.Return404AsNil = p
	return b
}

func (b *CommonFetchSpecBuilder) WithIgnoreResponse(p bool) *CommonFetchSpecBuilder {
	b.IgnoreResponse = p
	return b
}

func (b *CommonFetchSpecBuilder) Build() CommonFetchSpec {
	rv := CommonFetchSpec{
		Pagination:     b.Pagination,
		Resource:       b.Resource,
		Query:          b.Query,
		AppContext:     b.AppContext,
		Return404AsNil: b.Return404AsNil,
		IgnoreResponse: b.IgnoreResponse,
	}

	return rv
}

type ObjectFetchSpecBuilder[T any] struct {
	CommonFetchSpecBuilder
	ValueFactory Supplier[T]
}

type ObjectListFetchSpecBuilder[T any] struct {
	CommonFetchSpecBuilder
	ValueFactory Supplier[[]T]
}

type ObjectUpsertSpecBuilder[T any] struct {
	ObjectFetchSpecBuilder[T]
	Upsert T
}

type ObjectExchangeSpecBuilder[TIn, TOut any] struct {
	ObjectFetchSpecBuilder[TOut]
	Body TIn
}

func (b *ObjectFetchSpecBuilder[T]) WithValueFactory(p Supplier[T]) *ObjectFetchSpecBuilder[T] {
	b.ValueFactory = p
	return b
}

func (b *ObjectListFetchSpecBuilder[T]) WithValueFactory(p Supplier[[]T]) *ObjectListFetchSpecBuilder[T] {
	b.ValueFactory = p
	return b
}

func (b *ObjectUpsertSpecBuilder[T]) WithUpsert(p T) *ObjectUpsertSpecBuilder[T] {
	b.Upsert = p
	return b
}
func (b *ObjectExchangeSpecBuilder[TIn, TOut]) WithBody(p TIn) *ObjectExchangeSpecBuilder[TIn, TOut] {
	b.Body = p
	return b
}

func (b *ObjectFetchSpecBuilder[T]) Build() ObjectFetchSpec[T] {
	rv := ObjectFetchSpec[T]{
		ValueFactory:    b.ValueFactory,
		CommonFetchSpec: b.CommonFetchSpecBuilder.Build(),
	}
	return rv
}

func (b *ObjectListFetchSpecBuilder[T]) Build() ObjectListFetchSpec[T] {
	rv := ObjectListFetchSpec[T]{
		ValueFactory:    b.ValueFactory,
		CommonFetchSpec: b.CommonFetchSpecBuilder.Build(),
	}
	return rv
}

func (b *ObjectUpsertSpecBuilder[T]) Build() ObjectUpsertSpec[T] {
	rv := ObjectUpsertSpec[T]{
		Upsert:          b.Upsert,
		ObjectFetchSpec: b.ObjectFetchSpecBuilder.Build(),
	}
	return rv
}

func (b *ObjectExchangeSpecBuilder[TIn, TOut]) Build() ObjectExchangeSpec[TIn, TOut] {
	rv := ObjectExchangeSpec[TIn, TOut]{
		Body:            b.Body,
		ObjectFetchSpec: b.ObjectFetchSpecBuilder.Build(),
	}
	return rv
}
