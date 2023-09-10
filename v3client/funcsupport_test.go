package v3client

import (
	"context"
	"encoding/json"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type ClientSupplierFunc[TOut any] func(context.Context) (TOut, error)
type ClientSupplierFuncLocator[TOut any] func(Client) ClientSupplierFunc[TOut]

type ClientConsumerFunc[TIn any] func(context.Context, TIn) error
type ClientConsumerFuncLocator[TIn any] func(Client) ClientConsumerFunc[TIn]

type ClientBiConsumerFunc[TIdent, TIn any] func(context.Context, TIdent, TIn) error
type ClientBiConsumerFuncLocator[TIdent, TIn any] func(Client) ClientBiConsumerFunc[TIdent, TIn]

type ClientArraySupplierFunc[TOut any] func(context.Context) ([]TOut, error)
type ClientArraySupplierFuncLocator[TOut any] func(Client) ClientArraySupplierFunc[TOut]

type ClientFilteredArraySupplierFunc[TOut any] func(context.Context, map[string]string) ([]TOut, error)
type ClientFilteredArraySupplierFuncLocator[TOut any] func(Client) ClientFilteredArraySupplierFunc[TOut]

type ClientExchangeFunc[TIn, TOut any] func(context.Context, TIn) (TOut, error)
type ClientExchangeFuncLocator[TIn, TOut any] func(Client) ClientExchangeFunc[TIn, TOut]

type ClientBoolExchangeFunc[TIn, TOut any] func(context.Context, TIn) (TOut, bool, error)
type ClientBoolExchangeFuncLocator[TIn, TOut any] func(Client) ClientBoolExchangeFunc[TIn, TOut]

type ClientDualExchangeFunc[TInA, TInB, TOut any] func(context.Context, TInA, TInB) (TOut, error)
type ClientDualExchangeFuncLocator[TInA, TInB, TOut any] func(Client) ClientDualExchangeFunc[TInA, TInB, TOut]

type RequestMatcherFunc func(request *http.Request) bool

func IgnoreBiConsumerError[T1, T2 any](f BiConsumerCanErr[T1, T2]) BiConsumer[T1, T2] {
	return func(t1 T1, t2 T2) {
		_ = f(t1, t2)
	}
}

func IgnoreError[TIn, TOut any](f func(in TIn) (TOut, error)) func(argA TIn) TOut {
	return func(argA TIn) TOut {
		rv, _ := f(argA)
		return rv
	}
}

func IgnoreSupplyErr[T any](f SupplierCanErr[T]) Supplier[T] {
	return func() T {
		v, _ := f()
		return v
	}
}

var readAllFully func(io.Reader) []byte
var marshalJson func(interface{}) []byte
var unmarshalJson BiConsumer[[]byte, any]

func clone[T any](in T, f Supplier[T]) T {
	str := marshalJson(in)
	cloned := f()
	unmarshalJson(str, &cloned)

	return cloned
}

func cloneWithModification[T any](in T, c Consumer[*T]) T {
	str := marshalJson(in)
	cloned := new(T)

	unmarshalJson(str, &cloned)
	c(cloned)

	return *cloned
}

func cloneAllWithModification[T any](in []T, c Consumer[*T]) []T {
	cloned := make([]T, len(in))
	str := marshalJson(in)
	unmarshalJson(str, &cloned)

	for i, v := range cloned {
		c(&v)
		cloned[i] = v
	}

	return cloned
}

func uc(in []masherytypes.Application, c Consumer[*masherytypes.Application]) []masherytypes.Application {
	cloned := make([]masherytypes.Application, len(in))

	str := marshalJson(in)
	unmarshalJson(str, &cloned)

	for idx, v := range cloned {
		c(&v)
		cloned[idx] = v
	}

	return cloned
}

func init() {
	readAllFully = IgnoreError(io.ReadAll)
	marshalJson = IgnoreError(json.Marshal)
	unmarshalJson = IgnoreBiConsumerError(json.Unmarshal)
}

func TestCloning(t *testing.T) {
	app := masherytypes.Application{AddressableV3Object: masherytypes.AddressableV3Object{Id: "ABC"}}
	cloned := cloneWithModification(app, func(t1 *masherytypes.Application) {
		t1.AdsSystem = "Boo"
		t1.Id = t1.Id + "-mod"
	})

	assert.Equal(t, "ABC-mod", cloned.Id)
	assert.Equal(t, "Boo", cloned.AdsSystem)

	cloned2 := cloneWithModification(app, func(t1 *masherytypes.Application) {
		t1.AddressableV3Object = masherytypes.AddressableV3Object{Id: "Swapped"}
	})

	assert.Equal(t, "Swapped", cloned2.Id)
}

func TestArrayCloning(t *testing.T) {
	app := []masherytypes.Application{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "ABC"}},
	}
	cloned := cloneAllWithModification(app, func(t1 *masherytypes.Application) {
		t1.AdsSystem = "Cloned!"
	})

	for _, v := range cloned {
		assert.Equal(t, "Cloned!", v.AdsSystem)
	}

	cloned2 := uc(app, func(t1 *masherytypes.Application) {
		t1.AdsSystem = "Cloned!"
	})

	for _, v := range cloned2 {
		assert.Equal(t, "Cloned!", v.AdsSystem)
	}
}
