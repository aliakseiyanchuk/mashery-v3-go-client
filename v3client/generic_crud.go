package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
	"strings"
)

const fieldsContextKey = contextKeyType("fields.v3.client.mashery")

type GenericCRUDDecorator[ParentIdent, TIdent, T any] struct {
	ValueSupplier      transport.Supplier[T]
	ValueArraySupplier transport.Supplier[[]T]
	ResourceFor        func(TIdent) (string, error)
	// ResourceForUpsert Return a resource for upsert operation.
	ResourceForUpsert func(T) (string, error)
	ResourceForParent func(ParentIdent) (string, error)
	DefaultFields     []string
	GetFields         func(context.Context) []string
	Pagination        transport.PaginationType

	AcceptObjectIdent BiConsumer[TIdent, *T]
	AcceptParentIdent BiConsumer[ParentIdent, *T]
	AcceptIdentFrom   BiConsumer[T, *T]
}

// Dummy int supplier
func intSupplier() int {
	return 0
}

func ReturnFields(ctx context.Context, v []string) context.Context {
	return context.WithValue(ctx, fieldsContextKey, v)
}

func DefaultGetFieldsFromContext(defaultFields []string) func(context.Context) []string {
	return func(ctx context.Context) []string {
		if v := ctx.Value(fieldsContextKey); v != nil {
			if vArr, ok := v.([]string); ok {
				return vArr
			}
		}

		return defaultFields
	}
}

type GenericCRUD[TParent, TIdent, T any] struct {
	AppContext string
	Decorator  GenericCRUDDecorator[TParent, TIdent, T]

	querySupplier    transport.ContextSupplier[url.Values]
	fixedFieldsQuery url.Values

	doGet      func(ctx context.Context, spec transport.ObjectFetchSpec[T], c *transport.HttpTransport) (T, bool, error)
	doCreate   func(ctx context.Context, opCtx transport.ObjectUpsertSpec[T], c *transport.HttpTransport) (T, error)
	doUpdate   func(ctx context.Context, opCtx transport.ObjectUpsertSpec[T], c *transport.HttpTransport) (T, error)
	doDelete   func(ctx context.Context, opCtx transport.ObjectFetchSpec[T], c *transport.HttpTransport) error
	doCount    func(ctx context.Context, opCtx transport.ObjectListFetchSpec[T], c *transport.HttpTransport) (int64, error)
	doFetchAll func(ctx context.Context, opCtx transport.ObjectListFetchSpec[T], c *transport.HttpTransport) ([]T, error)
}

func NewCRUD[TParent, TIdent, T any](appContext string, d *GenericCRUDDecorator[TParent, TIdent, T]) *GenericCRUD[TParent, TIdent, T] {
	rv := GenericCRUD[TParent, TIdent, T]{
		AppContext: appContext,
		Decorator:  *d,
		doGet:      transport.GetObject[T],
		doCreate:   transport.CreateObject[T],
		doUpdate:   transport.UpdateObject[T],
		doDelete:   transport.DeleteObject[T],
		doCount:    transport.Count[T],
		doFetchAll: transport.FetchAll[T],
	}

	// If the decorator defined dynamic fields, it will be wrapped to produce dynamic
	// values. Otherwise, fixed values will be used.
	if d.GetFields != nil {
		rv.querySupplier = func(ctx context.Context) url.Values {
			if fields := d.GetFields(ctx); len(fields) > 0 {
				return url.Values{
					"fields": []string{strings.Join(fields, ",")},
				}
			} else {
				return nil
			}
		}
	} else {
		if len(d.DefaultFields) > 0 {
			rv.fixedFieldsQuery = url.Values{
				"fields": []string{strings.Join(d.DefaultFields, ",")},
			}
		} else {
			rv.fixedFieldsQuery = nil
		}

		rv.querySupplier = func(_ context.Context) url.Values {
			return rv.fixedFieldsQuery
		}
	}

	return &rv
}

type CRUDGetter[TIdent, T any] func(ctx context.Context, ident TIdent, c *transport.HttpTransport) (T, bool, error)
type CRUDCreator[TParent, T any] func(ctx context.Context, ident TParent, upsert T, c *transport.HttpTransport) (T, error)
type CRUDAllFetcher[TParent, T any] func(ctx context.Context, ident TParent, c *transport.HttpTransport) ([]T, error)
type CRUDFilteredFetcher[TParent, T any] func(ctx context.Context, ident TParent, filter map[string]string, c *transport.HttpTransport) ([]T, error)
type CRUDFilteredCounter[TParent, T any] func(ctx context.Context, ident TParent, filter map[string]string, c *transport.HttpTransport) (int64, error)

