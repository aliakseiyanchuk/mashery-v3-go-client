package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestFetchAll[TIdent, TPayload any](t *testing.T, ident TIdent, expRV []TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TIdent, []TPayload]) {
	autoTestSuccessfulFetchAll(t, ident, expRV, rmv, locator)
	autoTestFetchAllOn404(t, ident, rmv, locator)
	autoTestFetchAllOn403(t, ident, rmv, locator)
	autoTestFetchAllWithDeveloperOverRate(t, ident, expRV, rmv, locator)
}

func autoTestSuccessfulFetchAll[TIdent, TPayload any](t *testing.T, ident TIdent, expRV []TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TIdent, []TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), ident)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}

func autoTestFetchAllOn404[TIdent, TPayload any](t *testing.T, ident TIdent, rmv BuildVisitor, locator ClientExchangeFuncLocator[TIdent, []TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), ident)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	assert.Equal(t, 0, len(rvVal))

	wm.AssertExpectations(t)
}

func autoTestFetchAllOn403[TIdent, TPayload any](t *testing.T, ident TIdent, rmv BuildVisitor, locator ClientExchangeFuncLocator[TIdent, []TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), ident)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestFetchAllWithDeveloperOverRate[TIdent, TPayload any](t *testing.T, ident TIdent, expRV []TPayload, rmv BuildVisitor, locator ClientExchangeFuncLocator[TIdent, []TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), ident)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, rvVal))

	wm.AssertExpectations(t)
}
