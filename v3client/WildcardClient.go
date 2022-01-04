package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"golang.org/x/sync/semaphore"
	"net/http"
	"net/url"
	"time"
)

type WildcardClientImpl struct {
	transport *transport.HttpTransport
}

func (w WildcardClientImpl) FetchAny(ctx context.Context, resource string, qs *url.Values) (*http.Response, error) {
	destResource := resource
	if qs != nil && len(*qs) > 0 {
		destResource += "?" + qs.Encode()
	}

	return w.transport.Fetch(ctx, destResource)
}

func (w WildcardClientImpl) DeleteAny(ctx context.Context, resource string) (*http.Response, error) {
	return w.transport.Delete(ctx, resource)
}

func (w WildcardClientImpl) PostAny(ctx context.Context, resource string, body interface{}) (*http.Response, error) {
	return w.transport.Post(ctx, resource, body)
}

func (w WildcardClientImpl) PutAny(ctx context.Context, resource string, body interface{}) (*http.Response, error) {
	return w.transport.Put(ctx, resource, body)
}

// NewWildcardClient creates a "wildcard" client, which will auto-apply access tokens and will throttle the
// calls with the specified QPS.
func NewWildcardClient(p V3AccessTokenProvider, qps int64, travelTimeComp time.Duration) WildcardClient {
	impl := transport.HttpTransport{
		MashEndpoint: "https://api.mashery.com/v3/rest",
		Authorizer:   p,
		Sem:          semaphore.NewWeighted(qps),
		HttpClient: &http.Client{
			Timeout: time.Second * 60,
		},
		AvgNetLatency: travelTimeComp,
	}

	rv := WildcardClientImpl{
		transport: &impl,
	}

	return &rv
}
