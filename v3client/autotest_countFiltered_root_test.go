package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestRootFilteredCount(t *testing.T, params map[string]string, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[map[string]string, int64]) {
	autoTestSuccessfulFilteredRootCount(t, params, expectedCount, rmv, locator)
	autoTestRootFilteredCountOn404(t, params, rmv, locator)
	autoTestRootFilteredCountOn403(t, params, rmv, locator)
	autoTestRootFilteredCountWithDeveloperOverRate(t, params, expectedCount, rmv, locator)
}

func autoTestSuccessfulFilteredRootCount(t *testing.T, params map[string]string, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[map[string]string, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), params)
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootFilteredCountOn404(t *testing.T, params map[string]string, rmv BuildVisitor, locator ClientExchangeFuncLocator[map[string]string, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), params)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootFilteredCountOn403(t *testing.T, params map[string]string, rmv BuildVisitor, locator ClientExchangeFuncLocator[map[string]string, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), params)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestRootFilteredCountWithDeveloperOverRate(t *testing.T, params map[string]string, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[map[string]string, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), params)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}
