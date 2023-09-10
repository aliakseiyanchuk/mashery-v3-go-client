package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestCreateServiceCache(t *testing.T) {
	serviceIdent := masherytypes.ServiceIdentifier{}
	serviceIdent.ServiceId = "service-id"

	expBody := masherytypes.ServiceCache{
		CacheTtl: 4.5,
	}

	expRv := cloneWithModification(expBody, func(t1 *masherytypes.ServiceCache) {
		t1.ParentServiceId = serviceIdent
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/cache").
			WithMethod("post").
			RequestingNoFilters().
			Matching(PayloadMatcher(expBody)).
			WillReturnJsonOf(expBody)
	}

	autoTestCreate(t,
		serviceIdent,
		expBody,
		expRv,
		mockVisitor,
		func(client Client) ClientDualExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.ServiceCache, masherytypes.ServiceCache] {
			return client.CreateServiceCache
		},
	)
}

func TestGetServiceCache(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	onTheWire := masherytypes.ServiceCache{
		CacheTtl: 4.54,
	}
	expRV := cloneWithModification(onTheWire, func(t1 *masherytypes.ServiceCache) {
		t1.ParentServiceId = serviceId
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/cache").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(onTheWire)
	}

	autoTestGet(t,
		serviceId,
		expRV,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.ServiceCache] {
			return cl.GetServiceCache
		},
	)
}

func TestUpdateServiceCache(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}
	onTheWire := masherytypes.ServiceCache{
		CacheTtl: 3.453,
	}

	expReturn := cloneWithModification(onTheWire, func(t1 *masherytypes.ServiceCache) { t1.ParentServiceId = serviceId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/cache").
			WithMethod("put").
			RequestingNoFields().
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expReturn,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceCache, masherytypes.ServiceCache] {
			return client.UpdateServiceCache
		},
	)
}

func TestDeleteServiceCache(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/cache").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		serviceId,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceIdentifier] {
			return cl.DeleteServiceCache
		},
	)
}
