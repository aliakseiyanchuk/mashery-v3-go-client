package v3client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type V3OAuthHelper struct {
	client        transport.HttpExecutor
	TokenEndpoint string
}

type OAuthHelperParams struct {
	transport.HTTPClientParams
	MasheryTokenEndpoint string
}

func (ohp *OAuthHelperParams) FillDefaults() {
	if ohp.TLSConfig == nil && !ohp.TLSConfigDelegateSystem {
		ohp.TLSConfig = transport.DefaultTLSConfig()
	}

	if ohp.Timeout == 0 {
		ohp.Timeout = time.Second * 30
	}
	if len(ohp.MasheryTokenEndpoint) == 0 {
		ohp.MasheryTokenEndpoint = MasheryTokenEndpoint
	}
}

// NewOAuthHelper creates an instance of a helper that could be used directly
func NewOAuthHelper(params OAuthHelperParams) *V3OAuthHelper {
	params.FillDefaults()

	rv := V3OAuthHelper{
		client:        params.CreateHttpExecutor(),
		TokenEndpoint: params.MasheryTokenEndpoint,
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
		return nil, errors.New(fmt.Sprintf("server call was rejected (%s)", err.Error()))
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("server returned unexpected error code %d with message %s", resp.StatusCode, bodyText))
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read response body (%s)", err.Error()))
	}

	procResp := masherytypes.TimedAccessTokenResponse{
		Obtained: time.Now(),
		QPS:      creds.MaxQPS,
	}
	err = json.Unmarshal(bodyText, &procResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not unmarshal access token response (%s)", err.Error()))
	}

	procResp.ServerTime = ResponseDate(resp)

	return &procResp, nil
}
