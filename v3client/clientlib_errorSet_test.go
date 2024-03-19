package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetErrorSet(t *testing.T) {
	endpointId := masherytypes.ErrorSetIdentifier{}
	endpointId.ServiceId = "service-id"
	endpointId.ErrorSetId = "error-set-id"

	onTheWireReturn := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "error-set-od",
			Name: "error-set-name",
		},
		Type: "abc",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/errorSets/error-set-id").
			WithMethod("get").
			RequestingFields(MasheryErrorSetFields).
			WillReturnJsonOf(onTheWireReturn)
	}

	expRv := cloneWithModification(onTheWireReturn, func(t1 *masherytypes.ErrorSet) {
		t1.ParentServiceId = endpointId.ServiceIdentifier
	})

	autoTestGet(t,
		endpointId,
		expRv,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet] {
			return cl.GetErrorSet
		},
	)
}

func TestCreateErrorSet(t *testing.T) {
	serviceIdent := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	expRvCreate := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{Id: "", Name: "endpoint-name"},
		Type:                "abc",
		ParentServiceId:     masherytypes.ServiceIdentifier{ServiceId: "service-id"},
	}

	expRvCreateOnTheWire := cloneWithModification(expRvCreate, func(t1 *masherytypes.ErrorSet) { t1.ParentServiceId = masherytypes.ServiceIdentifier{} })
	apiResponseJson := cloneWithModification(expRvCreateOnTheWire, func(t1 *masherytypes.ErrorSet) { t1.Id = "error-set-id" })
	expRvEndpointReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.ErrorSet) { t1.ParentServiceId = serviceIdent })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/errorSets").
			WithMethod("post").
			RequestingFields(MasheryErrorSetFields).
			Matching(PayloadMatcher(expRvCreateOnTheWire)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		serviceIdent,
		expRvCreate,
		expRvEndpointReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.ErrorSet, masherytypes.ErrorSet] {
			return cl.CreateErrorSet
		},
	)
}

func TestUpdateErrorSet(t *testing.T) {
	expRvUpdatePayload := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "error-set-id",
			Name: "error-set-name",
		},
		Type: "abc",

		ParentServiceId: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
	}

	onTheWire := cloneWithModification(expRvUpdatePayload, func(t1 *masherytypes.ErrorSet) { t1.ParentServiceId = masherytypes.ServiceIdentifier{} })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/errorSets/error-set-id").
			WithMethod("put").
			RequestingFields(MasheryErrorSetFields).
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expRvUpdatePayload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.ErrorSet, masherytypes.ErrorSet] {
			return client.UpdateErrorSet
		},
	)
}

func TestUpdateErrorSetMessage(t *testing.T) {
	errorSetId := masherytypes.ErrorSetIdentifier{}
	errorSetId.ServiceId = "service-id"
	errorSetId.ErrorSetId = "endpoint-id"

	expRvUpdatePayload := masherytypes.MasheryErrorMessage{
		Id:   "MessageId",
		Code: 404,

		ParentErrorSet: masherytypes.ErrorSetIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
			ErrorSetId:        "error-set-id",
		},
	}

	onTheWire := cloneWithModification(expRvUpdatePayload, func(t1 *masherytypes.MasheryErrorMessage) {
		t1.ParentErrorSet = masherytypes.ErrorSetIdentifier{}
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/errorSets/error-set-id/errorMessages/MessageId").
			WithMethod("put").
			RequestingFields(MasheryErrorSetFields).
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expRvUpdatePayload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.MasheryErrorMessage, masherytypes.MasheryErrorMessage] {
			return client.UpdateErrorSetMessage
		},
	)
}

func TestDeleteErrorSet(t *testing.T) {
	postIdent := masherytypes.ErrorSetIdentifier{
		ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "service-id"},
		ErrorSetId:        "error-set-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/errorSets/error-set-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ErrorSetIdentifier] {
			return cl.DeleteErrorSet
		},
	)
}
