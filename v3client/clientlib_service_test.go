package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetService(t *testing.T) {
	memberId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	expRvService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "service-id",
			Name: "service-name",
		},
		Description: "desc-1",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id").
			WithMethod("get").
			RequestingFields(MasheryServiceFields).
			WillReturnJsonOf(expRvService)
	}

	autoTestGet(t,
		memberId,
		expRvService,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.Service] {
			return cl.GetService
		},
	)
}

func TestCreateService(t *testing.T) {

	expRvCreatePayload := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{Id: "", Name: "service-name"},
		Description:         "desc-1",
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.Service) { t1.Id = "service-id" })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services").
			WithMethod("post").
			RequestingFields(MasheryServiceFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestRootCreate(t,
		expRvCreatePayload,
		apiResponseJson,
		mockVisitor,
		func(cl Client) ClientExchangeFunc[masherytypes.Service, masherytypes.Service] {
			return cl.CreateService
		},
	)
}

func TestUpdateService(t *testing.T) {
	payload := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "service-id",
			Name: "service-name",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id").
			WithMethod("put").
			RequestingFields(MasheryServiceFields).
			Matching(PayloadMatcher(payload)).
			WillReturnJsonOf(payload)
	}

	autoTestUpdate(t,
		payload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.Service, masherytypes.Service] {
			return client.UpdateService
		},
	)
}

func TestDeleteService(t *testing.T) {
	postIdent := masherytypes.ServiceIdentifier{
		ServiceId: "service-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceIdentifier] {
			return cl.DeleteService
		},
	)
}

func TestListServices(t *testing.T) {
	mockedResponse := []masherytypes.Service{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "service-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services").
			WithMethod("get").
			RequestingFields(MasheryServiceFields).
			WillReturnJsonOf(mockedResponse)
	}

	autoTestRootFetchAll(t,
		mockedResponse,
		mockVisitor,
		func(client Client) ClientArraySupplierFunc[masherytypes.Service] {
			return client.ListServices
		},
	)
}

func TestListServicesFiltered(t *testing.T) {
	mockedResponse := []masherytypes.Service{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "service-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services").
			WithMethod("get").
			RequestingFields(MasheryServiceFields).
			FilteredOn("a", "b").
			WillReturnJsonOf(mockedResponse)
	}

	autoTestRootFetchFiltered(t,
		map[string]string{"a": "b"},
		mockedResponse,
		mockVisitor,
		func(client Client) ClientFilteredArraySupplierFunc[masherytypes.Service] {
			return client.ListServicesFiltered
		},
	)
}

func TestCountServices(t *testing.T) {
	mockedResponse := []masherytypes.Service{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "service-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services").
			WithMethod("get").
			RequestingNoFields().
			FilteredOn("a", "b").
			WillReturnJsonOf(mockedResponse).
			WillReturnTotalCount(98)
	}

	autoTestRootFilteredCount(t,
		map[string]string{"a": "b"},
		98,
		mockVisitor,
		func(client Client) ClientExchangeFunc[map[string]string, int64] {
			return client.CountServices
		},
	)
}
