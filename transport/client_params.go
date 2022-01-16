package transport

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HTTPClientParams struct {
	TLSConfig *tls.Config
	Timeout   time.Duration

	ProxyServer          *url.URL
	ProxyAuthType        string
	ProxyAuthCredentials string
}

type httpProxyFunction func(*http.Request) (*url.URL, error)

func (p *HTTPClientParams) determineProxyAuth() http.Header {
	if len(p.ProxyAuthType) > 0 && len(p.ProxyAuthCredentials) > 0 {
		rv := http.Header{}
		rv.Add("Proxy-Authorization", fmt.Sprintf("%s %s", p.ProxyAuthType, p.ProxyAuthCredentials))

		return rv
	} else {
		return nil
	}
}

func (p *HTTPClientParams) determineProxyMethod() httpProxyFunction {
	if p.ProxyServer != nil {
		return http.ProxyFromEnvironment
	} else {
		return p.yieldFixedProxyServer
	}
}

func (p *HTTPClientParams) yieldFixedProxyServer(_ *http.Request) (*url.URL, error) {
	return p.ProxyServer, nil
}

func (p *HTTPClientParams) CreateClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    p.TLSConfig,
			Proxy:              p.determineProxyMethod(),
			ProxyConnectHeader: p.determineProxyAuth(),
		},
		Timeout: p.Timeout,
	}
}
