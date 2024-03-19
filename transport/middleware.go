package transport

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"time"
)

const (
	SendErrorOn404 = ".send.error.on.404"
	RetryOn400     = ".retry.on.400"
	LeafExecutor   = ".leaf.executor"
)

func ThrottleFunc(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(c.DelayBeforeCall()):
		return next(ctx, c)
	}
}

func BackOffOnDeveloperOverQPSFunc(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	for i := 0; i < 10; i++ {
		// If there's an error, return
		if wr, err := next(ctx, c); err != nil {
			return wr, err
		} else if wr.StatusCode == 403 {
			// Retry if developer over QPS has been received.
			if str := wr.Header.Get("X-Mashery-Error-Code"); str == "ERR_403_DEVELOPER_OVER_QPS" {
				d := time.Duration(1+i) * time.Second
				time.Sleep(d)
				continue
			} else {
				return wr, err
			}
		} else {
			return wr, err
		}
	}

	return nil, errors.New("operation unsuccessful after all available retries")
}

func BreakOnDeveloperOverRateFunc(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	if wr, err := next(ctx, c); err != nil {
		return wr, err
	} else if str := wr.Header.Get("X-Mashery-Error-Code"); wr.StatusCode == 403 && str == "ERR_403_DEVELOPER_OVER_RATE" {
		// Break calling if developer call has been exhausted
		if str := wr.Header.Get("X-Mashery-Error-Code"); str == "ERR_403_DEVELOPER_OVER_RATE" {
			return wr, errors.New("further operations are impossible until developer rate is reset")
		} else {
			return wr, err
		}
	} else {
		return wr, err
	}
}

func boolKey(ctx context.Context, key string) bool {
	if v := ctx.Value(key); v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
	}

	return false
}

func ErrorOn404Func(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	if wr, err := next(ctx, c); err != nil {
		return wr, err
	} else if wr.StatusCode == 404 && boolKey(ctx, SendErrorOn404) {
		return wr, errors.New("error code 404 is not an expected response to this request")
	} else {
		return wr, err
	}
}

func RetryOn400Func(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	if boolKey(ctx, RetryOn400) {
		for i := 0; i < 5; i++ {
			if wr, err := next(ctx, c); err != nil {
				return wr, err
			} else if wr.StatusCode == 400 {
				time.Sleep(time.Second*3 + time.Duration(i*2))
				continue
			} else {
				return wr, err
			}
		}

		return nil, errors.New("attempts to perform this operation were unsuccessful after all possible retries")
	} else {
		return next(ctx, c)
	}
}

func UnmarshalServerError(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	if wr, err := next(ctx, c); err != nil {
		return wr, err
	} else if wr.StatusCode > 299 && wr.StatusCode != 404 {
		return wr, v3BasicError(wr)
	} else {
		return wr, err
	}
}

func EnsureBodyWasRead(ctx context.Context, c *HttpTransport, next MiddlewareFunc) (*WrappedResponse, error) {
	if wr, err := next(ctx, c); err != nil {
		return wr, err
	} else if _, readErr := wr.Body(); readErr != nil {
		return wr, &errwrap.WrappedError{
			Cause: readErr,
		}
	} else {
		return wr, err
	}
}

func ExecuteFunction(ctx context.Context, c *HttpTransport) (*WrappedResponse, error) {
	if f := ctx.Value(LeafExecutor); f != nil {
		if mf, ok := f.(MiddlewareFunc); ok {
			return mf(ctx, c)
		}
	}

	return nil, errors.New("no executable function in this context")
}
