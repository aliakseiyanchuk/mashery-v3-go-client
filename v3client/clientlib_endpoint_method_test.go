package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetEndpointMethod(t *testing.T) {
	endpointMethodId := masherytypes.ServiceEndpointMethodIdentifier{}
	endpointMethodId.ServiceId = "service-id"
	endpointMethodId.EndpointId = "endpoint-id"
	endpointMethodId.MethodId = "method-id"

	expRvMethod := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   "method-id",
				Name: "method-name",
			},
		},
		ParentEndpointId: endpointMethodId.ServiceEndpointIdentifier,
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id").
			WithMethod("get").
			RequestingFields(MasheryMethodsFields).
			WillReturnJsonOf(expRvMethod)
	}

	autoTestGet(t,
		endpointMethodId,
		expRvMethod,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod] {
			return cl.GetEndpointMethod
		},
	)
}

func tSupplyServiceEndpointMethod() masherytypes.ServiceEndpointMethod {
	return masherytypes.ServiceEndpointMethod{}
}

func TestCreateEndpointMethod(t *testing.T) {
	endpointId := masherytypes.ServiceEndpointIdentifier{
		EndpointId:        "endpoint-id",
		ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
	}

	expRvCreatePayload := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{Name: "method-name"},
		},
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.ServiceEndpointMethod) { t1.Id = "endpoint-id" })

	expRvEndpointReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.ServiceEndpointMethod) { t1.ParentEndpointId = endpointId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("post").
			RequestingFields(MasheryMethodsFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		endpointId,
		expRvCreatePayload,
		expRvEndpointReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethod, masherytypes.ServiceEndpointMethod] {
			return cl.CreateEndpointMethod
		},
	)
}

func TestUpdateEndpointMethod(t *testing.T) {
	expRvUpdatePayload := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   "method-id",
				Name: "method-name",
			},
		},
		ParentEndpointId: masherytypes.ServiceEndpointIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
			EndpointId:        "endpoint-id",
		},
	}

	onTheWire := cloneWithModification(expRvUpdatePayload, func(t1 *masherytypes.ServiceEndpointMethod) {
		t1.ParentEndpointId = masherytypes.ServiceEndpointIdentifier{}
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id").
			WithMethod("put").
			RequestingFields(MasheryMethodsFields).
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expRvUpdatePayload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointMethod, masherytypes.ServiceEndpointMethod] {
			return client.UpdateEndpointMethod
		},
	)
}

func TestDeleteEndpointMethod(t *testing.T) {
	postIdent := masherytypes.ServiceEndpointMethodIdentifier{
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{

			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
			EndpointId:        "endpoint-id",
		},
		MethodId: "method-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods/method-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceEndpointMethodIdentifier] {
			return cl.DeleteEndpointMethod
		},
	)
}

func TestCountEndpointsMethodsOf(t *testing.T) {
	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.ServiceEndpointMethod{{}}).
			WillReturnTotalCount(37)
	}

	autoTestCount(t,
		masherytypes.ServiceEndpointIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
			EndpointId:        "endpoint-id",
		},
		37,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointIdentifier, int64] {
			return client.CountEndpointsMethodsOf
		},
	)
}

func TestListEndpointMethodsWithFullInfo(t *testing.T) {
	mockedResponse := []masherytypes.ServiceEndpointMethod{
		{
			BaseMethod: masherytypes.BaseMethod{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "endpoint-id"},
			},
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("get").
			RequestingFields(MasheryMethodsFields).
			WillReturnJsonOf(mockedResponse)
	}

	identCtx := masherytypes.ServiceEndpointIdentifier{
		ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		EndpointId:        "endpoint-id",
	}

	expPop := cloneAllWithModification(mockedResponse,
		func(t1 *masherytypes.ServiceEndpointMethod) {
			t1.ParentEndpointId = identCtx
		})

	autoTestFetchAll(t,
		identCtx,
		expPop,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointIdentifier, []masherytypes.ServiceEndpointMethod] {
			return client.ListEndpointMethodsWithFullInfo
		},
	)
}

func TestListEndpointsMethods(t *testing.T) {
	mockedResponse := []masherytypes.AddressableV3Object{{Id: "method-id"}}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(mockedResponse)
	}

	autoTestFetchAll(t,
		masherytypes.ServiceEndpointIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
			EndpointId:        "endpoint-id",
		},
		mockedResponse,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ServiceEndpointIdentifier, []masherytypes.AddressableV3Object] {
			return client.ListEndpointMethods
		},
	)
}
