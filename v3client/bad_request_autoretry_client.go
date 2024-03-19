package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

type ExchangeFunc[TRequest, TResponse any] func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, error)
type ExchangeBoolFunc[TRequest, TResponse any] func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, bool, error)
type BiExchangeFunc[TRequest1, TRequest2, TResponse any] func(ctx context.Context, req1 TRequest1, req TRequest2, c *transport.HttpTransport) (TResponse, error)

func autoRetryBadRequest[TRequest, TResponse any](rawFunc ExchangeFunc[TRequest, TResponse]) ExchangeFunc[TRequest, TResponse] {
	return func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, error) {
		passCtx := context.WithValue(ctx, transport.RetryOn400, true)
		return rawFunc(passCtx, req, c)
	}
}

func autoRetryBadBiRequest[TIdent, TRequest, TResponse any](rawFunc BiExchangeFunc[TIdent, TRequest, TResponse]) BiExchangeFunc[TIdent, TRequest, TResponse] {
	return func(ctx context.Context, ident TIdent, req TRequest, c *transport.HttpTransport) (TResponse, error) {
		passCtx := context.WithValue(ctx, transport.RetryOn400, true)
		return rawFunc(passCtx, ident, req, c)
	}
}

func autoRetryBadGetRequest[TRequest, TResponse any](rawFunc ExchangeBoolFunc[TRequest, TResponse]) ExchangeBoolFunc[TRequest, TResponse] {
	return func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, bool, error) {
		passCtx := context.WithValue(ctx, transport.RetryOn400, true)
		return rawFunc(passCtx, req, c)
	}
}

func autoRetryBiBadRequest[TRequest1, TRequest2, TResponse any](rawFunc BiExchangeFunc[TRequest1, TRequest2, TResponse]) BiExchangeFunc[TRequest1, TRequest2, TResponse] {
	return func(ctx context.Context, req1 TRequest1, req2 TRequest2, c *transport.HttpTransport) (TResponse, error) {
		passCtx := context.WithValue(ctx, transport.RetryOn400, true)
		return rawFunc(passCtx, req1, req2, c)
	}
}
