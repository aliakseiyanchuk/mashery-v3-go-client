package transport_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCalculatingADelay(t *testing.T) {
	params := transport.HttpTransport{
		Mutex:  &sync.Mutex{},
		MaxQPS: 2,
	}

	assert.Equal(t, time.Duration(0), params.DelayBeforeCall(), "A")
	assert.Equal(t, time.Duration(0), params.DelayBeforeCall(), "B")
	assert.Equal(t, time.Second, params.DelayBeforeCall(), "C")
	assert.Equal(t, time.Second, params.DelayBeforeCall(), "D")
	assert.Equal(t, time.Second*2, params.DelayBeforeCall(), "E")
	assert.Equal(t, time.Second*2, params.DelayBeforeCall(), "F")

	time.Sleep(time.Second * 3)
	assert.Equal(t, time.Duration(0), params.DelayBeforeCall(), "G")
}
