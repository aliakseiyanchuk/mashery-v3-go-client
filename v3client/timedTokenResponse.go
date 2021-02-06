package v3client

import (
	"math"
	"time"
)

// Timed access token response, suitable for storing in a log file.
type TimedAccessTokenResponse struct {
	Obtained time.Time `json:"obtained"`
	AccessTokenResponse
}

func (a AccessTokenResponse) ObtainedNow() *TimedAccessTokenResponse {
	rv := TimedAccessTokenResponse{
		AccessTokenResponse: a,
		Obtained:            time.Now(),
	}
	return &rv
}

func (t *TimedAccessTokenResponse) Expired() bool {
	now := time.Now()
	secondsDiff := now.Unix() - t.Obtained.Unix()

	if secondsDiff > int64(math.Round(float64(t.ExpiresIn)*0.95)) {
		return true
	}

	return false
}

func (t *TimedAccessTokenResponse) ExpiryTime() time.Time {
	return time.Unix(t.Obtained.Unix()+int64(t.ExpiresIn), 0)
}

// Returns number of seconds that are still left in this access tokens.
func (t *TimedAccessTokenResponse) TimeLeft() int {
	diff := t.Obtained.Unix() + int64(t.ExpiresIn) - time.Now().Unix()
	if diff > 0 {
		return int(diff)
	} else {
		return 0
	}
}
