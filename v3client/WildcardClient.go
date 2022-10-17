package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

type WildcardClientImpl struct {
	transport *transport.HttpTransport
}

func (w *WildcardClientImpl) FetchAny(ctx context.Context, resource string, qs *url.Values) (*transport.WrappedResponse, error) {
	destResource := resource
	if qs != nil && len(*qs) > 0 {
		destResource += "?" + qs.Encode()
	}

	return w.transport.Fetch(ctx, destResource)
}

func (w *WildcardClientImpl) DeleteAny(ctx context.Context, resource string) (*transport.WrappedResponse, error) {
	return w.transport.Delete(ctx, resource)
}

func (w *WildcardClientImpl) PostAny(ctx context.Context, resource string, body interface{}) (*transport.WrappedResponse, error) {
	return w.transport.Post(ctx, resource, body)
}

func (w *WildcardClientImpl) PutAny(ctx context.Context, resource string, body interface{}) (*transport.WrappedResponse, error) {
	return w.transport.Put(ctx, resource, body)
}

func (w *WildcardClientImpl) Close(ctx context.Context) {
	w.transport.HttpClient.CloseIdleConnections()
}

// NewWildcardClient creates a "wildcard" client, which will auto-apply access tokens and will throttle the
// calls with the specified QPS.
func NewWildcardClient(params Params) WildcardClient {
	params.FillDefaults()
	impl := createHTTPTransport(params)

	rv := WildcardClientImpl{
		transport: &impl,
	}

	return &rv
}