func (crud *GenericCRUD[TParent, TIdent, T]) Get(ctx context.Context, ident TIdent, c *transport.HttpTransport) (T, bool, error) {
	if resourceURL, err := crud.Decorator.ResourceFor(ident); err != nil {
		return crud.StubValue(), false, err
	} else {
		fetchSpecBuilder := transport.ObjectFetchSpecBuilder[T]{}
		fetchSpecBuilder.
			WithValueFactory(crud.Decorator.ValueSupplier).
			WithResource(resourceURL).
			WithQuery(crud.querySupplier(ctx)).
			WithAppContext(crud.AppContext).
			WithReturn404AsNil(true)

		if get, exist, fetchErr := crud.doGet(ctx, fetchSpecBuilder.Build(), c); fetchErr != nil {
			return crud.StubValue(), false, fetchErr
		} else {
			if crud.Decorator.AcceptObjectIdent != nil {
				crud.Decorator.AcceptObjectIdent(ident, &get)
			}

			return get, exist, fetchErr
		}

	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) StubValue() T {
	return crud.Decorator.ValueSupplier()
}

func (crud *GenericCRUD[TParent, TIdent, T]) Exists(ctx context.Context, ident TIdent, c *transport.HttpTransport) (bool, error) {
	if resourceURL, err := crud.Decorator.ResourceFor(ident); err != nil {
		fetchSpecBuilder := transport.ObjectFetchSpecBuilder[T]{}
		fetchSpecBuilder.
			WithValueFactory(crud.Decorator.ValueSupplier).
			WithResource(resourceURL).
			WithAppContext(crud.AppContext).
			WithReturn404AsNil(true)

		return transport.Exists(ctx, fetchSpecBuilder.Build(), c)
	} else {
		return false, err
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) Create(ctx context.Context, ident TParent, upsert T, c *transport.HttpTransport) (T, error) {
	if resourceURL, err := crud.Decorator.ResourceForParent(ident); err != nil {
		return crud.StubValue(), err
	} else {
		fetchSpecBuilder := transport.ObjectUpsertSpecBuilder[T]{}
		fetchSpecBuilder.
			WithUpsert(upsert).
			WithValueFactory(crud.Decorator.ValueSupplier).
			WithResource(resourceURL).
			WithQuery(crud.querySupplier(ctx)).
			WithAppContext(crud.AppContext)

		if create, createErr := crud.doCreate(ctx, fetchSpecBuilder.Build(), c); createErr != nil {
			return crud.StubValue(), createErr
		} else {
			if crud.Decorator.AcceptParentIdent != nil {
				crud.Decorator.AcceptParentIdent(ident, &create)
			}
			return create, createErr
		}
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) Update(ctx context.Context, upsert T, c *transport.HttpTransport) (T, error) {
	if resourceURL, err := crud.Decorator.ResourceForUpsert(upsert); err != nil {
		return crud.StubValue(), err
	} else {
		objectUpsertSpecBuilder := transport.ObjectUpsertSpecBuilder[T]{}
		objectUpsertSpecBuilder.
			WithUpsert(upsert).
			WithValueFactory(crud.Decorator.ValueSupplier).
			WithResource(resourceURL).
			WithQuery(crud.querySupplier(ctx)).
			WithAppContext(crud.AppContext)

		if update, updateErr := crud.doUpdate(ctx, objectUpsertSpecBuilder.Build(), c); updateErr != nil {
			return crud.StubValue(), updateErr
		} else {
			if crud.Decorator.AcceptIdentFrom != nil {
				crud.Decorator.AcceptIdentFrom(upsert, &update)
			}
			return update, updateErr
		}
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) Delete(ctx context.Context, id TIdent, c *transport.HttpTransport) error {
	if resourceURL, err := crud.Decorator.ResourceFor(id); err != nil {
		return err
	} else {
		objectUpsertSpecBuilder := transport.ObjectFetchSpecBuilder[T]{}
		objectUpsertSpecBuilder.
			WithValueFactory(crud.Decorator.ValueSupplier).
			WithIgnoreResponse(true).
			WithResource(resourceURL).
			WithAppContext(crud.AppContext).
			WithIgnoreResponse(true)

		return crud.doDelete(ctx, objectUpsertSpecBuilder.Build(), c)
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) Count(ctx context.Context, id TParent, c *transport.HttpTransport) (int64, error) {
	return crud.CountFiltered(ctx, id, nil, c)
}

func (crud *GenericCRUD[TParent, TIdent, T]) CountFiltered(ctx context.Context, id TParent, filter map[string]string, c *transport.HttpTransport) (int64, error) {
	if resourceURL, err := crud.Decorator.ResourceForParent(id); err != nil {
		return 0, err
	} else {
		objectUpsertSpecBuilder := transport.ObjectListFetchSpecBuilder[T]{}
		objectUpsertSpecBuilder.
			WithValueFactory(crud.Decorator.ValueArraySupplier).
			WithResource(resourceURL).
			WithQuery(crud.toFilterQuery(filter)).
			WithAppContext(crud.AppContext).
			WithReturn404AsNil(false).
			WithPagination(transport.NotRequired)

		return crud.doCount(ctx, objectUpsertSpecBuilder.Build(), c)
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) FetchAllAsAddressable(ctx context.Context, id TParent, c *transport.HttpTransport) ([]masherytypes.AddressableV3Object, error) {
	return crud.FetchAllAsAddressableFiltered(ctx, id, nil, c)
}

func (crud *GenericCRUD[TParent, TIdent, T]) FetchAllAsAddressableFiltered(ctx context.Context, id TParent, filter map[string]string, c *transport.HttpTransport) ([]masherytypes.AddressableV3Object, error) {
	if resourceURL, err := crud.Decorator.ResourceForParent(id); err != nil {
		return nil, err
	} else {
		objectListSpecBuilder := transport.ObjectListFetchSpecBuilder[masherytypes.AddressableV3Object]{}
		objectListSpecBuilder.
			WithValueFactory(func() []masherytypes.AddressableV3Object {
				return []masherytypes.AddressableV3Object{}
			}).
			WithResource(resourceURL).
			WithQuery(crud.toFilterQuery(filter)).
			WithAppContext(crud.AppContext).
			WithPagination(crud.Decorator.Pagination)

		return transport.FetchAll(ctx, objectListSpecBuilder.Build(), c)
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) FetchAll(ctx context.Context, id TParent, c *transport.HttpTransport) ([]T, error) {
	return crud.FetchFiltered(ctx, id, nil, c)
}

func (crud *GenericCRUD[TParent, TIdent, T]) FetchFiltered(ctx context.Context, id TParent, filter map[string]string, c *transport.HttpTransport) ([]T, error) {
	if resourceURL, err := crud.Decorator.ResourceForParent(id); err != nil {
		return nil, err
	} else {
		objectListSpecBuilder := transport.ObjectListFetchSpecBuilder[T]{}
		objectListSpecBuilder.
			WithValueFactory(crud.Decorator.ValueArraySupplier).
			WithResource(resourceURL).
			WithQuery(crud.querySupplier(ctx)).
			WithMergedQuery(crud.toFilterQuery(filter)).
			WithAppContext(crud.AppContext).
			WithPagination(crud.Decorator.Pagination)

		if all, fetchAllErr := crud.doFetchAll(ctx, objectListSpecBuilder.Build(), c); fetchAllErr != nil {
			return nil, fetchAllErr
		} else {
			if len(all) > 0 && crud.Decorator.AcceptParentIdent != nil {
				for idx, v := range all {
					crud.Decorator.AcceptParentIdent(id, &v)
					all[idx] = v
				}
			}

			return all, fetchAllErr
		}
	}
}

func (crud *GenericCRUD[TParent, TIdent, T]) toFilterQuery(filter map[string]string) url.Values {
	if len(filter) > 0 {
		srchAtoms := make([]string, len(filter))

		var idx = 0
		for k, v := range filter {
			srchAtoms[idx] = fmt.Sprintf("%s:%s", k, v)
			idx++
		}

		return url.Values{
			"filter": []string{strings.Join(srchAtoms, ",")},
		}
	} else {
		return nil
	}
}

func GetWithFields[TIdent, T any](fields []string, getter CRUDGetter[TIdent, T]) CRUDGetter[TIdent, T] {
	return func(ctx context.Context, ident TIdent, c *transport.HttpTransport) (T, bool, error) {
		return getter(ReturnFields(ctx, fields), ident, c)
	}
}

// RootFetcher builds a fetch-all function that wraps a parent context object. Useful to shorten the signature
// where the object retrieved is already root, or to enforce a particular parent object context
func RootFetcher[TParent, T any](fetcher CRUDAllFetcher[TParent, T], rootIdent TParent) func(ctx context.Context, c *transport.HttpTransport) ([]T, error) {
	return func(ctx context.Context, c *transport.HttpTransport) ([]T, error) {
		return fetcher(ctx, rootIdent, c)
	}
}

func RootCreator[TParent, T any](creator CRUDCreator[TParent, T], rootIdent TParent) func(ctx context.Context, t T, c *transport.HttpTransport) (T, error) {
	return func(ctx context.Context, t T, c *transport.HttpTransport) (T, error) {
		return creator(ctx, rootIdent, t, c)
	}
}

func RootFilteredFetcher[TParent, T any](fetcher CRUDFilteredFetcher[TParent, T], rootIdent TParent) func(context.Context, map[string]string, *transport.HttpTransport) ([]T, error) {
	return func(ctx context.Context, filter map[string]string, c *transport.HttpTransport) ([]T, error) {
		return fetcher(ctx, rootIdent, filter, c)
	}
}

func RootFilteredCounter[TParent, T any](counter CRUDFilteredCounter[TParent, T], rootIdent TParent) func(context.Context, map[string]string, *transport.HttpTransport) (int64, error) {
	return func(ctx context.Context, filter map[string]string, c *transport.HttpTransport) (int64, error) {
		return counter(ctx, rootIdent, filter, c)
	}
}
