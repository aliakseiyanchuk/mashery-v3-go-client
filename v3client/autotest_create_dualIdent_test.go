package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func autoTestCreateDualIdent[TParent1, TParent2, TPayload any](t *testing.T, pIdent1 TParent1, pIdent2 TParent2, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent1, TParent2, *TPayload]) {
	autoTestSuccessfulCreateDualIdent(t, pIdent1, pIdent2, expRV, rmv, locator)
	autoTestCreateDualIdentOn404(t, pIdent1, pIdent2, rmv, locator)
	autoTestCreateDualIdentOn403(t, pIdent1, pIdent2, rmv, locator)
	autoTestCreateDualIdentWithDeveloperOverRate(t, pIdent1, pIdent2, expRV, rmv, locator)
}

func autoTestSuccessfulCreateDualIdent[TParent1, TParent2, TPayload any](t *testing.T, pIdent1 TParent1, pIdent2 TParent2, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent1, TParent2, *TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedData()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent1, pIdent2)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, *rvVal))

	wm.AssertExpectations(t)
}

func autoTestCreateDualIdentOn404[TParent1, TParent2, TPayload any](t *testing.T, pIdent TParent1, pIdent2 TParent2, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent1, TParent2, *TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned404()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, pIdent2)
	assert.Nil(t, err)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestCreateDualIdentOn403[TParent1, TParent2, TPayload any](t *testing.T, pIdent TParent1, pIdent2 TParent2, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent1, TParent2, *TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturned403()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, pIdent2)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	assert.True(t, strings.Index(err.Error(), "Not Authorized") > 0)
	assert.Nil(t, rvVal)

	wm.AssertExpectations(t)
}

func autoTestCreateDualIdentWithDeveloperOverRate[TParent1, TParent2, TPayload any](t *testing.T, pIdent TParent1, pIdent2 TParent2, expRV TPayload, rmv BuildVisitor, locator ClientDualExchangeFuncLocator[TParent1, TParent2, *TPayload]) {
	cl, wm := RequestMockBuilder(rmv).MockReturnedDeveloperOverRate()
	f := locator(cl)

	rvVal, err := f(context.TODO(), pIdent, pIdent2)
	assert.Nil(t, err)
	assert.NotNil(t, rvVal)
	assert.True(t, reflect.DeepEqual(expRV, *rvVal))

	wm.AssertExpectations(t)
}
