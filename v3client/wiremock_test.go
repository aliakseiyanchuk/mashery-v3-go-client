package v3client

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type WireMock struct {
	mock.Mock
	http.Client
}

func (wm *WireMock) Do(req *http.Request) (*http.Response, error) {
	args := wm.Called(req)

	rv := args.Get(0)
	if rv != nil {
		response := rv.(*http.Response)
		response.Request = req

		return response, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

type RequestMatcher struct {
	Matchers []func(*http.Request) bool
	Response http.Response
}

func (rmb *RequestMatcher) clientForMock(wm *WireMock) (Client, *WireMock) {
	p := Params{
		MashEndpoint: "http://localhost",
		HTTPClientParams: transport.HTTPClientParams{
			ExplicitHttpExecutor: wm,
		},
	}

	return NewHttpClient(p), wm
}

func (rmb *RequestMatcher) autoRetryBadRequestClientForMock(wm *WireMock) (Client, *WireMock) {
	p := Params{
		MashEndpoint: "http://localhost",
		HTTPClientParams: transport.HTTPClientParams{
			ExplicitHttpExecutor: wm,
		},
	}

	return NewHttpClientWithBadRequestAutoRetries(p), wm
}

func (rmb *RequestMatcher) MockReturnedData() (Client, *WireMock) {
	wm := WireMock{}
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&rmb.Response, nil).
		Once()

	return rmb.clientForMock(&wm)
}

func (rmb *RequestMatcher) MockBadRequestFollowedByReturnedData() (Client, *WireMock) {
	badResponse := http.Response{
		StatusCode: 400,
		Status:     "Bad Request",
		Body:       io.NopCloser(strings.NewReader("sample 400 response")),
	}
	wm := WireMock{}
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&badResponse, nil).
		Once()
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&rmb.Response, nil).
		Once()

	return rmb.autoRetryBadRequestClientForMock(&wm)
}

func (rmb *RequestMatcher) MockReturned404() (Client, *WireMock) {
	resp404 := http.Response{
		Status:     "Not Found",
		StatusCode: 404,
	}

	wm := WireMock{}
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&resp404, nil).
		Once()

	return rmb.clientForMock(&wm)
}

func (rmb *RequestMatcher) MockReturned403() (Client, *WireMock) {
	resp404 := http.Response{
		Status:     "Not Authorized",
		StatusCode: 403,
		Body:       io.NopCloser(strings.NewReader("<h1>Not Authorized</h1>")),
	}

	wm := WireMock{}
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&resp404, nil).
		Once()

	return rmb.clientForMock(&wm)
}

func (rmb *RequestMatcher) MockReturnedDeveloperOverRate() (Client, *WireMock) {
	resp403OverRate := http.Response{
		Status:     "Developer Over Rate",
		StatusCode: 403,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("<h1>Developer Over Rate</h1>")),
	}

	resp403OverRate.Header.Set("X-Mashery-Error-Code", "ERR_403_DEVELOPER_OVER_QPS")

	wm := WireMock{}
	wm.
		On("Do", mock.MatchedBy(rmb.Match)).
		Return(&resp403OverRate, nil).
		Once()

	wm.On("Do", mock.MatchedBy(rmb.Match)).
		Return(&rmb.Response, nil).
		Once()

	return rmb.clientForMock(&wm)
}

func (rmb *RequestMatcher) Match(r *http.Request) bool {
	for _, f := range rmb.Matchers {
		if !f(r) {
			return false
		}
	}

	return true
}

type BuildVisitor func(*RequestMatcher)

func RequestMockBuilder(v BuildVisitor) *RequestMatcher {
	rv := RequestMatcher{
		Response: http.Response{
			Header:     http.Header{},
			Status:     "OK",
			StatusCode: 200,
		},
	}

	if v != nil {
		v(&rv)
	}

	return &rv
}

func (rmb *RequestMatcher) verifyPath(path string, r *http.Request) bool {
	before, _, found := strings.Cut(r.URL.Path, "?")
	rv := false
	if found {
		rv = before == path
		return rv
	} else {
		rv = path == r.URL.Path
	}

	if !rv {
		fmt.Printf("Unexpected path:\n- need: `%s`\n- have: `%s`\n", path, r.URL.Path)
	}

	return rv
}

func (rmb *RequestMatcher) ForRequestPath(path string) *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, func(request *http.Request) bool {
		return rmb.verifyPath(path, request)
	})

	return rmb
}

