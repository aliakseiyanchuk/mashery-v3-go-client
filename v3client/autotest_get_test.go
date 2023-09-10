package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestGet[TId, IPayload any](t *testing.T, id TId, expRV IPayload, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, IPayload]) {
	autoTestSuccessfulGet(t, id, expRV, rmv, locator)
	autoTestGetOn404(t, id, rmv, locator)
	autoTestGetOn403(t, id, rmv, locator)
	autoTestGetWithDeveloperOverRate(t, id, expRV, rmv, locator)
}

func autoTestSuccessfulGet[TId, TPayload any](t *testing.T, id TId, expRV TPayload, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, exist, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestSuccessfulGetWithBadRequestAutoRetry[TId, TPayload any](t *testing.T, id TId, expRV TPayload, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockBadRequestFollowedByReturnedData()
	f := locator(cl)

	rvVal, exist, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestGetOn404[TId, TPayload any](t *testing.T, id TId, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	_, exist, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.False(t, exist)

	wm.AssertExpectations(t)
}

func autoTestGetOn403[TId, TPayload any](t *testing.T, id TId, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	_, exist, err := f(context.TODO(), id)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.False(t, exist)

	wm.AssertExpectations(t)
}

func autoTestGetWithDeveloperOverRate[TId, TPayload any](t *testing.T, id TId, expRV TPayload, rmv BuildVisitor, locator ClientBoolExchangeFuncLocator[TId, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, exist, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
