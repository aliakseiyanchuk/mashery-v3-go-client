package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestRootGet[IPayload any](t *testing.T, expRV IPayload, rmv BuildVisitor, locator ClientSupplierFuncLocator[IPayload]) {
	autoTestSuccessfulRootGet(t, expRV, rmv, locator)
	autoTestRootGetOn404(t, rmv, locator)
	autoTestRootGetOn403(t, rmv, locator)
	autoTestRootGetWithDeveloperOverRate(t, expRV, rmv, locator)
}

func autoTestSuccessfulRootGet[TPayload any](t *testing.T, expRV TPayload, rmv BuildVisitor, locator ClientSupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootGetOn404[TPayload any](t *testing.T, rmv BuildVisitor, locator ClientSupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootGetOn403[TPayload any](t *testing.T, rmv BuildVisitor, locator ClientSupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootGetWithDeveloperOverRate[TPayload any](t *testing.T, expRV TPayload, rmv BuildVisitor, locator ClientSupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
