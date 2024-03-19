package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func autoTestConsume[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientConsumerFuncLocator[TPayload]) {
	autoTestSuccessfulConsume(t, upsertPayload, rmv, locator)
	autoTestConsumeOn404(t, upsertPayload, rmv, locator)
	autoTestConsumeOn403(t, upsertPayload, rmv, locator)
	autoTestConsumeWithDeveloperOverRate(t, upsertPayload, rmv, locator)
}

func autoTestSuccessfulConsume[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientConsumerFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}

func autoTestConsumeOn404[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientConsumerFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, strings.Index(err.Error(), "error code 404 is not an expected response to this request") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestConsumeOn403[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientConsumerFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	err := f(context.TODO(), upsertPayload)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "not authorized") > 0)
	//assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestConsumeWithDeveloperOverRate[TPayload any](t *testing.T, upsertPayload TPayload, rmv BuildVisitor, locator ClientConsumerFuncLocator[TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	err := f(context.TODO(), upsertPayload)
	assert.Nil(t, err)

	wm.AssertExpectations(t)
}
