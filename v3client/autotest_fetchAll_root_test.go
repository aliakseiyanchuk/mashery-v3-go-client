package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestRootFetchAll[IPayload any](t *testing.T, expRV []IPayload, rmv BuildVisitor, locator ClientArraySupplierFuncLocator[IPayload]) {
	autoTestSuccessfulRootFetchAll(t, expRV, rmv, locator)
	autoTestRootFetchAllOn404(t, rmv, locator)
	autoTestRootFetchAllOn403(t, rmv, locator)
	autoTestRootFetchAllWithDeveloperOverRate(t, expRV, rmv, locator)
}

func autoTestSuccessfulRootFetchAll[TPayload any](t *testing.T, expRV []TPayload, rmv BuildVisitor, locator ClientArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootFetchAllOn404[TPayload any](t *testing.T, rmv BuildVisitor, locator ClientArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") > 0)
	assert.Equal(t, 0, len(rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootFetchAllOn403[TPayload any](t *testing.T, rmv BuildVisitor, locator ClientArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootFetchAllWithDeveloperOverRate[TPayload any](t *testing.T, expRV []TPayload, rmv BuildVisitor, locator ClientArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
