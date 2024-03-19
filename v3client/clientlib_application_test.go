package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	_ "github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"testing"
)

func TestGetApplication(t *testing.T) {
	expRvApp := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "AppId",
			Name: "AppName",
		},
		Username: "Foo",
		Eav:      masherytypes.EAV{},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/applications/AppId").
			WithMethod("get").
			RequestingNoFields().
			//RequestingFields(applicationFields).
			WillReturnJsonOf(expRvApp)
	}

	autoTestGet(t,
		masherytypes.ApplicationIdentifier{ApplicationId: "AppId"},
		expRvApp,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ApplicationIdentifier, masherytypes.Application] {
			return cl.GetApplication
		},
	)
}

func TestGetFullApplication(t *testing.T) {
	expRvApp := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "AppId-Full",
			Name: "AppName",
		},
		Username:    "Foo",
		PackageKeys: &[]masherytypes.ApplicationPackageKey{
			// TODO this needs to be fixed
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/applications/AppId-Full").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestGet(t,
		masherytypes.ApplicationIdentifier{ApplicationId: "AppId-Full"},
		expRvApp,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ApplicationIdentifier, masherytypes.Application] {
			return cl.GetFullApplication
		},
	)
}

func TestCreateApplication(t *testing.T) {
	postPayload := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Name: "AppName1",
		},
		Username: "Foo1",
	}

	expRvApp := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "AppId1",
			Name: "AppName1",
		},
		Username: "Foo1",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members/member-id/applications").
			WithMethod("post").
			RequestingNoFields().
			Matching(PayloadMatcher(postPayload)).
			WillReturnJsonOf(expRvApp)
	}

	autoTestCreate(t,
		masherytypes.MemberIdentifier{MemberId: "member-id"},
		postPayload,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.MemberIdentifier, masherytypes.Application, masherytypes.Application] {
			return cl.CreateApplication
		},
	)
}

//func TestUpdateApplication(t *testing.T) {
//	postPayload := masherytypes.Application{
//		AddressableV3Object: masherytypes.AddressableV3Object{
//			Name: "AppName1",
//		},
//		Username: "Foo1",
//	}
//
//	returnPayload := masherytypes.Application{
//		AddressableV3Object: masherytypes.AddressableV3Object{
//			Id:   "app-id-update",
//			Name: "AppName1",
//		},
//		Username: "Foo1",
//	}
//
//	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
//		matcher.
//			ForRequestPath("/applications/app-id-update").
//			WithMethod("put").
//			RequestingNoFields().
//			Matching(PayloadMatcher(postPayload)).
//			WillReturnJsonOf(returnPayload)
//	}
//
//	autoTestUpdate(t,
//		postPayload,
//		mockVisitor,
//		func(cl Client) ClientExchangeFunc[masherytypes.Application, masherytypes.Application] {
//			return cl.UpdateApplication
//		},
//	)
//
//}

func TestDeleteApplication(t *testing.T) {
	postIdent := masherytypes.ApplicationIdentifier{
		ApplicationId: "app-id-delete",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/applications/app-id-delete").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ApplicationIdentifier] {
			return cl.DeleteApplication
		},
	)
}
