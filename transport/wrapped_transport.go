package transport

import (
	"net/http"
	"sync"
)

// WrappedResponse Wraps the response so that calling applications can safely read the body multiple times.
type WrappedResponse struct {
	Response   *http.Response
	StatusCode int
	Header     http.Header

	once sync.Once

	readBody  []byte
	readError error
}

type WrappedRequest struct {
	Request *http.Request
	Body    interface{}
}

func (wr *WrappedResponse) Body() ([]byte, error) {
	wr.once.Do(func() {
		wr.readBody, wr.readError = ReadResponseBody(wr.Response)
	})

	return wr.readBody, wr.readError
}
