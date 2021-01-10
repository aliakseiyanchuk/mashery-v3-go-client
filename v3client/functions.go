package v3client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const AccessTokenEnv = "MASHERY_V3_TOKEN"
const AreaIdEnv = "MASHERY_AREA_ID"
const ApiKeyEnv = "MASHERY_V3API_KEY"
const ApiKeySecretEnv = "MASHERY_V3API_SECRET"
const UserNameEnv = "MASHERY_USER"
const UserPassEnv = "MASHERY_PASS"

const userSettingsFile = ".mashery-v3-credentials"

func ReadSavedV3TokenData(fp string) (*TimedAccessTokenResponse, error) {
	_, err := os.Stat(fp)
	if err != nil || os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s does not exist", fp))
	}

	dat, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("File %s could not be read (%s)", fp, err))
	}

	rv := TimedAccessTokenResponse{}
	err = json.Unmarshal(dat, &rv)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("File %f is not valid json: %s", fp, err))
	}

	return &rv, nil
}

func readUserHomeCredentials(m *MasheryV3Credentials) {
	if u, err := os.UserHomeDir(); err == nil {
		userSetFile := filepath.Join(u, userSettingsFile)
		tryReadCredentialSettingsFromYamlFile(m, userSetFile)
	}
}

func readCurrentProcessCredentials(m *MasheryV3Credentials) {
	if u, err := os.Getwd(); err == nil {
		userSetFile := filepath.Join(u, userSettingsFile)
		tryReadCredentialSettingsFromYamlFile(m, userSetFile)
	}
}

func tryReadCredentialSettingsFromYamlFile(m *MasheryV3Credentials, file string) {
	if _, err := os.Stat(file); err == nil || os.IsExist(err) {
		if dat, err := ioutil.ReadFile(file); err == nil {
			var fsCreds MasheryV3Credentials
			if err = yaml.Unmarshal(dat, &fsCreds); err == nil {
				m.Inherit(&fsCreds)
			}
		}
	}
	// else: the settings file doesn't exist.
}

func DeriveAccessCredentials(custFile string) MasheryV3Credentials {
	creds := MasheryV3Credentials{
		AreaId:   os.Getenv(AreaIdEnv),
		ApiKey:   os.Getenv(ApiKeyEnv),
		Secret:   os.Getenv(ApiKeySecretEnv),
		Username: os.Getenv(UserNameEnv),
		Password: os.Getenv(UserPassEnv),
	}

	readUserHomeCredentials(&creds)
	readCurrentProcessCredentials(&creds)
	tryReadCredentialSettingsFromYamlFile(&creds, custFile)

	return creds
}

func SavedAccessTokenFile() string {
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
