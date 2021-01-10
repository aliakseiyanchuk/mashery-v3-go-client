package mashery_v3_go_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type FileSystemTokenProvider struct {
	Response TimedAccessTokenResponse
}

func (f FileSystemTokenProvider) AccessToken() (string, error) {
	if f.Response.Expired() {
		return "", errors.New("Token has expired")
	} else {
		return f.Response.AccessToken, nil
	}
}

func NewFileSystemTokenProvider() (V3AccessTokenProvider, error) {
	return NewFileSystemTokenProviderFrom(SavedAccessTokenFile())
}

func NewFileSystemTokenProviderFrom(path string) (V3AccessTokenProvider, error) {
	var retErr error = nil
	var retVal V3AccessTokenProvider = nil

	if _, err := os.Stat(path); err == nil || os.IsExist(err) {
		rv := FileSystemTokenProvider{}
		if data, err := ioutil.ReadFile(path); err == nil {
			if err = json.Unmarshal(data, &(rv.Response)); err == nil {
				retVal = &rv
			} else {
				retErr = errors.New(fmt.Sprintf("data of file %s could not be unmarshalled (%s)", path, err))
			}
		} else {
			retErr = errors.New(fmt.Sprintf("File %s could not be read", path))
		}
	} else {
		retErr = errors.New("saved token file does not exist")
	}

	if retErr != nil {
		return nil, retErr
	} else {
		return retVal, nil
	}
}
