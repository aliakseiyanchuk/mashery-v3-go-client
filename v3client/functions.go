package v3client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const AccessTokenEnv = "MASHERY_V3_TOKEN"
const AreaIdEnv = "MASHERY_AREA_ID"
const ApiKeyEnv = "MASHERY_V3API_KEY"
const ApiKeySecretEnv = "MASHERY_V3API_SECRET"
const UserNameEnv = "MASHERY_USER"
const UserPassEnv = "MASHERY_PASS"

const userSettingsFile = ".mashery-v3-credentials"

func ReadSavedV3TokenData(fp string) (*masherytypes.TimedAccessTokenResponse, error) {
	_, err := os.Stat(fp)
	if err != nil || os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s does not exist", fp))
	}

	dat, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("File %s could not be read (%s)", fp, err))
	}

	rv := masherytypes.TimedAccessTokenResponse{}
	err = json.Unmarshal(dat, &rv)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("File %s is not valid json: %s", fp, err))
	}

	return &rv, nil
}

func DefaultCredentialsFile() string {
	if u, err := os.UserHomeDir(); err == nil {
		return filepath.Join(u, userSettingsFile)
	} else {
		return userSettingsFile
	}
}

func EncryptInPlace(path string, pass string) error {
	if len(pass) != 32 {
		return errors.New("password should be 32 characters")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &errwrap.WrappedError{Context: "reading", Cause: err}
	}

	if cphr, err := aes.NewCipher([]byte(pass)); err != nil {
		return &errwrap.WrappedError{Context: "obtaining cipher", Cause: err}
	} else if gcm, err := cipher.NewGCM(cphr); err != nil {
		return &errwrap.WrappedError{Context: "failed to initialize counter", Cause: err}
	} else {
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return &errwrap.WrappedError{Context: "cannot initialize nonce", Cause: err}
		} else {
			return ioutil.WriteFile(path, gcm.Seal(nonce, nonce, data, nil), 0777)
		}
	}
}

func ReadCiphertext(fileName string, pass string) ([]byte, error) {
	if len(pass) != 32 {
		return []byte{}, errors.New("password must be exactly 32 characters")
	}

	if ciphertext, err := ioutil.ReadFile(fileName); err != nil {
		return []byte{}, &errwrap.WrappedError{Context: "reading source file", Cause: err}
	} else if chr, err := aes.NewCipher([]byte(pass)); err != nil {
		return []byte{}, &errwrap.WrappedError{Context: "initializing cipher", Cause: err}
	} else if gcmDecrypt, err := cipher.NewGCM(chr); err != nil {
		return []byte{}, &errwrap.WrappedError{Context: "initializing counter", Cause: err}
	} else {
		nonceSize := gcmDecrypt.NonceSize()
		nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]

		dat, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
		return dat, err
	}
}

func tryReadCredentialSettingsFromYamlFile(m *MasheryV3Credentials, file, pass string) {
	if _, err := os.Stat(file); err == nil || os.IsExist(err) {
		if dat, err := ReadCiphertext(file, pass); err == nil {
			var fsCreds MasheryV3Credentials
			if err = yaml.Unmarshal(dat, &fsCreds); err == nil {
				m.Inherit(&fsCreds)
			}
		} else {
			fmt.Println(err)
		}
	}
	// else: the settings file doesn't exist.
}

// DeriveAccessCredentials derives credentials from all applicable sources, including the command line
// The sequence of derivation is:
// - Environment variables, overridden by
// - User settings file, overridden by
// - Credentials file in the working directory, overriden by
// - Command line arguments
func DeriveAccessCredentials(customFile, filePass string, fallbackCreds *MasheryV3Credentials) MasheryV3Credentials {
	creds := MasheryV3Credentials{
		AreaId:   os.Getenv(AreaIdEnv),
		ApiKey:   os.Getenv(ApiKeyEnv),
		Secret:   os.Getenv(ApiKeySecretEnv),
		Username: os.Getenv(UserNameEnv),
		Password: os.Getenv(UserPassEnv),
	}

	tryReadCredentialSettingsFromYamlFile(&creds, customFile, filePass)

	if fallbackCreds != nil {
		creds.Inherit(fallbackCreds)
	}

	return creds
}
