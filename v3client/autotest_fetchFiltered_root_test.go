package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestRootFetchFiltered[TPayload any](t *testing.T, param map[string]string, expRV []TPayload, rmv BuildVisitor, locator ClientFilteredArraySupplierFuncLocator[TPayload]) {
	autoTestSuccessfulRootFetchFiltered(t, param, expRV, rmv, locator)
	autoTestRootFetchFilteredOn404(t, param, rmv, locator)
	autoTestRootFetchFilteredOn403(t, param, rmv, locator)
	autoTestRootFetchFilteredWithDeveloperOverRate(t, param, expRV, rmv, locator)
}

func autoTestSuccessfulRootFetchFiltered[TPayload any](t *testing.T, param map[string]string, expRV []TPayload, rmv BuildVisitor, locator ClientFilteredArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootFetchFilteredOn404[TPayload any](t *testing.T, param map[string]string, rmv BuildVisitor, locator ClientFilteredArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), param)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") > 0)
	assert.Equal(t, 0, len(rvVal))

	wm.AssertExpectations(t)
}

func autoTestRootFetchFilteredOn403[TPayload any](t *testing.T, param map[string]string, rmv BuildVisitor, locator ClientFilteredArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), param)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootFetchFilteredWithDeveloperOverRate[TPayload any](t *testing.T, param map[string]string, expRV []TPayload, rmv BuildVisitor, locator ClientFilteredArraySupplierFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
