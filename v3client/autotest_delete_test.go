package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestDelete[TIdent any](t *testing.T, ident TIdent, rmv BuildVisitor, locator BiConsumerCanErrLocator[context.Context, TIdent]) {
	autoTestSuccessfulDelete(t, ident, rmv, locator)
	autoTestDeleteOn404(t, ident, rmv, locator)
	autoTestDeleteOn403(t, ident, rmv, locator)
	autoTestDeleteWithDeveloperOverRate(t, ident, rmv, locator)
}

func autoTestSuccessfulDelete[TIdent any](t *testing.T, ident TIdent, rmv BuildVisitor, locator BiConsumerCanErrLocator[context.Context, TIdent]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	err := f(context.TODO(), ident)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}

func autoTestDeleteOn404[TIdent any](t *testing.T, ident TIdent, rmv BuildVisitor, locator BiConsumerCanErrLocator[context.Context, TIdent]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	err := f(context.TODO(), ident)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)

	wm.AssertExpectations(t)
}

func autoTestDeleteOn403[TIdent any](t *testing.T, ident TIdent, rmv BuildVisitor, locator BiConsumerCanErrLocator[context.Context, TIdent]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	err := f(context.TODO(), ident)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)

	wm.AssertExpectations(t)
}

func autoTestDeleteWithDeveloperOverRate[TIdent any](t *testing.T, ident TIdent, rmv BuildVisitor, locator BiConsumerCanErrLocator[context.Context, TIdent]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	err := f(context.TODO(), ident)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}
