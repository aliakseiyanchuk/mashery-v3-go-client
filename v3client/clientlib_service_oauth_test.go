package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestCreateServiceOAuthSecurityProfile(t *testing.T) {
	serviceIdent := masherytypes.ServiceIdentifier{}
	serviceIdent.ServiceId = "service-id"

	expBody := masherytypes.MasheryOAuth{
		AccessTokenTtlEnabled: true,
	}

	expRv := cloneWithModification(expBody, func(t1 *masherytypes.MasheryOAuth) {
		t1.ParentService = serviceIdent
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/securityProfile/oauth").
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
		func(client Client) ClientDualExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth, masherytypes.MasheryOAuth] {
			return client.CreateServiceOAuthSecurityProfile
		},
	)
}

func TestGetMasheryOAuth(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	onTheWire := masherytypes.MasheryOAuth{
		AccessTokenTtlEnabled: true,
	}
	expRV := cloneWithModification(onTheWire, func(t1 *masherytypes.MasheryOAuth) {
		t1.ParentService = serviceId
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/securityProfile/oauth").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(onTheWire)
	}

	autoTestGet(t,
		serviceId,
		expRV,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth] {
			return cl.GetServiceOAuthSecurityProfile
		},
	)
}

func TestUpdateServiceOAuthSecurityProfile(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}
	onTheWire := masherytypes.MasheryOAuth{
		AccessTokenTtlEnabled: true,
	}

	expReturn := cloneWithModification(onTheWire, func(t1 *masherytypes.MasheryOAuth) { t1.ParentService = serviceId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/securityProfile/oauth").
			WithMethod("put").
			RequestingNoFields().
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expReturn,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.MasheryOAuth, masherytypes.MasheryOAuth] {
			return client.UpdateServiceOAuthSecurityProfile
		},
	)
}

func TestDeleteServiceOAuthSecurityProfile(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/securityProfile/oauth").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		serviceId,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceIdentifier] {
			return cl.DeleteServiceOAuthSecurityProfile
		},
	)
}
