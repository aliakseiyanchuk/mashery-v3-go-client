package mashery_v3_go_client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

type ClientCredentialsProvider struct {
	creds MasheryV3Credentials

	TokenEndpoint string
	Response      *TimedAccessTokenResponse

	client *http.Client
}

func NewClientCredentialsProvider(creds MasheryV3Credentials) *ClientCredentialsProvider {
	return NewLiveCredentialsProviderUsing(creds, "https://api.mashery.com/v3/token")
}

func NewLiveCredentialsProviderUsing(creds MasheryV3Credentials, endp string) *ClientCredentialsProvider {
	retVal := ClientCredentialsProvider{
		creds:         creds,
		TokenEndpoint: endp,
		client:        &http.Client{},
	}

	return &retVal
}

func (c *MasheryV3Credentials) Inherit(other *MasheryV3Credentials) {
	if len(other.AreaId) > 0 {
		c.AreaId = other.AreaId
	}
	if len(other.ApiKey) > 0 {
		c.ApiKey = other.ApiKey

		// Max QPS wil be inherited only if the key is supplied.
		if other.MaxQPS > 0 {
			c.MaxQPS = other.MaxQPS
		}
	}

	if len(other.Secret) > 0 {
		c.Secret = other.Secret
	}

	if len(other.Username) > 0 {
		c.Username = other.Username
	}

	if len(other.Password) > 0 {
		c.Password = other.Password
	}
}

func (lcp *ClientCredentialsProvider) Refresh() error {
	if lcp.Response != nil {
		data := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {lcp.Response.RefreshToken},
		}

		return lcp.postForToken(data)
	} else {
		return errors.New("cannot refresh without having a previous response")
	}
}

func (lcp *ClientCredentialsProvider) doRetrieve() error {
	if lcp.Response == nil || lcp.Response.Expired() {

		data := url.Values{
			"grant_type": {"password"},
			"username":   {lcp.creds.Username},
			"password":   {lcp.creds.Password},
			"scope":      {lcp.creds.AreaId},
		}

		return lcp.postForToken(data)
	}

	if lcp.Response == nil {
		return errors.New("Login failed to establish a response")
	} else {
		return nil
	}
}

func (lcp *ClientCredentialsProvider) postForToken(data url.Values) error {

	req, _ := http.NewRequest("POST", lcp.TokenEndpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(lcp.creds.ApiKey, lcp.creds.Secret)

	defer req.Body.Close()
	resp, err := lcp.client.Do(req)
	if err != nil {
		return errors.New("Failed to retrieve response from the server")
	}
	if resp.StatusCode != 200 {
		return errors.New("Server returned unexpected error code")
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Failed to read response body")
	}

	procResp := TimedAccessTokenResponse{
		Obtained: time.Now(),
	}
	err = json.Unmarshal(bodyText, &procResp)
	if err != nil {
		return errors.New("Could not unmarshal access token response")
	}

	lcp.Response = &procResp

	return nil
}

func (lcp *ClientCredentialsProvider) AccessToken() (string, error) {
	err := lcp.doRetrieve()
	if err != nil {
		return "", err
	} else {
		return lcp.Response.AccessToken, nil
	}
}

func (lcp *ClientCredentialsProvider) TokenData() (*TimedAccessTokenResponse, error) {
	err := lcp.doRetrieve()
	if err != nil {
		return nil, err
	} else {
		return lcp.Response, nil
	}
}
