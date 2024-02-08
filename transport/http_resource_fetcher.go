package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TokenFeederResponse struct {
	// Token contains token that can be used for authentication and authorization
	Token string `json:"access_token"`
	// Expiry contains the expiry time of this token.
	Expiry string `json:"expiry"`
	// ExpiryEpoch contains epoch time the token will expire and will no longer be usable
	ExpiryEpoch int64 `json:"expiry_epoch"`
}

type VaultTokenFeederResponse struct {
	Data TokenFeederResponse `json:"data"`
}

type ReceivedFeederResponse struct {
	received    TokenFeederResponse
	timeFetched time.Time
}

func (rfr *ReceivedFeederResponse) IsExpired(debounce time.Duration) bool {
	now := time.Now().Unix()
	return rfr.received.ExpiryEpoch <= now ||
		rfr.timeFetched.Unix()+int64(debounce.Seconds()) < now
}

type VaultToken string

func NewVaultTokenResourceAuthorizer(url string, token VaultToken) Authorizer {
	rv := HttpResourceFetcher{
		client: &http.Client{},
		url:    url,
		headers: map[string]string{
			"X-Vault-Token": string(token),
		},
		debounceCacheTime: time.Second * 20,
		parser:            parseVaultResponse,
	}

	return &rv
}

type HttpResourceFetcher struct {
	client *http.Client

	cachedResponse    ReceivedFeederResponse
	url               string
	headers           map[string]string
	debounceCacheTime time.Duration

	parser func([]byte, *TokenFeederResponse) error
}

func (h *HttpResourceFetcher) fetch() error {
	if req, reqErr := http.NewRequest("GET", h.url, nil); reqErr != nil {
		return reqErr
	} else {
		for hdr, hdrVal := range h.headers {
			req.Header.Set(hdr, hdrVal)
		}

		if resp, respErr := h.client.Do(req); respErr != nil {
			return respErr
		} else if resp.StatusCode > 299 {
			return errors.New(fmt.Sprintf("http resource fetcher received an unexpected code %d while attempting to read token from %s", resp.StatusCode, h.url))
		} else {
			wr := &WrappedResponse{Response: resp}
			if respData, rcvErr := wr.Body(); rcvErr != nil {
				return rcvErr
			} else {
				rv := ReceivedFeederResponse{
					timeFetched: time.Now(),
				}

				jsonErr := h.parser(respData, &rv.received)
				h.cachedResponse = rv
				return jsonErr
			}
		}
	}
}

func standardParserFunc(body []byte, resp *TokenFeederResponse) error {
	return json.Unmarshal(body, resp)
}

func parseVaultResponse(body []byte, resp *TokenFeederResponse) error {
	vaultStruct := VaultTokenFeederResponse{}
	jsonErr := json.Unmarshal(body, &vaultStruct)

	if jsonErr == nil {
		resp.Token = vaultStruct.Data.Token
		resp.Expiry = vaultStruct.Data.Expiry
		resp.ExpiryEpoch = vaultStruct.Data.ExpiryEpoch
	}

	return jsonErr
}

func (h *HttpResourceFetcher) HeaderAuthorization(ctx context.Context) (map[string]string, error) {
	rv := map[string]string{}

	// If the context is already closed; don't do anything else about it.
	if ctx.Err() != nil {
		return rv, ctx.Err()
	}

	if h.cachedResponse.IsExpired(h.debounceCacheTime) {
		if fetchErr := h.fetch(); fetchErr != nil {
			return rv, fetchErr
		}
	}

	rv["Authorization"] = fmt.Sprintf("Bearer %s", h.cachedResponse.received.Token)
	return rv, nil
}

func (h *HttpResourceFetcher) QueryStringAuthorization(_ context.Context) (map[string]string, error) {
	return nil, errors.New("unsupported operation")
}

func (h *HttpResourceFetcher) Close() {
	h.client.CloseIdleConnections()
}
