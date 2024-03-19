package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestCreate[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	autoTestSuccessfulCreate(t, pIdent, upsertPayload, expRV, rmv, locator)
	autoTestCreateOn404(t, pIdent, upsertPayload, rmv, locator)
	autoTestCreateOn403(t, pIdent, upsertPayload, rmv, locator)
	autoTestCreateWithDeveloperOverRate(t, pIdent, upsertPayload, expRV, rmv, locator)
}

func autoTestSuccessfulCreate[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestSuccessfulCreateWithBadRequestAutoRetry[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockBadRequestFollowedByReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestCreateOn404[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	_, err := f(context.TODO(), pIdent, upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestCreateOn403[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	_, err := f(context.TODO(), pIdent, upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestCreateWithDeveloperOverRate[TParent, TPayload any](t *testing.T, pIdent TParent, upsertPayload TPayload, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent, TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
