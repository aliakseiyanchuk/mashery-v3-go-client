package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestRootCount(t *testing.T, expectedCount int64, rmv BuildVisitor, locator ClientSupplierFuncLocator[int64]) {
	autoTestSuccessfulRootCount(t, expectedCount, rmv, locator)
	autoTestRootCountOn404(t, rmv, locator)
	autoTestRootCountOn403(t, rmv, locator)
	autoTestRootCountWithDeveloperOverRate(t, expectedCount, rmv, locator)
}

func autoTestSuccessfulRootCount(t *testing.T, expectedCount int64, rmv BuildVisitor, locator ClientSupplierFuncLocator[int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootCountOn404(t *testing.T, rmv BuildVisitor, locator ClientSupplierFuncLocator[int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "->no such resource") > 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootCountOn403(t *testing.T, rmv BuildVisitor, locator ClientSupplierFuncLocator[int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootCountWithDeveloperOverRate(t *testing.T, expectedCount int64, rmv BuildVisitor, locator ClientSupplierFuncLocator[int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}
