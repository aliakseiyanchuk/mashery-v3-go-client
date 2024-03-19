package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetMember(t *testing.T) {
	memberId := masherytypes.MemberIdentifier{MemberId: "member-id"}

	expRvMember := masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "member-id",
			Name: "member-name",
		},
		Address1: "addr-1",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members/member-id").
			WithMethod("get").
			RequestingFields(memberFields).
			WillReturnJsonOf(expRvMember)
	}

	autoTestGet(t,
		memberId,
		expRvMember,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.MemberIdentifier, masherytypes.Member] {
			return cl.GetMember
		},
	)
}

func TestGetFullMember(t *testing.T) {
	memberId := masherytypes.MemberIdentifier{MemberId: "member-id"}

	expRvMember := masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "member-id",
			Name: "member-name",
		},
		Address1: "addr-1",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members/member-id").
			WithMethod("get").
			RequestingFields(memberDeepFields).
			WillReturnJsonOf(expRvMember)
	}

	autoTestGet(t,
		memberId,
		expRvMember,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.MemberIdentifier, masherytypes.Member] {
			return cl.GetFullMember
		},
	)
}

func TestCreateMember(t *testing.T) {

	expRvCreatePayload := masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{Id: "", Name: "member-name"},
		Address1:            "addr-1",
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.Member) { t1.Id = "member-id" })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members").
			WithMethod("post").
			RequestingFields(memberFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestRootCreate(t,
		expRvCreatePayload,
		apiResponseJson,
		mockVisitor,
		func(cl Client) ClientExchangeFunc[masherytypes.Member, masherytypes.Member] {
			return cl.CreateMember
		},
	)
}

//func TestUpdateMember(t *testing.T) {
//	payload := masherytypes.Member{
//		AddressableV3Object: masherytypes.AddressableV3Object{
//			Id:   "member-id",
//			Name: "member-name",
//		},
//	}
//
//	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
//		matcher.
//			ForRequestPath("/members/member-id").
//			WithMethod("put").
//			RequestingFields(memberFields).
//			Matching(PayloadMatcher(payload)).
//			WillReturnJsonOf(payload)
//	}
//
//	autoTestUpdate(t,
//		payload,
//		mockVisitor,
//		func(client Client) ClientExchangeFunc[masherytypes.Member, masherytypes.Member] {
//			return client.UpdateMember
//		},
//	)
//}

func TestDeleteMember(t *testing.T) {
	postIdent := masherytypes.MemberIdentifier{
		MemberId: "member-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members/member-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.MemberIdentifier] {
			return cl.DeleteMember
		},
	)
}

func TestListMembers(t *testing.T) {
	mockedResponse := []masherytypes.Member{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "member-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/members").
			WithMethod("get").
			RequestingFields(memberFields).
			WillReturnJsonOf(mockedResponse)
	}

	autoTestRootFetchAll(t,
		mockedResponse,
		mockVisitor,
		func(client Client) ClientArraySupplierFunc[masherytypes.Member] {
			return client.ListMembers
		},
	)
}
