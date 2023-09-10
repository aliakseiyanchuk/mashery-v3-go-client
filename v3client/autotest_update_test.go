package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestUpdate[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	autoTestSuccessfulUpdate(t, upsertPayload, rmv, locator)
	autoTestUpdateOn404(t, upsertPayload, rmv, locator)
	autoTestUpdateOn403(t, upsertPayload, rmv, locator)
	autoTestUpdateWithDeveloperOverRate(t, upsertPayload, rmv, locator)
}

func autoTestSuccessfulUpdate[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(upsertPayload, rvVal))

	wm.AssertExpectations(t)
}

func autoTestUpdateOn404[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	_, err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestUpdateOn403[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	_, err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestUpdateWithDeveloperOverRate[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TPayload, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(upsertPayload, rvVal))

	wm.AssertExpectations(t)
}
