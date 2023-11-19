package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestCount[TId any](t *testing.T, id TId, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, int64]) {
	autoTestSuccessfulCount(t, id, expectedCount, rmv, locator)
	autoTestCountOn404(t, id, rmv, locator)
	autoTestCountOn403(t, id, rmv, locator)
	autoTestCountWithDeveloperOverRate(t, id, expectedCount, rmv, locator)
}

func autoTestSuccessfulCount[TId any](t *testing.T, id TId, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}

func autoTestCountOn404[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestCountOn403[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	assert.Equal(t, int64(-1), rvVal)

	wm.AssertExpectations(t)
}

func autoTestCountWithDeveloperOverRate[TId any](t *testing.T, id TId, expectedCount int64, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, int64]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.Equal(t, expectedCount, rvVal)

	wm.AssertExpectations(t)
}