func (rmb *RequestMatcher) WithMethod(method string) *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, func(request *http.Request) bool {
		b := strings.ToLower(method) == strings.ToLower(request.Method)
		if !b {
			fmt.Printf("incorrect method: %s expected, %s supplied\n", method, request.Method)
		}
		return b
	})

	return rmb
}

func PayloadMatcher[T any](exp T) RequestMatcherFunc {
	return func(request *http.Request) bool {
		v := new(T)

		bodyCopy := IgnoreSupplyErr(request.GetBody)()
		b := readAllFully(bodyCopy)
		unmarshalJson(b, &v)

		equal := reflect.DeepEqual(exp, *v)
		if !equal {
			fmt.Println("Payload doesn't not match the expected")
			fmt.Println("Requested:")
			fmt.Println(string(b))
			fmt.Println("Expected")
			fmt.Println(string(marshalJson(exp)))
		}
		return equal
	}
}

func (rmb *RequestMatcher) Matching(f RequestMatcherFunc) *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, f)
	return rmb
}

func (rmb *RequestMatcher) RequestingNoFields() *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, matchMissingQueryField("fields"))
	return rmb
}

func (rmb *RequestMatcher) RequestingNoFilters() *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, matchMissingQueryField("filter"))
	return rmb
}

func matchMissingQueryField(fld string) func(request *http.Request) bool {
	return func(request *http.Request) bool {
		qs := request.URL.Query()
		if qs != nil {
			rv := !qs.Has(fld)
			if !rv {
				fmt.Printf("Request specifies query string parameter `%s`, which isn't exected here", fld)
			}
			return rv
		} else {
			return true
		}
	}
}

func (rmb *RequestMatcher) FilteredOn(key, value string) *RequestMatcher {
	return rmb.FilteringUsing(map[string]string{key: value})
}

func (rmb *RequestMatcher) FilteringUsing(filters map[string]string) *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, func(request *http.Request) bool {
		qs := request.URL.Query()
		if qs.Has("filter") {
			filterStr := qs.Get("filter")
			parsedFilters := map[string]string{}

			for _, filter := range strings.Split(filterStr, ",") {
				kv := strings.Split(filter, "=")
				parsedFilters[kv[0]] = kv[1]
			}
			return reflect.DeepEqual(filters, parsedFilters)
		} else {
			fmt.Println("Missing query parameter filter in this request")
			return false
		}
	})

	return rmb
}

func (rmb *RequestMatcher) RequestingFields(fields []string) *RequestMatcher {
	rmb.Matchers = append(rmb.Matchers, func(request *http.Request) bool {
		qs := request.URL.Query()
		if qs.Has("fields") {
			fieldsStr := strings.Split(qs.Get("fields"), ",")

			// Search fields -> passed
		outer:
			for _, f := range fields {
				for _, v := range fieldsStr {
					if f == v {
						continue outer
					}
				}

				fmt.Printf("missing key %s\n", f)
				return false
			}

			// Search fields -> passed
		outerReversed:
			for _, f := range fieldsStr {
				for _, v := range fields {
					if f == v {
						continue outerReversed
					}
				}

				fmt.Printf("unnessary key %s passed to the back-end", f)
				return false
			}

			return true
		} else {
			fmt.Println("this request should request fields, but none are specified")
			return false
		}

	})

	return rmb
}

func (rmb *RequestMatcher) WillReturnHeader(header, value string) *RequestMatcher {
	rmb.Response.Header.Set(header, value)
	return rmb
}

func (rmb *RequestMatcher) WillReturnTotalCount(num int) *RequestMatcher {
	rmb.Response.Header.Set("X-Total-Count", strconv.Itoa(num))
	return rmb
}

func (rmb *RequestMatcher) WillReturnStatus(status string, statusCode int) *RequestMatcher {
	rmb.Response.Status = status
	rmb.Response.StatusCode = statusCode

	return rmb
}

func (rmb *RequestMatcher) WillReturnBody(json string) *RequestMatcher {
	rmb.Response.Body = io.NopCloser(strings.NewReader(json))
	return rmb
}

func (rmb *RequestMatcher) WillReturnUnspecified() *RequestMatcher {
	rmb.Response.Body = io.NopCloser(strings.NewReader("[\"\"]"))
	return rmb
}

func (rmb *RequestMatcher) WillReturnJsonOf(js interface{}) *RequestMatcher {
	return rmb.WillReturnBody(string(marshalJson(js)))
}
