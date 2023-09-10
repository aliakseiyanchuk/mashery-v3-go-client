package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestExists[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, bool]) {
	autoTestSuccessfulExists(t, id, rmv, locator)
	autoTestExistsOn404(t, id, rmv, locator)
	autoTestExistsOn403(t, id, rmv, locator)
	autoTestExistsWithDeveloperOverRate(t, id, rmv, locator)
}

func autoTestSuccessfulExists[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, bool]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.True(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestExistsOn404[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, bool]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.False(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestExistsOn403[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, bool]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	assert.False(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestExistsWithDeveloperOverRate[TId any](t *testing.T, id TId, rmv BuildVisitor, locator ClientExchangeFuncLocator[TId, bool]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), id)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, rvVal)

	wm.AssertExpectations(t)
}
