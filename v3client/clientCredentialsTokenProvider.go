package v3client

import (
	"crypto/tls"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"net/http"
	"time"
)

type MasheryV3Credentials struct {
	AreaId   string `yaml:"areaId,omitempty"`
	ApiKey   string `yaml:"apiKey,omitempty"`
	Secret   string `yaml:"secret,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	MaxQPS   int    `yaml:"maxQPS,omitempty"`
}

func (m *MasheryV3Credentials) FullySpecified() bool {
	return len(m.AreaId) > 0 &&
		len(m.ApiKey) > 0 &&
		len(m.Secret) > 0 &&
		len(m.Username) > 0 &&
		len(m.Password) > 0
}

//------------------------------------------------------------------------
// Abstract credentials provider

type ClientCredentialsProvider struct {
	V3OAuthHelper

	Response            *masherytypes.TimedAccessTokenResponse
	tokenFile           string
	credentials         MasheryV3Credentials
	comm                chan int
	asyncRefreshRunning bool

	postRefreshAction func()
}

func NewClientCredentialsProvider(credentials MasheryV3Credentials, tlsCfg *tls.Config) *ClientCredentialsProvider {
	return NewLiveCredentialsProviderUsing(credentials, MasheryTokenEndpoint, tlsCfg)
}

func NewLiveCredentialsProviderUsing(credentials MasheryV3Credentials, endpoint string, tlsCfg *tls.Config) *ClientCredentialsProvider {
	if tlsCfg == nil {
		panic("nil tls configuration is not allowed")
	}

	retVal := ClientCredentialsProvider{
		V3OAuthHelper: V3OAuthHelper{
			TokenEndpoint: endpoint,
			client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsCfg,
				},
				Timeout: time.Second * 30,
			},
		},

		credentials: credentials,
		comm:        make(chan int),
	}

	return &retVal
}

func (m *MasheryV3Credentials) Inherit(other *MasheryV3Credentials) {
	if len(other.AreaId) > 0 {
		m.AreaId = other.AreaId
	}
	if len(other.ApiKey) > 0 {
		m.ApiKey = other.ApiKey

		// Max QPS wil be inherited only if the key is supplied.
		if other.MaxQPS > 0 {
			m.MaxQPS = other.MaxQPS
		}
	}

	if len(other.Secret) > 0 {
		m.Secret = other.Secret
	}

	if len(other.Username) > 0 {
		m.Username = other.Username
	}

	if len(other.Password) > 0 {
		m.Password = other.Password
	}
}

func (lcp *ClientCredentialsProvider) OnPostRefresh(f func()) {
	lcp.postRefreshAction = f
}

func (lcp *ClientCredentialsProvider) EnsureRefresh() {
	if !lcp.asyncRefreshRunning {
		go lcp.doEnsureRefresh()
		lcp.asyncRefreshRunning = true
	}
}

func (lcp *ClientCredentialsProvider) Close() {
	// Send the command to the refresh thread if it is running
	if lcp.asyncRefreshRunning {
		lcp.comm <- 1
	}
}

func (lcp *ClientCredentialsProvider) doEnsureRefresh() {
	for lcp.Response != nil {
		waitDur := lcp.Response.ExpiryTime().Sub(time.Now())
		waitDur -= time.Minute * 5

		select {
		case <-lcp.comm:
			lcp.asyncRefreshRunning = false
			return
		case <-time.After(waitDur):
			err := lcp.Refresh()
			if err != nil {
				lcp.asyncRefreshRunning = false
				return
			}
		}

		if lcp.postRefreshAction != nil {
			lcp.postRefreshAction()
		}
	}
}

func ResponseDate(resp *http.Response) time.Time {
	if val := resp.Header.Get("Date"); len(val) > 0 {
		if t, err := time.Parse(time.RFC1123, val); err == nil {
			return t
		}
	}

	return time.Unix(0, 0)
}

func (lcp *ClientCredentialsProvider) Refresh() error {
	if lcp.Response == nil {
		return errors.New("cannot refresh without having a previous response")
	} else if lcp.Response.Expired() {
		return errors.New("refresh token has already expired")
	}

	resp, err := lcp.ExchangeRefreshToken(&lcp.credentials, lcp.Response.RefreshToken)
	lcp.Response = resp

	return err
}

func (lcp *ClientCredentialsProvider) TokenData() (*masherytypes.TimedAccessTokenResponse, error) {
	if lcp.Response != nil && !lcp.Response.Expired() {
		return lcp.Response, nil
	} else {
		resp, err := lcp.RetrieveAccessTokenFor(&lcp.credentials)
		lcp.Response = resp

		return resp, err
	}
}

func (lcp *ClientCredentialsProvider) AccessToken() (string, error) {
	if dat, err := lcp.TokenData(); err != nil {
		return "", err
	} else if dat == nil {
		return "", errors.New("empty token data returned while trying to provide access token")
	} else {
		return dat.AccessToken, nil
	}
}

func (lcp *ClientCredentialsProvider) HeaderAuthorization() (map[string]string, error) {
	var token string
	var err error
	var rv map[string]string

	if token, err = lcp.AccessToken(); err == nil {
		rv["Authorization"] = token
	}

	return rv, err
}

func (lcp *ClientCredentialsProvider) QueryStringAuthorization() (map[string]string, error) {
	var emptyMap map[string]string

	return emptyMap, nil
}
