package transport

import (
	"context"
)

// ExchangeListener receive a notification of a raw exchange: what was actually sent to Mashery, and which response
// was received
type ExchangeListener func(ctx context.Context, req *WrappedRequest, res *WrappedResponse, err error)
