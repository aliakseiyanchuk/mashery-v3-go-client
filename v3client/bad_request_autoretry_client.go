package v3client

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"time"
)

type ExchangeFunc[TRequest, TResponse any] func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, error)
type ExchangeBoolFunc[TRequest, TResponse any] func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, bool, error)
type BiExchangeFunc[TRequest1, TRequest2, TResponse any] func(ctx context.Context, req1 TRequest1, req TRequest2, c *transport.HttpTransport) (TResponse, error)

func autoRetryBadRequest[TRequest, TResponse any](rawFunc ExchangeFunc[TRequest, TResponse]) ExchangeFunc[TRequest, TResponse] {
	return func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, error) {
		var resp TResponse
		var err error

		for i := 0; i < 5; i++ {
			if resp, err = rawFunc(ctx, req, c); err != nil {
				if we, ok := unwrapUndeterminedError(err); ok {
					if we.Code == 400 {
						time.Sleep(time.Second*3 + time.Duration(i*2))
						continue
					}
				}
			} else if err == nil {
				break
			}
		}

		return resp, err
	}
}

func autoRetryBadBiRequest[TIdent, TRequest, TResponse any](rawFunc BiExchangeFunc[TIdent, TRequest, TResponse]) BiExchangeFunc[TIdent, TRequest, TResponse] {
	return func(ctx context.Context, ident TIdent, req TRequest, c *transport.HttpTransport) (TResponse, error) {
		var resp TResponse
		var err error

		for i := 0; i < 5; i++ {
			if resp, err = rawFunc(ctx, ident, req, c); err != nil {
				if we, ok := unwrapUndeterminedError(err); ok {
					if we.Code == 400 {
						time.Sleep(time.Second*3 + time.Duration(i*2))
						continue
					}
				}
			} else if err == nil {
				break
			}
		}

		return resp, err
	}
}

func autoRetryBadGetRequest[TRequest, TResponse any](rawFunc ExchangeBoolFunc[TRequest, TResponse]) ExchangeBoolFunc[TRequest, TResponse] {
	return func(ctx context.Context, req TRequest, c *transport.HttpTransport) (TResponse, bool, error) {
		var resp TResponse
		var err error
		var exists bool

		for i := 0; i < 5; i++ {
			if resp, exists, err = rawFunc(ctx, req, c); err != nil {
				if we, ok := unwrapUndeterminedError(err); ok {
					if we.Code == 400 {
						time.Sleep(time.Second*3 + time.Duration(i*2))
						continue
					}
				}
			} else if err == nil {
				break
			}
		}

		return resp, exists, err
	}
}

func autoRetryBiBadRequest[TRequest1, TRequest2, TResponse any](rawFunc BiExchangeFunc[TRequest1, TRequest2, TResponse]) BiExchangeFunc[TRequest1, TRequest2, TResponse] {
	return func(ctx context.Context, req1 TRequest1, req2 TRequest2, c *transport.HttpTransport) (TResponse, error) {
		var resp TResponse
		var err error

		for i := 0; i < 5; i++ {
			if resp, err = rawFunc(ctx, req1, req2, c); err != nil {
				if we, ok := unwrapUndeterminedError(err); ok {
					if we.Code == 400 {
						time.Sleep(time.Second*3 + time.Duration(i*2))
						continue
					}
				}
			} else if err == nil {
				break
			}
		}

		return resp, err
	}
}

func unwrapUndeterminedError(err error) (*masherytypes.V3UndeterminedError, bool) {
	var wrappedError *errwrap.WrappedError
	if errors.As(err, &wrappedError) {
		var unspecifiedError *masherytypes.V3UndeterminedError
		if errors.As(wrappedError.Cause, &unspecifiedError) {
			return unspecifiedError, true
		}
	}

	return nil, false
}
