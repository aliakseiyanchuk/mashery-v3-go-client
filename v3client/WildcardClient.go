package v3client

import (
	"context"
	"golang.org/x/sync/semaphore"
	"net/http"
	"net/url"
	"time"
)

type WildcardClientImpl struct {
	transport *HttpTransport
}

func (w WildcardClientImpl) FetchAny(ctx context.Context, resource string, qs *url.Values) (*http.Response, error) {
	destResource := resource
	if qs != nil && len(*qs) > 0 {
		destResource += qs.Encode()
	}

	return w.transport.fetch(ctx, destResource)
}

func (w WildcardClientImpl) DeleteAny(ctx context.Context, resource string) (*http.Response, error) {
	return w.transport.delete(ctx, resource)
}

func (w WildcardClientImpl) PostAny(ctx context.Context, resource string, body interface{}) (*http.Response, error) {
	return w.transport.post(ctx, resource, body)
}

func (w WildcardClientImpl) PutAny(ctx context.Context, resource string, body interface{}) (*http.Response, error) {
	return w.transport.put(ctx, resource, body)
}

// NewWildcardClient creates a "wildcard" client, which will auto-apply access tokens and will throttle the
// calls with the specified QPS.
func NewWildcardClient(p V3AccessTokenProvider, qps int64, travelTimeComp time.Duration) WildcardClient {
	impl := HttpTransport{
		mashEndpoint:  "https://api.mashery.com/v3/rest",
		tokenProvider: p,
		sem:           semaphore.NewWeighted(qps),
		httpCl:        &http.Client{},
		avgNetLatency: travelTimeComp,
	}

	rv := WildcardClientImpl{
		transport: &impl,
	}

	return &rv
}
