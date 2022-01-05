package v3client

import (
	"encoding/json"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func DefaultSavedAccessTokenFilePath() string {
	u, e := user.Current()
	if e == nil {
		return filepath.Join(u.HomeDir, tokenFile)
	} else {
		wd, e := os.Getwd()
		if e != nil {
			panic("Process with no user and working directory? How can it be")
		}
		return filepath.Join(wd, tokenFile)
	}
}

func PersistV3TokenResponse(dat *masherytypes.TimedAccessTokenResponse, path string) error {
	if stat, err := os.Stat(path); (err == nil || os.IsExist(err)) && stat.IsDir() {
		return errors.New("cannot persis a file into existing directory")
	}

	if m, err := json.Marshal(dat); err == nil {
		return ioutil.WriteFile(path, m, 0644)
	} else {
		return err
	}
}

func LoadV3TokenResponse(path string) (*masherytypes.TimedAccessTokenResponse, error) {
	if stat, err := os.Stat(path); (err == nil || os.IsExist(err)) && !stat.IsDir() {
		if dat, err := ioutil.ReadFile(path); err == nil {
			resp := masherytypes.TimedAccessTokenResponse{}
			err = json.Unmarshal(dat, &resp)
			return &resp, err
		} else {
			return nil, err
		}
	} else {
		return nil, nil
	}
}
