package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestBiConsume[TIdent, TPayload any](t *testing.T, ident TIdent, upsertPayload TPayload, rmv BuildVisitor, locator ClientBiConsumerFuncLocator[TIdent, TPayload]) {
	autoTestSuccessfulBiConsume(t, ident, upsertPayload, rmv, locator)
	autoTestBiConsumeOn404(t, ident, upsertPayload, rmv, locator)
	autoTestBiConsumeOn403(t, ident, upsertPayload, rmv, locator)
	autoTestBiConsumeWithDeveloperOverRate(t, ident, upsertPayload, rmv, locator)
}

func autoTestSuccessfulBiConsume[TIdent, TPayload any](t *testing.T, ident TIdent, upsertPayload TPayload, rmv BuildVisitor, locator ClientBiConsumerFuncLocator[TIdent, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	err := f(context.TODO(), ident, upsertPayload)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}

func autoTestBiConsumeOn404[TIdent, TPayload any](t *testing.T, ident TIdent, upsertPayload TPayload, rmv BuildVisitor, locator ClientBiConsumerFuncLocator[TIdent, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	err := f(context.TODO(), ident, upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") >= 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestBiConsumeOn403[TIdent, TPayload any](t *testing.T, ident TIdent, upsertPayload TPayload, rmv BuildVisitor, locator ClientBiConsumerFuncLocator[TIdent, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	err := f(context.TODO(), ident, upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestBiConsumeWithDeveloperOverRate[TIdent, TPayload any](t *testing.T, ident TIdent, upsertPayload TPayload, rmv BuildVisitor, locator ClientBiConsumerFuncLocator[TIdent, TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	err := f(context.TODO(), ident, upsertPayload)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}
