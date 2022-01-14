package v3client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type V3OAuthHelper struct {
	client        *http.Client
	TokenEndpoint string
}

// NewOAuthHelper creates an instance of a helper that could be used directly
func NewOAuthHelper(tlsCfg *tls.Config) *V3OAuthHelper {
	rv := V3OAuthHelper{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
		},
		TokenEndpoint: MasheryTokenEndpoint,
	}

	return &rv
}

func (lcp *V3OAuthHelper) RetrieveAccessTokenFor(creds *MasheryV3Credentials) (*masherytypes.TimedAccessTokenResponse, error) {
	data := url.Values{
		"grant_type": {"password"},
		"username":   {creds.Username},
		"password":   {creds.Password},
		"scope":      {creds.AreaId},
	}

	return lcp.postForToken(data, creds)
}

func (lcp *V3OAuthHelper) ExchangeRefreshToken(creds *MasheryV3Credentials, refreshToken string) (*masherytypes.TimedAccessTokenResponse, error) {

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	return lcp.postForToken(data, creds)
}

func (lcp *V3OAuthHelper) postForToken(data url.Values, creds *MasheryV3Credentials) (*masherytypes.TimedAccessTokenResponse, error) {

	req, _ := http.NewRequest("POST", lcp.TokenEndpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(creds.ApiKey, creds.Secret)

	defer req.Body.Close()
	resp, err := lcp.client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to retrieve response from the server")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Server returned unexpected error code")
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read response body")
	}

	procResp := masherytypes.TimedAccessTokenResponse{
		Obtained: time.Now(),
		QPS:      creds.MaxQPS,
	}
	err = json.Unmarshal(bodyText, &procResp)
	if err != nil {
		return nil, errors.New("Could not unmarshal access token response")
	}

	procResp.ServerTime = ResponseDate(resp)

	return &procResp, nil
}
