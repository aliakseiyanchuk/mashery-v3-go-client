package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestRootCreate[TPayload any](t *testing.T, upsertPayload TPayload, expRV TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	autoTestRootAsymmetricCreate[TPayload, TPayload](t, upsertPayload, expRV, rmv, locator)
}

func autoTestRootAsymmetricCreate[TPayload, TReturn any](t *testing.T, upsertPayload TPayload, expRV TReturn, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TReturn]) {
	autoTestSuccessfulRootCreate(t, upsertPayload, expRV, rmv, locator)
	autoTestRootCreateOn404(t, upsertPayload, rmv, locator)
	autoTestRootCreateOn403(t, upsertPayload, rmv, locator)
	autoTestRootCreateWithDeveloperOverRate(t, upsertPayload, expRV, rmv, locator)
}

func autoTestSuccessfulRootCreate[TPayload, TReturn any](t *testing.T, upsertPayload TPayload, expRV TReturn, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TReturn]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootCreateOn404[TPayload, TReturn any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TReturn]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	_, err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootCreateOn403[TPayload, TReturn any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TReturn]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	_, err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootCreateWithDeveloperOverRate[TPayload, TReturn any](t *testing.T, upsertPayload TPayload, expRV TReturn, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TReturn]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
