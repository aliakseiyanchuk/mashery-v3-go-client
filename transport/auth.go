package transport

import "context"

type Authorizer interface {
	HeaderAuthorization(ctx context.Context) (map[string]string, error)
	QueryStringAuthorization(ctx context.Context) (map[string]string, error)
	Close()
}
