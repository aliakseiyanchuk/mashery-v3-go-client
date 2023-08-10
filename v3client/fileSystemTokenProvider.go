package v3client

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"os"
	"time"
)

// FileSystemTokenProvider
// An access token provider that  serve access token that has been pre-saved and persisted. The file will be periodically
// checked for the modification. The provider will retain the most recent successfully read response.
type FileSystemTokenProvider struct {
	FixedTokenProvider

	path               string
	Response           *masherytypes.TimedAccessTokenResponse
	lastFSCheck        time.Time
	sourceLastModified time.Time
	syncInterval       time.Duration
}

func (f *FileSystemTokenProvider) HeaderAuthorization(_ context.Context) (map[string]string, error) {
	f.checkFileSync()

	return map[string]string{
		"Authorization": f.header,
	}, nil
}

func (f *FileSystemTokenProvider) AccessToken(ctx context.Context) (string, error) {
	f.checkFileSync()

	if f.Response == nil {
		return "", errors.New("no saved token data found")
	} else if f.Response.Expired() {
		return "", errors.New("saved token has already expired")
	}

	return f.FixedTokenProvider.AccessToken(ctx)
}

func (f *FileSystemTokenProvider) checkFileSync() {
	now := time.Now()
	if now.Sub(f.lastFSCheck) > f.syncInterval {
		if info, err := os.Stat(f.path); (err == nil || os.IsNotExist(err)) &&
			!info.IsDir() &&
			f.sourceLastModified.Before(info.ModTime()) {

			if resp, err := ReadSavedV3TokenData(f.path); err == nil {
				f.Response = resp
				f.UpdateToken(resp.AccessToken)
			}

			f.sourceLastModified = info.ModTime()
		}
	}

	f.lastFSCheck = now
}

func (f *FileSystemTokenProvider) Close() {
	// Do nothing
}

func NewFileSystemTokenProvider() V3AccessTokenProvider {
	return NewFileSystemTokenProviderFrom(DefaultSavedAccessTokenFilePath())
}

func NewFileSystemTokenProviderFrom(path string) V3AccessTokenProvider {
	syncInterval, _ := time.ParseDuration("1m")

	return &FileSystemTokenProvider{
		path:               path,
		lastFSCheck:        time.Unix(0, 0),
		sourceLastModified: time.Unix(0, 0),
		syncInterval:       syncInterval,
	}
}
