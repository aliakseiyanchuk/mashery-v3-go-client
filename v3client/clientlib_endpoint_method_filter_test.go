package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetEndpointMethodFilter(t *testing.T) {
	filterId := masherytypes.ServiceEndpointMethodFilterIdentifier{}
	filterId.ServiceId = "service-id"
	filterId.EndpointId = "endpoint-id"
	filterId.MethodId = "method-id"
	filterId.FilterId = "filter-id"

	mockAPIReturn := masherytypes.ServiceEndpointMethodFilter{}
	mockAPIReturn.Id = "filter-id"
	mockAPIReturn.Name = "filter-name"
	mockAPIReturn.XmlFilterFields = "a,b,c"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters/filter-id").
			WithMethod("get").
			RequestingFields(MasheryResponseFilterFields).
			WillReturnJsonOf(mockAPIReturn)
	}

	expReturn := cloneWithModification(mockAPIReturn, func(t1 *masherytypes.ServiceEndpointMethodFilter) {
		t1.ServiceEndpointMethod = filterId.ServiceEndpointMethodIdentifier
	})

	autoTestGet(t,
		filterId,
		expReturn,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter] {
			return cl.GetEndpointMethodFilter
		},
	)
}

func TestCreateEndpointMethodFilter(t *testing.T) {
	methodId := masherytypes.ServiceEndpointMethodIdentifier{}
	methodId.ServiceId = "service-id"
	methodId.EndpointId = "endpoint-id"
	methodId.MethodId = "method-id"

	expRvCreatePayload := masherytypes.ServiceEndpointMethodFilter{}
	expRvCreatePayload.Name = "filter-name"
	expRvCreatePayload.JsonFilterFields = "a,b,c"

	apiResponseJson := cloneWithModification(expRvCreatePayload,
		func(t1 *masherytypes.ServiceEndpointMethodFilter) { t1.Id = "filter-id" })

	expRvEndpointReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.ServiceEndpointMethodFilter) { t1.ServiceEndpointMethod = methodId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters").
			WithMethod("post").
			RequestingFields(MasheryResponseFilterFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		methodId,
		expRvCreatePayload,
		expRvEndpointReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilter, masherytypes.ServiceEndpointMethodFilter] {
			return cl.CreateEndpointMethodFilter
		},
	)
}

func TestUpdateEndpointMethodFilter(t *testing.T) {
	methodId := masherytypes.ServiceEndpointMethodIdentifier{}
	methodId.ServiceId = "service-id"
	methodId.EndpointId = "endpoint-id"
	methodId.MethodId = "method-id"

	expRvUpdatePayload := masherytypes.ServiceEndpointMethodFilter{}
	expRvUpdatePayload.Id = "filter-id"
	expRvUpdatePayload.Name = "filter-name"
	expRvUpdatePayload.JsonFilterFields = "a,b,c,"

	expRvReturn := cloneWithModification(expRvUpdatePayload,
		func(t1 *masherytypes.ServiceEndpointMethodFilter) { t1.ServiceEndpointMethod = methodId },
	)

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters/filter-id").
			WithMethod("put").
			RequestingFields(MasheryResponseFilterFields).
			Matching(PayloadMatcher(expRvUpdatePayload)).
			WillReturnJsonOf(expRvUpdatePayload)
	}

	autoTestUpdate(t,
		expRvReturn,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointMethodFilter, masherytypes.ServiceEndpointMethodFilter] {
			return client.UpdateEndpointMethodFilter
		},
	)
}

func TestDeleteEndpointMethodFilter(t *testing.T) {
	filterId := masherytypes.ServiceEndpointMethodFilterIdentifier{}
	filterId.ServiceId = "service-id"
	filterId.EndpointId = "endpoint-id"
	filterId.MethodId = "method-id"
	filterId.FilterId = "filter-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters/filter-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		filterId,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceEndpointMethodFilterIdentifier] {
			return cl.DeleteEndpointMethodFilter
		},
	)
}

func TestCountEndpointsMethodsFiltersOf(t *testing.T) {
	methodId := masherytypes.ServiceEndpointMethodIdentifier{}
	methodId.ServiceId = "service-id"
	methodId.EndpointId = "endpoint-id"
	methodId.MethodId = "method-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.PackagePlanServiceEndpointMethodFilter{{}}).
			WillReturnTotalCount(3)
	}

	autoTestCount(t,
		methodId,
		3,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointMethodIdentifier, int64] {
			return client.CountEndpointsMethodsFiltersOf
		},
	)
}

func TestListEndpointMethodFiltersWithFullInfo(t *testing.T) {
	rvObj := masherytypes.ServiceEndpointMethodFilter{}
	rvObj.Id = "filter-id"
	rvObj.Name = "filter-name"
	rvObj.JsonFilterFields = "a,b,c"

	mockedResponse := []masherytypes.ServiceEndpointMethodFilter{
		rvObj,
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters").
			WithMethod("get").
			RequestingFields(MasheryResponseFilterFields).
			WillReturnJsonOf(mockedResponse)
	}

	methodId := masherytypes.ServiceEndpointMethodIdentifier{}
	methodId.ServiceId = "service-id"
	methodId.EndpointId = "endpoint-id"
	methodId.MethodId = "method-id"

	expPop := cloneAllWithModification(mockedResponse,
		func(t1 *masherytypes.ServiceEndpointMethodFilter) {
			t1.ServiceEndpointMethod = methodId
		})

	autoTestFetchAll(t,
		methodId,
		expPop,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointMethodIdentifier, []masherytypes.ServiceEndpointMethodFilter] {
			return client.ListEndpointMethodFiltersWithFullInfo
		},
	)
}

func TestListEndpointMethodFilters(t *testing.T) {
	mockedResponse := []masherytypes.AddressableV3Object{{Id: "filter-id"}}

	methodId := masherytypes.ServiceEndpointMethodIdentifier{}
	methodId.ServiceId = "service-id"
	methodId.EndpointId = "endpoint-id"
	methodId.MethodId = "method-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilters").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(mockedResponse)
	}

	autoTestFetchAll(t,
		methodId,
		mockedResponse,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointMethodIdentifier, []masherytypes.AddressableV3Object] {
			return client.ListEndpointMethodFilters
		},
	)
}
