package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetEndpoint(t *testing.T) {
	endpointId := masherytypes.ServiceEndpointIdentifier{}
	endpointId.ServiceId = "service-id"
	endpointId.EndpointId = "endpoint-id"

	expRvEndpoint := masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "endpoint-od",
			Name: "endpoint-name",
		},
		ApiKeyValueLocationKey: "apikey",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id").
			WithMethod("get").
			RequestingFields(MasheryEndpointFields).
			WillReturnJsonOf(expRvEndpoint)
	}

	autoTestGet(t,
		endpointId,
		expRvEndpoint,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint] {
			return cl.GetEndpoint
		},
	)
}

func TestCreateEndpoint(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	//endpointId := masherytypes.ServiceEndpointIdentifier{
	//	ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
	//	EndpointId: "endpoint-id",
	//}

	expRvCreatePayload := masherytypes.Endpoint{
		AddressableV3Object:    masherytypes.AddressableV3Object{Id: "", Name: "endpoint-name"},
		ApiKeyValueLocationKey: "apikey",
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.Endpoint) { t1.Id = "endpoint-id" })

	expRvEndpointReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.Endpoint) { t1.ParentServiceId = serviceId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints").
			WithMethod("post").
			RequestingFields(MasheryEndpointFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		serviceId,
		expRvCreatePayload,
		expRvEndpointReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.Endpoint, masherytypes.Endpoint] {
			return cl.CreateEndpoint
		},
	)
}

func TestCreateEndpointAutoRetryOnBadRequest(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	expRvCreatePayload := masherytypes.Endpoint{
		AddressableV3Object:    masherytypes.AddressableV3Object{Id: "", Name: "endpoint-name"},
		ApiKeyValueLocationKey: "apikey",
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.Endpoint) { t1.Id = "endpoint-id" })

	expRvEndpointReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.Endpoint) { t1.ParentServiceId = serviceId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints").
			WithMethod("post").
			RequestingFields(MasheryEndpointFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestSuccessfulCreateWithBadRequestAutoRetry(t,
		serviceId,
		expRvCreatePayload,
		expRvEndpointReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.Endpoint, masherytypes.Endpoint] {
			return cl.CreateEndpoint
		},
	)
}

func TestUpdateEndpoint(t *testing.T) {
	expRvUpdatePayload := masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "endpoint-id",
			Name: "endpoint-name",
		},
		ApiKeyValueLocationKey: "apikey",
		ParentServiceId:        masherytypes.ServiceIdentifier{ServiceId: "service-id"},
	}

	onTheWire := cloneWithModification(expRvUpdatePayload, func(t1 *masherytypes.Endpoint) { t1.ParentServiceId = masherytypes.ServiceIdentifier{} })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id").
			WithMethod("put").
			RequestingFields(MasheryEndpointFields).
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expRvUpdatePayload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.Endpoint, masherytypes.Endpoint] {
			return client.UpdateEndpoint
		},
	)
}

func TestDeleteEndpoint(t *testing.T) {
	postIdent := masherytypes.ServiceEndpointIdentifier{
		ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		EndpointId:        "endpoint-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceEndpointIdentifier] {
			return cl.DeleteEndpoint
		},
	)
}

func TestCountEndpointsOf(t *testing.T) {
	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.Endpoint{{}}).
			WillReturnTotalCount(34)
	}

	autoTestCount(t,
		masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		34,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceIdentifier, int64] {
			return client.CountEndpointsOf
		},
	)
}

func TestListEndpointsWithFullInfo(t *testing.T) {
	mockedResponse := []masherytypes.Endpoint{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "endpoint-id"}},
	}
	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints").
			WithMethod("get").
			RequestingFields(MasheryEndpointFields).
			WillReturnJsonOf(mockedResponse)
	}

	identCtx := masherytypes.ServiceIdentifier{ServiceId: "service-id"}
	expPop := cloneAllWithModification(mockedResponse,
		func(t1 *masherytypes.Endpoint) {
			t1.ParentServiceId = identCtx
		})

	autoTestFetchAll(t,
		masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		expPop,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceIdentifier, []masherytypes.Endpoint] {
			return client.ListEndpointsWithFullInfo
		},
	)
}

func TestListEndpoints(t *testing.T) {
	mockedResponse := []masherytypes.AddressableV3Object{{Id: "endpoint-id"}}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(mockedResponse)
	}

	autoTestFetchAll(t,
		masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		mockedResponse,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceIdentifier, []masherytypes.AddressableV3Object] {
			return client.ListEndpoints
		},
	)
}
